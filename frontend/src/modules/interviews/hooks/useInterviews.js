import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { interviewsApi } from '@/services/api'

export const INTERVIEWS_QUERY_KEY = ['interviews']

export function useInterviewsList(params) {
  return useQuery({
    queryKey: [...INTERVIEWS_QUERY_KEY, params],
    queryFn: () => interviewsApi.list(params),
  })
}

export function useCreateInterview() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: interviewsApi.create,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: INTERVIEWS_QUERY_KEY }),
  })
}
