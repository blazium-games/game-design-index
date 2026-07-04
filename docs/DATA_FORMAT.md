# Data format (schema 1.2)

Neutral field guide for contributors to the game mechanics index.

## Schema versions

| Document | Version | Notes |
|----------|---------|-------|
| Mechanic entries | **1.2** | All entries in `library/mechanics.json` |
| Gameplay maps | **1.1** (or 1.2) | Maps migrate independently; new `gdd_outline` / `mechanic_relationships` optional |

## Gameplay map (`data/source/maps/{slug}.json`)

| Field | Required | Notes |
|-------|----------|-------|
| `schema_version` | yes | `"1.1"` or `"1.2"` |
| `slug` | yes | kebab-case stable ID |
| `title` | yes | Display title |
| `map_type` | yes | `game` or `genre` |
| `subject` | yes | `name`, `genres[]`, optional `influences[]` |
| `narrative` | yes | `description`, `core_loop`, `skills_tested[]` |
| `signature_gameplay` | yes | 3–7 mechanic slugs |
| `mechanics` | yes | Bindings: `mechanic_slug`, `role`, `map_notes`, `domain` |
| `gdd_outline` | no | GDD seed: `overview`, `core_loop`, `player_goals[]`, `constraints[]`, `progression_notes`, `combat_notes`, `economy_notes` |
| `mechanic_relationships[]` | no | Map-local edges: `from_mechanic`, `to_mechanic`, `relationship`, `notes` |

### Binding roles

- `signature` — DNA mechanic for this game
- `supporting` — important but not identity-defining
- `common` — genre-standard presence

### Flavors (mechanic entries)

Each mechanic has exactly one **flavor**: `action`, `adventure`, or `strategy`.

### Domains

`locomotion`, `combat`, `progression`, `economy`, `level`, `session`

## Mechanic entry (`library/mechanics.json`)

### Required (1.1 core)

| Field | Required |
|-------|----------|
| `slug`, `name`, `flavor`, `domain`, `summary`, `tags[]` | yes |
| `requirements[]` | recommended |
| `synergies[]`, `conflicts[]` | optional |
| `signature_of.games[]` | when known |

### Schema 1.2 additions (GDD / AI)

| Field | Purpose |
|-------|---------|
| `flavor_rationale` | Why the flavor classification fits (Category Insights) |
| `design_guidance.when_to_use` | When to apply this mechanic |
| `design_guidance.where_to_use` | Which design layer (boss roster, macro routing, etc.) |
| `design_guidance.when_to_avoid[]` | Contexts to avoid |
| `design_guidance.designer_notes` | Tuning / implementation notes |
| `synergy_notes[]` | `{ slug, note }` — synergies with designer prose |
| `player_experience` | What the player optimizes for |
| `complexity` | `S` / `M` / `L` — public design effort |
| `examples[]` | `{ label, description, map_slug? }` |
| `relationship_model` | Network mechanics: `type` (`directed_cycle` \| `graph`), `edges[]` |
| `agent_context.summary_for_agents` | One-paragraph agent synopsis |
| `agent_context.gdd_prompt` | Ready-made GDD section seed |
| `agent_context.implementation_checklist[]` | Build-order checklist |
| `media.snippet_path` | Internal asset path; **stripped on public export** |

Keep `parameter_knobs`, `requirements`, `synergies` (slug index), `conflicts`, and `anti_patterns` as first-class design fields.

## Quality metadata

`metadata.quality_tier`: `curated` | `catalog` | `template` | `stub`

## Public export

The export pipeline strips `implementation`, `media`, and map provenance URLs before publishing.

Multi-format exports (Markdown, YAML, XML, plain text) are documented in [FORMATS.md](FORMATS.md).
