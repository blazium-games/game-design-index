import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'
import { api } from '../api/client'
import { ExportDropdown } from '../components/ExportDropdown'
import { SuggestEditLink } from '../components/Layout'

export function MechanicDetailPage() {
  const { slug = '' } = useParams()
  const { data: mech } = useQuery({
    queryKey: ['mechanic', slug],
    queryFn: () => api.fetchMechanic(slug),
    enabled: !!slug,
  })
  const { data: mechMaps } = useQuery({
    queryKey: ['mechanic-maps', slug],
    queryFn: async () => {
      const idx = await api.fetchMechanicToMaps()
      return idx.mechanics[slug] ?? []
    },
    enabled: !!slug,
  })
  const { data: varMechs } = useQuery({
    queryKey: ['variable-to-mechanics', slug],
    queryFn: async () => {
      const idx = await api.fetchVariableToMechanics()
      const out: string[] = []
      for (const [v, mechs] of Object.entries(idx.variables)) {
        if (mechs.includes(slug)) out.push(v)
      }
      return out.sort()
    },
    enabled: !!slug,
  })
  const { data: menuIndex } = useQuery({
    queryKey: ['ui-menus-index'],
    queryFn: api.fetchUIMenusIndex,
  })
  const { data: relatedMenus } = useQuery({
    queryKey: ['menus-for-mechanic', slug, menuIndex?.length],
    queryFn: async () => {
      if (!menuIndex) return []
      const menus = await Promise.all(menuIndex.map((row) => api.fetchUIMenu(row.slug)))
      return menus.filter((m) => m.related_mechanics?.includes(slug))
    },
    enabled: !!slug && !!menuIndex,
  })

  if (!mech) return <p>Loading…</p>

  return (
    <div>
      <div className="detail-header">
        <h1>{mech.name}</h1>
        <div className="detail-actions">
          <ExportDropdown kind="mechanic" slug={slug} entity={mech} />
          <SuggestEditLink slug={slug} kind="mechanic" />
        </div>
      </div>
      <p className="meta">
        {mech.domain} · flavor: {mech.flavor}
      </p>
      <section>
        <h2>Summary</h2>
        <p>{mech.summary}</p>
      </section>
      {mech.player_experience && (
        <section>
          <h2>Player experience</h2>
          <p>{mech.player_experience}</p>
        </section>
      )}
      {mech.flavor_rationale && (
        <section>
          <h2>Flavor rationale</h2>
          <p>{mech.flavor_rationale}</p>
        </section>
      )}
      {mech.tags && mech.tags.length > 0 && (
        <section>
          <h2>Tags</h2>
          <div className="chips">
            {mech.tags.map((t) => (
              <span key={t} className="chip">
                {t}
              </span>
            ))}
          </div>
        </section>
      )}
      {varMechs && varMechs.length > 0 && (
        <section>
          <h2>Related variables</h2>
          <div className="chips">
            {varMechs.map((v) => (
              <Link key={v} className="chip" to={`/variables/${v}`}>
                {v}
              </Link>
            ))}
          </div>
        </section>
      )}
      {relatedMenus && relatedMenus.length > 0 && (
        <section>
          <h2>Related UI menus</h2>
          <div className="chips">
            {relatedMenus.map((m) => (
              <Link key={m.slug} className="chip" to={`/ui-menus/${m.slug}`}>
                {m.slug}
              </Link>
            ))}
          </div>
        </section>
      )}
      {mech.signature_of?.games && mech.signature_of.games.length > 0 && (
        <section>
          <h2>Signature games</h2>
          <ul>
            {mech.signature_of.games.map((g) => (
              <li key={g}>{g}</li>
            ))}
          </ul>
        </section>
      )}
      {mech.synergies && mech.synergies.length > 0 && (
        <section>
          <h2>Synergies</h2>
          <div className="chips">
            {mech.synergies.map((s) => (
              <Link key={s} className="chip" to={`/mechanics/${s}`}>
                {s}
              </Link>
            ))}
          </div>
        </section>
      )}
      {mech.design_guidance?.when_to_use && (
        <section>
          <h2>When to use</h2>
          <p>{mech.design_guidance.when_to_use}</p>
        </section>
      )}
      <section>
        <h2>Featured in maps</h2>
        <div className="chips">
          {(mechMaps ?? mech.featured_in ?? []).slice(0, 40).map((s) => (
            <Link key={s} className="chip" to={`/games/${s}`}>
              {s}
            </Link>
          ))}
        </div>
      </section>
    </div>
  )
}
