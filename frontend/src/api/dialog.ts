import { OpenDialog, SaveDialog } from '../../wailsjs/go/main/App'

type Filter = { name: string; extensions: string[] }

// 对齐 Tauri plugin-dialog 的 open/save 签名，供 views 透明替换 import 路径。
// 取消选择返回 null（Wails 端取消返回空字符串，此处归一为 null）。
// multiple 字段保留以兼容原调用，Wails OpenDialog 为单选（原项目均 multiple:false）。

type OpenOpts = { multiple?: boolean; filters?: Filter[] }

export async function openDialog(opts: OpenOpts = {}): Promise<string | null> {
  const p = await OpenDialog('', opts.filters ?? [])
  return p || null
}

type SaveOpts = { filters?: Filter[]; defaultPath?: string }

export async function saveDialog(opts: SaveOpts = {}): Promise<string | null> {
  const p = await SaveDialog('', opts.defaultPath ?? '', opts.filters ?? [])
  return p || null
}
