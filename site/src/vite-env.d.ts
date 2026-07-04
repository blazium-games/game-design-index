/// <reference types="vite/client" />

interface ModelContextTool {
  name: string
  description: string
  inputSchema?: object
  execute: (args: Record<string, unknown>) => Promise<{ content: Array<{ type: string; text: string }> }>
}

interface ModelContext {
  registerTool(tool: ModelContextTool, options?: { signal?: AbortSignal }): Promise<void>
}

interface Document {
  modelContext: ModelContext
}
