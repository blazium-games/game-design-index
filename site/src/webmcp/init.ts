import { registerIndexTools } from './tools'
import type { WebMCPDeps } from './types'

export async function setupWebMCP(deps: WebMCPDeps): Promise<() => void> {
  const controller = new AbortController()

  try {
    if (!('modelContext' in document) || !document.modelContext) {
      const { initializeWebMCPPolyfill } = await import('@mcp-b/webmcp-polyfill')
      initializeWebMCPPolyfill({ installTestingShim: false })
    }

    await registerIndexTools(deps, { signal: controller.signal })
  } catch (err) {
    console.warn('WebMCP unavailable:', err)
  }

  return () => controller.abort()
}
