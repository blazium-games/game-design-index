import { ORG_NAME, SITE_NAME, SITE_URL, TWITTER_URL } from '../types'

export function organizationJsonLd() {
  return {
    '@context': 'https://schema.org',
    '@type': 'Organization',
    name: ORG_NAME,
    url: 'https://blazium.app',
    sameAs: [TWITTER_URL, SITE_URL, 'https://github.com/blazium-games'],
  }
}

export function webSiteJsonLd() {
  return {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    name: SITE_NAME,
    url: SITE_URL,
    description:
      'An open index of video game mechanics, gameplay decomposition, design skills, and genre recipes.',
    publisher: { '@type': 'Organization', name: ORG_NAME },
    potentialAction: {
      '@type': 'SearchAction',
      target: `${SITE_URL}/games?q={search_term_string}`,
      'query-input': 'required name=search_term_string',
    },
  }
}

export function datasetJsonLd(opts: {
  name: string
  description: string
  recordCount?: number
}) {
  return {
    '@context': 'https://schema.org',
    '@type': 'Dataset',
    name: opts.name,
    description: opts.description,
    url: SITE_URL,
    license: 'https://opensource.org/licenses/MIT',
    creator: { '@type': 'Organization', name: ORG_NAME },
    ...(opts.recordCount != null ? { size: `${opts.recordCount} records` } : {}),
  }
}

export function definedTermJsonLd(opts: {
  name: string
  description: string
  url: string
  termSet?: string
}) {
  return {
    '@context': 'https://schema.org',
    '@type': 'DefinedTerm',
    name: opts.name,
    description: opts.description,
    url: opts.url,
    inDefinedTermSet: opts.termSet ?? SITE_NAME,
  }
}

export function learningResourceJsonLd(opts: {
  name: string
  description: string
  url: string
}) {
  return {
    '@context': 'https://schema.org',
    '@type': 'LearningResource',
    name: opts.name,
    description: opts.description,
    url: opts.url,
    learningResourceType: 'game design mechanic',
    provider: { '@type': 'Organization', name: ORG_NAME },
  }
}

export function videoGameJsonLd(opts: {
  name: string
  description: string
  url: string
  genres?: string[]
}) {
  return {
    '@context': 'https://schema.org',
    '@type': 'VideoGame',
    name: opts.name,
    description: opts.description,
    url: opts.url,
    ...(opts.genres?.length ? { genre: opts.genres } : {}),
  }
}

export function breadcrumbJsonLd(items: { name: string; url: string }[]) {
  return {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    itemListElement: items.map((item, i) => ({
      '@type': 'ListItem',
      position: i + 1,
      name: item.name,
      item: item.url,
    })),
  }
}

export function faqPageJsonLd(faqs: { question: string; answer: string }[]) {
  return {
    '@context': 'https://schema.org',
    '@type': 'FAQPage',
    mainEntity: faqs.map((faq) => ({
      '@type': 'Question',
      name: faq.question,
      acceptedAnswer: { '@type': 'Answer', text: faq.answer },
    })),
  }
}
