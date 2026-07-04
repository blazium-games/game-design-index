import {
  DEFAULT_OG_IMAGE,
  ORG_NAME,
  SITE_NAME,
  SITE_URL,
  TWITTER_HANDLE,
} from '../types'

export interface PageMetaInput {
  title: string
  description: string
  path: string
  ogType?: string
  image?: string
  jsonLd?: Record<string, unknown> | Record<string, unknown>[]
  noIndex?: boolean
}

const DEFAULT_DESCRIPTION =
  'Open MIT-licensed index of video game mechanics, gameplay maps, design skills, variables, and genre recipes for game designers and AI agents.'

export function buildCanonical(path: string): string {
  const normalized = path.startsWith('/') ? path : `/${path}`
  const base = SITE_URL.replace(/\/$/, '')
  if (normalized === '/' || normalized === '') return `${base}/`
  return `${base}${normalized}`
}

export function buildOgImageUrl(image?: string): string {
  if (!image) return DEFAULT_OG_IMAGE
  if (image.startsWith('http')) return image
  const base = SITE_URL.replace(/\/$/, '')
  return `${base}${image.startsWith('/') ? image : `/${image}`}`
}

export function resolvePageMeta(input: Partial<PageMetaInput> & { path: string }): PageMetaInput {
  const title = input.title ?? SITE_NAME
  return {
    title,
    description: input.description ?? DEFAULT_DESCRIPTION,
    path: input.path,
    ogType: input.ogType ?? 'website',
    image: buildOgImageUrl(input.image),
    jsonLd: input.jsonLd,
    noIndex: input.noIndex,
  }
}

export function pageTitle(entityTitle: string): string {
  return `${entityTitle} · ${SITE_NAME}`
}

export { DEFAULT_DESCRIPTION, ORG_NAME, SITE_NAME, SITE_URL, TWITTER_HANDLE }
