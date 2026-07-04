import { GeneratePasswords } from '../../wailsjs/go/main/App'

export type PasswordOptions = {
  length: number
  upper: boolean
  lower: boolean
  digit: boolean
  symbol: boolean
  excludeAmbiguous: boolean
}

export const passwordApi = {
  generate: (opts: PasswordOptions, count: number) => GeneratePasswords(opts, count),
}
