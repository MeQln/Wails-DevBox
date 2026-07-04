import { HashText, HashBytes, HashFile } from '../../wailsjs/go/main/App'

export type HashResult = {
  size: number
  md5: string
  sha1: string
  sha256: string
  sha384: string
  sha512: string
}

export const hashApi = {
  text: (text: string) => HashText(text),
  bytes: (bytes: number[]) => HashBytes(bytes),
  file: (path: string) => HashFile(path),
}
