# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目定位

**DevBox** 是一个本地开发者工具集（DevToys 风格）的 **Wails v2** 桌面应用（Go 后端 + Vue 3 前端）：左侧导航 + 右侧工具页。已实现二维码、URL、Base64（文本/图片）、JSON/SQL/XML 格式化、端口管理、连通性测试、WebSocket、哈希校验、密码、UUID 等工具页，其余导航项落到统一的 Coming Soon 占位页。

**产品名固定为 DevBox**：`wails.json#name`/`outputfilename`、`package.json#name`、桌面窗口标题 `DevBox · 开发工具箱`、macOS bundle id `com.meqln.devbox` 全部一致。改文案时不要改回 fs-wails / fs-tauri。

> 本仓库由同工作空间的 Tauri 版（`../fs-tauri/DevBox`）迁移而来：前端 UI / 样式 / 路由 / store 整体保留，后端 Rust 命令改写为 Go 方法，前端交互层从 Tauri invoke / plugin 改为 Wails binding / runtime。Go 单测（`backend_test.go`）对齐原 Rust 单测，验证行为等价。

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

## 架构概览

### 前后端分工

```
前端 frontend/src/  ──▶  src/api/*.ts  ──▶  wailsjs/go/main/App.*  ──▶  Go 方法（App struct，tools_*.go）
```

**核心契约**：工具的真实逻辑下沉到 Go，前端通过 `api/*.ts` 单一入口调用 Wails 生成的 binding。**修改行为时两侧都要改**：Go 方法签名 + 字符集 + 单测；TS 端类型与参数。绝不在前端 fallback（如 `encodeURIComponent`），那会让两端行为漂移。

例：URL 编解码的 percent-encode 字符集在 `tools_url.go` 的 `urlEncodeOne`（safe set `A-Za-z0-9-_.!~*'()`），与 JS `encodeURIComponent` 字面对齐。

### 后端结构（Wails v2）

- `main.go`：`wails.Run` 启动 + 窗口配置（标题 / 尺寸 / 浅色背景）+ `Bind: []interface{}{app}`
- `app.go`：`App` struct（含 `ctx`）+ `NewApp` + `startup` + 文件读写 binding（`ReadFile`/`ReadTextFile`/`WriteFile`）
- `tools_*.go`：按工具拆分，所有命令作为 `*App` 方法（同 `package main`）：`tools_url/base64/hash/uuid/password/qrcode/port/net/dialog.go`
- `kill_unix.go` / `kill_windows.go`：`kill_port` 的平台特定实现（build tags 拆分）
- `backend_test.go`：Go 单测，对齐原 Rust 单测

**新增命令**：在对应 `tools_*.go` 加 `func (a *App) Xxx(...) (...)` 方法，然后 `wails generate module`（或 `wails dev`/`wails build` 会自动）重新生成 `frontend/wailsjs/go/main/App.*` 与 `models.ts`。前端在 `src/api/` 加封装。

### 前端交互层映射（从 Tauri 迁移）

| 原 Tauri | 现 Wails |
|---|---|
| `invoke('cmd')` | `wailsjs/go/main/App` 的方法（`api/*.ts` 封装） |
| `plugin-clipboard-manager` | `wailsjs/runtime/runtime` 的 `ClipboardGetText`/`ClipboardSetText`（`api/clipboard.ts`） |
| `plugin-dialog` open/save | Go `OpenDialog`/`SaveDialog` binding（`api/dialog.ts` 封装，签名对齐原 open/save） |
| `plugin-fs` readFile/readTextFile/writeFile | Go `ReadFile`/`ReadTextFile`/`WriteFile` binding（`api/fs.ts` 封装，签名对齐） |
| `api/event` listen(`ping:line`) | `wailsjs/runtime/runtime` 的 `EventsOn`；Go 端 `runtime.EventsEmit` |

**`[]byte` 在 Wails v2 中以 `Array<number>` 传输**（与原 Tauri 一致），前端字节处理逻辑（`Array.from(new Uint8Array(...))` 等）无需改。

### 单向依赖

```
views
  ├─▶ components/nav (AsideNav / NavGroup / NavItem)
  ├─▶ components/ui  (Switch / PillBtn / CodeArea)
  ├─▶ stores/nav     (useNavStore — 单一全局 store)
  └─▶ api/*          (Wails binding 封装)
```

**不要反向引用**：`components/ui/*` 不知道路由 / store 存在；`stores/nav.ts` 不知道任何具体工具页。

### 路由结构

vue-router 嵌套在 `AppShell` 之下（`frontend/src/router/index.ts`）。`/tools/:id` 命中占位页 `PlaceholderView`（通过 `useNavStore.findLabel(id)` 反查工具名，未命中显示「未知工具」）。

