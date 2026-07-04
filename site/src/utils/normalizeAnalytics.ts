import type {
  AnalyticsOverview,
  AnalyticsSnapshot,
  EntityEnrichmentStats,
  GenreDomainHeatmap,
  MenuFlowStats,
} from '../types'

export const EMPTY_ENRICHMENT_STATS: EntityEnrichmentStats = {
  enrichment_complete: 0,
  enrichment_needs_info: 0,
  with_map_bindings: 0,
}

function normalizeEnrichmentStats(raw?: Partial<EntityEnrichmentStats>): EntityEnrichmentStats {
  return {
    by_category: raw?.by_category ?? [],
    by_scope: raw?.by_scope ?? [],
    by_menu_type: raw?.by_menu_type ?? [],
    by_layer: raw?.by_layer ?? [],
    enrichment_complete: raw?.enrichment_complete ?? 0,
    enrichment_needs_info: raw?.enrichment_needs_info ?? 0,
    with_map_bindings: raw?.with_map_bindings ?? 0,
  }
}

function normalizeOverview(raw?: Partial<AnalyticsOverview>): AnalyticsOverview {
  return {
    map_count: raw?.map_count ?? 0,
    game_count: raw?.game_count ?? 0,
    genre_recipe_count: raw?.genre_recipe_count ?? 0,
    mechanic_count: raw?.mechanic_count ?? 0,
    variable_count: raw?.variable_count ?? 0,
    menu_count: raw?.menu_count ?? 0,
    skill_count: raw?.skill_count,
    mechanic_enrichment_pct: raw?.mechanic_enrichment_pct,
    skill_enrichment_pct: raw?.skill_enrichment_pct,
    variable_enrichment_pct: raw?.variable_enrichment_pct ?? 0,
    menu_enrichment_pct: raw?.menu_enrichment_pct ?? 0,
    avg_signature_count: raw?.avg_signature_count ?? 0,
  }
}

function normalizeHeatmap(raw?: Partial<GenreDomainHeatmap>): GenreDomainHeatmap {
  return {
    genres: raw?.genres ?? [],
    domains: raw?.domains ?? [],
    cells: raw?.cells ?? [],
  }
}

function normalizeMenuFlow(raw?: Partial<MenuFlowStats>): MenuFlowStats {
  return {
    edge_counts_by_relationship: raw?.edge_counts_by_relationship ?? [],
    top_hubs: raw?.top_hubs ?? [],
  }
}

export function analyticsMissingEnrichmentStats(raw: Partial<AnalyticsSnapshot>): boolean {
  return raw.mechanic_stats == null || raw.skill_stats == null
}

export function normalizeAnalyticsSnapshot(raw: Partial<AnalyticsSnapshot>): AnalyticsSnapshot {
  return {
    schema_version: raw.schema_version ?? '1.0',
    overview: normalizeOverview(raw.overview),
    quality_tiers: raw.quality_tiers ?? [],
    games_by_decade: raw.games_by_decade ?? [],
    top_genres: raw.top_genres ?? [],
    mechanic_domains: raw.mechanic_domains ?? [],
    mechanic_flavors: raw.mechanic_flavors ?? [],
    mechanic_complexity: raw.mechanic_complexity ?? [],
    top_mechanics: raw.top_mechanics ?? [],
    genre_domain_heatmap: normalizeHeatmap(raw.genre_domain_heatmap),
    cooccurrence: raw.cooccurrence ?? [],
    mechanic_stats: normalizeEnrichmentStats(raw.mechanic_stats),
    skill_stats: normalizeEnrichmentStats(raw.skill_stats),
    variable_stats: normalizeEnrichmentStats(raw.variable_stats),
    menu_stats: normalizeEnrichmentStats(raw.menu_stats),
    variable_mechanic_pairs: raw.variable_mechanic_pairs ?? [],
    menu_flow: normalizeMenuFlow(raw.menu_flow),
    signature_distribution: raw.signature_distribution ?? [],
    insights: raw.insights ?? [],
  }
}
