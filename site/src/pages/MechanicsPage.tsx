import { useQuery } from '@tanstack/react-query'
import { useMemo } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api/client'
import { useFilters } from '../context/Filters'

export function MechanicsPage() {
  const { mechanics, setMechanics } = useFilters()
  const { data: rows } = useQuery({ queryKey: ['mechanics-index'], queryFn: api.fetchMechanicsIndex })

  const filtered = useMemo(() => {
    if (!rows) return []
    return rows
      .filter((m) => !mechanics.domain || m.domain === mechanics.domain)
      .filter((m) => !mechanics.flavor || m.flavor === mechanics.flavor)
      .filter((m) => !mechanics.tag || m.tags?.includes(mechanics.tag))
      .filter((m) => {
        if (!mechanics.query) return true
        const q = mechanics.query.toLowerCase()
        return m.name.toLowerCase().includes(q) || m.slug.includes(q)
      })
  }, [rows, mechanics])

  return (
    <div>
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
      </div>
      <p className="meta">{filtered.length} mechanics</p>
      <div className="card-grid">
        {filtered.map((m) => (
          <Link key={m.slug} className="card" to={`/mechanics/${m.slug}`}>
            <h3>{m.name}</h3>
            <p className="meta">
              {m.domain} · {m.flavor} · {m.featured_count} maps
            </p>
          </Link>
        ))}
      </div>
    </div>
  )
}
