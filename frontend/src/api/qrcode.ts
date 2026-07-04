import { QrEncode, QrDecode } from '../../wailsjs/go/main/App'

export const qrApi = {
  encode: (text: string) => QrEncode(text),
  decode: (image: number[]) => QrDecode(image),
}
