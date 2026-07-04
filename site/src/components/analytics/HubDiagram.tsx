import { Link } from 'react-router-dom'
import type { MenuFlowStats } from '../../types'
import { BarBreakdown } from './BarBreakdown'

export function HubDiagram({ menuFlow }: { menuFlow: MenuFlowStats }) {
  const hubs = menuFlow.top_hubs
  const maxDeg = hubs.reduce((m, h) => Math.max(m, h.total_degree), 1)

  return (
    <div className="hub-diagram">
      <div className="hub-nodes">
        {hubs.map((hub, i) => {
          const scale = 0.6 + (hub.total_degree / maxDeg) * 0.4
          const angle = (i / Math.max(hubs.length, 1)) * Math.PI * 2 - Math.PI / 2
          const radius = 100
          const cx = 120 + Math.cos(angle) * radius
          const cy = 120 + Math.sin(angle) * radius
          return (
            <Link
              key={hub.slug}
              to={`/ui-menus/${hub.slug}`}
              className="hub-node"
              style={{
                left: `${cx}px`,
                top: `${cy}px`,
                transform: `translate(-50%, -50%) scale(${scale})`,
              }}
              title={`${hub.name}: ${hub.in_degree} in, ${hub.out_degree} out`}
            >
              <span className="hub-node-name">{hub.name}</span>
              <span className="hub-node-deg">
                ↓{hub.in_degree} ↑{hub.out_degree}
              </span>
            </Link>
          )
        })}
        <svg className="hub-spokes" viewBox="0 0 240 240" aria-hidden>
          <circle cx="120" cy="120" r="100" fill="none" stroke="#2a2e38" strokeWidth="1" />
          {hubs.map((_, i) => {
            const angle = (i / Math.max(hubs.length, 1)) * Math.PI * 2 - Math.PI / 2
            const x2 = 120 + Math.cos(angle) * 100
            const y2 = 120 + Math.sin(angle) * 100
            return (
              <line
                key={i}
                x1="120"
                y1="120"
                x2={x2}
                y2={y2}
                stroke="#3a3f4b"
                strokeWidth="1"
                strokeDasharray="4 3"
              />
            )
          })}
        </svg>
      </div>
      <div className="hub-table-wrap">
        <table className="table">
          <thead>
            <tr>
              <th>Menu</th>
              <th>In</th>
              <th>Out</th>
              <th>Total</th>
            </tr>
          </thead>
          <tbody>
            {hubs.map((h) => (
              <tr key={h.slug}>
                <td>
                  <Link to={`/ui-menus/${h.slug}`}>{h.name}</Link>
                </td>
                <td>{h.in_degree}</td>
                <td>{h.out_degree}</td>
                <td>{h.total_degree}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {menuFlow.edge_counts_by_relationship.length > 0 && (
        <div className="hub-edges-chart">
          <h3>Edge relationships</h3>
          <BarBreakdown data={menuFlow.edge_counts_by_relationship} layout="horizontal" maxBars={12} />
        </div>
      )}
    </div>
  )
}
