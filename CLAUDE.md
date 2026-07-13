# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目定位

**DevBox** 是一个本地开发者工具集（DevToys 风格）的 **Wails v2** 桌面应用（Go 后端 + Vue 3 前端）。左侧导航 + 右侧工具页结构，共 20 个已实现的工具页 + 兜底占位页。

**产品名固定为 DevBox**：`wails.json#name`/`outputfilename`、`package.json#name`、桌面窗口标题 `DevBox · 开发工具箱`、macOS bundle id `com.meqln.devbox` 全部一致。改文案时不要改回 fs-wails / fs-tauri。

## 常用命令

```bash
# 桌面应用
wails dev                     # 启动 Wails 桌面窗口（Vite dev server + Go，热重载）
wails build                   # 产出 build/bin/DevBox.app（macOS）

# 仅前端（纯 UI 调试，但 binding 调用会失败 — 桌面 IPC 需要 Wails 进程）
cd frontend && pnpm dev       # Vite dev server @ 127.0.0.1:5173
cd frontend && pnpm build     # vue-tsc -b && vite build

# 类型检查 / 单测
cd frontend && pnpm exec vue-tsc --noEmit   # 前端类型检查；提交前必须 exit 0
go test ./...                                # Go 后端单测（URL/Base64/Hash/UUID/Password）
```

**没有** ESLint / Prettier / Vitest / Playwright — 这是 spec 明确禁止的（YAGNI），不要新增。

## 构建策略

**仅构建当前平台（macOS/darwin），不交叉编译。**

- `wails build` 只在 macOS 上构建 macOS `.app` 包，不做 `-platform windows/amd64` 等交叉编译。
- 修改了 `tools_kill_windows.go` 等仅 Windows 平台相关的代码时，不要求在 macOS 上编译验证 —— 平台相关代码通过 Go build tags（`//go:build windows` / `//go:build darwin`）隔离，在 macOS 上本就不会编译。
- 如需验证 Windows 端修改，依赖 CI 或手动在 Windows 上 `wails build`。

## 架构概览

### 前后端分工

```
前端 frontend/src/  ──▶  src/api/*.ts  ──▶  wailsjs/go/main/App.*  ──▶  Go 方法（App struct，tools_*.go）
```

**核心契约**：需要跨语言一致性或调用系统能力的工具逻辑下沉到 Go，前端通过 `api/*.ts` 单一入口调用 Wails 生成的 binding。**修改 Go 行为时两侧都要改**：Go 方法签名 + 字符集 + 单测；TS 端类型与参数。绝不在前端 fallback（如 `encodeURIComponent`），那会让两端行为漂移。

**例外（纯前端工具，无 Go binding）**：JSON / SQL / XML 格式化（用 `sql-formatter` 与原生 `JSON` API）、WebSocket 客户端（浏览器 `WebSocket` API）、文本处理（转义/文本比对/Markdown）、设置页主题/配色。这些工具没有 `api/*.ts` 封装，也没有对应 `tools_*.go`，新增类似纯展示型工具时遵循同一约定。

例：URL 编解码的 percent-encode 字符集在 `tools_url.go` 的 `urlEncodeOne`（safe set `A-Za-z0-9-_.!~*'()`），与 JS `encodeURIComponent` 字面对齐。

### 后端结构（Wails v2）

- `main.go`：`wails.Run` 启动 + 窗口配置（标题 / 尺寸 / 浅色背景）+ `Bind: []interface{}{app}`
- `app.go`：`App` struct（含 `ctx`）+ `NewApp` + `startup` + 文件读写 binding（`ReadFile`/`ReadTextFile`/`WriteFile`）
- `tools_*.go`：按工具拆分，所有命令作为 `*App` 方法（同 `package main`）
- `kill_unix.go` / `kill_windows.go`：`kill_port` 的平台特定实现（build tags 拆分）
- `backend_test.go`：Go 单测，对齐原 Rust 单测

当前暴露的 `*App` 方法（修改签名后需 `wails generate module` 重新生成 binding）：

