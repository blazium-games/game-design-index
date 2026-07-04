export type EnrichmentStatus = 'complete' | 'needs_info'

export function resolveEnrichmentStatus(raw?: string): EnrichmentStatus {
  return raw === 'complete' ? 'complete' : 'needs_info'
}

export function formatEnrichmentLabel(status: EnrichmentStatus): string {
  return status === 'complete' ? 'complete' : 'needs info'
}
