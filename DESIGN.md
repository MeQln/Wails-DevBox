# DevBox UI 设计约束

> 本文是 UI 视觉与交互的单一真源。技术架构、构建、命名层级见 [`CLAUDE.md`](CLAUDE.md)。

视觉与交互的最终参照是 [`prototype/index.html`](prototype/index.html)（DevToys 风原型）。当原型与本文冲突时，**以原型为准**，并修订本文使其重新对齐。

---

## 一、视觉

颜色 / 圆角 / 字体栈集中在 `src/styles/tokens.css`（CSS 变量），是唯一真源。Tailwind 仅暴露常用子集（`bg-aside`、`text-ink-2`、`rounded-md` 等），剩余在 scoped CSS 里直接用 `var(--xxx)`。

新增 / 调整 token 时**先改 `tokens.css`**，再决定是否暴露到 Tailwind。

### 颜色

**主色表（浅色模式 `:root`，未设 `[data-color]`）**：

| 类别 | 变量 | 值 | 用途 |
|---|---|---|---|
| 背景 | `--bg` | `#f6f5f3` | 页面背景（主区外露处） |
| 表面 | `--surface` | `#ffffff` | 主区卡片、segmented 容器 |
| 侧栏 | `--aside` | `#ffffff` | 侧栏背景（始终白色） |
| 卡片 | `--card` | `#f3f2ef` | 设置项行、seg-btn 激活底、PillBtn 底 |
| 卡片 2 | `--card-2` | `#ffffff` | CodeArea 背景、工具页行内卡背景 |
| 分割线 | `--rule` | `#e3e0d9` | 主分割线（侧栏右边框、视图边框） |
| 次分割 | `--rule-soft` | `#ecebe6` | CodeArea gutter 右边框 |
| 文字主 | `--ink` | `#1d1d1f` | 主内容文字 |
| 文字 2 | `--ink-2` | `#4a4a4d` | 标题、次要文字 |
| 文字 3 | `--ink-3` | `#6f6e72` | 说明文字、搜索图标 |
| 文字 4 | `--ink-4` | `#9a9893` | 占位符、次要图标、空态提示 |
| 文字 5 | `--ink-5` | `#c2c0bb` | 滚动条滑块 |

**语义色**：

| 变量 | 值 | 用途 |
|---|---|---|
| `--amber` | `#e8a534` | 警告强调（如「有更新」bulb） |
| `--amber-d` | `#b8821e` | 警告强调深色变体 |
| `--link` | `#2769d0` | 链接文字、编号 |
| `--ok` | `#4f9a59` | 成功状态（文本比对添加行、有效徽标） |
| `--warn` | `#d97a3b` | 警告状态（文本比对删除行、错误徽标、错误栏） |
| `--danger` | `#d33` | 危险状态（当前未使用） |

**配色系统（5 色）**：通过 `[data-color]` 选择器覆盖 `--aside-top`/`--aside-2`/`--aside-3`/`--border-accent`/`--link`/`--accent` 六项变量。浅色与深色模式各 5 套完整值。

| 配色 | 默认(蓝) | 浅紫 | 绿 | 玫瑰 | 青 |
|---|---|---|---|---|---|
| `--aside-top` | `#f3f6fd` | `#f3eefb` | `#eef5ed` | `#fdf0f0` | `#edf5f3` |
| `--aside-2` | `#dce8f9` | `#ece4f7` | `#dbe8d8` | `#f5dcdc` | `#d4e8e4` |
| `--aside-3` | `#c8d8f0` | `#dcd0ee` | `#c8d8c4` | `#ecc8c8` | `#bcd8d4` |
| `--border-accent` | `#c4d4eb` | `#d4c8ec` | `#c8d4c4` | `#ecc8c8` | `#bcd4d0` |
| `--accent` | `#5b8cff` | `#8b6cf0` | `#5cb86c` | `#d96a6a` | `#5cb8b8` |

