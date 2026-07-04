import { useMemo } from 'react'
import { useLocation } from 'react-router-dom'
import { PageMeta } from './PageMeta'
import type { PageMetaInput } from './meta'

export function usePageMeta(input: Omit<Partial<PageMetaInput>, 'path'> & { path?: string }) {
  const location = useLocation()
  const path = input.path ?? location.pathname
  return useMemo(
    () => ({ ...input, path }),
    [path, input.title, input.description, input.ogType, input.image, input.noIndex],
  )
}

export function DocumentMeta(props: Partial<PageMetaInput> & { path?: string }) {
  const location = useLocation()
  const resolved = usePageMeta({ ...props, path: props.path ?? location.pathname })
  return <PageMeta {...resolved} />
}
