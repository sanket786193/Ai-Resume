/**
 * Single source of truth for route paths.
 */
export const ROUTES = {
  HOME: '/',
  // Auth
  LOGIN: '/login',
  REGISTER: '/register',
  REGISTER_HR: '/register/hr',
  // HR
  HR_DASHBOARD: '/hr',
  HR_JOBS: '/hr/jobs',
  HR_JOB_CREATE: '/hr/jobs/new',
  HR_JOB_EDIT: '/hr/jobs/:jobId/edit',
  HR_JOB_APPLICANTS: '/hr/jobs/:jobId/applicants',
  HR_APPLICANT_DETAIL: '/hr/jobs/:jobId/applicants/:applicationId',
  HR_APPLICANT_CONTACT: '/hr/jobs/:jobId/applicants/:applicationId/contact',
  HR_APPLICANT_APPLICATION: '/hr/jobs/:jobId/applicants/:applicationId/application',
  HR_APPLICANT_RESUME: '/hr/jobs/:jobId/applicants/:applicationId/resume',
  HR_APPLICANT_AI: '/hr/jobs/:jobId/applicants/:applicationId/ai',
  HR_ATS: '/hr/ats',
  HR_INTERVIEWS: '/hr/interviews',
  HR_OFFERS: '/hr/offers',
  // Candidate
  CANDIDATE_DASHBOARD: '/candidate',
  CANDIDATE_APPLICATIONS: '/candidate/applications',
  // Public
  PUBLIC_JOBS: '/jobs',
  PUBLIC_JOB_DETAIL: '/jobs/:jobId',
  PUBLIC_APPLY: '/jobs/:jobId/apply',
}

export const ROUTE_BUILDERS = {
  hrJobEdit: (jobId) => `/hr/jobs/${jobId}/edit`,
  hrJobApplicants: (jobId) => `/hr/jobs/${jobId}/applicants`,
  hrApplicantDetail: (jobId, applicationId, tab = 'contact') =>
    `/hr/jobs/${jobId}/applicants/${applicationId}/${tab}`,
  publicJobDetail: (jobId) => `/jobs/${jobId}`,
  publicApply: (jobId) => `/jobs/${jobId}/apply`,
}