**深色模式**：通过 `[data-theme='dark']` 覆盖，侧栏保持 `--aside: #ffffff`，主区 `--bg`/`--surface` 变深，`--ink` 不变（深色侧栏外露区域使用深色背景）。深色模式下 `--aside-top` 等 6 项配色变量随 `[data-theme='dark'][data-color='xxx']` 组合选择器切换。

### 圆角

| 变量 | 值 | 用途 |
|---|---|---|
| `--r-sm` | 6px | 预览图背景、清除按钮 |
| `--r-md` | 10px | 配置卡、CodeArea、搜索框、segmented 容器 |
| `--r-lg` | 14px | （预留） |

按钮 / 菜单项独立用 `8px`（不在 token 中）以贴合原型。

### 字体

```
sans  Inter Tight, PingFang SC, Noto Sans SC, -apple-system, BlinkMacSystemFont, sans-serif
mono  JetBrains Mono, ui-monospace, Menlo, Consolas, monospace
serif Inter Tight, PingFang SC, Noto Sans SC, -apple-system, sans-serif  （H1 用）
```

不引入 woff2 字体文件，依赖系统回退。

### 过渡

- 标准 `.15s`（hover / 颜色切换 / 边框过渡）
- 滑动 `.22s cubic-bezier(.2, .7, .2, 1)`（开关滑块）
- 放大 `.15s`（色环按钮 hover: `scale(1.15)`）

---

## 二、整体布局

```
┌──────────────────────────────────────────────────────────────┐
│ window  100vw × 100vh    grid-cols-[280px 1fr]               │
├────────────┬─────────────────────────────────────────────────┤
│ AsideNav   │ main  bg-surface  flex-col  px-8 pt-[22px] pb-8 │
│ 280px      │  background: linear-gradient(90deg, #fff,       │
│ grid-rows  │   var(--aside-top))                             │
│ [auto 1fr  │   ├─ router-view  (flex-col flex-1 min-h-0)     │
│  auto]     │     │                                           │
│            │     └─ page-head → section-title → config →      │
│ background │        CodeArea × 2 (flex-1 平分剩余空间)        │
│ linear      │                                                │
│ gradient    │     工具页通用模板：                             │
│ (180deg,    │       header.page-head (H1)                    │
│  var(--aside├─ .section-title: "配置" / "输入" / "输出"      │
│   -top),    │       .config (flex-col gap-4px) ← 可选的配置区  │
│  var(--aside│       CodeArea (flex-1, 输入)                   │
│  ))         │       CodeArea (flex-1, 输出 readonly)          │
└────────────┴─────────────────────────────────────────────────┘

最小窗口  880 × 600     默认窗口  1100 × 760
```

**自适应**：宽度由 `[280px 1fr]` 网格驱动；高度由 `flex-col flex-1` 链传递至 `<router-view>` 内的所有 flex 容器，最终由两个 `<CodeArea class="flex-1">` 平分主区剩余空间——**不要给 CodeArea 加 `min-height`**，会破坏平分。图像处理视图使用 `position: relative` + 进度条覆盖。

---

## 三、组件视觉规范

### 侧栏 AsideNav

```
┌─────────────────────┐
│ 搜索框              │   auto
├─────────────────────┤
│ 导航组 / 项         │   1fr
│ ├─ g-system         │
│ ├─ g-codec          │
│ ├─ g-format         │
│ ├─ g-net            │
│ ├─ g-gen            │
│ ├─ g-img            │
│ └─ g-text           │
├─────────────────────┤
│ 设置                │   auto
└─────────────────────┘
```

