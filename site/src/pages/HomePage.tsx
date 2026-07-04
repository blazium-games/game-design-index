import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { api } from '../api/client'
import { DISCORD_URL, REPO_URL } from '../types'
import { DEFAULT_DESCRIPTION, SITE_NAME } from '../seo/meta'
import { datasetJsonLd, organizationJsonLd, webSiteJsonLd } from '../seo/jsonld'
import { DocumentMeta } from '../seo/usePageMeta'

export function HomePage() {
  const { data: catalog } = useQuery({ queryKey: ['catalog'], queryFn: api.fetchCatalog })
  const { data: analytics } = useQuery({ queryKey: ['analytics'], queryFn: api.fetchAnalytics })
  const { data: maps } = useQuery({ queryKey: ['maps-index'], queryFn: api.fetchMapsIndex })

  const curated =
    maps?.filter((m) => m.quality_tier === 'curated' && m.map_type === 'game').slice(0, 12) ?? []

  const hasWebMCP = typeof document !== 'undefined' && 'modelContext' in document

  return (
    <div>
      <DocumentMeta
        title={SITE_NAME}
        description={DEFAULT_DESCRIPTION}
        path="/"
        jsonLd={[
          webSiteJsonLd(),
          organizationJsonLd(),
          datasetJsonLd({
            name: SITE_NAME,
            description: DEFAULT_DESCRIPTION,
            recordCount: catalog?.game_count,
          }),
        ]}
      />
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
          <a className="btn github" href={REPO_URL} target="_blank" rel="noreferrer">
            GitHub
          </a>
          <a className="btn discord" href={DISCORD_URL} target="_blank" rel="noreferrer">
            Join Discord
          </a>
          <a className="btn" href={`${REPO_URL}/issues/new/choose`} target="_blank" rel="noreferrer">
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
          {catalog.skill_count != null && catalog.skill_count > 0 && (
            <div className="stat">
              <Link to="/skills">
                <strong>{catalog.skill_count}</strong>
              </Link>
              <span>design skills</span>
            </div>
          )}
          {analytics?.overview.mechanic_enrichment_pct != null && (
            <div className="stat">
              <Link to="/explore/analytics">
                <strong>{analytics.overview.mechanic_enrichment_pct}%</strong>
              </Link>
              <span>mechanic enrichment</span>
            </div>
          )}
          <div className="stat">
            <strong>{catalog.genre_count}</strong>
            <span>genre recipes</span>
          </div>
          <div className="stat">
            <Link to="/changelog">
              <strong>v{catalog.release_version}</strong>
            </Link>
            <span>
              API {catalog.schema_version} · <Link to="/changelog">changelog</Link>
            </span>
          </div>
        </section>
      )}

      <section className="cta-banner">
        <h2>Explore corpus analytics</h2>
        <p>Charts, correlations, genre×domain heatmaps, and menu-flow hubs from the full index.</p>
        <div className="cta-actions">
          <Link className="btn" to="/explore/analytics">
            View analytics
          </Link>
          <Link className="btn secondary" to="/explore/cooccurrence">
            Co-occurrence table
          </Link>
        </div>
      </section>

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
