import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { Layout } from './components/Layout'
import { FilterProvider } from './context/Filters'
import { ApiDocsPage } from './pages/ApiDocsPage'
import { ContributePage } from './pages/ContributePage'
import { CooccurrencePage } from './pages/CooccurrencePage'
import { GameDetailPage } from './pages/GameDetailPage'
import { GamesPage } from './pages/GamesPage'
import { GenreDetailPage, GenresPage } from './pages/GenresPage'
import { HomePage } from './pages/HomePage'
import { MechanicDetailPage } from './pages/MechanicDetailPage'
import { MechanicsPage } from './pages/MechanicsPage'
import { WebMcpDocsPage } from './pages/WebMcpDocsPage'

const basename = import.meta.env.BASE_URL.replace(/\/$/, '') || '/'

const router = createBrowserRouter(
  [
    {
      path: '/',
      element: <Layout />,
      children: [
        { index: true, element: <HomePage /> },
        { path: 'games', element: <GamesPage /> },
        { path: 'games/:slug', element: <GameDetailPage /> },
        { path: 'mechanics', element: <MechanicsPage /> },
        { path: 'mechanics/:slug', element: <MechanicDetailPage /> },
        { path: 'genres', element: <GenresPage /> },
        { path: 'genres/:slug', element: <GenreDetailPage /> },
        { path: 'explore/cooccurrence', element: <CooccurrencePage /> },
        { path: 'contribute', element: <ContributePage /> },
        { path: 'docs/api', element: <ApiDocsPage /> },
        { path: 'docs/webmcp', element: <WebMcpDocsPage /> },
      ],
    },
  ],
  { basename: basename === '/' ? undefined : basename },
)

const queryClient = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <FilterProvider>
        <RouterProvider router={router} />
      </FilterProvider>
    </QueryClientProvider>
  )
}
