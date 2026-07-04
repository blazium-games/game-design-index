import { useEffect } from 'react'
import { useLocation } from 'react-router-dom'

export function useHashScroll(enabled = true) {
  const location = useLocation()

  useEffect(() => {
    if (!enabled || !location.hash) return
    const id = location.hash.slice(1)
    const timer = window.setTimeout(() => {
      document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }, 100)
    return () => window.clearTimeout(timer)
  }, [enabled, location.hash, location.pathname])
}
