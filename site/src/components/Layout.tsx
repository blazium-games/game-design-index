import { Link, Outlet } from 'react-router-dom'
import { REPO_URL } from '../types'
import { WebMCPBridge } from '../WebMCPBridge'

export function Layout() {
  return (
    <div className="app">
      <WebMCPBridge />
      <header className="header">
        <Link to="/" className="logo">
          Game Mechanics Index
        </Link>
        <nav>
          <Link to="/games">Games</Link>
          <Link to="/mechanics">Mechanics</Link>
          <Link to="/variables">Variables</Link>
          <Link to="/ui-menus">UI Menus</Link>
          <Link to="/genres">Genres</Link>
          <Link to="/explore/cooccurrence">Co-occurrence</Link>
          <Link to="/contribute">Contribute</Link>
          <Link to="/docs/api">API</Link>
          <Link to="/docs/webmcp">WebMCP</Link>
        </nav>
      </header>
      <main className="main">
        <Outlet />
      </main>
      <footer className="footer">
        <span>MIT License · Blazium Games</span>
        <a href={REPO_URL} target="_blank" rel="noreferrer">
          GitHub
        </a>
        <a href={`${REPO_URL}/issues/new/choose`} target="_blank" rel="noreferrer">
          Open an Issue
        </a>
      </footer>
    </div>
  )
}

export function SuggestEditLink({
  slug,
  kind,
}: {
  slug: string
  kind: 'game' | 'mechanic' | 'genre' | 'variable' | 'ui-menu'
}) {
  const title = encodeURIComponent(`Data correction: ${slug}`)
  const body = encodeURIComponent(
    `**${kind} slug:** \`${slug}\`\n\n**What should change:**\n\n`,
  )
  const url = `${REPO_URL}/issues/new?template=data-correction.yml&title=${title}&body=${body}`
  return (
    <a className="suggest-edit" href={url} target="_blank" rel="noreferrer">
      Suggest an edit
    </a>
  )
}