| 文件 | 方法 |
|---|---|
| `app.go` | `ReadFile` / `ReadTextFile` / `WriteFile` |
| `tools_url.go` | `UrlEncode` / `UrlDecode` |
| `tools_base64.go` | `Base64Encode` / `Base64Decode` |
| `tools_hash.go` | `HashText` / `HashBytes` / `HashFile` |
| `tools_uuid.go` | `GenerateUuids` |
| `tools_password.go` | `GeneratePasswords` |
| `tools_qrcode.go` | `QrEncode` / `QrDecode` |
| `tools_port.go` | `ListPorts` / `KillPort` |
| `tools_net.go` | `PingHost` / `CheckPort` |
| `tools_dialog.go` | `OpenDialog` / `SaveDialog` |
| `tools_image.go` | `ImageRead` / `ImageConvert` / `ImageCompress` |

**新增命令**：在对应 `tools_*.go` 加 `func (a *App) Xxx(...) (...)` 方法，然后 `wails generate module`（或 `wails dev`/`wails build` 会自动）重新生成 `frontend/wailsjs/go/main/App.*` 与 `models.ts`。前端在 `src/api/` 加封装。

### 系统托盘

系统托盘使用 CGo + macOS Cocoa API 实现（`tray_darwin.go` + `tray_darwin_helper.go`），通过 `dispatch_async` 在主线程初始化 NSStatusBar。

- **`tray.go`**：`App.initTray(ctx)` 在 `App.startup` 回调中调用，读取 `build/appicon.png` 图标，调用平台 `startTray` 并启动菜单事件分发 goroutine。
- **`tray_darwin.go`**：darwin 平台实现，`startTray` 将图标数据传给 CGo 层，`startTrayActionHandler` 启动 goroutine 监听 `actionCh` 分发显示/隐藏/退出事件。
- **`tray_darwin_helper.go`**：CGo + ObjC 实现，用 `dispatch_async` 在主线程创建 `NSStatusItem`、设置图标和菜单。菜单点击通过 `handleTrayAction` 回调发送到 Go channel。
- **`tray_other.go`**：非 darwin 平台 stub，`startTray` 和 `startTrayActionHandler` 均为空操作，防止 goroutine 泄漏。
- **`main.go`**：`HideWindowOnClose: true` 使关闭窗口时隐藏而非退出，符合托盘应用惯例。
- 托盘为跨平台行为（macOS + Windows）。macOS 菜单栏显示，Windows 通知区域显示（待实现）。

### 前端交互层（Wails binding）

| 功能 | 实现 |
|---|---|
| 后端命令调用 | `wailsjs/go/main/App` 的方法（`api/*.ts` 封装） |
| 剪贴板 | `wailsjs/runtime/runtime` 的 `ClipboardGetText`/`ClipboardSetText`（`api/clipboard.ts`） |
| 文件对话框 | Go `OpenDialog`/`SaveDialog` binding（`api/dialog.ts` 封装） |
| 文件读写 | Go `ReadFile`/`ReadTextFile`/`WriteFile` binding（`api/fs.ts` 封装） |
| 事件监听 | `wailsjs/runtime/runtime` 的 `EventsOn`；Go 端 `runtime.EventsEmit` |

**`[]byte` 在 Wails v2 中以 `Array<number>` 传输**，前端字节处理逻辑无需特殊处理。

### 工具页一览

| 分组 | 工具 | 后端 Go | 前端纯展现 |
|---|---|---|---|
| 系统工具 | 端口管理 | ✅ `ListPorts`/`KillPort` | — |
| 编解码器 | 二维码、URL、Base64 文本/图片 | ✅ | — |
| 格式化工具 | JSON、SQL、XML | — | ✅ |
| 测试工具 | 连通性测试、WebSocket | ✅ PingHost/CheckPort | WebSocket 纯前端 |
| 生成器 | 哈希校验、密码生成、UUID | ✅ | — |
| 图像处理 | 格式转换、图片压缩 | ✅ `ImageRead`/`ImageConvert`/`ImageCompress` | — |
| 文本处理 | 转义/反转义、文本比对、Markdown 预览 | — | ✅ |
| 设置 | 主题/配色切换 | — | ✅ |