- 三段：搜索框 → 导航 → footer，由 grid `[auto 1fr auto]` 行划分。
- 背景：`linear-gradient(180deg, var(--aside-top), var(--aside))`，`--rule` 强制为 `var(--border-accent)` 使侧栏边框匹配配色。
- 搜索框 `32px` 高、`transparent` 背景、`1px solid var(--aside-2)` 描边、`var(--r-md)` 圆角、`gap: 6px`。focus 时 `border-color: var(--aside-3)` + `box-shadow 0 0 0 3px rgba(0,0,0,0.04)`。
- 搜索图标（放大镜）`14×14`，色 `var(--ink-3)`。输入框 `font-size: 13px`，占位符色 `var(--ink-4)`。
- 清除按钮（叉号）`16×16`，hover 时色 `var(--ink-2)` + `bg var(--aside-2)`。
- footer 与导航之间用 `border-t` 分割，`--rule` 色（已强制为 `--border-accent`）。
- 搜索无结果时显示「未找到匹配的工具」，居中，`padding: 24px 8px`，色 `var(--ink-4)`，字号 `12.5px`。
- 搜索支持 `Cmd/Ctrl+F` 聚焦搜索框。
- 滚动条：宽 `6px`，滑块 `var(--ink-5)` 圆角 `3px`，hover `var(--ink-4)`。

### NavItem（菜单项）

```
┌──────────────────────────┐
│ [SVG/glyph]  label       │  32px, 8px 圆角
└──────────────────────────┘
```

- 高度 `32px`，圆角 `8px`，左右内边距 `6px`，上下外边距 `my-px`（1px）。
- 网格 `[22px | 1fr]`，gap `1px`：左槽放 SVG `16×16` / glyph `18×18`，右栏 label（`truncate`）。
- **active 状态**：`linear-gradient(180deg, var(--aside-2), var(--aside-3))` 背景 + `inset 0 0 0 1px rgba(0,0,0,0.04)` 内描 + `color: var(--ink)` + `font-weight: 500`。
- **hover**：`bg-aside-2`，色 `text-ink-2`。
- **glyph 占位图标**：`18×18`，4px 圆角，背景 `color-mix(in srgb, var(--ink) 5%, transparent)`，等宽字体 `9px uppercase`，色 `var(--ink-3)`。
- `hasUpdate` 字段保留（`qrcode` 项标记）但暂时无视觉呈现。

### NavGroup（分组头）

```
┌────────────────────────────────┐
│ [SVG]  label              [▼] │  36px, 8px 圆角
└────────────────────────────────┘
```

- 高度 `36px`，网格 `[22px 1fr 16px]`：左槽 SVG `16×16`、中间 label（`text-ink-3`）、右侧箭头 `<path d="M6 9l6 6 6-6"/>` `14×14`。
- 箭头过渡 `transition-transform`：展开态 `rotate(0)`，折叠态 `-rotate-90`。
- **hover**：`bg-aside-2`。
- 子项使用 `v-show="expanded"` 控制显隐，左内边距 `pl-[22px]` 缩进层级。

### 顶层节点之间的分割线

`AsideNav` 循环渲染时，当前节点 `type` 与前一个不同则插入 `<hr>` 分隔（`my-1.5 mx-1.5 border-0 border-t`）。集团与集团之间不画线。

### 图标

侧栏图标统一从 `src/components/nav/icons.ts` 的 `ICONS` 表取，lucide 风 SVG inner markup（`24×24` viewBox，stroke 风，`1.5/2px` 描边），通过 `v-html` 注入 `<svg width="16" height="16" stroke="currentColor">`。新增图标只需在 `icons.ts` 加 key，store 中相应节点写 `icon: 'newkey'`。

### Switch（滑动开关）

```
标签文本 [○───]  关闭
标签文本 [───●]  开启
```

- 容器 `row-ctl`：`display: flex; align-items: center; gap: 8px; font-size: 12.5px; color: var(--ink-3)`。
- 开关尺寸 `44 × 24px`，圆角 `999px`。
- 关闭态背景 `var(--aside-3)`，开启态背景 `var(--accent)`（随配色变化）。
- 滑块 `20 × 20px`（`top: 2px; left: 2px`），`background: var(--surface)` + `box-shadow 0 1px 2px rgba(0,0,0,0.18)`。
- 开启态滑块 `translateX(20px)`，过渡 `transform .22s cubic-bezier(.2,.7,.2,1)`，背景过渡 `.15s`。
- **事件绑定只在 `<label>`**，`<input>` 仅作可访问性挂载（`tabindex="-1"; pointer-events: none; opacity: 0`）；不要给 input 加 `@click`，会导致 label-input 自动 click 链路双触发。

