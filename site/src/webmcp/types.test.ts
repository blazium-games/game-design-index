import { describe, expect, it } from 'vitest'
import { TOOL_NAMES } from './types'

describe('WebMCP tool catalog', () => {
  it('lists registered tools including analytics', () => {
    expect(TOOL_NAMES).toHaveLength(25)
    expect(TOOL_NAMES).toContain('get-analytics')
    expect(TOOL_NAMES).toContain('get-game')
    expect(TOOL_NAMES).toContain('compose-design-brief')
    expect(TOOL_NAMES).toContain('list-skills')
    expect(TOOL_NAMES).toContain('get-skill')
  })
})