### 单向依赖

```
views
  ├─▶ components/nav (AsideNav / NavGroup / NavItem)
  ├─▶ components/ui  (Switch / PillBtn)
  ├─▶ stores/nav     (useNavStore — 导航真源)
  ├─▶ stores/theme   (useThemeStore — 主题 / 明暗 / 配色)
  └─▶ api/*          (Wails binding 封装)
```

**不要反向引用**：`components/ui/*` 不知道路由 / store 存在；`stores/nav.ts` 不知道任何具体工具页；`stores/theme.ts` 只管 `data-theme`/`data-color` 与 localStorage，不依赖导航。

### 路由结构

vue-router 嵌套在 `AppShell` 之下（`frontend/src/router/index.ts`）。所有工具页均有具名路由（`/tools/:id`），兜底路由 `/tools/:id` 命中占位页 `PlaceholderView`（通过 `useNavStore.findLabel(id)` 反查工具名，未命中显示「未知工具」）。

### 样式分层

- **Token 真源**：`frontend/src/styles/tokens.css`（CSS 变量；颜色 / 圆角 / 字体栈都在这里）。浅色在 `:root`，深色在 `[data-theme='dark']` 覆盖同一组变量。配色通过 `[data-color='blue'|'purple'|'green'|'rose'|'teal'|'warm']` 选择器切换 6 套 `--aside-top`/`--border-accent`/`--accent`/`--link` 变量。
- **公用组件样式**：`frontend/src/styles/common.css`。包含 `.page-head`、`.section-title`、`.config`/`.row`、`.text-area`、`.input`、`.btn`/`.btn-primary`/`.btn-danger`、`.error-bar`、`.progress-bar` 等基础控件样式。各视图只需引用 class 名，无需重复定义。
- **Tailwind**：`tailwind.config.ts` 仅暴露**常用子集**（`bg-aside`、`text-ink-2`、`rounded-md` 等）通过 `var(--xxx)` 引用 token；剩余 token 在 scoped CSS 里直接 `var()` 用。
- **scoped CSS**：仅用于视图特有布局和自定义元件（Switch 滑动开关、PillBtn 工具按钮）的样式覆盖。

新颜色 / 圆角 / 字体先加到 `tokens.css`（浅深两套 + 6 色配色都要给），再决定是否暴露到 Tailwind。Naive UI 作容器层（`n-config-provider` + `n-message-provider` 于 `App.vue`，`n-dialog-provider` 于 `AppShell`）与主题映射（深色传 `darkTheme`），**不**用 `n-switch` / `n-button` 替换自定义组件。

### 主题与配色

`stores/theme.ts`：
- `mode`（`'light' | 'dark'`）持久化到 `localStorage`（key `devbox-theme`）
- `color`（`'blue' | 'purple' | 'green' | 'rose' | 'teal' | 'warm'`）持久化到 `localStorage`（key `devbox-color`）
- `watch` 同步 `<html data-theme>`/`<html data-color>` 与存储
- `initTheme()` 在 `main.ts` 的 `app.mount()` 之前同步调用，避免首屏闪烁
- 设置页 `SettingsView.vue` 用分段按钮切换主题、色环按钮切换配色
- 新增依赖配色的元件一律走 `var(--xxx)`，不要硬编码颜色

### 导航数据

`frontend/src/stores/nav.ts` 中的 `NAV_DATA` 是导航真源（7 个分组 + 子项），`FOOT_DATA` 是侧栏底部项（当前为 `settings`）。每项有 `id`，路由 `to` 由 id 派生 (`/tools/${id}`)。`icon` 字段引用 `src/components/nav/icons.ts` 的 ICONS 表（lucide 风 SVG inner markup，通过 `v-html` 注入）；`glyph` 字段是纯文本/emoji 占位图标。新增图标只需在 `icons.ts` 加 key + 在 store 里写 `icon: 'newkey'`。`useNavStore` 还提供 `filteredItems`（按 label/id 搜索过滤）与 `findLabel`（id → 工具名反查）。

