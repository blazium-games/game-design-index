import { Helmet } from 'react-helmet-async'
import { TWITTER_HANDLE } from '../types'
import type { PageMetaInput } from './meta'
import { buildCanonical, resolvePageMeta } from './meta'

export function PageMeta(props: Partial<PageMetaInput> & { path: string }) {
  const meta = resolvePageMeta(props)
  const canonical = buildCanonical(meta.path)
  const jsonLdBlocks = meta.jsonLd
    ? Array.isArray(meta.jsonLd)
      ? meta.jsonLd
      : [meta.jsonLd]
    : []

  return (
    <Helmet>
      <title>{meta.title}</title>
      <meta name="description" content={meta.description} />
      <link rel="canonical" href={canonical} />
      {meta.noIndex && <meta name="robots" content="noindex" />}

      <meta property="og:title" content={meta.title} />
      <meta property="og:description" content={meta.description} />
      <meta property="og:url" content={canonical} />
      <meta property="og:type" content={meta.ogType ?? 'website'} />
      <meta property="og:site_name" content="Game Design Index" />
      <meta property="og:image" content={meta.image} />
      <meta property="og:image:width" content="1200" />
      <meta property="og:image:height" content="630" />
      <meta property="og:locale" content="en_US" />

      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:site" content={TWITTER_HANDLE} />
      <meta name="twitter:creator" content={TWITTER_HANDLE} />
      <meta name="twitter:title" content={meta.title} />
      <meta name="twitter:description" content={meta.description} />
      <meta name="twitter:image" content={meta.image} />

      <meta name="theme-color" content="#1a1d24" />

      {jsonLdBlocks.map((block, i) => (
        <script key={i} type="application/ld+json">
          {JSON.stringify(block)}
        </script>
      ))}
    </Helmet>
  )
}
