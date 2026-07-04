import { describe, expect, it } from 'vitest'
import { buildCanonical, buildOgImageUrl, pageTitle } from './meta'
import { SITE_URL, TWITTER_HANDLE } from '../types'

describe('buildCanonical', () => {
  it('builds root canonical URL', () => {
    expect(buildCanonical('/')).toBe(`${SITE_URL}/`)
  })

  it('builds detail page canonical URL', () => {
    expect(buildCanonical('/mechanics/melee-primary-combat')).toBe(
      `${SITE_URL}/mechanics/melee-primary-combat`,
    )
  })
})

describe('buildOgImageUrl', () => {
  it('returns absolute default OG image URL', () => {
    expect(buildOgImageUrl()).toBe(`${SITE_URL}/og-default.png`)
  })

  it('preserves absolute image URLs', () => {
    expect(buildOgImageUrl('https://example.com/img.png')).toBe('https://example.com/img.png')
  })
})

describe('pageTitle', () => {
  it('formats entity titles with site name', () => {
    expect(pageTitle('Melee Primary Combat')).toBe('Melee Primary Combat · Game Design Index')
  })
})

describe('twitter constants', () => {
  it('uses @BlaziumGames handle', () => {
    expect(TWITTER_HANDLE).toBe('@BlaziumGames')
  })
})
