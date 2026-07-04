import { Link } from 'react-router-dom'

export type RankedRow = {
  slug?: string
  label: string
  count: number
}

export function RankedTable({
  rows,
  labelHeader = 'Item',
  countHeader = 'Count',
  linkPrefix,
}: {
  rows: RankedRow[]
  labelHeader?: string
  countHeader?: string
  linkPrefix?: string
}) {
  if (!rows.length) return <p className="meta">No data</p>
  return (
    <table className="table ranked-table">
      <thead>
        <tr>
          <th className="rank-col">#</th>
          <th>{labelHeader}</th>
          <th className="count-col">{countHeader}</th>
        </tr>
      </thead>
      <tbody>
        {rows.map((row, i) => (
          <tr key={row.slug ?? row.label}>
            <td className="rank-col">{i + 1}</td>
            <td>
              {linkPrefix && row.slug ? (
                <Link to={`${linkPrefix}/${row.slug}`}>{row.label}</Link>
              ) : (
                row.label
              )}
            </td>
            <td className="count-col">{row.count.toLocaleString()}</td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
