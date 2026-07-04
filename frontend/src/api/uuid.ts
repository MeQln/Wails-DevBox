import { GenerateUuids } from '../../wailsjs/go/main/App'

export const uuidApi = {
  generate: (version: number, count: number, uppercase: boolean, hyphen: boolean) =>
    GenerateUuids(version, count, uppercase, hyphen),
}
