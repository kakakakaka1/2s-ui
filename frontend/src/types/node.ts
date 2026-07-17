// A managed node: another 2s-ui panel this one monitors over its apiv2.
// Connection config lives in the DB; live status (NodeStatus) is an in-memory
// snapshot on the server, delivered as `nodesStatus` in api/load responses.

export interface Node {
  id: number
  enable: boolean
  name: string
  baseUrl: string
  webPath: string
  // Write-only: the server never echoes the token back, it reports tokenSet.
  token?: string
  insecure: boolean
  certPin?: string
  desc?: string
  lastSeen: number
  tokenSet?: boolean
  // Phase 2 sync state
  dirty?: boolean
  lastSync?: number
}

export type NodeState = 'online' | 'offline' | 'core-stopped'

export interface NodeStatus {
  state: NodeState
  latency: number
  cpu: number
  mem: { current: number; total: number }
  appVersion: string
  coreVersion: string
  error?: string
  checkedAt: number
  lastOnline: number
}

export const defaultNode: Node = {
  id: 0,
  enable: true,
  name: '',
  baseUrl: '',
  webPath: '/app/',
  token: '',
  insecure: false,
  certPin: '',
  desc: '',
  lastSeen: 0,
}
