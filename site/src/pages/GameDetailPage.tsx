import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'
import { api } from '../api/client'
import { SuggestEditLink } from '../components/Layout'

export function GameDetailPage() {
  const { slug = '' } = useParams()
  const { data: map } = useQuery({
    queryKey: ['map', slug],
    queryFn: () => api.fetchMap(slug),
    enabled: !!slug,
  })

  if (!map) return <p>Loading…</p>

  return (
    <div>
      <div className="detail-header">
        <h1>{map.subject.name}</h1>
        <SuggestEditLink slug={slug} kind="game" />
      </div>
      <p className="meta">
        {map.subject.genres?.join(' · ')} · {map.metadata?.quality_tier ?? 'template'}
      </p>
      <section>
        <h2>Description</h2>
        <p>{map.narrative.description}</p>
        {map.narrative.core_loop && (
          <>
            <h3>Core loop</h3>
            <p>{map.narrative.core_loop}</p>
          </>
        )}
      </section>
      <section>
        <h2>Signature mechanics</h2>
        <div className="chips">
          {map.signature_gameplay.map((s) => (
            <Link key={s} className="chip" to={`/mechanics/${s}`}>
              {s}
            </Link>
          ))}
        </div>
      </section>
      <section>
        <h2>Mechanic bindings</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Mechanic</th>
              <th>Role</th>
              <th>Domain</th>
              <th>Notes</th>
            </tr>
          </thead>
          <tbody>
            {map.mechanics.map((b) => (
              <tr key={b.mechanic_slug}>
                <td>
                  <Link to={`/mechanics/${b.mechanic_slug}`}>{b.mechanic_slug}</Link>
                </td>
                <td>{b.role}</td>
                <td>{b.domain}</td>
                <td>{b.map_notes || b.expression}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
      {map.design_intent?.theme_tags && (
        <section>
          <h2>Theme tags</h2>
          <div className="chips">
            {map.design_intent.theme_tags.map((t) => (
              <span key={t} className="chip">
                {t}
              </span>
            ))}
          </div>
        </section>
      )}
    </div>
  )
}
