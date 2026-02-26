import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { applicationsApi } from '@/services/api'

export const MY_APPLICATIONS_QUERY_KEY = ['candidate', 'applications']

/** Phase 3: refetch every 30s for "real-time" status feel. */
const APPLICATIONS_REFETCH_MS = 30_000

export function useMyApplications(candidateId, options = {}) {
  return useQuery({
    queryKey: [...MY_APPLICATIONS_QUERY_KEY, candidateId],
    queryFn: () => applicationsApi.myApplications(candidateId),
    enabled: !!candidateId,
    refetchInterval: options.refetchInterval ?? APPLICATIONS_REFETCH_MS,
    ...options,
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
