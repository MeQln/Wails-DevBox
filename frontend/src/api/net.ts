import { PingHost, CheckPort } from '../../wailsjs/go/main/App'

export type PortCheckResult = {
  host: string
  port: number
  open: boolean
  latency_ms: number
  message: string
}

export const netApi = {
  ping: (host: string) => PingHost(host),
  checkPort: (host: string, port: number) => CheckPort(host, port),
}
