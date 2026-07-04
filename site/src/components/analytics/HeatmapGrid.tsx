import type { GenreDomainHeatmap } from '../../types'

function cellIntensity(count: number, max: number): number {
  if (max <= 0 || count <= 0) return 0
  return 0.15 + (count / max) * 0.85
}

export function HeatmapGrid({ heatmap }: { heatmap: GenreDomainHeatmap }) {
  const { genres, domains, cells } = heatmap
  if (!genres.length || !domains.length) return <p className="meta">No heatmap data</p>

  const lookup = new Map<string, number>()
  let max = 0
  for (const c of cells) {
    lookup.set(`${c.genre}\0${c.domain}`, c.count)
    if (c.count > max) max = c.count
  }

  return (
    <div className="heatmap-wrap">
      <div
        className="heatmap-grid"
        style={{ gridTemplateColumns: `100px repeat(${domains.length}, minmax(48px, 1fr))` }}
      >
        <div className="heatmap-corner" />
        {domains.map((d) => (
          <div key={d} className="heatmap-col-head">
            {d}
          </div>
        ))}
        {genres.map((genre) => (
          <div key={genre} className="heatmap-row" style={{ display: 'contents' }}>
            <div className="heatmap-row-head">{genre}</div>
            {domains.map((domain) => {
              const count = lookup.get(`${genre}\0${domain}`) ?? 0
              const alpha = cellIntensity(count, max)
              return (
                <div
                  key={`${genre}-${domain}`}
                  className="heatmap-cell"
                  style={{ background: `rgba(126, 184, 255, ${alpha})` }}
                  title={`${genre} × ${domain}: ${count} maps`}
                >
                  <span>{count > 0 ? count : ''}</span>
                </div>
              )
            })}
          </div>
        ))}
      </div>
      <p className="meta heatmap-legend">
        Cell value = number of {heatmap.genres.length} top-genre games binding at least one mechanic in
        that domain. Darker blue = higher count (max {max}).
      </p>
    </div>
  )
}
