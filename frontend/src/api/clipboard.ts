import { ClipboardGetText, ClipboardSetText } from '../../wailsjs/runtime/runtime'

// 对齐 Tauri plugin-clipboard-manager：read 返回 Promise<string>，write 返回 Promise。
// Wails ClipboardSetText 返回 Promise<boolean>，此处忽略返回值以保持原签名。
export const clipboardApi = {
  read: (): Promise<string> => ClipboardGetText(),
  write: (text: string): Promise<boolean> => ClipboardSetText(text),
}
