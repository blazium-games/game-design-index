import { describe, expect, it } from 'vitest'
import { formatChartLabel, formatSlugLabel, truncateLabel } from './formatChartLabel'

describe('formatChartLabel', () => {
  it('maps domain and flavor keys', () => {
    expect(formatChartLabel('progression')).toBe('Progression')
    expect(formatChartLabel('action')).toBe('Action')
  })

  it('maps complexity and quality keys', () => {
    expect(formatChartLabel('unset')).toBe('Not classified')
    expect(formatChartLabel('S')).toBe('Small')
    expect(formatChartLabel('curated')).toBe('Curated')
  })

  it('passes through decade and mechanic names', () => {
    expect(formatChartLabel('2010s')).toBe('2010s')
    expect(formatChartLabel('Environmental Trigger Switch')).toBe('Environmental Trigger Switch')
  })
})

describe('truncateLabel', () => {
  it('truncates with ellipsis', () => {
    expect(truncateLabel('Environmental Trigger Switch', 20)).toBe('Environmental Trigg…')
  })
})

describe('formatSlugLabel', () => {
  it('title-cases hyphenated slugs', () => {
    expect(formatSlugLabel('permadeath-system')).toBe('Permadeath System')
  })
})
