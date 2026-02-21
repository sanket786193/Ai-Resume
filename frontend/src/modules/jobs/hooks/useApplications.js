import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { applicationsApi } from '@/services/api'

export const MY_APPLICATIONS_QUERY_KEY = ['candidate', 'applications']

export function useMyApplications() {
  return useQuery({
    queryKey: MY_APPLICATIONS_QUERY_KEY,
    queryFn: applicationsApi.myApplications,
  })
}

export function useApplyToJob(jobId) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload) => applicationsApi.apply(jobId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MY_APPLICATIONS_QUERY_KEY })
    },
  })
}
