import { describe, expect, it } from 'vitest'
import { TOOL_NAMES } from './types'

describe('WebMCP tool catalog', () => {
  it('lists 18 tools', () => {
    expect(TOOL_NAMES).toHaveLength(18)
    expect(TOOL_NAMES).toContain('get-game')
    expect(TOOL_NAMES).toContain('compose-design-brief')
    expect(TOOL_NAMES).toContain('navigate')
  })
})
