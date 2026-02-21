/**
 * API service. All backend calls live here. No React, no hooks.
 */
import { http } from './http'

// ---------- Auth ----------
export const authApi = {
  login: (payload) => http.post('/auth/login', payload),
  register: (payload) => http.post('/auth/register', payload),
  me: () => http.get('/auth/me'),
  logout: () => http.post('/auth/logout'),
}

// ---------- Jobs ----------
export const jobsApi = {
  list: (params) => {
    const q = params ? `?${new URLSearchParams(params)}` : ''
    return http.get(`/jobs${q}`)
  },
  getById: (id) => http.get(`/jobs/${id}`),
  create: (payload) => http.post('/jobs', payload),
  update: (id, payload) => http.patch(`/jobs/${id}`, payload),
  close: (id) => http.patch(`/jobs/${id}/close`),
  getApplicants: (jobId) => http.get(`/jobs/${jobId}/applicants`),
}

// ---------- Candidates (ATS) ----------
export const candidatesApi = {
  list: (params) => {
    const q = params ? `?${new URLSearchParams(params)}` : ''
    return http.get(`/candidates${q}`)
  },
  getById: (id) => http.get(`/candidates/${id}`),
  updateStatus: (id, status) => http.patch(`/candidates/${id}/status`, { status }),
  getByJob: (jobId) => http.get(`/jobs/${jobId}/applicants`),
}

// ---------- Applications (candidate-facing) ----------
export const applicationsApi = {
  apply: (jobId, payload) => {
    if (payload?.resume instanceof File) {
      const form = new FormData()
      form.append('name', payload.name ?? '')
      form.append('email', payload.email ?? '')
      form.append('resume', payload.resume)
      return http.postForm(`/jobs/${jobId}/apply`, form)
    }
    return http.post(`/jobs/${jobId}/apply`, { name: payload.name, email: payload.email })
  },
  myApplications: () => http.get('/candidate/applications'),
}

// ---------- Interviews ----------
export const interviewsApi = {
  list: (params) => {
    const q = params ? `?${new URLSearchParams(params)}` : ''
    return http.get(`/interviews${q}`)
  },
  create: (payload) => http.post('/interviews', payload),
  update: (id, payload) => http.patch(`/interviews/${id}`, payload),
}

// ---------- Offers ----------
export const offersApi = {
  list: (params) => {
    const q = params ? `?${new URLSearchParams(params)}` : ''
    return http.get(`/offers${q}`)
  },
  create: (payload) => http.post('/offers', payload),
  sendSelection: (candidateId) => http.post(`/candidates/${candidateId}/select`),
  sendRejection: (candidateId) => http.post(`/candidates/${candidateId}/reject`),
}
