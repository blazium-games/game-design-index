export function InsightCallout({ insights }: { insights: string[] }) {
  if (!insights.length) return null
  return (
    <section className="chart-card">
      <h2>Correlation insights</h2>
      <p className="meta">
        Auto-generated callouts from corpus statistics — use these as starting points for deeper
        exploration.
      </p>
      <ol className="insight-list">
        {insights.map((text, i) => (
          <li key={i}>{text}</li>
        ))}
      </ol>
    </section>
  )
}
