import { describe, expect, it } from 'vitest'
import { formatEntity } from './entityFormat'

describe('entityFormat', () => {
  const variable = {
    slug: 'health',
    name: 'Health',
    summary: 'Primary survivability pool.',
    category: 'stat',
    scope: 'player',
    value_kind: 'integer',
    shared_rationale: 'Universal failure boundary.',
    tags: ['combat'],
  }

  it('formats variable as json', () => {
    const out = formatEntity('variable', variable, 'json')
    expect(out).toContain('"slug": "health"')
  })

  it('formats variable as markdown', () => {
    const out = formatEntity('variable', variable, 'md')
    expect(out).toContain('# Health')
    expect(out).toContain('Shared Rationale')
  })

  it('formats variable as xml', () => {
    const out = formatEntity('variable', variable, 'xml')
    expect(out).toContain('<game_variable')
    expect(out).toContain('slug="health"')
  })
})
