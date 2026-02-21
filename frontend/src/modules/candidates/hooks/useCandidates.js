import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { candidatesApi } from '@/services/api'
import { jobsApi } from '@/services/api'

export const CANDIDATES_QUERY_KEY = ['candidates']
export const CANDIDATE_QUERY_KEY = (id) => ['candidates', id]
export const JOB_APPLICANTS_QUERY_KEY = (jobId) => ['jobs', jobId, 'applicants']

export function useCandidatesList(params) {
  return useQuery({
    queryKey: [...CANDIDATES_QUERY_KEY, params],
    queryFn: () => candidatesApi.list(params),
  })
}

export function useCandidate(id) {
  return useQuery({
    queryKey: CANDIDATE_QUERY_KEY(id),
    queryFn: () => candidatesApi.getById(id),
    enabled: !!id,
  })
}

export function useJobApplicants(jobId) {
  return useQuery({
    queryKey: JOB_APPLICANTS_QUERY_KEY(jobId),
    queryFn: () => jobsApi.getApplicants(jobId),
    enabled: !!jobId,
  })
}

export function useUpdateCandidateStatus() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }) => candidatesApi.updateStatus(id, status),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: CANDIDATES_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: CANDIDATE_QUERY_KEY(id) })
      queryClient.invalidateQueries({ predicate: (q) => q.queryKey[0] === 'jobs' && q.queryKey[2] === 'applicants' })
    },
  })
}
