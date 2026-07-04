import { describe, expect, it } from 'vitest'
import { formatEnrichmentLabel, resolveEnrichmentStatus } from './enrichmentStatus'

describe('enrichmentStatus', () => {
  it('defaults unknown values to needs_info', () => {
    expect(resolveEnrichmentStatus(undefined)).toBe('needs_info')
    expect(resolveEnrichmentStatus('')).toBe('needs_info')
    expect(resolveEnrichmentStatus('bogus')).toBe('needs_info')
  })

  it('recognizes complete', () => {
    expect(resolveEnrichmentStatus('complete')).toBe('complete')
  })

  it('formats human-readable labels', () => {
    expect(formatEnrichmentLabel('needs_info')).toBe('needs info')
    expect(formatEnrichmentLabel('complete')).toBe('complete')
  })
})
