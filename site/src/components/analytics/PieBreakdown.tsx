import { Cell, Legend, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts'
import type { CountRow } from '../../types'
import { formatChartLabel } from '../../utils/formatChartLabel'
import { CHART_COLORS, CHART_THEME } from './chartTheme'

type ChartDatum = { name: string; displayName: string; value: number }

export function PieBreakdown({
  data,
  totalLabel,
}: {
  data: CountRow[]
  totalLabel?: string
}) {
  if (!data.length) return <p className="meta">No data</p>

  const chartData: ChartDatum[] = data.map((d) => ({
    name: d.label,
    displayName: formatChartLabel(d.label),
    value: d.count,
  }))
  const total = chartData.reduce((s, d) => s + d.value, 0)
  const dominant = chartData.reduce((a, b) => (b.value > a.value ? b : a), chartData[0])
  const dominantPct = total > 0 ? (dominant.value / total) * 100 : 0

  if (dominantPct >= 95) {
    return (
      <div className="chart-pie-wrap">
        <p className="chart-stat-callout">
          {totalLabel ??
            `${total} items — ${dominant.displayName.toLowerCase()} (${dominantPct.toFixed(0)}% of total)`}
        </p>
      </div>
    )
  }

  return (
    <div className="chart-pie-wrap">
      <ResponsiveContainer width="100%" height={300}>
        <PieChart margin={{ top: 8, right: 8, bottom: 8, left: 8 }}>
          <Pie
            data={chartData}
            dataKey="value"
            nameKey="displayName"
            cx="50%"
            cy="45%"
            innerRadius={45}
            outerRadius={70}
            paddingAngle={2}
          >
            {chartData.map((_, i) => (
              <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />
            ))}
          </Pie>
          <Tooltip
            content={({ active, payload }) => {
              if (!active || !payload?.[0]) return null
              const row = payload[0].payload as ChartDatum
              const value = Number(payload[0].value ?? 0)
              const pct = total > 0 ? ((value / total) * 100).toFixed(0) : '0'
              return (
                <div
                  style={{
                    background: CHART_THEME.tooltipBg,
                    border: `1px solid ${CHART_THEME.tooltipBorder}`,
                    padding: '8px 10px',
                    fontSize: 12,
                    color: '#e8e6e3',
                  }}
                >
                  <div>{row.displayName}</div>
                  <div>
                    {value} ({pct}%)
                  </div>
                </div>
              )
            }}
          />
          <Legend
            layout="horizontal"
            verticalAlign="bottom"
            wrapperStyle={{ fontSize: '0.85rem', color: CHART_THEME.text, paddingTop: 12 }}
            formatter={(value: string, entry) => {
              const row = entry.payload as ChartDatum | undefined
              const count = row?.value ?? 0
              const pct = total > 0 ? ((count / total) * 100).toFixed(0) : '0'
              return `${value} · ${count} (${pct}%)`
            }}
          />
        </PieChart>
      </ResponsiveContainer>
    </div>
  )
}
