import { registerExploreTools } from './explore'
import { registerQueryTools } from './query'
import { registerUiTools } from './ui'
import type { WebMCPDeps } from '../types'

export async function registerIndexTools(deps: WebMCPDeps, opts: { signal?: AbortSignal }) {
  await registerQueryTools(deps, opts)
  await registerExploreTools(deps, opts)
  await registerUiTools(deps, opts)
}
