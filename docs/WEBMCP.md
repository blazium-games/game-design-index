# WebMCP agent tools

The live site registers [WebMCP](https://github.com/webmachinelearning/webmcp) tools on `document.modelContext` so browser-integrated AI agents can query the Game Design Index with structured schemas.

## Quick start

1. Open [the live index](https://blazium-games.github.io/game-design-index/).
2. In DevTools: `document.modelContext.getTools()` — confirm tools are registered.
3. Paste a task prompt (see [Cursor workflow](#cursor-workflow)) and let the agent call tools.

## Cursor workflow

1. Open the index URL in Cursor's integrated browser.
2. Paste the user prompt below into chat.
3. The agent should call `document.modelContext` tools — not scrape HTML.

**User prompt** (copy from `/docs/webmcp#cursor` on the site):

```text
You are on the Game Design Index. Use document.modelContext WebMCP tools (not HTML scraping).
Task: Find curated metroidvania-style games, list their signature mechanics, and draft a one-page design brief combining the top two references.
Tools to use: search-index, get-game, compose-design-brief.
```

## Static API fallback (no browser)

```bash
curl -s https://blazium-games.github.io/game-design-index/api/v1/catalog.json
curl -s https://blazium-games.github.io/game-design-index/api/v1/mechanics/boss-weakness-network.json
curl -s https://blazium-games.github.io/game-design-index/formats/v1/mechanics/boss-weakness-network.md
```

## Copy-paste: project rule

See [docs/snippets/webmcp-cursor-rule.md](snippets/webmcp-cursor-rule.md) or paste from `/docs/webmcp#rules`.

## Copy-paste: system context

```text
Game Design Index — 1389 games, 248 mechanics, schema 1.2.
WebMCP tools on document.modelContext return { content: [{ type: "text", text: "<json>" }] }.
Catalog: get-catalog. Search: search-index { query }. Game: get-game { slug }. Mechanic: get-mechanic { slug }. Formatted: get-mechanic-formatted { slug, format: "md" }. GDD seed: compose-design-brief { ref_slugs: ["hollow-knight", "dead-cells"] }.
```

## Runtime

- Native WebMCP when the browser supports it (polyfill chunk is not loaded)
- [`@mcp-b/webmcp-polyfill`](https://www.npmjs.com/package/@mcp-b/webmcp-polyfill) otherwise (lazy-loaded after first paint; testing shim disabled in production)

Tools return MCP-style payloads: `{ content: [{ type: "text", text: "<json>" }] }`.

## Tool catalog

| Tool | Description |
|------|-------------|
| `get-catalog` | Corpus stats and schema version |
| `search-index` | Search games and mechanics |
| `list-games` | Filter game index rows |
| `get-game` | Full gameplay map by slug |
| `list-mechanics` | Filter mechanic index rows |
| `get-mechanic` | Full mechanic entry |
| `get-mechanic-formatted` | Mechanic as md/yaml/txt from `formats/v1/` |
| `list-genres` | Genre recipe list |
| `get-genre` | Genre recipe map |
| `get-mechanic-maps` | Maps featuring a mechanic |
| `get-cooccurrence` | Top mechanic pairs |
| `get-similar-games` | Games sharing signatures |
| `compose-design-brief` | Merge 1–4 reference games |
| `list-tags` | Tag vocabulary |
| `list-variables` | Filter game variable catalog entries |
| `get-variable` | Full variable entry by slug |
| `list-ui-menus` | Filter UI menu catalog entries |
| `get-ui-menu` | Full UI menu entry by slug |
| `list-skills` | Filter design skill catalog entries |
| `get-skill` | Full design skill entry by slug |
| `get-analytics` | Pre-computed corpus analytics and enrichment stats |
| `navigate` | Navigate catalog UI |
| `filter-games-view` | Apply games filters + navigate |
| `filter-mechanics-view` | Apply mechanics filters |
| `open-contribute` | Open contribution flow |

Full schemas: browse to `/docs/webmcp` on the live site or read `site/public/webmcp-tools.json`.

## Example agent prompts

- "What are Hollow Knight's signature mechanics?"
- "List curated metroidvania games"
- "Compose a design brief from hades and dead-cells"
- "Fetch boss-weakness-network as Markdown for a GDD section"

## Verify tools

```javascript
document.modelContext.getTools().then((t) => console.log(t.map((x) => x.name)))
```

## Testing

```powershell
cd site
npm run test:webmcp
```

With the polyfill testing shim enabled locally, use `navigator.modelContextTesting.executeTool(...)`.

## Troubleshooting

**Browser extension noise** — Warnings from `contentscript.js`, `ObjectMultiplex`, or `MaxListenersExceededWarning` usually come from wallet or agent extensions (e.g. MetaMask), not this site. Test in an incognito window with extensions disabled to see only site messages.

**CSP / unsafe-eval** — The production bundle does not use `eval()` or `new Function()`. The site sets a strict Content-Security-Policy without `'unsafe-eval'`. If you still see a CSP violation attributed to `index-*.js`, confirm it reproduces with extensions off. Native Chrome WebMCP uses `document.modelContext` directly; the polyfill is only downloaded when native support is missing.

## License

MIT — same as the index data.
