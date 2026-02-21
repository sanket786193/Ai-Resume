/**
 * API service. All backend calls go through here. Aligned with Go API routes.
 */
import { http } from './http'

// ---------- Auth ----------
export const authApi = {
  login: (payload) => http.post('/auth/login', payload),
  registerCandidate: (payload) => http.post('/auth/register/candidate', payload),
  registerHR: (payload) => http.post('/auth/register/hr', payload),
  me: () => http.get('/auth/me'),
  refresh: (payload) => http.post('/auth/refresh', payload),
  logout: (payload) => http.post('/auth/logout', payload ?? {}),
}

// ---------- Jobs (public: list/detail; HR: list/create/update/delete/close/applications) ----------
export const jobsApi = {
  list: (params) => {
    const q = params ? `?${new URLSearchParams(params)}` : ''
    return http.get(`/jobs${q}`)
  },
  listForHR: (params) => {
    const q = params ? `?${new URLSearchParams(params)}` : ''
    return http.get(`/api/hr/jobs${q}`)
  },
  getById: (id) => http.get(`/jobs/${id}`),
  create: (payload) => http.post('/api/hr/jobs', payload),
  update: (id, payload) => http.put(`/api/hr/jobs/${id}`, payload),
  publish: (id) => http.post(`/api/hr/jobs/${id}/publish`),
  close: (id) => http.post(`/api/hr/jobs/${id}/close`),
  delete: (id) => http.delete(`/api/hr/jobs/${id}`),
  getApplicants: (jobId) => http.get(`/api/hr/jobs/${jobId}/applications`),
}

// ---------- Candidates / ATS (HR: applications list and status) ----------
export const candidatesApi = {
  list: (params) => {
    const q = params ? `?${new URLSearchParams(params)}` : ''
    return http.get(`/api/hr/applications${q}`)
  },
  getById: (id) => http.get(`/api/hr/applications/${id}`),
  getByJob: (jobId) => http.get(`/api/hr/jobs/${jobId}/applications`),
  updateStatus: (id, status) => http.put(`/api/hr/applications/${id}/status`, { status }),
}

// ---------- Applications (candidate-facing; requires candidate_id from auth) ----------
export const applicationsApi = {
  myApplications: (candidateId) => http.get(`/api/candidates/${candidateId}/applications`),
  apply: (candidateId, payload) => http.post(`/api/candidates/${candidateId}/applications`, { job_id: payload.jobId, resume_id: payload.resumeId }),
  getApplicationStatus: (candidateId, jobId) => http.get(`/api/candidates/${candidateId}/applications/${jobId}/status`),
}

// ---------- Resumes (candidate-facing) ----------
export const resumesApi = {
  list: (candidateId) => http.get(`/api/candidates/${candidateId}/resumes`),
  upload: (candidateId, file) => {
    const fd = new FormData()
    fd.append('file', file)
    return http.postForm(`/api/candidates/${candidateId}/resumes/upload`, fd)
  },
}

// ---------- Interviews (HR) ----------
export const interviewsApi = {
  list: (params) => {
    const q = params ? `?${new URLSearchParams(params)}` : ''
    return http.get(`/api/hr/interviews${q}`)
  },
  getById: (id) => http.get(`/api/hr/interviews/${id}`),
  create: (payload) => http.post('/api/hr/interviews', payload),
  update: (id, payload) => http.put(`/api/hr/interviews/${id}`, payload),
}

// ---------- Offers (HR create; candidate accept/reject) ----------
export const offersApi = {
  list: (params) => {
    const q = params ? `?${new URLSearchParams(params)}` : ''
    return http.get(`/api/hr/offers${q}`)
  },
  getById: (id) => http.get(`/api/hr/offers/${id}`),
  create: (payload) => http.post('/api/hr/offers', payload),
  accept: (candidateId, offerId) => http.post(`/api/candidates/${candidateId}/offers/${offerId}/accept`),
  reject: (candidateId, offerId) => http.post(`/api/candidates/${candidateId}/offers/${offerId}/reject`),
}
