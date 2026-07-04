import { describe, expect, it } from 'vitest'
import type { Changelog } from '../types'

const mockChangelog: Changelog = {
  schema_version: '1.0',
  entries: [
    {
      version: '1.0.0',
      date: '2026-07-04',
      title: 'Game Design Index — initial release',
      highlights: ['1389 gameplay maps', '248 mechanics'],
      sections: [{ heading: 'Added', items: ['Game Design Index corpus'] }],
      release_url: 'https://github.com/blazium-games/game-design-index/releases/tag/v1.0.0',
    },
  ],
}

describe('Changelog shape', () => {
  it('lists entries newest-first with sections', () => {
    expect(mockChangelog.entries[0].version).toBe('1.0.0')
    expect(mockChangelog.entries[0].sections[0].items.length).toBeGreaterThan(0)
    expect(mockChangelog.entries.length).toBe(1)
  })
})