### 样式三层

- **真源**：`frontend/src/styles/tokens.css`（26 个 CSS 变量；颜色 / 圆角 / 字体栈都在这里）
- **Tailwind**：`tailwind.config.ts` 仅暴露**常用子集**（`bg-aside`、`text-ink-2`、`rounded-md` 等）通过 `var(--xxx)` 引用 token；剩余 token 在 scoped CSS 里直接 `var()` 用
- **scoped CSS**：仅用于原型自定义元件（Switch 滑动开关、CodeArea gutter、PillBtn）

新颜色 / 圆角 / 字体先加到 `tokens.css`，再决定是否暴露到 Tailwind。Naive UI 仅作 `<n-config-provider>` + `<n-message-provider>` 容器，**不**用 `n-switch` / `n-button` 替换原型自定义组件。

### 导航数据

`frontend/src/stores/nav.ts` 中的 `NAV_DATA` 是导航真源。每项有 `id`，路由 `to` 由 id 派生 (`/tools/${id}`)。`icon` 字段引用 `src/components/nav/icons.ts` 的 ICONS 表（lucide 风 SVG inner markup，通过 `v-html` 注入）；新增图标只需在 `icons.ts` 加 key + 在 store 里写 `icon: 'newkey'`。

## 关键约束

### Wails 项目结构

- 前端在 `frontend/` 子目录（Wails 标准），`wails.json` 的 `frontend:install`/`frontend:build`/`frontend:dev:watcher` 用 pnpm。
- `main.go` 用 `//go:embed all:frontend/dist` 嵌入前端产物，故 `wails build` 前必须先有 `frontend/dist`（`wails build`/`wails dev` 会自动跑 `pnpm build`）。
- 修改 Go 方法签名后必须重新生成 binding（`wails generate module`），否则前端 `wailsjs/go/main/App.*` 过期。
- `frontend/wailsjs/` 是自动生成目录，**不要手改**。

### pnpm 11 严格策略

`frontend/pnpm-workspace.yaml` 同时含 `onlyBuiltDependencies` 与 `allowBuilds` 两份字段（pnpm 11 strict 策略要求）。**不要删任何一个** — 删了 `pnpm install` 会因 `ERR_PNPM_IGNORED_BUILDS` 失败。

### 错误处理反原则

不预先发明错误。三套策略各得其所：
- **Go `UrlDecode` 失败**：返回原文，不报错（`urlDecodeOne` 用 `url.PathUnescape`，失败回退原文，与原 Rust 一致）
- **前端 watcher binding 失败**：`try/catch` 静默，保留上次 output，**不打扰用户**
- **剪贴板失败**：才用 `n-message.error('复制失败')` toast

`frontend/src/views/UrlView.vue` 等的 watcher race token (`reqId`) 模式不要简化 — 输入连打时防止旧响应覆盖新结果。

### 第三方依赖

- `github.com/google/uuid`（v4/v7）、`github.com/skip2/go-qrcode`（QR 生成，取模块矩阵手写 SVG）、`github.com/tuotoo/qrcode`（QR 解码）、`github.com/shirou/gopsutil/v3`（端口列表 / 进程名 / 结束进程）、`golang.org/x/image`（bmp/tiff/webp 解码）、标准库 `crypto/*`、`encoding/base64`。
- **已知偏差**：QR 解码库 `tuotoo/qrcode` 与原 Rust `rqrr` 行为可能略有差异；TGA 图片格式未注册解码器（其他格式 PNG/JPEG/GIF/BMP/TIFF/WEBP 均支持）；`gopsutil` 在 macOS 获取其他用户进程的 PID 受系统权限限制（与原 netstat2 同样限制）。

### Commit message 规范

中文 + `类型: 简短描述` 格式：`feat: …`、`fix: …`、`chore: …`、`refactor: …`、`docs: …`、`test: …`。**不要**生成英文 commit message 或加 `Co-Authored-By` 之类标记 — 仓库历史不带它们。

## 文档与设计

- **UI 视觉与交互真源**：`DESIGN.md` + `prototype/index.html`（DevToys 风原型）。当原型与本文冲突时以原型为准。`prototype/` 是视觉参照，**不要修改也不要从这里直接搬代码**。
- `docs/`：从 Tauri 版继承的设计 spec / plan，技术栈部分已过时（仍写 Tauri），仅作历史参考。

## 不做项（spec 明确禁止）

国际化、自动更新、系统托盘、ESLint / Prettier、前端单测 / E2E、Tailwind 暗色配置。除非用户明确要求引入，否则不要主动加。
