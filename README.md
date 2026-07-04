# DevBox · 开发工具箱

本地开发者工具集（DevToys 风格）桌面应用：左侧导航 + 右侧工具页。基于 **Wails v2（Go 后端 + Vue 3 前端）** 构建。

> 本仓库由同工作空间的 Tauri 版（`fs-tauri/DevBox`）迁移而来，前端 UI / 样式 / 路由 / store 整体保留，后端由 Rust 改写为 Go，交互层从 Tauri invoke / plugin 改为 Wails binding / runtime。

## 功能工具

二维码、URL 编解码、Base64（文本 / 图片）、JSON / SQL / XML 格式化、端口管理、连通性测试（Ping / TCP 端口）、WebSocket、哈希校验（MD5/SHA1/SHA256/SHA384/SHA512）、密码生成、UUID（v4/v7）。

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | Vue 3 + Pinia + Vue Router + Naive UI + Tailwind CSS + Vite |
| 后端 | Go + Wails v2 |
| 通信 | Wails binding（`frontend/wailsjs/go/main/App`）+ runtime（事件 / 剪贴板） |

## 常用命令

```bash
# 桌面应用（开发模式，热重载）
wails dev

# 产出安装包到 build/bin/
wails build

# 仅前端（纯 UI 调试，但 binding 调用会失败 — 桌面 IPC 需要 Wails 进程）
cd frontend && pnpm dev

# 类型检查 / 单测
cd frontend && pnpm exec vue-tsc --noEmit       # 前端类型检查
go test ./...                                    # Go 后端单测（对齐原 Rust 单测）
```

## 架构

```
前端 frontend/src/  ──▶  src/api/*.ts  ──▶  wailsjs/go/main/App.*  ──▶  Go 方法（App struct，tools_*.go）
```

- **后端命令**挂在 `App` struct（`app.go` + `tools_*.go`），按工具拆分文件，同 `package main`。
- **前端 API 层** `src/api/*.ts` 是单一调用入口，内部转发到 Wails 生成的 binding，对外签名与原 Tauri 版一致。
- **剪贴板 / 事件**用本地 `wailsjs/runtime/runtime`；**文件对话框**由 Go 端 `OpenDialog`/`SaveDialog` binding 包装；**文件读写**由 `ReadFile`/`ReadTextFile`/`WriteFile` binding 提供，`src/api/fs.ts` 封装对齐原 Tauri plugin-fs 签名。

更多架构与约束见 [`CLAUDE.md`](CLAUDE.md)，UI 视觉与交互见 [`DESIGN.md`](DESIGN.md)。
