import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { jobsApi } from '@/services/api'

export const JOBS_QUERY_KEY = ['jobs']
export const JOB_QUERY_KEY = (id) => ['jobs', id]
export const JOB_APPLICANTS_QUERY_KEY = (jobId) => ['jobs', jobId, 'applicants']

export function useJobsList(params) {
  return useQuery({
    queryKey: [...JOBS_QUERY_KEY, params],
    queryFn: () => jobsApi.list(params),
  })
}

export function useJob(id, options = {}) {
  return useQuery({
    queryKey: JOB_QUERY_KEY(id),
    queryFn: () => jobsApi.getById(id),
    enabled: !!id && (options.enabled !== false),
  })
}

export function useJobApplicants(jobId) {
  return useQuery({
    queryKey: JOB_APPLICANTS_QUERY_KEY(jobId),
    queryFn: () => jobsApi.getApplicants(jobId),
    enabled: !!jobId,
  })
}

export function useCreateJob() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: jobsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: JOBS_QUERY_KEY })
    },
  })
}

export function useUpdateJob(id) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload) => jobsApi.update(id, payload),
    onSuccess: (_, __, context) => {
      queryClient.invalidateQueries({ queryKey: JOBS_QUERY_KEY })
      if (id) queryClient.invalidateQueries({ queryKey: JOB_QUERY_KEY(id) })
    },
  })
}

export function useCloseJob(id) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => jobsApi.close(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: JOBS_QUERY_KEY })
      if (id) queryClient.invalidateQueries({ queryKey: JOB_QUERY_KEY(id) })
    },
  })
}
