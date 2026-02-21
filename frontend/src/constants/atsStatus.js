/**
 * ATS pipeline status enum. Single source of truth for candidate stages.
 */
export const ATS_STATUS = {
  APPLIED: 'APPLIED',
  SCREENING: 'SCREENING',
  SHORTLISTED: 'SHORTLISTED',
  INTERVIEW: 'INTERVIEW',
  REJECTED: 'REJECTED',
  HIRED: 'HIRED',
}

/** Ordered list for pipeline columns (excludes terminal states for main board). */
export const ATS_PIPELINE_ORDER = [
  ATS_STATUS.APPLIED,
  ATS_STATUS.SCREENING,
  ATS_STATUS.SHORTLISTED,
  ATS_STATUS.INTERVIEW,
  ATS_STATUS.HIRED,
  ATS_STATUS.REJECTED,
]

/** Statuses that appear as Kanban columns (active pipeline). */
export const ATS_KANBAN_STATUSES = [
  ATS_STATUS.APPLIED,
  ATS_STATUS.SCREENING,
  ATS_STATUS.SHORTLISTED,
  ATS_STATUS.INTERVIEW,
  ATS_STATUS.HIRED,
  ATS_STATUS.REJECTED,
]

export const ATS_STATUS_LABEL = {
  [ATS_STATUS.APPLIED]: 'Applied',
  [ATS_STATUS.SCREENING]: 'Screening',
  [ATS_STATUS.SHORTLISTED]: 'Shortlisted',
  [ATS_STATUS.INTERVIEW]: 'Interview',
  [ATS_STATUS.REJECTED]: 'Rejected',
  [ATS_STATUS.HIRED]: 'Hired',
}

/** Badge variant per status for candidate/HR UI (default | secondary | outline | destructive) */
export const ATS_STATUS_BADGE_VARIANT = {
  [ATS_STATUS.APPLIED]: 'secondary',
  [ATS_STATUS.SCREENING]: 'default',
  [ATS_STATUS.SHORTLISTED]: 'default',
  [ATS_STATUS.INTERVIEW]: 'default',
  [ATS_STATUS.REJECTED]: 'destructive',
  [ATS_STATUS.HIRED]: 'default',
}
