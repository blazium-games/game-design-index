import { useQuery } from '@tanstack/react-query'
import { useMemo } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { api } from '../api/client'
import { useFilters } from '../context/Filters'
import { pageTitle } from '../seo/meta'
import { DocumentMeta } from '../seo/usePageMeta'

export function GamesPage() {
  const [params] = useSearchParams()
  const { games, setGames } = useFilters()
  const genre = params.get('genre') ?? games.genre
  const tier = params.get('tier') ?? games.tier
  const query = params.get('q') ?? games.query

  const { data: rows } = useQuery({ queryKey: ['maps-index'], queryFn: api.fetchMapsIndex })

  const filtered = useMemo(() => {
    if (!rows) return []
    return rows
      .filter((m) => m.map_type === 'game')
      .filter((m) => !genre || m.genres?.some((g) => g.toLowerCase() === genre.toLowerCase()))
      .filter((m) => !tier || m.quality_tier === tier)
      .filter((m) => {
        if (!query) return true
        const q = query.toLowerCase()
        return m.name.toLowerCase().includes(q) || m.slug.includes(q)
      })
  }, [rows, genre, tier, query])

  return (
    <div>
      <DocumentMeta
        title={pageTitle('Games')}
        description="Browse gameplay maps and signature mechanics for indexed video games."
        path="/games"
      />
      <h1>Game index</h1>
      <div className="filters">
        <input
          placeholder="Search games…"
          value={query}
          onChange={(e) => setGames({ query: e.target.value })}
        />
        <select value={tier} onChange={(e) => setGames({ tier: e.target.value })}>
          <option value="">All tiers</option>
          <option value="curated">Curated</option>
          <option value="template">Template</option>
        </select>
        <input
          placeholder="Genre filter"
          value={genre}
          onChange={(e) => setGames({ genre: e.target.value })}
        />
      </div>
      <p className="meta">{filtered.length} games</p>
      <div className="card-grid">
        {filtered.map((m) => (
          <Link key={m.slug} className="card" to={`/games/${m.slug}`}>
            <h3>{m.name}</h3>
            <p className="meta">
              {m.genres?.join(' · ')} · {m.quality_tier} · {m.signature_count} sig
            </p>
          </Link>
        ))}
      </div>
    </div>
  )
}
