import { formatEnrichmentLabel, resolveEnrichmentStatus } from '../utils/enrichmentStatus'

export function EnrichmentStatusChip({ status }: { status?: string }) {
  const resolved = resolveEnrichmentStatus(status)
  return (
    <span className={`chip ${resolved === 'needs_info' ? 'chip-warn' : ''}`}>
      {formatEnrichmentLabel(resolved)}
    </span>
  )
}
