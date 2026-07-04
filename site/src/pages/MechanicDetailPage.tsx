import { useQuery } from '@tanstack/react-query'
import { Link, useParams } from 'react-router-dom'
import { api } from '../api/client'
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

  if (!mech) return <p>Loading…</p>

  return (
    <div>
      <div className="detail-header">
        <h1>{mech.name}</h1>
        <SuggestEditLink slug={slug} kind="mechanic" />
      </div>
      <p className="meta">
        {mech.domain} · flavor: {mech.flavor}
      </p>
      <section>
        <h2>Summary</h2>
        <p>{mech.summary}</p>
      </section>
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
