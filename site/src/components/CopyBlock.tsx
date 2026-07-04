import { useState } from 'react'

export function CopyBlock({ label, children }: { label?: string; children: string }) {
  const [copied, setCopied] = useState(false)
  async function copy() {
    await navigator.clipboard.writeText(children)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }
  return (
    <div className="copy-block">
      {label && <div className="copy-block-label">{label}</div>}
      <pre>{children}</pre>
      <button type="button" className="copy-btn" onClick={() => void copy()}>
        {copied ? 'Copied' : 'Copy'}
      </button>
    </div>
  )
}
