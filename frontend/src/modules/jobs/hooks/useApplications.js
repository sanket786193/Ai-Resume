import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { applicationsApi } from '@/services/api'

export const MY_APPLICATIONS_QUERY_KEY = ['candidate', 'applications']

export function useMyApplications(candidateId) {
  return useQuery({
    queryKey: [...MY_APPLICATIONS_QUERY_KEY, candidateId],
    queryFn: () => applicationsApi.myApplications(candidateId),
    enabled: !!candidateId,
  })
}

export function useApplyToJob(jobId, candidateId) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload) => applicationsApi.apply(candidateId, { jobId, resumeId: payload.resumeId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MY_APPLICATIONS_QUERY_KEY })
    },
  })
}
