import { useQuery } from '@tanstack/react-query'
import { useMemo } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api/client'
import { EnrichmentStatusChip } from '../components/EnrichmentStatusChip'
import { useFilters } from '../context/Filters'
import { resolveEnrichmentStatus } from '../utils/enrichmentStatus'
import { pageTitle } from '../seo/meta'
import { DocumentMeta } from '../seo/usePageMeta'

export function MechanicsPage() {
  const { mechanics, setMechanics } = useFilters()
  const { data: rows } = useQuery({ queryKey: ['mechanics-index'], queryFn: api.fetchMechanicsIndex })

  const filtered = useMemo(() => {
    if (!rows) return []
    return rows
      .filter((m) => !mechanics.domain || m.domain === mechanics.domain)
      .filter((m) => !mechanics.flavor || m.flavor === mechanics.flavor)
      .filter((m) => !mechanics.tag || m.tags?.includes(mechanics.tag))
      .filter(
        (m) =>
          !mechanics.enrichmentStatus ||
          resolveEnrichmentStatus(m.enrichment_status) === mechanics.enrichmentStatus,
      )
      .filter((m) => {
        if (!mechanics.query) return true
        const q = mechanics.query.toLowerCase()
        return m.name.toLowerCase().includes(q) || m.slug.includes(q)
      })
  }, [rows, mechanics])

  return (
    <div>
      <DocumentMeta
        title={pageTitle('Mechanics')}
        description="Reusable game design mechanics catalog with design guidance and enrichment status."
        path="/mechanics"
      />
      <h1>Mechanic index</h1>
      <div className="filters">
        <input
          placeholder="Search mechanics…"
          value={mechanics.query}
          onChange={(e) => setMechanics({ query: e.target.value })}
        />
        <select value={mechanics.domain} onChange={(e) => setMechanics({ domain: e.target.value })}>
          <option value="">All domains</option>
          <option value="locomotion">locomotion</option>
          <option value="combat">combat</option>
          <option value="progression">progression</option>
          <option value="economy">economy</option>
          <option value="level">level</option>
          <option value="session">session</option>
        </select>
        <select value={mechanics.flavor} onChange={(e) => setMechanics({ flavor: e.target.value })}>
          <option value="">All flavors</option>
          <option value="action">action</option>
          <option value="adventure">adventure</option>
          <option value="strategy">strategy</option>
        </select>
        <select
          value={mechanics.enrichmentStatus ?? ''}
          onChange={(e) => setMechanics({ enrichmentStatus: e.target.value })}
        >
          <option value="">All enrichment</option>
          <option value="complete">complete</option>
          <option value="needs_info">needs info</option>
        </select>
      </div>
      <p className="meta">{filtered.length} mechanics</p>
      <table className="table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Domain</th>
            <th>Flavor</th>
            <th>Maps</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map((m) => (
            <tr key={m.slug}>
              <td>
                <Link to={`/mechanics/${m.slug}`}>{m.name}</Link>
              </td>
              <td>{m.domain}</td>
              <td>{m.flavor}</td>
              <td>{m.featured_count}</td>
              <td>
                <EnrichmentStatusChip status={m.enrichment_status} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
