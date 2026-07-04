import { UrlEncode, UrlDecode } from '../../wailsjs/go/main/App'

export const urlApi = {
  encode: (text: string, multiline: boolean) => UrlEncode(text, multiline),
  decode: (text: string, multiline: boolean) => UrlDecode(text, multiline),
}
