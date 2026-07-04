import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'
import { api } from '../api/client'
import { EmptyField } from '../components/EmptyField'
import { ExportDropdown } from '../components/ExportDropdown'
import { SuggestEditLink } from '../components/Layout'

export function UIMenuDetailPage() {
  const { slug = '' } = useParams()
  const { data: menu } = useQuery({
    queryKey: ['ui-menu', slug],
    queryFn: () => api.fetchUIMenu(slug),
    enabled: !!slug,
  })
  const { data: menuMaps } = useQuery({
    queryKey: ['menu-maps', slug],
    queryFn: async () => {
      const idx = await api.fetchMenuToMaps()
      return idx.menus[slug] ?? []
    },
    enabled: !!slug,
  })
  const { data: flow } = useQuery({
    queryKey: ['menu-flow'],
    queryFn: api.fetchMenuFlowEdges,
  })

  if (!menu) return <p>Loading…</p>

  const outEdges = flow?.edges.filter((e) => e.from_menu === slug) ?? []
  const inEdges = flow?.edges.filter((e) => e.to_menu === slug) ?? []

  return (
    <div>
      <div className="detail-header">
        <h1>{menu.name}</h1>
        <div className="detail-actions">
          <ExportDropdown kind="menu" slug={slug} entity={menu} />
          <SuggestEditLink slug={slug} kind="ui-menu" />
        </div>
      </div>
      <p className="meta">
        {menu.menu_type} · {menu.layer}
        {menu.input_context ? ` · ${menu.input_context}` : ''}
      </p>
      <section>
        <h2>Summary</h2>
        <p>{menu.summary}</p>
      </section>
      <section>
        <h2>Shared rationale</h2>
        {menu.shared_rationale ? (
          <p>{menu.shared_rationale}</p>
        ) : (
          <EmptyField slug={slug} field="shared_rationale" kind="ui-menu" />
        )}
      </section>
      {menu.typical_actions && menu.typical_actions.length > 0 && (
        <section>
          <h2>Typical actions</h2>
          <div className="chips">
            {menu.typical_actions.map((a) => (
              <span key={a} className="chip">
                {a}
              </span>
            ))}
          </div>
        </section>
      )}
      {menu.related_mechanics && menu.related_mechanics.length > 0 && (
        <section>
          <h2>Related mechanics</h2>
          <div className="chips">
            {menu.related_mechanics.map((m) => (
              <Link key={m} className="chip" to={`/mechanics/${m}`}>
                {m}
              </Link>
            ))}
          </div>
        </section>
      )}
      {menu.related_variables && menu.related_variables.length > 0 && (
        <section>
          <h2>Related variables</h2>
          <div className="chips">
            {menu.related_variables.map((v) => (
              <Link key={v} className="chip" to={`/variables/${v}`}>
                {v}
              </Link>
            ))}
          </div>
        </section>
      )}
      {menu.related_menus && menu.related_menus.length > 0 && (
        <section>
          <h2>Related menus</h2>
          <div className="chips">
            {menu.related_menus.map((rel) => (
              <Link key={rel.slug} className="chip" to={`/ui-menus/${rel.slug}`}>
                {rel.slug}
                {rel.relationship ? ` (${rel.relationship})` : ''}
              </Link>
            ))}
          </div>
        </section>
      )}
      {(outEdges.length > 0 || inEdges.length > 0) && (
        <section>
          <h2>Menu flow</h2>
          <ul>
            {outEdges.map((e) => (
              <li key={`${e.from_menu}-${e.to_menu}`}>
                → <Link to={`/ui-menus/${e.to_menu}`}>{e.to_menu}</Link>
                {e.relationship ? ` (${e.relationship})` : ''}
              </li>
            ))}
            {inEdges.map((e) => (
              <li key={`in-${e.from_menu}-${e.to_menu}`}>
                ← <Link to={`/ui-menus/${e.from_menu}`}>{e.from_menu}</Link>
                {e.relationship ? ` (${e.relationship})` : ''}
              </li>
            ))}
          </ul>
        </section>
      )}
      <section>
        <h2>Featured in maps</h2>
        {(menuMaps ?? menu.featured_in ?? []).length > 0 ? (
          <div className="chips">
            {(menuMaps ?? menu.featured_in ?? []).map((s) => (
              <Link key={s} className="chip" to={`/games/${s}`}>
                {s}
              </Link>
            ))}
          </div>
        ) : (
          <p className="meta">No maps yet.</p>
        )}
      </section>
    </div>
  )
}
