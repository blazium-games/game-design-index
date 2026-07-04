import { useCallback, useState } from 'react'
import { useLocation } from 'react-router-dom'
import { SITE_URL } from '../types'

export function ShareSectionLink({
  sectionId,
  label = 'Copy section link',
}: {
  sectionId: string
  label?: string
}) {
  const location = useLocation()
  const [copied, setCopied] = useState(false)

  const copy = useCallback(async () => {
    const base = SITE_URL.replace(/\/$/, '')
    const path = location.pathname.replace(/\/$/, '') || '/'
    const url = `${base}${path === '/' ? '' : path}#${sectionId}`
    try {
      await navigator.clipboard.writeText(url)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 2000)
    } catch {
      window.prompt('Copy link:', url)
    }
  }, [location.pathname, sectionId])

  return (
    <button type="button" className="share-section-link" onClick={copy} title={label}>
      {copied ? 'Copied' : 'Share'}
    </button>
  )
}
