import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from './api/client'
import { useFilters } from './context/Filters'
import { setupWebMCP } from './webmcp/init'

export function WebMCPBridge() {
  const navigate = useNavigate()
  const { setGames, setMechanics } = useFilters()

  useEffect(() => {
    let teardown: (() => void) | undefined

    void setupWebMCP({ api, navigate, setGames, setMechanics }).then((fn) => {
      teardown = fn
    })

    return () => teardown?.()
  }, [navigate, setGames, setMechanics])

  return null
}
