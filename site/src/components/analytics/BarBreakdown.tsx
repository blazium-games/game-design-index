import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import type { CountRow } from '../../types'
import { formatChartLabel, truncateLabel } from '../../utils/formatChartLabel'
import { CHART_THEME } from './chartTheme'

type ChartRow = CountRow & { displayLabel: string }

export function BarBreakdown({
  data,
  dataKey = 'count',
  layout = 'vertical',
  maxBars = 20,
  formatLabels = true,
}: {
  data: CountRow[]
  dataKey?: string
  layout?: 'vertical' | 'horizontal'
  maxBars?: number
  formatLabels?: boolean
}) {
  const rows: ChartRow[] = data.slice(0, maxBars).map((row) => ({
    ...row,
    displayLabel: formatLabels ? formatChartLabel(row.label) : row.label,
  }))
  if (!rows.length) return <p className="meta">No data</p>

  const vertical = layout === 'vertical'
  const longest = rows.reduce((m, r) => Math.max(m, r.displayLabel.length), 0)
  const axisWidth = vertical
    ? Math.min(280, Math.max(120, Math.ceil(longest * 6.5)))
    : 8

  return (
    <ResponsiveContainer width="100%" height={Math.max(220, rows.length * (vertical ? 36 : 24))}>
      <BarChart
        data={rows}
        layout={vertical ? 'vertical' : 'horizontal'}
        margin={{
          top: 8,
          right: 16,
          left: vertical ? axisWidth + 8 : 8,
          bottom: vertical ? 8 : 60,
        }}
      >
        <CartesianGrid stroke={CHART_THEME.grid} strokeDasharray="3 3" />
        {vertical ? (
          <>
            <XAxis
              type="number"
              stroke={CHART_THEME.text}
              tick={{ fill: CHART_THEME.text, fontSize: 11 }}
            />
            <YAxis
              type="category"
              dataKey="displayLabel"
              width={axisWidth}
              stroke={CHART_THEME.text}
              tick={{ fill: CHART_THEME.text, fontSize: 11 }}
              tickFormatter={(v: string) => truncateLabel(v, 28)}
            />
            <Bar dataKey={dataKey} fill="#7eb8ff" radius={[0, 4, 4, 0]} />
          </>
        ) : (
          <>
            <XAxis
              dataKey="displayLabel"
              stroke={CHART_THEME.text}
              tick={{ fill: CHART_THEME.text, fontSize: 10 }}
              angle={-35}
              textAnchor="end"
              height={70}
              tickFormatter={(v: string) => truncateLabel(v, 16)}
            />
            <YAxis stroke={CHART_THEME.text} tick={{ fill: CHART_THEME.text, fontSize: 11 }} />
            <Bar dataKey={dataKey} fill="#7eb8ff" radius={[4, 4, 0, 0]} />
          </>
        )}
        <Tooltip
          content={({ active, payload }) => {
            if (!active || !payload?.[0]) return null
            const row = payload[0].payload as ChartRow
            return (
              <div
                style={{
                  background: CHART_THEME.tooltipBg,
                  border: `1px solid ${CHART_THEME.tooltipBorder}`,
                  borderRadius: 6,
                  padding: '8px 10px',
                  fontSize: 12,
                  color: '#e8e6e3',
                }}
              >
                <div>{row.displayLabel}</div>
                <div>{row.count}</div>
              </div>
            )
          }}
        />
      </BarChart>
    </ResponsiveContainer>
  )
}
