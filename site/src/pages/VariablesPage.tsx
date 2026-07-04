import { useQuery } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api/client'

export function VariablesPage() {
  const [query, setQuery] = useState('')
  const [category, setCategory] = useState('')
  const { data: rows } = useQuery({ queryKey: ['variables-index'], queryFn: api.fetchVariablesIndex })

  const filtered = useMemo(() => {
    if (!rows) return []
    return rows
      .filter((v) => !category || v.category === category)
      .filter((v) => {
        if (!query) return true
        const q = query.toLowerCase()
        return v.name.toLowerCase().includes(q) || v.slug.includes(q)
      })
  }, [rows, query, category])

  return (
    <div>
      <h1>Game variables</h1>
      <p className="meta">
        Tracked state patterns (health, currency, slots) shared across games. All catalog entries
        are listed — empty fields can be enriched via Issues.
      </p>
      <div className="filters">
        <input placeholder="Search variables…" value={query} onChange={(e) => setQuery(e.target.value)} />
        <select value={category} onChange={(e) => setCategory(e.target.value)}>
          <option value="">All categories</option>
          <option value="stat">stat</option>
          <option value="resource">resource</option>
          <option value="currency">currency</option>
          <option value="slot">slot</option>
          <option value="meter">meter</option>
          <option value="counter">counter</option>
          <option value="flag">flag</option>
        </select>
      </div>
      <p className="meta">{filtered.length} variables</p>
      <table className="table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Category</th>
            <th>Scope</th>
            <th>Maps</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {filtered.map((v) => (
            <tr key={v.slug}>
              <td>
                <Link to={`/variables/${v.slug}`}>{v.name}</Link>
              </td>
              <td>{v.category}</td>
              <td>{v.scope}</td>
              <td>{v.featured_count}</td>
              <td>
                <span className={`chip ${v.enrichment_status === 'needs_info' ? 'chip-warn' : ''}`}>
                  {v.enrichment_status === 'complete' ? 'complete' : 'needs info'}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
