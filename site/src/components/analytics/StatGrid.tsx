import type { AnalyticsOverview } from '../../types'

export function StatGrid({ overview }: { overview: AnalyticsOverview }) {
  const varPct = overview.variable_enrichment_pct ?? 0
  const menuPct = overview.menu_enrichment_pct ?? 0
  const stats = [
    { value: overview.game_count, label: 'games' },
    { value: overview.mechanic_count, label: 'mechanics' },
    { value: overview.variable_count, label: 'variables' },
    { value: overview.menu_count, label: 'UI menus' },
    { value: overview.genre_recipe_count, label: 'genre recipes' },
    { value: overview.avg_signature_count.toFixed(1), label: 'avg signatures / game' },
    { value: `${varPct}%`, label: 'variable enrichment' },
    { value: `${menuPct}%`, label: 'menu enrichment' },
  ]
  if (overview.mechanic_enrichment_pct != null) {
    stats.splice(6, 0, {
      value: `${overview.mechanic_enrichment_pct}%`,
      label: 'mechanic enrichment',
    })
  }
  if (overview.skill_enrichment_pct != null) {
    stats.splice(overview.mechanic_enrichment_pct != null ? 7 : 6, 0, {
      value: `${overview.skill_enrichment_pct}%`,
      label: 'skill enrichment',
    })
  }
  return (
    <section className="stats analytics-stats">
      {stats.map((s) => (
        <div key={s.label} className="stat">
          <strong>{s.value}</strong>
          <span>{s.label}</span>
        </div>
      ))}
    </section>
  )
}
