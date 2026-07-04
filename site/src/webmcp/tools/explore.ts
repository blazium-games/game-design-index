import { jsonResponse } from '../responses'
import type { WebMCPDeps } from '../types'

type RegisterOpts = { signal?: AbortSignal }

export async function registerExploreTools(deps: WebMCPDeps, opts: RegisterOpts) {
  const mc = document.modelContext

  await mc.registerTool(
    {
      name: 'get-cooccurrence',
      description: 'Return top mechanic pairs that co-occur on the same gameplay maps.',
      inputSchema: {
        type: 'object',
        properties: {
          limit: { type: 'number', description: 'Max pairs (default 50)' },
          min_count: { type: 'number', description: 'Minimum co-occurrence count' },
        },
      },
      async execute(args) {
        const { limit = 50, min_count = 0 } = args as { limit?: number; min_count?: number }
        const raw = await deps.api.fetchCooccurrence(500)
        const pairs = raw.filter((p) => p.count >= min_count).slice(0, limit)
        return jsonResponse(pairs)
      },
    },
    opts,
  )

  await mc.registerTool(
    {
      name: 'get-similar-games',
      description: 'Find games sharing signature mechanics with a reference game.',
      inputSchema: {
        type: 'object',
        properties: {
          slug: { type: 'string' },
          limit: { type: 'number' },
        },
        required: ['slug'],
      },
      async execute(args) {
        const { slug, limit = 10 } = args as { slug: string; limit?: number }
        const ref = await deps.api.fetchMap(slug)
        const sigs = new Set(ref.signature_gameplay)
        const index = await deps.api.fetchMapsIndex()
        const scored: { slug: string; shared: number }[] = []
        for (const row of index) {
          if (row.slug === slug || row.map_type !== 'game') continue
          const full = await deps.api.fetchMap(row.slug)
          const shared = full.signature_gameplay.filter((s) => sigs.has(s)).length
          if (shared > 0) scored.push({ slug: row.slug, shared })
        }
        scored.sort((a, b) => b.shared - a.shared)
        return jsonResponse(scored.slice(0, limit))
      },
    },
    opts,
  )

  await mc.registerTool(
    {
      name: 'compose-design-brief',
      description: 'Merge 1-4 reference games into a design brief with shared signatures and synergies.',
      inputSchema: {
        type: 'object',
        properties: {
          ref_slugs: {
            type: 'array',
            items: { type: 'string' },
            description: '1-4 game map slugs',
          },
        },
        required: ['ref_slugs'],
      },
      async execute(args) {
        const { ref_slugs } = args as { ref_slugs: string[] }
        const slugs = ref_slugs.slice(0, 4)
        const maps = await Promise.all(slugs.map((s) => deps.api.fetchMap(s)))
        const sigSeen = new Map<string, { slug: string; notes: string[] }>()
        for (const m of maps) {
          for (const sig of m.signature_gameplay) {
            const entry = sigSeen.get(sig) ?? { slug: sig, notes: [] }
            const note = m.mechanics.find((b) => b.mechanic_slug === sig)?.map_notes
            if (note) entry.notes.push(note)
            sigSeen.set(sig, entry)
          }
        }
        const cooc = await deps.api.fetchCooccurrence(200)
        const sigList = [...sigSeen.keys()]
        const pairs = cooc.filter(
          (p) => sigList.includes(p.mechanic_a) && sigList.includes(p.mechanic_b),
        )
        return jsonResponse({
          references: maps.map((m) => ({
            slug: m.slug,
            name: m.subject.name,
            signatures: m.signature_gameplay,
          })),
          signatures: [...sigSeen.values()],
          synergy_pairs: pairs.slice(0, 20),
        })
      },
    },
    opts,
  )
}
