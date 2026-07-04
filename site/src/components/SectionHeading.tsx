import type { ReactNode } from 'react'
import { ShareSectionLink } from './ShareSectionLink'

export function SectionHeading({
  sectionId,
  children,
}: {
  sectionId: string
  children: ReactNode
}) {
  return (
    <h2>
      {children}
      <ShareSectionLink sectionId={sectionId} />
    </h2>
  )
}
