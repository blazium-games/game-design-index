import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'
import { api } from '../api/client'
import { ExportDropdown } from '../components/ExportDropdown'
import { SuggestEditLink } from '../components/Layout'
import { REPO_URL } from '../types'

export function GameDetailPage() {
  const { slug = '' } = useParams()
  const { data: map } = useQuery({
    queryKey: ['map', slug],
    queryFn: () => api.fetchMap(slug),
    enabled: !!slug,
  })

  if (!map) return <p>Loading…</p>

  const hasVars = (map.variables?.length ?? 0) > 0
  const hasMenus = (map.ui_menus?.length ?? 0) > 0

  return (
    <div>
      <div className="detail-header">
        <h1>{map.subject.name}</h1>
        <div className="detail-actions">
          <ExportDropdown kind="game" slug={slug} entity={map} />
          <SuggestEditLink slug={slug} kind="game" />
        </div>
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
      <section>
        <h2>Variable bindings</h2>
        {hasVars ? (
          <table className="table">
            <thead>
              <tr>
                <th>Variable</th>
                <th>Role</th>
                <th>Expression</th>
              </tr>
            </thead>
            <tbody>
              {map.variables!.map((vb) => (
                <tr key={vb.variable_slug}>
                  <td>
                    <Link to={`/variables/${vb.variable_slug}`}>{vb.variable_slug}</Link>
                  </td>
                  <td>{vb.role}</td>
                  <td>{vb.expression || vb.map_notes}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <p className="meta">
            No variable bindings documented.{' '}
            <a href={`${REPO_URL}/issues/new?template=map-variable-binding.yml`} target="_blank" rel="noreferrer">
              Add bindings
            </a>
          </p>
        )}
      </section>
      <section>
        <h2>UI menu bindings</h2>
        {hasMenus ? (
          <table className="table">
            <thead>
              <tr>
                <th>Menu</th>
                <th>Role</th>
                <th>Notes</th>
              </tr>
            </thead>
            <tbody>
              {map.ui_menus!.map((mb) => (
                <tr key={mb.menu_slug}>
                  <td>
                    <Link to={`/ui-menus/${mb.menu_slug}`}>{mb.menu_slug}</Link>
                  </td>
                  <td>{mb.role}</td>
                  <td>{mb.map_notes}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <p className="meta">
            No UI menu bindings documented.{' '}
            <a href={`${REPO_URL}/issues/new?template=map-variable-binding.yml`} target="_blank" rel="noreferrer">
              Add bindings
            </a>
          </p>
        )}
      </section>
      {map.variable_relationships && map.variable_relationships.length > 0 && (
        <section>
          <h2>Variable relationships</h2>
          <ul>
            {map.variable_relationships.map((rel) => (
              <li key={`${rel.from_variable}-${rel.to_variable}`}>
                <Link to={`/variables/${rel.from_variable}`}>{rel.from_variable}</Link> →{' '}
                <Link to={`/variables/${rel.to_variable}`}>{rel.to_variable}</Link>
                {rel.relationship ? ` (${rel.relationship})` : ''}
                {rel.notes ? ` — ${rel.notes}` : ''}
              </li>
            ))}
          </ul>
        </section>
      )}
      {map.menu_flow && map.menu_flow.length > 0 && (
        <section>
          <h2>Menu flow</h2>
          <ul>
            {map.menu_flow.map((edge) => (
              <li key={`${edge.from_menu}-${edge.to_menu}`}>
                <Link to={`/ui-menus/${edge.from_menu}`}>{edge.from_menu}</Link> →{' '}
                <Link to={`/ui-menus/${edge.to_menu}`}>{edge.to_menu}</Link>
                {edge.relationship ? ` (${edge.relationship})` : ''}
              </li>
            ))}
          </ul>
        </section>
      )}
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
