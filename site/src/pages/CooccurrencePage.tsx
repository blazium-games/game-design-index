import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { api } from '../api/client'

export function CooccurrencePage() {
  const { data: pairs } = useQuery({
    queryKey: ['cooccurrence'],
    queryFn: () => api.fetchCooccurrence(100),
  })
  return (
    <div>
      <h1>Mechanic co-occurrence</h1>
      <p className="meta">Top pairs that appear together on gameplay maps</p>
      <table className="table">
        <thead>
          <tr>
            <th>Mechanic A</th>
            <th>Mechanic B</th>
            <th>Count</th>
          </tr>
        </thead>
        <tbody>
          {(pairs ?? []).map((p) => (
            <tr key={`${p.mechanic_a}-${p.mechanic_b}`}>
              <td>
                <Link to={`/mechanics/${p.mechanic_a}`}>{p.mechanic_a}</Link>
              </td>
              <td>
                <Link to={`/mechanics/${p.mechanic_b}`}>{p.mechanic_b}</Link>
              </td>
              <td>{p.count}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
