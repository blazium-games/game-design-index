import {
  CartesianGrid,
  ResponsiveContainer,
  Scatter,
  ScatterChart,
  Tooltip,
  XAxis,
  YAxis,
  ZAxis,
} from 'recharts'
import { Link } from 'react-router-dom'
import type { CooccurrenceAnalyticsRow } from '../../types'
import { formatSlugLabel } from '../../utils/formatChartLabel'
import { CHART_THEME } from './chartTheme'

export function ScatterLift({ pairs }: { pairs: CooccurrenceAnalyticsRow[] }) {
  if (!pairs.length) return <p className="meta">No co-occurrence data</p>

  const data = pairs.map((p) => ({
    x: p.count,
    y: p.lift,
    z: p.lift,
    label: `${p.mechanic_a} + ${p.mechanic_b}`,
    mechanic_a: p.mechanic_a,
    mechanic_b: p.mechanic_b,
  }))

  const topLift = [...pairs].sort((a, b) => b.lift - a.lift).slice(0, 5)

  return (
    <div>
      <ResponsiveContainer width="100%" height={300}>
        <ScatterChart margin={{ top: 16, right: 16, bottom: 8, left: 8 }}>
          <CartesianGrid stroke={CHART_THEME.grid} strokeDasharray="3 3" />
          <XAxis
            type="number"
            dataKey="x"
            name="Count"
            stroke={CHART_THEME.text}
            tick={{ fill: CHART_THEME.text, fontSize: 11 }}
            label={{
              value: 'How often pairs appear together',
              position: 'bottom',
              fill: CHART_THEME.text,
              fontSize: 11,
            }}
          />
          <YAxis
            type="number"
            dataKey="y"
            name="Lift"
            stroke={CHART_THEME.text}
            tick={{ fill: CHART_THEME.text, fontSize: 11 }}
            label={{
              value: 'How much more common than chance',
              angle: -90,
              position: 'insideLeft',
              fill: CHART_THEME.text,
              fontSize: 11,
            }}
          />
          <ZAxis type="number" dataKey="z" range={[40, 200]} />
          <Tooltip
            content={({ payload }) => {
              if (!payload?.[0]) return null
              const d = payload[0].payload as (typeof data)[0]
              return (
                <div
                  style={{
                    background: CHART_THEME.tooltipBg,
                    border: `1px solid ${CHART_THEME.tooltipBorder}`,
                    padding: '8px 10px',
                    borderRadius: 6,
                    fontSize: 12,
                  }}
                >
                  <div>
                    <Link to={`/mechanics/${d.mechanic_a}`}>{formatSlugLabel(d.mechanic_a)}</Link>
                    {' + '}
                    <Link to={`/mechanics/${d.mechanic_b}`}>{formatSlugLabel(d.mechanic_b)}</Link>
                  </div>
                  <div>Maps together: {d.x}</div>
                  <div>Lift: {d.y.toFixed(2)}× expected</div>
                </div>
              )
            }}
          />
          <Scatter data={data} fill="#7eb8ff" fillOpacity={0.75} />
        </ScatterChart>
      </ResponsiveContainer>
      <p className="meta">
        Lift = observed co-occurrence ÷ expected if independent. Points top-right are frequent and
        over-represented together.
      </p>
      <ul className="scatter-top-lift">
        {topLift.map((p) => (
          <li key={`${p.mechanic_a}-${p.mechanic_b}`}>
            <Link to={`/mechanics/${p.mechanic_a}`}>{formatSlugLabel(p.mechanic_a)}</Link>
            {' + '}
            <Link to={`/mechanics/${p.mechanic_b}`}>{formatSlugLabel(p.mechanic_b)}</Link>
            {' — '}
            {p.lift.toFixed(2)}× expected ({p.count} maps)
          </li>
        ))}
      </ul>
    </div>
  )
}
