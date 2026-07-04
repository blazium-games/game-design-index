import { useQuery } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api/client'

export function UIMenusPage() {
  const [query, setQuery] = useState('')
  const [layer, setLayer] = useState('')
  const { data: rows } = useQuery({ queryKey: ['ui-menus-index'], queryFn: api.fetchUIMenusIndex })

  const filtered = useMemo(() => {
    if (!rows) return []
    return rows
      .filter((m) => !layer || m.layer === layer)
      .filter((m) => {
        if (!query) return true
        const q = query.toLowerCase()
        return m.name.toLowerCase().includes(q) || m.slug.includes(q)
      })
  }, [rows, query, layer])

  return (
    <div>
      <h1>UI menus</h1>
      <p className="meta">
        Reusable screen patterns (pause, inventory, stage select) and how they connect to game
        state.
      </p>
      <div className="filters">
        <input placeholder="Search menus…" value={query} onChange={(e) => setQuery(e.target.value)} />
        <select value={layer} onChange={(e) => setLayer(e.target.value)}>
          <option value="">All layers</option>
          <option value="meta">meta</option>
          <option value="in_game">in_game</option>
          <option value="combat_overlay">combat_overlay</option>
        </select>
      </div>
      <p className="meta">{filtered.length} menus</p>
      <div className="card-grid">
        {filtered.map((m) => (
          <Link key={m.slug} className="card" to={`/ui-menus/${m.slug}`}>
            <h3>{m.name}</h3>
            <p className="meta">
              {m.menu_type} · {m.layer} · {m.featured_count} maps
            </p>
            <span className={`chip ${m.enrichment_status === 'needs_info' ? 'chip-warn' : ''}`}>
              {m.enrichment_status === 'complete' ? 'complete' : 'needs info'}
            </span>
          </Link>
        ))}
      </div>
    </div>
  )
}
