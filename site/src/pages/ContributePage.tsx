import { Link } from 'react-router-dom'
import { REPO_URL } from '../types'

export function ContributePage() {
  return (
    <div>
      <h1>Contribute</h1>
      <p>We are actively looking for additional game data. Help expand the index via GitHub.</p>
      <section>
        <h2>What to add</h2>
        <ul>
          <li>
            <strong>New game map</strong> — subject, narrative, signatures, mechanic bindings with
            map_notes
          </li>
          <li>
            <strong>Enrich existing game</strong> — upgrade quality toward curated with
            game-specific notes
          </li>
          <li>
            <strong>New mechanic</strong> — schema 1.1 entry with domain, tags, summary
          </li>
          <li>
            <strong>Game variables</strong> — enrich catalog entries (shared rationale, player
            focus, typical range) via enrich-variable issues
          </li>
          <li>
            <strong>UI menus</strong> — document screen patterns and map bindings via enrich-ui-menu
            issues
          </li>
          <li>
            <strong>Map bindings</strong> — add variables[] and ui_menus[] to gameplay maps
          </li>
          <li>
            <strong>Corrections</strong> — fix genres, synergies, or signature lists
          </li>
        </ul>
      </section>
      <section>
        <h2>Quality tiers</h2>
        <ul>
          <li>
            <strong>curated</strong> — hand-enriched, game-specific mechanic notes
          </li>
          <li>
            <strong>template</strong> — genre-derived starter maps awaiting enrichment
          </li>
        </ul>
      </section>
      <div className="cta-actions">
        <a className="btn" href={`${REPO_URL}/issues/new/choose`} target="_blank" rel="noreferrer">
          Open an Issue
        </a>
        <a className="btn secondary" href={`${REPO_URL}/blob/main/CONTRIBUTING.md`} target="_blank" rel="noreferrer">
          CONTRIBUTING.md
        </a>
      </div>
      <p>
        See also <Link to="/docs/api">API docs</Link> and <Link to="/docs/webmcp">WebMCP tools</Link>.
      </p>
    </div>
  )
}
