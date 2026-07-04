export const REPO = 'blazium-games/game-mechanics-index'
export const REPO_URL = `https://github.com/${REPO}`

export interface Catalog {
  schema_version: string
  license: string
  release_version: string
  map_count: number
  mechanic_count: number
  genre_count: number
  game_count: number
  genres: string[]
  domains: string[]
  flavors: string[]
}

export interface MapIndexRow {
  slug: string
  title: string
  name: string
  map_type: string
  genres?: string[]
  quality_tier: string
  signature_count: number
  year?: number
}

export interface MechanicIndexRow {
  slug: string
  name: string
  domain: string
  flavor: string
  tags?: string[]
  featured_count: number
  signature_games_count: number
}

export interface GenreIndexRow {
  slug: string
  title: string
  name: string
}

export interface SearchRow {
  type: 'game' | 'mechanic'
  slug: string
  title: string
  genres?: string[]
  tags?: string[]
  tier?: string
}

export interface GameplayMap {
  slug: string
  title: string
  map_type: string
  subject: { name: string; genres?: string[]; influences?: string[] }
  narrative: { description: string; core_loop?: string; skills_tested?: string[] }
  signature_gameplay: string[]
  mechanics: Array<{
    mechanic_slug: string
    role: string
    domain?: string
    map_notes?: string
    expression?: string
  }>
  metadata?: { quality_tier?: string }
  relationships?: { genre_maps?: string[]; influence_slugs?: string[] }
  design_intent?: { theme_tags?: string[]; player_fantasy?: string; design_pillars?: string[] }
  context?: { year?: number; platforms?: string[] }
}

export interface MechanicEntry {
  slug: string
  name: string
  flavor: string
  domain: string
  summary: string
  tags?: string[]
  synergies?: string[]
  conflicts?: string[]
  signature_of?: { games?: string[]; genres?: string[] }
  featured_in?: string[]
}

export interface CooccurrencePair {
  mechanic_a: string
  mechanic_b: string
  count: number
}
