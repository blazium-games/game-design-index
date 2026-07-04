import { describe, expect, it } from 'vitest'
import type { AnalyticsSnapshot } from '../types'
import {
  analyticsMissingEnrichmentStats,
  normalizeAnalyticsSnapshot,
} from './normalizeAnalytics'

describe('normalizeAnalyticsSnapshot', () => {
  it('fills missing mechanic_stats and skill_stats from legacy payloads', () => {
    const raw: Partial<AnalyticsSnapshot> = {
      schema_version: '1.0',
      overview: {
        map_count: 0,
        game_count: 10,
        genre_recipe_count: 0,
        mechanic_count: 248,
        variable_count: 44,
        menu_count: 23,
        variable_enrichment_pct: 0,
        menu_enrichment_pct: 0,
        avg_signature_count: 0,
      },
      variable_stats: {
        enrichment_complete: 5,
        enrichment_needs_info: 39,
        with_map_bindings: 8,
      },
      menu_stats: {
        enrichment_complete: 7,
        enrichment_needs_info: 16,
        with_map_bindings: 5,
      },
    }
    expect(analyticsMissingEnrichmentStats(raw)).toBe(true)
    const snap = normalizeAnalyticsSnapshot(raw)
    expect(snap.mechanic_stats.enrichment_complete).toBe(0)
    expect(snap.skill_stats.by_category).toEqual([])
    expect(snap.overview.game_count).toBe(10)
    expect(snap.insights).toEqual([])
  })

  it('preserves enrichment stats when present', () => {
    const raw = {
      mechanic_stats: {
        enrichment_complete: 10,
        enrichment_needs_info: 238,
        with_map_bindings: 0,
        by_category: [{ label: 'combat', count: 87 }],
      },
      skill_stats: {
        enrichment_complete: 15,
        enrichment_needs_info: 0,
        with_map_bindings: 0,
      },
    }
    expect(analyticsMissingEnrichmentStats(raw)).toBe(false)
    const snap = normalizeAnalyticsSnapshot(raw)
    expect(snap.mechanic_stats.enrichment_complete).toBe(10)
    expect(snap.mechanic_stats.by_category?.[0].label).toBe('combat')
  })
})