### PillBtn（工具按钮）

- 高度 `30px`，圆角 `8px`，左右内边距 `10px`（icon-only 变体宽 `32px`，`padding: 0`，居中）。
- 背景 `var(--card)`，hover `var(--aside-2)`，active `var(--aside-3)`。
- 内嵌 SVG 统一 `14 × 14px`。
- `:deep(svg)` 穿透生效，确保所有 slot 内 svg 尺寸一致。
- `:disabled` 态：`opacity: .45; cursor: not-allowed; pointer-events: none`。

### CodeArea（输入 / 输出区）

```
┌─────────────────────────────────────┐
│  1 │ text content here...            │  --rule 边框
│  2 │ more text                       │  --r-md 圆角
│  3 │                                 │  --card-2 背景
└─────────────────────────────────────┘
```

- 容器 `.io`：`border: 1px solid var(--rule)`、圆角 `var(--r-md)`、背景 `var(--card-2)`、网格 `[56px 1fr]`、`overflow: hidden`、`margin-bottom: 12px`。
- gutter（左行号列）：`background: linear-gradient(180deg, transparent, rgba(0,0,0,0.012))`、右边框 `1px solid var(--rule-soft)`、`padding: 12px 8px`、右对齐、字体 `var(--mono) 13px` / `line-height: 1.85`、字色 `var(--link)`、`white-space: pre`、`user-select: none`、`overflow: hidden`。
- content（右内容列）：`padding: 12px 14px`、字体 `var(--mono) 13px` / `line-height: 1.85`、`white-space: pre-wrap`、`word-break: break-all`、`overflow: auto`。
- `readonly === true` 渲染 `<pre>`、否则 `<textarea>`（不统一为 `textarea readonly`，贴合原型 HTML 结构）。

### 配置卡片（config / row）

**工具页配置区**（通用模板）：

```html
<div class="section-title"><span>配置</span></div>
<div class="config">
  <div class="row">...</div>
</div>
```

`.config`：
- 背景 `color-mix(in srgb, var(--aside-2) 6%, var(--card-2))`、`1px solid var(--border-accent)`、圆角 `var(--r-md)`、`padding: 6px`、`flex-col gap: 4px`。
- **设置页例外**：`.config` 背景 `transparent`，`.row` 背景 `transparent` + `1px solid var(--border-accent)`（带边框的行）。

`.row`：
- 网格 `[44px 1fr auto]`、`gap: 12px`、`min-height: 64px`、`padding: 14px 16px`、圆角 `8px`、`box-shadow: 0 1px 0 rgba(0,0,0,0.02)`。
- 背景 `var(--card-2)`（工具页）/ `transparent`（设置页）。
- `.row-icon`：`22×22px`、居中 flex、`color: var(--ink-2)`、内嵌 SVG `18×18`。
- `.row-title`：`font-size: 14px; font-weight: 500`。
- `.row-desc`：`font-size: 12.5px; color: var(--ink-3); margin-top: 2px`。

### 分段按钮（segmented）

```
┌──────────────────────────────┐
│ ┌──────┐ ┌──────┐ ┌──────┐ │  --aside 背景
│ │ 浅色  │ │ 深色  │ │      │ │  --r-md 圆角
│ └──────┘ └──────┘ └──────┘ │  padding: 3px; gap: 2px
└──────────────────────────────┘
```

- 容器：`display: inline-flex`、`background: var(--aside)`、圆角 `var(--r-md)`、`padding: 3px`、`gap: 2px`。
- 按钮：`padding: 6px 14px`、圆角 `7px`、`font-size: 13px`、`color: var(--ink-3)`、过渡 `background .15s color .15s box-shadow .15s`。
- hover `color: var(--ink)`。
- active：`background: var(--card)`、`color: var(--ink)`、`box-shadow: 0 1px 2px rgba(0,0,0,0.08)`。
- SVG 图标 `14×14`，`gap: 6px` 图文间距。

