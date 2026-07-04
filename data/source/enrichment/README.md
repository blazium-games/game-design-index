# Enrichment patches

Curated design pedagogy content (GDD guidance, agent context, skills links) is authored as **local patch files**, merged into canonical catalogs with `scripts/apply_enrichment`, then published via PRs that contain only the merged `library/` and `maps/` changes.

## Workflow

1. Open a GitHub issue using **Enrich mechanic** (or variable/menu templates for other entities).
2. Author a cohort patch file under `cohorts/` locally (see naming below). **Cohort JSON is gitignored** — it stays on your machine for review and re-application.
3. Run `go run ./scripts/apply_enrichment` from the repo root to merge patches into `library/` and `maps/`.
4. Open a PR with the **merged library and map changes only** (not the cohort file).

## Cohort file naming

Use `{YYYY-MM}-{short-label}-v{release}.json`, for example:

- `2026-07-pilot-v1.1.0.json` — initial pilot: 10 mechanics, 3 variables, 2 menus, 5 maps

Each file may include a `_meta` block describing the cohort. Patch keys:

| Key | Target |
|-----|--------|
| `mechanics` | Slug → fields merged into `library/mechanics.json` |
| `variables` | Slug → fields merged into `library/variables.json` |
| `ui_menus` | Slug → fields merged into `library/ui-menus.json` |
| `skills` | Slug → fields merged into `library/skills.json` |
| `skills_add` | New skill entries appended to `library/skills.json` |
| `maps` | Map slug → fields merged into `maps/{slug}.json` |
| `maps_add` | New map files written to `maps/` if absent |

See `patch-format.schema.json` for the documented shape.

## Apply script

```powershell
# Apply all cohorts in data/source/enrichment/cohorts/ (local files)
go run ./scripts/apply_enrichment

# Apply one cohort only
go run ./scripts/apply_enrichment -cohort 2026-07-pilot-v1.1.0.json
```

Bulk cohort generators (e.g. pedagogy or catalog fill scripts) are **local-only tooling** and are not committed to this repository. Author patches manually or with your own local scripts, then apply with `apply_enrichment`.

## Enrichment status

Mechanics are **complete** when both `design_guidance.when_to_use` and `agent_context.summary_for_agents` are populated (schema 1.3 pedagogy). Export and the site analytics page report coverage counts.

Gold-standard depth (pilot reference: `ability-driven-backtracking`) adds full `design_guidance`, `agent_context`, `design_exercises` with constraints, map-linked `examples`, `featured_in`, `synergies`, `synergy_notes`, and domain-tuned `parameter_knobs`.
