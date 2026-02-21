import { Outlet, useParams, NavLink, useNavigate } from 'react-router-dom'
import { PageHeader } from '@/components/common/PageHeader'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { useApplicationDetail } from '@/modules/ats/hooks/useApplicationDetail'
import { ROUTE_BUILDERS, ATS_STATUS_LABEL, ATS_STATUS_BADGE_VARIANT } from '@/constants'
import { IconArrowLeft } from '@/components/icons/ApplicantIcons'
import { cn } from '@/lib/utils'

const TAB_BASE = (jobId, applicationId) => `/hr/jobs/${jobId}/applicants/${applicationId}`

function ApplicantSidebar({ detail }) {
  if (!detail) return null
  const status = detail.status
  const statusLabel = status ? (ATS_STATUS_LABEL[status] ?? status) : null
  const variant = status ? (ATS_STATUS_BADGE_VARIANT[status] ?? 'secondary') : 'secondary'
  const appliedAt = detail.created_at
    ? new Date(detail.created_at).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })
    : null
  const hasScoresOrMeta = statusLabel || detail.job_title || appliedAt || detail.ats_score != null || detail.skill_match_pct != null || (detail.ai_summary && detail.ai_summary.trim() !== '')
  if (!hasScoresOrMeta) return null
  const aiSummaryShort = detail.ai_summary?.trim() ? (detail.ai_summary.trim().length > 120 ? `${detail.ai_summary.trim().slice(0, 120)}…` : detail.ai_summary.trim()) : null
  return (
    <Card className="sticky top-4 shadow-sm border-border/80 h-fit">
      <CardContent className="p-4 space-y-4">
        {statusLabel && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Status</span>
            <div className="mt-1.5">
              <Badge variant={variant} className="font-medium">
                {statusLabel}
              </Badge>
            </div>
          </div>
        )}
        {detail.job_title && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Job</span>
            <p className="text-sm font-medium mt-1 line-clamp-2">{detail.job_title}</p>
          </div>
        )}
        {appliedAt && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Applied</span>
            <p className="text-sm font-medium mt-1">{appliedAt}</p>
          </div>
        )}
        {detail.ats_score != null && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">ATS score</span>
            <p className="text-lg font-semibold mt-1 text-primary">{detail.ats_score}/100</p>
          </div>
        )}
        {detail.skill_match_pct != null && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Skill match</span>
            <p className="text-lg font-semibold mt-1 text-primary">{detail.skill_match_pct}%</p>
          </div>
        )}
        {aiSummaryShort && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">AI summary</span>
            <p className="text-sm mt-1 text-muted-foreground line-clamp-4">{aiSummaryShort}</p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

export function ApplicantDetailLayout() {
  const { jobId, applicationId } = useParams()
  const navigate = useNavigate()
  const { data: detail, isPending, isError } = useApplicationDetail(applicationId, { enabled: Boolean(applicationId) })

  const name = detail?.candidate_name ?? detail?.candidate_email ?? 'Applicant'
  const hasContact = name || detail?.candidate_email || detail?.candidate_phone || detail?.candidate_linkedin
  const hasApplication = detail?.status != null || detail?.created_at != null || detail?.job_title
  const hasResume = detail?.resume_file_name
  const hasAI =
    detail?.ats_score != null ||
    detail?.skill_match_score != null ||
    detail?.skill_match_pct != null ||
    detail?.qualified != null ||
    (detail?.ai_summary && detail.ai_summary !== '') ||
    (detail?.experience_match && detail.experience_match !== '') ||
    (Array.isArray(detail?.missing_skills) && detail.missing_skills.length > 0) ||
    (Array.isArray(detail?.experience_warnings) && detail.experience_warnings.length > 0) ||
    (Array.isArray(detail?.keyword_matches) && detail.keyword_matches.length > 0) ||
    (Array.isArray(detail?.semantic_matches) && detail.semantic_matches.length > 0)

  const backHref = ROUTE_BUILDERS.hrJobApplicants(jobId)

  if (isError) {
    return (
      <div className="max-w-4xl">
        <p className="text-sm text-destructive">Failed to load applicant details. Please try again.</p>
        <Button variant="outline" size="sm" className="mt-3" onClick={() => navigate(backHref)}>
          Back to applicants
        </Button>
      </div>
    )
  }

  const tabClass = ({ isActive }) =>
    cn(
      'px-4 py-2.5 rounded-md text-sm font-medium transition-colors',
      isActive
        ? 'bg-primary text-primary-foreground shadow-sm'
        : 'text-muted-foreground hover:text-foreground hover:bg-muted/60'
    )

  return (
    <div className="max-w-4xl">
      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mb-6">
        <PageHeader
          title={isPending ? <Skeleton className="h-8 w-48" /> : name}
          description={isPending ? null : 'Applicant details'}
          actions={
            <Button
              variant="outline"
              size="sm"
              className="shrink-0 gap-1.5 text-muted-foreground hover:text-foreground"
              onClick={() => navigate(backHref)}
            >
              <IconArrowLeft className="w-4 h-4" />
              Back to applicants
            </Button>
          }
        />
      </div>
      {!isPending && detail && (
        <nav
          className="flex flex-wrap gap-1 p-1 rounded-lg bg-muted/40 border border-border/50 mb-6 w-fit"
          aria-label="Applicant sections"
        >
          {hasContact && (
            <NavLink to={`${TAB_BASE(jobId, applicationId)}/contact`} end={false} className={tabClass}>
              Contact
            </NavLink>
          )}
          {hasApplication && (
            <NavLink to={`${TAB_BASE(jobId, applicationId)}/application`} className={tabClass}>
              Application
            </NavLink>
          )}
          {hasResume && (
            <NavLink to={`${TAB_BASE(jobId, applicationId)}/resume`} className={tabClass}>
              Resume
            </NavLink>
          )}
          {(hasResume || hasAI) && (
            <NavLink to={`${TAB_BASE(jobId, applicationId)}/ai`} className={tabClass}>
              AI feedback
            </NavLink>
          )}
        </nav>
      )}
      <div className="grid grid-cols-1 lg:grid-cols-[1fr_260px] gap-6 lg:gap-8">
        <div className="min-w-0">
          {isPending && <Skeleton className="h-48 w-full rounded-xl" />}
          {!isPending && detail && <Outlet context={{ detail, jobId }} />}
        </div>
        {!isPending && detail && (
          <aside className="hidden lg:block">
            <ApplicantSidebar detail={detail} />
          </aside>
        )}
      </div>
    </div>
  )
}