### 配色选择器（color-picker）

```
  ○  ○  ●  ○  ○    28×28 圆形色钮
```

- 容器：`display: inline-flex; gap: 8px; align-items: center`。
- 色钮：`28×28px`、`border-radius: 50%`、`background: var(--swatch)`（通过 `--swatch` 自定义属性传入）。
- 默认 `border: 2px solid transparent`。
- hover：`transform: scale(1.15)`。
- active：`border-color: var(--ink)` + `box-shadow: 0 0 0 2px var(--surface), 0 0 0 4px var(--ink)`（双层环，内环 `--surface` 作过渡）。

### 格式切换按钮组（fmt-group / type-group / view-group）

```
┌───────────────┐
│ [PNG][JPEG]…  │  间距 gap: 2px/4px
└───────────────┘
```

- 容器：`display: flex; gap: 2px`（type-group/view-group）或 `gap: 4px`（fmt-group），背景 `color-mix(in srgb, var(--aside-2) 12%, transparent)`、圆角 `6px`、`padding: 2px`。
- 按钮：`padding: 4px 12px`（type/view）或 `min-width: 52px; padding: 5px 8px`（fmt）、圆角 `4px`、`background: transparent`、`color: var(--ink-2)`、`font-size: 12.5px`。
- active：背景 `var(--card-2)`（type/view）或 `var(--accent)` + `#fff` + `border-color: var(--accent)`（fmt）、`box-shadow: 0 1px 2px rgba(0,0,0,0.06)`、`color: var(--ink)`。
- hover（非 active）：`color: var(--ink)`。

### 进度条（progress-bar）

```
▌▌▌▌▌▌▌▌▌▌▌▌▌   3px 高, --accent 色
```

- `position: absolute; top: 0; left: 0; right: 0; height: 3px; z-index: 10`。
- `background: linear-gradient(90deg, var(--accent) 30%, transparent 30%)` / `background-size: 200% 100%`。
- 动画 `1.2s ease infinite`：`background-position: 200% 0` → `-200% 0`。

### 通用按钮（btn）

- `padding: 7px 16px; border: 1px solid var(--border-accent); border-radius: var(--r-md)`。
- 背景 `transparent`、色 `var(--ink)`、`font-size: 13px`、`white-space: nowrap`。
- `:disabled`：`opacity: 0.5; cursor: not-allowed`。
- hover：`background: color-mix(in srgb, var(--aside-2) 10%, transparent)`。
- `btn-primary`：`background: var(--accent); color: #fff; border-color: var(--accent)`，hover `opacity: 0.85`。

### 页头（page-head）

```html
<header class="page-head"><h1>工具名</h1></header>
```

- `display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 18px`。
- `<h1>`：`font-family: var(--serif); font-size: 28px; font-weight: 500; letter-spacing: -0.015em`。

### 节标题（section-title）

```html
<div class="section-title"><span>配置</span><div class="section-actions">…</div></div>
```

- `display: flex; align-items: center; justify-content: space-between`。
- `font-size: 13.5px; font-weight: 500; color: var(--ink-2); margin: 12px 0 8px`（配置栏下 `margin: 6px 0 4px`）。
- `.section-actions`：`display: flex; gap: 4px; align-items: center`。

### 有效性徽标（badge）

```
输入 [有效]
输入 [无效]
```

- 位于 JSON 视图的 `section-title` 内。
- `.badge`：`font-size: 11px; font-weight: 600; padding: 2px 8px; border-radius: 999px; letter-spacing: 0.02em`。
- `.badge-ok`：`background: color-mix(in srgb, var(--ok) 14%, transparent); color: var(--ok)`。
- `.badge-err`：`background: color-mix(in srgb, var(--warn) 14%, transparent); color: var(--warn)`。

