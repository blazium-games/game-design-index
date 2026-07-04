import { REPO_URL } from '../types'

type EnrichKind = 'variable' | 'ui-menu'

export function EmptyField({
  slug,
  field,
  kind,
}: {
  slug: string
  field: string
  kind: EnrichKind
}) {
  const template = kind === 'variable' ? 'enrich-variable.yml' : 'enrich-ui-menu.yml'
  const title = encodeURIComponent(`Enrich ${kind === 'variable' ? 'variable' : 'UI menu'}: ${slug} — ${field}`)
  const body = encodeURIComponent(
    `**Slug:** \`${slug}\`\n**Field:** \`${field}\`\n\n**Suggested content:**\n\n`,
  )
  const url = `${REPO_URL}/issues/new?template=${template}&title=${title}&body=${body}`
  return (
    <p className="empty-field">
      <a href={url} target="_blank" rel="noreferrer">
        Add information here
      </a>
    </p>
  )
}
