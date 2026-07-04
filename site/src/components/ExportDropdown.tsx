import { useState } from 'react'
import type { GameVariable, GameplayMap, MechanicEntry, UIMenu } from '../types'
import { RELEASE_URL } from '../types'
import { formatEntity, type EntityKind, type ExportFormat } from '../utils/entityFormat'

const FORMATS: { id: ExportFormat; label: string; ext: string }[] = [
  { id: 'json', label: 'JSON', ext: 'json' },
  { id: 'md', label: 'Markdown', ext: 'md' },
  { id: 'yaml', label: 'YAML', ext: 'yaml' },
  { id: 'xml', label: 'XML', ext: 'xml' },
  { id: 'txt', label: 'Plain text', ext: 'txt' },
]

export function ExportDropdown({
  kind,
  slug,
  entity,
}: {
  kind: EntityKind
  slug: string
  entity: MechanicEntry | GameVariable | UIMenu | GameplayMap
}) {
  const [open, setOpen] = useState(false)

  function download(format: ExportFormat) {
    const content = formatEntity(kind, entity, format)
    const ext = FORMATS.find((f) => f.id === format)?.ext ?? format
    const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${slug}.${ext}`
    a.click()
    URL.revokeObjectURL(url)
    setOpen(false)
  }

  return (
    <div className="export-dropdown">
      <button type="button" className="btn secondary" onClick={() => setOpen((o) => !o)}>
        Export ▾
      </button>
      {open && (
        <div className="export-menu">
          {FORMATS.map((f) => (
            <button key={f.id} type="button" onClick={() => download(f.id)}>
              {f.label}
            </button>
          ))}
          <hr />
          <a href={RELEASE_URL} target="_blank" rel="noreferrer" onClick={() => setOpen(false)}>
            Bulk formats (release zip)
          </a>
        </div>
      )}
    </div>
  )
}
