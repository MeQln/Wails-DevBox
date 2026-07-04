import { ListPorts, KillPort } from '../../wailsjs/go/main/App'

export type PortEntry = {
  port: number
  pid: number
  process_name: string
  address: string
}

export const portApi = {
  list: () => ListPorts(),
  kill: (pid: number) => KillPort(pid),
}
