import { describe, expect, it } from 'vitest'
import type { AnalyticsSnapshot } from '../types'

const mockAnalytics: AnalyticsSnapshot = {
  schema_version: '1.0',
  overview: {
    map_count: 1389,
    game_count: 1373,
    genre_recipe_count: 16,
    mechanic_count: 248,
    variable_count: 44,
    menu_count: 23,
    variable_enrichment_pct: 12.5,
    menu_enrichment_pct: 30,
    mechanic_enrichment_pct: 4,
    skill_enrichment_pct: 100,
    avg_signature_count: 5.2,
  },
  quality_tiers: [{ label: 'curated', count: 10 }],
  games_by_decade: [{ label: '2010s', count: 100 }],
  top_genres: [{ label: 'Action', count: 200 }],
  mechanic_domains: [{ label: 'combat', count: 50 }],
  mechanic_flavors: [{ label: 'action', count: 80 }],
  mechanic_complexity: [{ label: 'M', count: 40 }],
  top_mechanics: [{ slug: 'jump', name: 'Jump', count: 500 }],
  genre_domain_heatmap: {
    genres: ['Action'],
    domains: ['combat'],
    cells: [{ genre: 'Action', domain: 'combat', count: 10 }],
  },
  cooccurrence: [{ mechanic_a: 'a', mechanic_b: 'b', count: 20, lift: 2.5 }],
  mechanic_stats: {
    enrichment_complete: 10,
    enrichment_needs_info: 238,
    with_map_bindings: 0,
    by_category: [{ label: 'combat', count: 50 }],
  },
  skill_stats: {
    enrichment_complete: 15,
    enrichment_needs_info: 0,
    with_map_bindings: 0,
    by_category: [{ label: 'cognitive', count: 5 }],
  },
  variable_stats: {
    enrichment_complete: 5,
    enrichment_needs_info: 39,
    with_map_bindings: 8,
    by_category: [{ label: 'stat', count: 10 }],
  },
  menu_stats: {
    enrichment_complete: 7,
    enrichment_needs_info: 16,
    with_map_bindings: 5,
    by_menu_type: [{ label: 'main', count: 3 }],
  },
  variable_mechanic_pairs: [{ variable: 'health', mechanic: 'depletable-health-pool' }],
  menu_flow: {
    edge_counts_by_relationship: [{ label: 'opens', count: 5 }],
    top_hubs: [{ slug: 'main-menu', name: 'Main Menu', in_degree: 0, out_degree: 3, total_degree: 3 }],
  },
  signature_distribution: [{ label: '5', count: 100 }],
  insights: ['Test insight one.', 'Test insight two.'],
}

describe('AnalyticsSnapshot shape', () => {
  it('has required sections for the analytics page', () => {
    expect(mockAnalytics.overview.game_count).toBeGreaterThan(0)
    expect(mockAnalytics.top_mechanics.length).toBeGreaterThan(0)
    expect(mockAnalytics.cooccurrence[0].lift).toBeGreaterThan(0)
    expect(mockAnalytics.insights.length).toBeGreaterThanOrEqual(2)
    expect(mockAnalytics.genre_domain_heatmap.cells.length).toBeGreaterThan(0)
  })
})
