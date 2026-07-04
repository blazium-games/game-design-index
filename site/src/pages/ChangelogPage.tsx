import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { api } from '../api/client'
import { RELEASE_URL } from '../types'
import { pageTitle } from '../seo/meta'
import { DocumentMeta } from '../seo/usePageMeta'

export function ChangelogPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['changelog'],
    queryFn: api.fetchChangelog,
  })

  if (isLoading) return <p className="meta">Loading changelog…</p>
  if (error || !data) return <p className="meta">Failed to load changelog.</p>

  return (
    <div className="changelog-page">
      <DocumentMeta
        title={pageTitle('Changelog')}
        description="Release history for the Game Design Index data and site."
        path="/changelog"
      />
      <section className="hero">
        <h1>Changelog</h1>
        <p>
          Release history for the Game Design Index data corpus and public site. Data releases are
          versioned separately from site UI deploys.
        </p>
        <p className="meta">
          <a href={RELEASE_URL} target="_blank" rel="noreferrer">
            Latest GitHub release
          </a>
          {' · '}
          <Link to="/explore/analytics">Corpus analytics</Link>
        </p>
      </section>

      <div className="changelog-entries">
        {data.entries.map((entry) => (
          <article key={entry.version} className="changelog-entry card">
            <header className="changelog-header">
              <span className="changelog-version">v{entry.version}</span>
              <time className="meta" dateTime={entry.date}>
                {entry.date}
              </time>
            </header>
            <h2>{entry.title}</h2>
            {entry.highlights && entry.highlights.length > 0 && (
              <div className="chips">
                {entry.highlights.map((h) => (
                  <span key={h} className="chip">
                    {h}
                  </span>
                ))}
              </div>
            )}
            {entry.sections.map((section) => (
              <div key={section.heading} className="changelog-section">
                <h3>{section.heading}</h3>
                <ul>
                  {section.items.map((item, i) => (
                    <li key={i}>{item}</li>
                  ))}
                </ul>
              </div>
            ))}
            {entry.release_url && (
              <p className="changelog-release-link">
                <a href={entry.release_url} target="_blank" rel="noreferrer">
                  Release notes &amp; download assets
                </a>
              </p>
            )}
          </article>
        ))}
      </div>

      <p className="meta changelog-footer">
        To propose changes for the next release, open an issue or pull request on GitHub.
      </p>
    </div>
  )
}
