import { isRouteErrorResponse, Link, useNavigate, useRouteError } from 'react-router-dom'

export function RouteErrorFallback() {
  const error = useRouteError()
  const navigate = useNavigate()

  let message = 'Something went wrong loading this page.'
  if (isRouteErrorResponse(error)) {
    message = error.statusText || message
  } else if (error instanceof Error) {
    message = error.message
  }

  return (
    <div className="error-fallback">
      <h1>Page error</h1>
      <p className="meta">{message}</p>
      <div className="cta-actions">
        <button type="button" className="btn" onClick={() => window.location.reload()}>
          Try again
        </button>
        <button type="button" className="btn secondary" onClick={() => navigate('/')}>
          Go home
        </button>
        <Link className="btn secondary" to="/explore/analytics">
          Analytics
        </Link>
      </div>
    </div>
  )
}
