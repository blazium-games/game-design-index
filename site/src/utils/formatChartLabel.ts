const LABEL_MAP: Record<string, string> = {
  locomotion: 'Locomotion',
  combat: 'Combat',
  progression: 'Progression',
  economy: 'Economy',
  level: 'Level',
  session: 'Session',
  action: 'Action',
  adventure: 'Adventure',
  strategy: 'Strategy',
  unset: 'Not classified',
  S: 'Small',
  M: 'Medium',
  L: 'Large',
  curated: 'Curated',
  catalog: 'Catalog',
  template: 'Template',
  stub: 'Stub',
  unknown: 'Unknown year',
}

/** Map corpus keys to human-readable chart labels. */
export function formatChartLabel(raw: string): string {
  if (!raw) return raw
  const lower = raw.toLowerCase()
  if (LABEL_MAP[lower]) return LABEL_MAP[lower]
  if (LABEL_MAP[raw]) return LABEL_MAP[raw]
  return raw
}

/** Truncate long axis labels with ellipsis. */
export function truncateLabel(text: string, max: number): string {
  if (text.length <= max) return text
  return `${text.slice(0, Math.max(0, max - 1))}…`
}

/** Turn mechanic slugs into readable title-case text. */
export function formatSlugLabel(slug: string): string {
  return slug
    .split('-')
    .map((w) => (w ? w[0].toUpperCase() + w.slice(1) : w))
    .join(' ')
}
