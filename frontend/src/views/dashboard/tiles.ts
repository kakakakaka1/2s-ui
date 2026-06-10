export const TILES = [
  { id: 'stat-clients', group: 'grpStats', k: 'onlineClients' },
  { id: 'stat-download', group: 'grpStats', k: 'download30' },
  { id: 'stat-upload', group: 'grpStats', k: 'upload30' },
  { id: 'stat-inbounds', group: 'grpStats', k: 'activeInbounds' },
  { id: 'traffic', group: 'grpCharts', k: 'traffic' },
  { id: 'protocol', group: 'grpCharts', k: 'protocolMix' },
  { id: 'resources', group: 'grpMonitors', k: 'systemResources' },
  { id: 'network', group: 'grpMonitors', k: 'networkThroughput' },
  { id: 'server', group: 'grpInfo', k: 'server' },
  { id: 'activity', group: 'grpInfo', k: 'activity' },
]

export const TILE_GROUPS = ['grpStats', 'grpCharts', 'grpMonitors', 'grpInfo']

export const TILES_KEY = '2sui-tiles'

export function defaultTiles(): Record<string, boolean> {
  return Object.fromEntries(TILES.map((t) => [t.id, true]))
}

export function loadTiles(): Record<string, boolean> {
  try {
    const raw = localStorage.getItem(TILES_KEY)
    if (raw) return { ...defaultTiles(), ...JSON.parse(raw) }
  } catch { /* corrupted value -> defaults */ }
  return defaultTiles()
}

export function saveTiles(map: Record<string, boolean>) {
  localStorage.setItem(TILES_KEY, JSON.stringify(map))
}
