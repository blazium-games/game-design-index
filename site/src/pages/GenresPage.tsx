import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'
import { api } from '../api/client'
import { ExportDropdown } from '../components/ExportDropdown'
import { SuggestEditLink } from '../components/Layout'
import { MapDetailSections } from '../components/MapDetailSections'
import { useHashScroll } from '../hooks/useHashScroll'
import { buildCanonical, pageTitle } from '../seo/meta'
import { breadcrumbJsonLd, definedTermJsonLd } from '../seo/jsonld'
import { DocumentMeta } from '../seo/usePageMeta'

export function GenresPage() {
  const { data: genres } = useQuery({ queryKey: ['genres'], queryFn: api.fetchGenresIndex })
  return (
    <div>
      <DocumentMeta
        title={pageTitle('Genres')}
        description="Genre recipe maps and common mechanic bindings."
        path="/genres"
      />
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

  useHashScroll(!!map)

  if (!map) return <p>Loading…</p>

  const canonical = buildCanonical(`/genres/${slug}`)
  const description = map.narrative.description || `Genre recipe map: ${map.subject.name}.`

  return (
    <div>
      <DocumentMeta
        title={pageTitle(map.subject.name)}
        description={description}
        path={`/genres/${slug}`}
        ogType="article"
        jsonLd={[
          definedTermJsonLd({
            name: map.subject.name,
            description,
            url: canonical,
            termSet: 'Genre Recipes',
          }),
          breadcrumbJsonLd([
            { name: 'Genres', url: buildCanonical('/genres') },
            { name: map.subject.name, url: canonical },
          ]),
        ]}
      />
      <div className="detail-header">
        <h1>{map.subject.name}</h1>
        <div className="detail-actions">
          <ExportDropdown kind="genre" slug={slug} entity={map} />
          <SuggestEditLink slug={slug} kind="genre" />
        </div>
      </div>
      <p className="meta">{map.subject.genres?.join(' · ')} · genre recipe</p>
      <MapDetailSections map={map} />
    </div>
  )
}
