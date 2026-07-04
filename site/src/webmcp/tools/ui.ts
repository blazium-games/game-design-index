import { REPO_URL } from '../../types'
import { jsonResponse } from '../responses'
import type { WebMCPDeps } from '../types'

type RegisterOpts = { signal?: AbortSignal }

export async function registerUiTools(deps: WebMCPDeps, opts: RegisterOpts) {
  const mc = document.modelContext

  await mc.registerTool(
    {
      name: 'navigate',
      description: 'Navigate the catalog UI to a game, mechanic, genre, or path.',
      inputSchema: {
        type: 'object',
        properties: {
          target: { type: 'string', enum: ['game', 'mechanic', 'genre', 'path'] },
          slug: { type: 'string' },
          path: { type: 'string' },
        },
        required: ['target'],
      },
      async execute(args) {
        const { target, slug, path } = args as {
          target: string
          slug?: string
          path?: string
        }
        let dest = '/'
        if (target === 'path' && path) dest = path
        else if (target === 'game' && slug) dest = `/games/${slug}`
        else if (target === 'mechanic' && slug) dest = `/mechanics/${slug}`
        else if (target === 'genre' && slug) dest = `/genres/${slug}`
        deps.navigate(dest)
        return jsonResponse({ navigated: dest })
      },
    },
    opts,
  )

  await mc.registerTool(
    {
      name: 'filter-games-view',
      description: 'Apply filters on the games index view and navigate there.',
      inputSchema: {
        type: 'object',
        properties: {
          genre: { type: 'string' },
          quality_tier: { type: 'string' },
          query: { type: 'string' },
        },
      },
      async execute(args) {
        const { genre, quality_tier, query } = args as {
          genre?: string
          quality_tier?: string
          query?: string
        }
        deps.setGames({ genre: genre ?? '', tier: quality_tier ?? '', query: query ?? '' })
        deps.navigate('/games')
        const rows = await deps.api.fetchMapsIndex()
        let count = rows.filter((m) => m.map_type === 'game').length
        if (genre)
          count = rows.filter(
            (m) =>
              m.map_type === 'game' &&
              m.genres?.some((g) => g.toLowerCase() === genre.toLowerCase()),
          ).length
        return jsonResponse({ matching_count: count, filters: { genre, quality_tier, query } })
      },
    },
    opts,
  )

  await mc.registerTool(
    {
      name: 'filter-mechanics-view',
      description: 'Apply filters on the Design Index view and navigate there.',
      inputSchema: {
        type: 'object',
        properties: {
          domain: { type: 'string' },
          flavor: { type: 'string' },
          tag: { type: 'string' },
          query: { type: 'string' },
        },
      },
      async execute(args) {
        const { domain, flavor, tag, query } = args as {
          domain?: string
          flavor?: string
          tag?: string
          query?: string
        }
        deps.setMechanics({
          domain: domain ?? '',
          flavor: flavor ?? '',
          tag: tag ?? '',
          query: query ?? '',
        })
        deps.navigate('/mechanics')
        return jsonResponse({ filters: { domain, flavor, tag, query } })
      },
    },
    opts,
  )

  await mc.registerTool(
    {
      name: 'open-contribute',
      description: 'Open the contribution flow for missing games, corrections, or new mechanics.',
      inputSchema: {
        type: 'object',
        properties: {
          kind: {
            type: 'string',
            enum: ['missing-game', 'data-correction', 'new-mechanic'],
          },
          slug: { type: 'string', description: 'Optional slug for correction issues' },
        },
        required: ['kind'],
      },
      async execute(args) {
        const { kind, slug } = args as { kind: string; slug?: string }
        if (kind === 'data-correction' && slug) {
          const url = `${REPO_URL}/issues/new?template=data-correction.yml&title=${encodeURIComponent('Data correction: ' + slug)}`
          window.open(url, '_blank')
          return jsonResponse({ opened: url })
        }
        deps.navigate('/contribute')
        return jsonResponse({ navigated: '/contribute', kind })
      },
    },
    opts,
  )
}
