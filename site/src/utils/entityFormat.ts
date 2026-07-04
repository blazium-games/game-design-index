import type { GameVariable, GameplayMap, MechanicEntry, UIMenu } from '../types'

export type ExportFormat = 'json' | 'md' | 'yaml' | 'xml' | 'txt'
export type EntityKind = 'mechanic' | 'variable' | 'menu' | 'game' | 'genre'

export function formatEntity(
  kind: EntityKind,
  entity: MechanicEntry | GameVariable | UIMenu | GameplayMap,
  format: ExportFormat,
): string {
  if (format === 'json') {
    return JSON.stringify(entity, null, 2) + '\n'
  }
  if (format === 'yaml') {
    return toYAML(entity)
  }
  const md =
    kind === 'mechanic'
      ? mechanicMarkdown(entity as MechanicEntry)
      : kind === 'variable'
        ? variableMarkdown(entity as GameVariable)
        : kind === 'menu'
          ? menuMarkdown(entity as UIMenu)
          : mapMarkdown(entity as GameplayMap)
  if (format === 'txt') {
    return md.replace(/^## /gm, '')
  }
  if (format === 'xml') {
    return kind === 'mechanic'
      ? mechanicXML(entity as MechanicEntry)
      : kind === 'variable'
        ? variableXML(entity as GameVariable)
        : kind === 'menu'
          ? menuXML(entity as UIMenu)
          : mapXML(entity as GameplayMap)
  }
  return md
}

function toYAML(obj: unknown): string {
  const lines: string[] = []
  walkYAML(obj, lines, 0)
  return lines.join('\n') + '\n'
}

function walkYAML(v: unknown, lines: string[], depth: number): void {
  const pad = '  '.repeat(depth)
  if (v === null || v === undefined) return
  if (Array.isArray(v)) {
    for (const item of v) {
      if (typeof item === 'object' && item !== null) {
        lines.push(`${pad}-`)
        walkYAML(item, lines, depth + 1)
      } else {
        lines.push(`${pad}- ${yamlScalar(item)}`)
      }
    }
    return
  }
  if (typeof v === 'object') {
    for (const [k, val] of Object.entries(v as Record<string, unknown>)) {
      if (val === undefined || val === null || val === '') continue
      if (typeof val === 'object') {
        lines.push(`${pad}${k}:`)
        walkYAML(val, lines, depth + 1)
      } else {
        lines.push(`${pad}${k}: ${yamlScalar(val)}`)
      }
    }
  }
}

function yamlScalar(v: unknown): string {
  const s = String(v)
  if (/[:#\n]/.test(s) || s.includes(' ')) return JSON.stringify(s)
  return s
}

function escapeXML(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

function mechanicMarkdown(e: MechanicEntry): string {
  const lines = [`# ${e.name}`, '', `**Flavor:** ${e.flavor} · **Domain:** ${e.domain}`, '']
  if (e.summary) lines.push('## Summary', '', e.summary, '')
  if (e.player_experience) lines.push('## Player Experience', '', e.player_experience, '')
  if (e.featured_in?.length) lines.push('## Featured In', '', ...e.featured_in.map((s) => `- ${s}`), '')
  if (e.synergies?.length) lines.push('## Synergies', '', ...e.synergies.map((s) => `- ${s}`), '')
  if (e.design_guidance?.when_to_use) lines.push('## When to Use', '', e.design_guidance.when_to_use, '')
  return lines.join('\n').trim() + '\n'
}

function variableMarkdown(v: GameVariable): string {
  const lines = [
    `# ${v.name}`,
    '',
    `**Category:** ${v.category} · **Scope:** ${v.scope} · **Value kind:** ${v.value_kind}`,
    '',
  ]
  if (v.summary) lines.push('## Summary', '', v.summary, '')
  if (v.shared_rationale) lines.push('## Shared Rationale', '', v.shared_rationale, '')
  if (v.player_focus) lines.push('## Player Focus', '', v.player_focus, '')
  if (v.typical_range) lines.push('## Typical Range', '', v.typical_range, '')
  if (v.related_mechanics?.length)
    lines.push('## Related Mechanics', '', ...v.related_mechanics.map((m) => `- ${m}`), '')
  if (v.featured_in?.length) lines.push('## Featured In', '', ...v.featured_in.map((s) => `- ${s}`), '')
  return lines.join('\n').trim() + '\n'
}

function menuMarkdown(m: UIMenu): string {
  const lines = [`# ${m.name}`, '', `**Type:** ${m.menu_type} · **Layer:** ${m.layer}`, '']
  if (m.summary) lines.push('## Summary', '', m.summary, '')
  if (m.shared_rationale) lines.push('## Shared Rationale', '', m.shared_rationale, '')
  if (m.typical_actions?.length)
    lines.push('## Typical Actions', '', ...m.typical_actions.map((a) => `- ${a}`), '')
  if (m.related_mechanics?.length)
    lines.push('## Related Mechanics', '', ...m.related_mechanics.map((x) => `- ${x}`), '')
  if (m.related_variables?.length)
    lines.push('## Related Variables', '', ...m.related_variables.map((x) => `- ${x}`), '')
  if (m.featured_in?.length) lines.push('## Featured In', '', ...m.featured_in.map((s) => `- ${s}`), '')
  return lines.join('\n').trim() + '\n'
}

function mapMarkdown(m: GameplayMap): string {
  const lines = [`# ${m.subject.name}`, '']
  if (m.narrative.description) lines.push('## Overview', '', m.narrative.description, '')
  if (m.narrative.core_loop) lines.push('## Core Loop', '', m.narrative.core_loop, '')
  if (m.signature_gameplay?.length)
    lines.push('## Signature Gameplay', '', ...m.signature_gameplay.map((s) => `- ${s}`), '')
  if (m.variables?.length) {
    lines.push('## Variable Bindings', '')
    for (const vb of m.variables) {
      lines.push(`- **${vb.variable_slug}** (${vb.role})${vb.expression ? `: ${vb.expression}` : ''}`)
    }
    lines.push('')
  }
  if (m.ui_menus?.length) {
    lines.push('## UI Menu Bindings', '')
    for (const mb of m.ui_menus) {
      lines.push(`- **${mb.menu_slug}** (${mb.role})${mb.map_notes ? `: ${mb.map_notes}` : ''}`)
    }
    lines.push('')
  }
  return lines.join('\n').trim() + '\n'
}

function mechanicXML(e: MechanicEntry): string {
  return `<?xml version="1.0" encoding="UTF-8"?>\n<mechanic slug="${e.slug}" name="${escapeXML(e.name)}" flavor="${e.flavor}" domain="${e.domain}">\n  <summary>${escapeXML(e.summary)}</summary>\n</mechanic>\n`
}

function variableXML(v: GameVariable): string {
  return `<?xml version="1.0" encoding="UTF-8"?>\n<game_variable slug="${v.slug}" name="${escapeXML(v.name)}" category="${v.category}" scope="${v.scope}" value_kind="${v.value_kind}">\n  <summary>${escapeXML(v.summary)}</summary>\n</game_variable>\n`
}

function menuXML(m: UIMenu): string {
  return `<?xml version="1.0" encoding="UTF-8"?>\n<ui_menu slug="${m.slug}" name="${escapeXML(m.name)}" menu_type="${m.menu_type}" layer="${m.layer}">\n  <summary>${escapeXML(m.summary)}</summary>\n</ui_menu>\n`
}

function mapXML(m: GameplayMap): string {
  return `<?xml version="1.0" encoding="UTF-8"?>\n<gameplay_map slug="${m.slug}" title="${escapeXML(m.title)}" name="${escapeXML(m.subject.name)}" map_type="${m.map_type}"/>\n`
}
