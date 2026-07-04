export const REPO = 'blazium-games/game-mechanics-index'
export const REPO_URL = `https://github.com/${REPO}`
export const RELEASE_URL = `https://github.com/${REPO}/releases/latest`

export interface Catalog {
  schema_version: string
  license: string
  release_version: string
  map_count: number
  mechanic_count: number
  variable_count?: number
  menu_count?: number
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

export interface VariableIndexRow {
  slug: string
  name: string
  category: string
  scope: string
  tags?: string[]
  featured_count: number
  enrichment_status: 'complete' | 'needs_info'
}

export interface MenuIndexRow {
  slug: string
  name: string
  menu_type: string
  layer: string
  tags?: string[]
  featured_count: number
  enrichment_status: 'complete' | 'needs_info'
}

export interface GenreIndexRow {
  slug: string
  title: string
  name: string
}

export interface SearchRow {
  type: 'game' | 'mechanic' | 'variable' | 'menu'
  slug: string
  title: string
  genres?: string[]
  tags?: string[]
  tier?: string
}

export interface MapVariableBinding {
  variable_slug: string
  role: string
  map_notes?: string
  expression?: string
  related_mechanics?: string[]
}

export interface MapUIMenuBinding {
  menu_slug: string
  role: string
  map_notes?: string
  opens_from?: string[]
  displays_variables?: string[]
  supports_mechanics?: string[]
}

export interface VariableRelationship {
  from_variable: string
  to_variable: string
  relationship?: string
  notes?: string
}

export interface MenuFlowEdge {
  from_menu: string
  to_menu: string
  relationship?: string
  notes?: string
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
  variables?: MapVariableBinding[]
  ui_menus?: MapUIMenuBinding[]
  variable_relationships?: VariableRelationship[]
  menu_flow?: MenuFlowEdge[]
  metadata?: { quality_tier?: string }
  relationships?: { genre_maps?: string[]; influence_slugs?: string[] }
  design_intent?: { theme_tags?: string[]; player_fantasy?: string; design_pillars?: string[] }
  context?: { year?: number; platforms?: string[] }
  gdd_outline?: { combat_notes?: string; economy_notes?: string; player_goals?: string[] }
  systems?: { failure_model?: string; economy_tightness?: string }
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
  flavor_rationale?: string
  player_experience?: string
  common_in?: string[]
  design_guidance?: {
    when_to_use?: string
    where_to_use?: string
    when_to_avoid?: string[]
    designer_notes?: string
  }
  examples?: Array<{ label?: string; description?: string; map_slug?: string }>
  synergy_notes?: Array<{ slug: string; note?: string }>
}

export interface GameVariable {
  slug: string
  name: string
  summary: string
  category: string
  scope: string
  value_kind: string
  shared_rationale?: string
  player_focus?: string
  typical_range?: string
  reset_behavior?: string
  related_mechanics?: string[]
  related_variables?: Array<{ slug: string; relationship?: string; notes?: string }>
  featured_in?: string[]
  common_in?: string[]
  tags?: string[]
  design_guidance?: MechanicEntry['design_guidance']
}

export interface UIMenu {
  slug: string
  name: string
  summary: string
  menu_type: string
  layer: string
  input_context?: string
  shared_rationale?: string
  typical_actions?: string[]
  related_mechanics?: string[]
  related_variables?: string[]
  related_menus?: Array<{ slug: string; relationship?: string; notes?: string }>
  featured_in?: string[]
  common_in?: string[]
  tags?: string[]
  design_guidance?: MechanicEntry['design_guidance']
}

export interface CooccurrencePair {
  mechanic_a: string
  mechanic_b: string
  count: number
}