## 关键约束

### Wails 项目结构

- 前端在 `frontend/` 子目录（Wails 标准），`wails.json` 的 `frontend:install`/`frontend:build`/`frontend:dev:watcher` 用 pnpm。
- `main.go` 用 `//go:embed all:frontend/dist` 嵌入前端产物，故 `wails build` 前必须先有 `frontend/dist`（`wails build`/`wails dev` 会自动跑 `pnpm build`）。
- 修改 Go 方法签名后必须重新生成 binding（`wails generate module`），否则前端 `wailsjs/go/main/App.*` 过期。
- `frontend/wailsjs/` 是自动生成目录，**不要手改**。
- `wails` CLI 版本必须 ≥ 2.12.0（与 `go.mod` 的 `github.com/wailsapp/wails/v2 v2.12.0` 匹配）。

### pnpm 11 严格策略

`frontend/pnpm-workspace.yaml` 同时含 `onlyBuiltDependencies` 与 `allowBuilds` 两份字段（pnpm 11 strict 策略要求）。**不要删任何一个** — 删了 `pnpm install` 会因 `ERR_PNPM_IGNORED_BUILDS` 失败。

### 错误处理反原则

不预先发明错误。三套策略各得其所：
- **Go `UrlDecode` 失败**：返回原文，不报错（`urlDecodeOne` 用 `url.PathUnescape`，失败回退原文，与原 Rust 一致）
- **前端 watcher binding 失败**：`try/catch` 静默，保留上次 output，**不打扰用户**
- **剪贴板失败**：才用 `n-message.error('复制失败')` toast

`frontend/src/views/UrlView.vue` 等的 watcher race token (`reqId`) 模式不要简化 — 输入连打时防止旧响应覆盖新结果。

### 第三方依赖

- Go：`github.com/google/uuid`（v4/v7）、`github.com/skip2/go-qrcode`（QR 生成，取模块矩阵手写 SVG）、`github.com/tuotoo/qrcode`（QR 解码）、`github.com/shirou/gopsutil/v3`（端口列表 / 进程名 / 结束进程）、`golang.org/x/image`（bmp/tiff/webp 解码）、`github.com/getlantern/systray`（系统托盘）、标准库 `crypto/*`、`encoding/base64`。
- 前端：`marked`（Markdown 渲染）、`sql-formatter`（SQL 格式化）、`naive-ui`、`pinia`、`vue-router`。
- **已知偏差**：QR 解码库 `tuotoo/qrcode` 与原 Rust `rqrr` 行为可能略有差异；ICO 格式未注册解码器（其他格式 PNG/JPEG/GIF/BMP/TIFF/WEBP 均支持）；WebP 编码暂不支持（`encodeImage` 无 WebP 编码器，前端 `FormatConversionView` 与 `ImageCompressionView` 均未暴露 WebP 选项）；`gopsutil` 在 macOS 获取其他用户进程的 PID 受系统权限限制（与原 netstat2 同样限制）。

### Commit message 规范

中文 + `类型: 简短描述` 格式：`feat: …`、`fix: …`、`chore: …`、`refactor: …`、`docs: …`、`test: …`。**不要**生成英文 commit message 或加 `Co-Authored-By` 之类标记 — 仓库历史不带它们。

## 文档与设计

- **UI 视觉与交互真源**：`DESIGN.md` + `prototype/index.html`（DevToys 风原型）。当原型与本文冲突时以原型为准。`prototype/` 是视觉参照，**不要修改也不要从这里直接搬代码**。
- `docs/`：从 Tauri 版继承的设计 spec / plan，技术栈部分已过时（仍写 Tauri），仅作历史参考。

## 不做项（spec 明确禁止）

国际化、自动更新、ESLint / Prettier、前端单测 / E2E。除非用户明确要求引入，否则不要主动加。（暗色主题、配色系统、系统托盘均已实现，不再是禁项。）