import { formatMechanic } from '../utils/mechanicFormat'
import { resolveEnrichmentStatus } from '../utils/enrichmentStatus'
import {
  analyticsMissingEnrichmentStats,
  normalizeAnalyticsSnapshot,
} from '../utils/normalizeAnalytics'
import type {
  AnalyticsSnapshot,
  Catalog,
  Changelog,
  CooccurrencePair,
  GameplayMap,
  GameVariable,
  GenreIndexRow,
  MapIndexRow,
  MechanicEntry,
  MechanicIndexRow,
  MenuIndexRow,
  MenuFlowEdge,
  SearchRow,
  UIMenu,
  VariableIndexRow,
  DesignSkill,
  SkillIndexRow,
} from '../types'

const BASE =
  import.meta.env.VITE_DATA_BASE_URL ??
  `${import.meta.env.BASE_URL}api/v1`

async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE}${path}`)
  if (!res.ok) throw new Error(`API ${path}: ${res.status}`)
  return res.json() as Promise<T>
}

function withEnrichmentStatus<T extends { enrichment_status?: string }>(
  rows: T[],
): (T & { enrichment_status: ReturnType<typeof resolveEnrichmentStatus> })[] {
  return rows.map((r) => ({
    ...r,
    enrichment_status: resolveEnrichmentStatus(r.enrichment_status),
  }))
}

export const api = {
  base: BASE,
  fetchCatalog: () => get<Catalog>('/catalog.json'),
  fetchAnalytics: async () => {
    const raw = await get<Partial<AnalyticsSnapshot>>('/analytics.json')
    const normalized = normalizeAnalyticsSnapshot(raw)
    return Object.assign(normalized, {
      _legacyEnrichmentStats: analyticsMissingEnrichmentStats(raw),
    }) as AnalyticsSnapshot & { _legacyEnrichmentStats?: boolean }
  },
  fetchChangelog: () => get<Changelog>('/changelog.json'),
  fetchMapsIndex: () => get<MapIndexRow[]>('/maps/index.json'),
  fetchMap: (slug: string) => get<GameplayMap>(`/maps/${slug}.json`),
  fetchMechanicsIndex: async () => withEnrichmentStatus(await get<MechanicIndexRow[]>('/mechanics/index.json')),
  fetchMechanic: (slug: string) => get<MechanicEntry>(`/mechanics/${slug}.json`),
  fetchVariablesIndex: async () => withEnrichmentStatus(await get<VariableIndexRow[]>('/variables/index.json')),
  fetchVariable: (slug: string) => get<GameVariable>(`/variables/${slug}.json`),
  fetchUIMenusIndex: async () => withEnrichmentStatus(await get<MenuIndexRow[]>('/ui-menus/index.json')),
  fetchUIMenu: (slug: string) => get<UIMenu>(`/ui-menus/${slug}.json`),
  fetchGenresIndex: () => get<GenreIndexRow[]>('/genres/index.json'),
  fetchGenre: (slug: string) => get<GameplayMap>(`/genres/${slug}.json`),
  fetchSearch: () => get<SearchRow[]>('/search.json'),
  fetchCooccurrence: async (limit = 100) => {
    const data = await get<{ pairs?: CooccurrencePair[] } | CooccurrencePair[]>(
      '/indexes/cooccurrence-top500.json',
    )
    const pairs = Array.isArray(data) ? data : (data.pairs ?? [])
    return pairs.slice(0, limit)
  },
  fetchMechanicToMaps: () =>
    get<{ mechanics: Record<string, string[]> }>('/indexes/mechanic-to-maps.json'),
  fetchVariableToMaps: () =>
    get<{ variables: Record<string, string[]> }>('/indexes/variable-to-maps.json'),
  fetchMenuToMaps: () =>
    get<{ menus: Record<string, string[]> }>('/indexes/menu-to-maps.json'),
  fetchVariableToMechanics: () =>
    get<{ variables: Record<string, string[]> }>('/indexes/variable-to-mechanics.json'),
  fetchMenuFlowEdges: () =>
    get<{ edges: MenuFlowEdge[] }>('/indexes/menu-flow-edges.json'),
  fetchTags: () => get<{ tags: string[] }>('/tags.json'),
  fetchVariableTags: () => get<{ tags: string[] }>('/variable-tags.json'),
  fetchMenuTags: () => get<{ tags: string[] }>('/menu-tags.json'),
  fetchSkillsIndex: async () => withEnrichmentStatus(await get<SkillIndexRow[]>('/skills/index.json')),
  fetchSkill: (slug: string) => get<DesignSkill>(`/skills/${slug}.json`),
  fetchSkillToMechanics: () =>
    get<{ skills: Record<string, string[]> }>('/indexes/skill-to-mechanics.json'),
  fetchMechanicFormatted: async (slug: string, format: 'md' | 'yaml' | 'txt' = 'md') => {
    const entry = await get<MechanicEntry>(`/mechanics/${slug}.json`)
    return formatMechanic(entry, format)
  },
}

export type ApiClient = typeof api
