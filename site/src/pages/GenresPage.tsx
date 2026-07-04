import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'
import { api } from '../api/client'
import { ExportDropdown } from '../components/ExportDropdown'
import { SuggestEditLink } from '../components/Layout'

export function GenresPage() {
  const { data: genres } = useQuery({ queryKey: ['genres'], queryFn: api.fetchGenresIndex })
  return (
    <div>
      <h1>Genre recipes</h1>
      <div className="card-grid">
        {(genres ?? []).map((g) => (
          <Link key={g.slug} className="card" to={`/genres/${g.slug}`}>
            <h3>{g.name}</h3>
            <p className="meta">{g.slug}</p>
          </Link>
        ))}
      </div>
    </div>
  )
}

export function GenreDetailPage() {
  const { slug = '' } = useParams()
  const { data: map } = useQuery({
    queryKey: ['genre', slug],
    queryFn: () => api.fetchGenre(slug),
    enabled: !!slug,
  })
  if (!map) return <p>Loading…</p>
  return (
    <div>
      <div className="detail-header">
        <h1>{map.subject.name}</h1>
        <div className="detail-actions">
          <ExportDropdown kind="genre" slug={slug} entity={map} />
          <SuggestEditLink slug={slug} kind="genre" />
        </div>
      </div>
      <p>{map.narrative.description}</p>
      <h2>Signature mechanics</h2>
      <div className="chips">
        {map.signature_gameplay.map((s) => (
          <Link key={s} className="chip" to={`/mechanics/${s}`}>
            {s}
          </Link>
        ))}
      </div>
    </div>
  )
}
