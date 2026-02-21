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
      const previous = queryClient.getQueryData(ATS_PIPELINE_QUERY_KEY)
      queryClient.setQueryData(ATS_PIPELINE_QUERY_KEY, (old) => {
        if (!old?.list) return old
        return {
          ...old,
          list: old.list.map((c) => (c.id === id ? { ...c, status } : c)),
        }
      })
      return { previous }
    },
    onError: (err, variables, context) => {
      if (context?.previous) {
        queryClient.setQueryData(ATS_PIPELINE_QUERY_KEY, context.previous)
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ATS_PIPELINE_QUERY_KEY })
    },
  })
}
