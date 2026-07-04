import type { NavigateFunction } from 'react-router-dom'
import type { ApiClient } from '../api/client'
import type { GameFilters, MechanicFilters } from '../context/Filters'

export interface WebMCPDeps {
  api: ApiClient
  navigate: NavigateFunction
  setGames: (f: Partial<GameFilters>) => void
  setMechanics: (f: Partial<MechanicFilters>) => void
}

export const TOOL_NAMES = [
  'get-catalog',
  'get-analytics',
  'search-index',
  'list-games',
  'get-game',
  'list-mechanics',
  'get-mechanic',
  'get-mechanic-formatted',
  'list-genres',
  'get-genre',
  'get-mechanic-maps',
  'get-cooccurrence',
  'get-similar-games',
  'compose-design-brief',
  'list-tags',
  'list-variables',
  'get-variable',
  'list-ui-menus',
  'get-ui-menu',
  'list-skills',
  'get-skill',
  'navigate',
  'filter-games-view',
  'filter-mechanics-view',
  'open-contribute',
] as const
