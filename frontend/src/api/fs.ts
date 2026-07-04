import { ReadFile, ReadTextFile, WriteFile } from '../../wailsjs/go/main/App'

// 对齐 Tauri plugin-fs 的签名，供 views 透明替换 import 路径。
// Wails ReadFile 返回 number[]，这里包成 Uint8Array 与 Tauri 行为一致。

export async function readFile(path: string): Promise<Uint8Array> {
  const arr = await ReadFile(path)
  return new Uint8Array(arr)
}

export function readTextFile(path: string): Promise<string> {
  return ReadTextFile(path)
}

export function writeFile(path: string, data: Uint8Array | number[]): Promise<void> {
  const arr = data instanceof Uint8Array ? Array.from(data) : data
  return WriteFile(path, arr)
}
