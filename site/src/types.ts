export const REPO = 'blazium-games/game-design-index'
export const REPO_URL = `https://github.com/${REPO}`
export const RELEASE_URL = `https://github.com/${REPO}/releases/latest`
export const DISCORD_URL = 'https://discord.gg/sZaf9KYzDp'
export const SITE_URL = 'https://blazium-games.github.io/game-design-index'
export const SITE_NAME = 'Game Design Index'
export const TWITTER_URL = 'https://x.com/BlaziumGames'
export const TWITTER_HANDLE = '@BlaziumGames'
export const ORG_NAME = 'Blazium Games'
export const DEFAULT_OG_IMAGE = `${SITE_URL}/og-default.png`

export interface Catalog {
  schema_version: string
  license: string
  release_version: string
  map_count: number
  mechanic_count: number
  variable_count?: number
  menu_count?: number
  skill_count?: number
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
  enrichment_status?: 'complete' | 'needs_info'
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

export interface SkillIndexRow {
  slug: string
  name: string
  category: string
  tags?: string[]
  mechanic_count: number
  enrichment_status: 'complete' | 'needs_info'
}

export interface DesignSkill {
  slug: string
  name: string
  summary: string
  category: string
  learning_outcome: string
  practice_activities?: string[]
  related_mechanics?: string[]
  related_variables?: string[]
  design_guidance?: {
    when_to_use?: string
    where_to_use?: string
    when_to_avoid?: string[]
    designer_notes?: string
  }
  tags?: string[]
}

export interface AgentContext {
  summary_for_agents?: string
  gdd_prompt?: string
  implementation_checklist?: string[]
}

export interface DesignExercise {
  prompt: string
  constraints?: string[]
  success_criteria?: string[]
}

export interface GenreIndexRow {
  slug: string
  title: string
  name: string
}

export interface SearchRow {
  type: 'game' | 'mechanic' | 'variable' | 'menu' | 'skill'
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
    phase?: string
    weight?: number
    depends_on?: string[]
  }>
  variables?: MapVariableBinding[]
  ui_menus?: MapUIMenuBinding[]
  variable_relationships?: VariableRelationship[]
  menu_flow?: MenuFlowEdge[]
  metadata?: { quality_tier?: string }
  relationships?: { genre_maps?: string[]; influence_slugs?: string[] }
  design_intent?: { theme_tags?: string[]; player_fantasy?: string; design_pillars?: string[] }
  context?: {
    year?: number
    platforms?: string[]
    dimension?: string
    perspective?: string
    world_structure?: string
    session_type?: string
  }
  gdd_outline?: {
    overview?: string
    core_loop?: string
    combat_notes?: string
    economy_notes?: string
    progression_notes?: string
    constraints?: string[]
    player_goals?: string[]
  }
  systems?: {
    failure_model?: string
    economy_tightness?: string
    metaprogression?: boolean
    pacing?: { early?: string; mid?: string; late?: string }
    primary_loop_phases?: Array<{ id: string; label: string; mechanics?: string[] }>
  }
  variants?: Array<{
    id: string
    label: string
    notes?: string
    add_signatures?: string[]
    drop_signatures?: string[]
  }>
  views?: Array<{ id: string; label: string; filter?: Record<string, string> }>
  skill_slugs?: string[]
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
  learning_objectives?: string[]
  design_exercises?: DesignExercise[]
  skills_developed?: string[]
  agent_context?: AgentContext
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
  agent_context?: AgentContext
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
  agent_context?: AgentContext
}

export interface CooccurrencePair {
  mechanic_a: string
  mechanic_b: string
  count: number
}

export interface CountRow {
  label: string
  count: number
}

export interface MechanicAdoptionRow {
  slug: string
  name: string
  count: number
}

export interface CooccurrenceAnalyticsRow {
  mechanic_a: string
  mechanic_b: string
  count: number
  lift: number
}

export interface HeatmapCell {
  genre: string
  domain: string
  count: number
}

export interface GenreDomainHeatmap {
  genres: string[]
  domains: string[]
  cells: HeatmapCell[]
}

export interface EntityEnrichmentStats {
  by_category?: CountRow[]
  by_scope?: CountRow[]
  by_menu_type?: CountRow[]
  by_layer?: CountRow[]
  enrichment_complete: number
  enrichment_needs_info: number
  with_map_bindings: number
}

export interface MenuHubRow {
  slug: string
  name: string
  in_degree: number
  out_degree: number
  total_degree: number
}

export interface MenuFlowStats {
  edge_counts_by_relationship: CountRow[]
  top_hubs: MenuHubRow[]
}

export interface VariableMechanicPair {
  variable: string
  mechanic: string
}

export interface AnalyticsOverview {
  map_count: number
  game_count: number
  genre_recipe_count: number
  mechanic_count: number
  variable_count: number
  menu_count: number
  skill_count?: number
  mechanic_enrichment_pct?: number
  skill_enrichment_pct?: number
  variable_enrichment_pct: number
  menu_enrichment_pct: number
  avg_signature_count: number
}

export interface AnalyticsSnapshot {
  schema_version: string
  overview: AnalyticsOverview
  quality_tiers: CountRow[]
  games_by_decade: CountRow[]
  top_genres: CountRow[]
  mechanic_domains: CountRow[]
  mechanic_flavors: CountRow[]
  mechanic_complexity: CountRow[]
  top_mechanics: MechanicAdoptionRow[]
  genre_domain_heatmap: GenreDomainHeatmap
  cooccurrence: CooccurrenceAnalyticsRow[]
  mechanic_stats: EntityEnrichmentStats
  skill_stats: EntityEnrichmentStats
  variable_stats: EntityEnrichmentStats
  menu_stats: EntityEnrichmentStats
  variable_mechanic_pairs: VariableMechanicPair[]
  menu_flow: MenuFlowStats
  signature_distribution: CountRow[]
  insights: string[]
}

export interface ChangelogSection {
  heading: string
  items: string[]
}

export interface ChangelogEntry {
  version: string
  date: string
  title: string
  highlights?: string[]
  sections: ChangelogSection[]
  release_url?: string
}

export interface Changelog {
  schema_version: string
  entries: ChangelogEntry[]
}