### Toast / Message

仅用 `useMessage()` 的 `success / error`，不引入 `useDialog` / `useNotification`。文案：「已复制」 / 「复制失败」。

---

## 四、文案规范

- 全中文（原型已确定的中英混排短语保留，如「Encoding / Decoding Multiline」）。
- 标题用全角斜杠 `/` 与半角空格组合，如「URL 编码 / 解码工具」「转义 / 反转义工具」「Markdown 预览工具」。
- 「即将上线」用于占位页（PlaceholderView）——仅 `image-format`/`image-compress`/`escape`/`list-cmp`/`md` 之外的未知路由才命中。
- **产品名固定 `DevBox`**：窗口标题 `DevBox · 开发工具箱`、页内不出现「fs-tauri」「fs-wails」字样。

---

## 五、交互约束

- 单输入 → 单输出，**输出区只读**（`<CodeArea readonly>`），用户不可编辑。
- 输入 / 开关变化触发 watcher，结果立即写回输出区（无防抖；本地 invoke 往返 < 1ms 不需要）。
- 导航项**点击都跳转**到 `/tools/:id`；`AsideNav` 的搜索框过滤、PillBtn 的「粘贴 / 读取文件 / 清空 / 复制」等按功能触发；进度条在图片处理时显示。
- 搜索框 `Esc` 清空 `query`；`Cmd/Ctrl+F` 聚焦搜索框。
- 文本比对：输入 A / B 自动触发 LCS 比对，支持交换 A/B 和复制结果。逐行模式含行内字符高亮，全文逐字模式全文字符级比对。
- Markdown 预览：分栏/编辑/预览三模式切换，支持粘贴 Markdown 和复制渲染后的 HTML。
- 图片处理：选择源文件 → 读图片信息 → 选择目标格式/质量 → 保存对话框 → 后端处理 → 结果通知。

### 错误表现

| 错误源 | 表现 |
|---|---|
| Go decode 非法输入 | 返回原文（用户视角相当于"无变化"） |
| 前端 watcher invoke 失败 | 静默，保留上次输出（不弹 toast，不打扰） |
| JSON 解析失败 | 保留上次有效输出，显示错误文本，标注「无效」徽标 |
| unescapeUnicode 非法代理对 | `String.fromCharCode` 处理（不崩溃） |
| 剪贴板权限 / 失败 | `n-message.error('复制失败')` 右上角 toast |
| 图片读取/转换/压缩失败 | `n-message.error('失败原因')` + 状态栏红色提示 |
| 深色模式 | `--aside` 保持白色（侧栏始终浅色），主区背景加深 |

详细策略与代码位置见 `CLAUDE.md`。

---

## 六、不做项（UI 范畴）

| 项 | 理由 |
|---|---|
| 国际化 i18n | 仅中文 |
| 响应式 / 移动端断点 | 桌面应用，最小窗口 880 已限定 |
| 动画曲线之外的过渡（弹簧、惯性等） | 已用 `cubic-bezier(.2, .7, .2, 1)` 与 `.15s` 标准 |
| ARIA 完整无障碍 | 仅基础 `aria-label` / `role`；不为屏幕阅读器做完整支持 |

---

## 七、修订规则

UI 改动时：

1. 改原型 `prototype/index.html` 或本文，二者**不能**同时落后于代码。
2. 颜色 / 圆角 / 字体先进 `tokens.css`，再决定 Tailwind 暴露。
3. 自定义元件（Switch / Pill / CodeArea / NavItem 等）保留 scoped CSS；不强求 Tailwind 化。
4. 新增 Naive UI 组件需评估是否破坏「原型优先」基线——默认**不**用 `n-switch` / `n-button` / `n-input`。
5. 配色系统变更时需同时覆盖 `:root`、`[data-color='xxx']`、`[data-theme='dark']`、`[data-theme='dark'][data-color='xxx']` 四层选择器。