import { useQuery } from '@tanstack/react-query'
import { candidatesApi } from '@/services/api'

export const APPLICATION_DETAIL_QUERY_KEY = (id) => ['applications', id]

/** True when we expect AI screening to run but don't have ai_processed_at yet. */
export function isAIScreeningPending(detail) {
  if (!detail?.id) return false
  if (detail.ai_processed_at) return false
  return Boolean(detail.resume_file_name)
}

export function useApplicationDetail(applicationId, options = {}) {
  const enabled = Boolean(applicationId) && (options.enabled !== false)
  return useQuery({
    queryKey: APPLICATION_DETAIL_QUERY_KEY(applicationId),
    queryFn: () => candidatesApi.getById(applicationId),
    enabled,
    refetchInterval: (query) => {
      const d = query.state.data
      if (!d) return false
      if (d.ai_processed_at) return false
      if (!d.resume_file_name) return false
      return 3000
    },
    ...options,
  })
}
