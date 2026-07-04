# Contributing

All changes must go through a **feature branch** and **pull request**. Do not push directly to `main`.

## Workflow

1. Create a branch from `main`:
   ```bash
   git checkout main
   git pull
   git checkout -b feat/your-change
   ```
2. Edit corpus files under `data/source/` (maps, library, schema) and/or site code under `site/`.
3. Validate:
   ```bash
   go test ./...
   go run ./cmd/export -version 0.1.0
   go run ./scripts/lint_public_copy
   cd site && npm run build && npm run test:webmcp
   ```
4. Open a pull request against `main`. CI runs `pr-validate.yml` on every PR.
5. Merge via GitHub after checks pass.

## Releases

- `main` deploys GitHub Pages on merge (see `.github/workflows/deploy-pages.yml`).
- Tagged releases are published via `.github/workflows/publish-release.yml`.
