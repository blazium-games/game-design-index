import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'
import { api } from '../api/client'
import { EmptyField } from '../components/EmptyField'
import { ExportDropdown } from '../components/ExportDropdown'
import { SuggestEditLink } from '../components/Layout'

export function VariableDetailPage() {
  const { slug = '' } = useParams()
  const { data: v } = useQuery({
    queryKey: ['variable', slug],
    queryFn: () => api.fetchVariable(slug),
    enabled: !!slug,
  })
  const { data: varMaps } = useQuery({
    queryKey: ['variable-maps', slug],
    queryFn: async () => {
      const idx = await api.fetchVariableToMaps()
      return idx.variables[slug] ?? []
    },
    enabled: !!slug,
  })
  const { data: varMechs } = useQuery({
    queryKey: ['variable-mechanics', slug],
    queryFn: async () => {
      const idx = await api.fetchVariableToMechanics()
      return idx.variables[slug] ?? []
    },
    enabled: !!slug,
  })
  const { data: menuIndex } = useQuery({
    queryKey: ['ui-menus-index'],
    queryFn: api.fetchUIMenusIndex,
  })
  const { data: relatedMenuDetails } = useQuery({
    queryKey: ['menus-for-variable', slug, menuIndex?.length],
    queryFn: async () => {
      if (!menuIndex) return []
      const menus = await Promise.all(menuIndex.map((row) => api.fetchUIMenu(row.slug)))
      return menus.filter((m) => m.related_variables?.includes(slug))
    },
    enabled: !!slug && !!menuIndex,
  })

  if (!v) return <p>Loading…</p>

  return (
    <div>
      <div className="detail-header">
        <h1>{v.name}</h1>
        <div className="detail-actions">
          <ExportDropdown kind="variable" slug={slug} entity={v} />
          <SuggestEditLink slug={slug} kind="variable" />
        </div>
      </div>
      <p className="meta">
        {v.category} · {v.scope} · {v.value_kind}
        {v.reset_behavior ? ` · resets ${v.reset_behavior}` : ''}
      </p>
      <section>
        <h2>Summary</h2>
        <p>{v.summary}</p>
      </section>
      <section>
        <h2>Shared rationale</h2>
        {v.shared_rationale ? <p>{v.shared_rationale}</p> : <EmptyField slug={slug} field="shared_rationale" kind="variable" />}
      </section>
      <section>
        <h2>Player focus</h2>
        {v.player_focus ? <p>{v.player_focus}</p> : <EmptyField slug={slug} field="player_focus" kind="variable" />}
      </section>
      <section>
        <h2>Typical range</h2>
        {v.typical_range ? <p>{v.typical_range}</p> : <EmptyField slug={slug} field="typical_range" kind="variable" />}
      </section>
      {v.tags && v.tags.length > 0 && (
        <section>
          <h2>Tags</h2>
          <div className="chips">
            {v.tags.map((t) => (
              <span key={t} className="chip">
                {t}
              </span>
            ))}
          </div>
        </section>
      )}
      <section>
        <h2>Related mechanics</h2>
        {(varMechs ?? v.related_mechanics ?? []).length > 0 ? (
          <div className="chips">
            {(varMechs ?? v.related_mechanics ?? []).map((m) => (
              <Link key={m} className="chip" to={`/mechanics/${m}`}>
                {m}
              </Link>
            ))}
          </div>
        ) : (
          <p className="meta">No mechanic links yet.</p>
        )}
      </section>
      {v.related_variables && v.related_variables.length > 0 && (
        <section>
          <h2>Related variables</h2>
          <div className="chips">
            {v.related_variables.map((rel) => (
              <Link key={rel.slug} className="chip" to={`/variables/${rel.slug}`}>
                {rel.slug}
                {rel.relationship ? ` (${rel.relationship})` : ''}
              </Link>
            ))}
          </div>
        </section>
      )}
      {relatedMenuDetails && relatedMenuDetails.length > 0 && (
        <section>
          <h2>Related UI menus</h2>
          <div className="chips">
            {relatedMenuDetails.map((m) => (
              <Link key={m.slug} className="chip" to={`/ui-menus/${m.slug}`}>
                {m.slug}
              </Link>
            ))}
          </div>
        </section>
      )}
      <section>
        <h2>Featured in maps</h2>
        {(varMaps ?? v.featured_in ?? []).length > 0 ? (
          <div className="chips">
            {(varMaps ?? v.featured_in ?? []).map((s) => (
              <Link key={s} className="chip" to={`/games/${s}`}>
                {s}
              </Link>
            ))}
          </div>
        ) : (
          <p className="meta">
            No maps yet.{' '}
            <a
              href={`https://github.com/blazium-games/game-mechanics-index/issues/new?template=map-variable-binding.yml`}
              target="_blank"
              rel="noreferrer"
            >
              Add a binding
            </a>
          </p>
        )}
      </section>
    </div>
  )
}
