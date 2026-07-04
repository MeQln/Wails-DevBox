import { Base64Encode, Base64Decode } from '../../wailsjs/go/main/App'

export const base64Api = {
  encode: (text: string, urlSafe: boolean) => Base64Encode(text, urlSafe),
  decode: (text: string, urlSafe: boolean) => Base64Decode(text, urlSafe),
}
