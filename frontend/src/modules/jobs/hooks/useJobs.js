import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { jobsApi } from '@/services/api'
import { ATS_PIPELINE_QUERY_KEY } from '@/modules/ats/hooks/useAtsPipeline'

export const JOBS_QUERY_KEY = ['jobs']
export const HR_JOBS_QUERY_KEY = ['jobs', 'hr']
export const JOB_QUERY_KEY = (id) => ['jobs', id]
export const JOB_APPLICANTS_QUERY_KEY = (jobId) => ['jobs', jobId, 'applicants']

export function useJobsList(params) {
  return useQuery({
    queryKey: [...JOBS_QUERY_KEY, params],
    queryFn: () => jobsApi.list(params),
  })
}

export function useHrJobsList(params) {
  return useQuery({
    queryKey: [...HR_JOBS_QUERY_KEY, params],
    queryFn: () => jobsApi.listForHR(params),
  })
}

export function useJob(id, options = {}) {
  const validId = id != null && String(id) !== '' && String(id) !== 'undefined'
  return useQuery({
    queryKey: JOB_QUERY_KEY(id),
    queryFn: () => jobsApi.getById(id),
    enabled: validId && (options.enabled !== false),
  })
}

export function useJobApplicants(jobId) {
  const validId = jobId != null && String(jobId) !== '' && String(jobId) !== 'undefined'
  return useQuery({
    queryKey: JOB_APPLICANTS_QUERY_KEY(jobId),
    queryFn: () => jobsApi.getApplicants(jobId),
    enabled: validId,
  })
}

export function useCreateJob() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: jobsApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: JOBS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: HR_JOBS_QUERY_KEY })
    },
  })
}

export function useUpdateJob(id) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload) => jobsApi.update(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: JOBS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: HR_JOBS_QUERY_KEY })
      if (id) queryClient.invalidateQueries({ queryKey: JOB_QUERY_KEY(id) })
    },
  })
}

export function usePublishJob() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id) => jobsApi.publish(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: JOBS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: HR_JOBS_QUERY_KEY })
      if (id) queryClient.invalidateQueries({ queryKey: JOB_QUERY_KEY(id) })
    },
  })
}

export function useCloseJob() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id) => jobsApi.close(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: JOBS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: HR_JOBS_QUERY_KEY })
      if (id) queryClient.invalidateQueries({ queryKey: JOB_QUERY_KEY(id) })
    },
  })
}

export function useDeleteJob() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id) => jobsApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: JOBS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: HR_JOBS_QUERY_KEY })
    },
  })
}

/** Phase 4: bulk add applicants to a job (HR). */
export function useBulkApply(jobId) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload) => jobsApi.bulkApply(jobId, payload),
    onSuccess: (_, __, context) => {
      queryClient.invalidateQueries({ queryKey: JOB_APPLICANTS_QUERY_KEY(jobId) })
      queryClient.invalidateQueries({ queryKey: ATS_PIPELINE_QUERY_KEY })
    },
  })
}
