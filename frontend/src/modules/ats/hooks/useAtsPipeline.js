import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { candidatesApi } from '@/services/api'

export const ATS_PIPELINE_QUERY_KEY = ['ats', 'pipeline']

export function useAtsPipeline(params = {}) {
  return useQuery({
    queryKey: [...ATS_PIPELINE_QUERY_KEY, params],
    queryFn: () => candidatesApi.list(params),
  })
}

export function useUpdateCandidateStatus() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }) => candidatesApi.updateStatus(id, status),
    onMutate: async ({ id, status }) => {
      await queryClient.cancelQueries({ queryKey: ATS_PIPELINE_QUERY_KEY })
      const previous = queryClient.getQueriesData({ queryKey: ATS_PIPELINE_QUERY_KEY })
      queryClient.setQueriesData({ queryKey: ATS_PIPELINE_QUERY_KEY }, (old) => {
        if (!old) return old
        const list = Array.isArray(old) ? old : old?.list ?? []
        const next = list.map((c) => (c.id === id ? { ...c, status } : c))
        return Array.isArray(old) ? next : { ...old, list: next }
      })
      return { previous }
    },
    onError: (err, variables, context) => {
      context?.previous?.forEach(([queryKey, data]) => {
        queryClient.setQueryData(queryKey, data)
      })
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ATS_PIPELINE_QUERY_KEY })
    },
  })
}
