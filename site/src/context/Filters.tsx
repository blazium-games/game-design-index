import { createContext, useContext, useState, type ReactNode } from 'react'

export interface GameFilters {
  genre: string
  tier: string
  query: string
}

export interface MechanicFilters {
  domain: string
  flavor: string
  tag: string
  query: string
  enrichmentStatus: string
}

interface FilterState {
  games: GameFilters
  mechanics: MechanicFilters
  setGames: (f: Partial<GameFilters>) => void
  setMechanics: (f: Partial<MechanicFilters>) => void
}

const FilterContext = createContext<FilterState | null>(null)

export function FilterProvider({ children }: { children: ReactNode }) {
  const [games, setGamesState] = useState<GameFilters>({ genre: '', tier: '', query: '' })
  const [mechanics, setMechState] = useState<MechanicFilters>({
    domain: '',
    flavor: '',
    tag: '',
    query: '',
    enrichmentStatus: '',
  })
  return (
    <FilterContext.Provider
      value={{
        games,
        mechanics,
        setGames: (f) => setGamesState((s) => ({ ...s, ...f })),
        setMechanics: (f) => setMechState((s) => ({ ...s, ...f })),
      }}
    >
      {children}
    </FilterContext.Provider>
  )
}

export function useFilters() {
  const ctx = useContext(FilterContext)
  if (!ctx) throw new Error('useFilters outside provider')
  return ctx
}
