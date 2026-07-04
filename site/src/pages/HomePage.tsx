import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { api } from '../api/client'

export function HomePage() {
  const { data: catalog } = useQuery({ queryKey: ['catalog'], queryFn: api.fetchCatalog })
  const { data: maps } = useQuery({ queryKey: ['maps-index'], queryFn: api.fetchMapsIndex })

  const curated =
    maps?.filter((m) => m.quality_tier === 'curated' && m.map_type === 'game').slice(0, 12) ?? []

  const hasWebMCP = typeof document !== 'undefined' && 'modelContext' in document

  return (
    <div>
      <section className="hero">
        <h1>An open index of video game mechanics</h1>
        <p>
          MIT-licensed gameplay decomposition — browse games, reusable mechanics, genre recipes, and
          relationships.
        </p>
        {hasWebMCP && (
          <Link className="badge" to="/docs/webmcp#cursor">
            WebMCP tools available
          </Link>
        )}
      </section>

      <section className="cta-banner">
        <h2>Help us index more games</h2>
        <p>Missing titles, wrong signatures, or thin notes — we welcome Issues and Pull Requests.</p>
        <div className="cta-actions">
          <a
            className="btn"
            href="https://github.com/blazium-games/game-mechanics-index/issues/new/choose"
            target="_blank"
            rel="noreferrer"
          >
            Open an Issue
          </a>
          <Link className="btn secondary" to="/contribute">
            How to contribute
          </Link>
        </div>
      </section>

      {catalog && (
        <section className="stats">
          <div className="stat">
            <strong>{catalog.game_count}</strong>
            <span>games</span>
          </div>
          <div className="stat">
            <strong>{catalog.mechanic_count}</strong>
            <span>mechanics</span>
          </div>
          {catalog.variable_count != null && (
            <div className="stat">
              <strong>{catalog.variable_count}</strong>
              <span>variables</span>
            </div>
          )}
          {catalog.menu_count != null && (
            <div className="stat">
              <strong>{catalog.menu_count}</strong>
              <span>UI menus</span>
            </div>
          )}
          <div className="stat">
            <strong>{catalog.genre_count}</strong>
            <span>genre recipes</span>
          </div>
          <div className="stat">
            <strong>v{catalog.release_version}</strong>
            <span>API {catalog.schema_version}</span>
          </div>
        </section>
      )}

      {catalog && (
        <section>
          <h2>Browse by genre</h2>
          <div className="chips">
            {catalog.genres.slice(0, 20).map((g) => (
              <Link key={g} className="chip" to={`/games?genre=${encodeURIComponent(g)}`}>
                {g}
              </Link>
            ))}
          </div>
        </section>
      )}

      <section>
        <h2>Curated highlights</h2>
        <div className="card-grid">
          {curated.map((m) => (
            <Link key={m.slug} className="card" to={`/games/${m.slug}`}>
              <h3>{m.name}</h3>
              <p className="meta">
                {m.genres?.join(' · ')} · {m.signature_count} signatures
              </p>
            </Link>
          ))}
        </div>
      </section>
    </div>
  )
}
