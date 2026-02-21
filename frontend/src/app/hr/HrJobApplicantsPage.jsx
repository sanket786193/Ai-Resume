import { useParams, useNavigate } from 'react-router-dom'
import { PageHeader } from '@/components/common/PageHeader'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { EmptyState } from '@/components/common/EmptyState'
import { useJob, useJobApplicants } from '@/modules/jobs/hooks/useJobs'
import { ATS_STATUS_LABEL, ATS_STATUS_BADGE_VARIANT, ROUTES, ROUTE_BUILDERS } from '@/constants'
import { cn } from '@/lib/utils'

function getDisplayName(a) {
  return a?.candidate_name ?? a?.name ?? a?.candidate_email ?? a?.email ?? 'Applicant'
}

function getEmail(a) {
  return a?.candidate_email ?? a?.email ?? ''
}

function getInitial(a) {
  const name = getDisplayName(a)
  if (name && name !== 'Applicant') return name.charAt(0).toUpperCase()
  const email = getEmail(a)
  return email ? email.charAt(0).toUpperCase() : '?'
}

export function HrJobApplicantsPage() {
  const { jobId } = useParams()
  const navigate = useNavigate()

  const validJobId = jobId != null && String(jobId) !== '' && String(jobId) !== 'undefined'
  const { data: job, isPending: jobLoading } = useJob(jobId)
  const { data: applicants, isPending: applicantsLoading } = useJobApplicants(jobId)

  const list = Array.isArray(applicants) ? applicants : applicants?.list ?? []

  const openDetail = (appId) => {
    navigate(ROUTE_BUILDERS.hrApplicantDetail(jobId, appId, 'contact'))
  }

  if (!validJobId) {
    return (
      <EmptyState
        title="Invalid job"
        description="This job link is invalid or the job was removed."
        actionLabel="Back to jobs"
        onAction={() => window.location.assign(ROUTES.HR_JOBS)}
      />
    )
  }

  if (jobLoading || !job) {
    return <Skeleton className="h-64 w-full rounded-lg" />
  }

  const jobTitle = job?.title ?? job?.Title ?? 'Job'
  return (
    <div>
      <PageHeader
        title={jobTitle}
        description={`Applicants for ${jobTitle}`}
      />
      {applicantsLoading && <Skeleton className="h-48 w-full rounded-lg" />}
      {!applicantsLoading && list.length === 0 && (
        <EmptyState title="No applicants yet" description="Applicants will appear here when candidates apply." />
      )}
      {!applicantsLoading && list.length > 0 && (
        <div className="space-y-3">
          {list.map((a) => {
            const appId = a.id ?? a.ID
            const status = a.status ?? a.Status
            const variant = ATS_STATUS_BADGE_VARIANT[status] ?? 'secondary'
            const appliedAt = a.created_at ? new Date(a.created_at).toLocaleDateString(undefined, { dateStyle: 'medium' }) : null
            return (
              <Card
                key={appId}
                className={cn(
                  'transition-colors cursor-pointer hover:bg-muted/50 hover:border-primary/20',
                  'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
                )}
                tabIndex={0}
                onClick={() => openDetail(appId)}
                onKeyDown={(e) => (e.key === 'Enter' || e.key === ' ') && openDetail(appId)}
              >
                <CardContent className="p-4 flex items-center gap-4">
                  <div
                    className="flex-shrink-0 w-11 h-11 rounded-full bg-primary/10 text-primary flex items-center justify-center text-base font-semibold"
                    aria-hidden
                  >
                    {getInitial(a)}
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="font-medium truncate">{getDisplayName(a)}</p>
                    {getEmail(a) && (
                      <p className="text-sm text-muted-foreground truncate">{getEmail(a)}</p>
                    )}
                    {appliedAt && (
                      <p className="text-xs text-muted-foreground mt-0.5">Applied {appliedAt}</p>
                    )}
                  </div>
                  <div className="flex items-center gap-2 shrink-0">
                    {a.ats_score != null && (
                      <span className="text-xs text-muted-foreground" title="ATS score">
                        {a.ats_score}/100
                      </span>
                    )}
                    <Badge variant={variant}>{ATS_STATUS_LABEL[status] ?? status}</Badge>
                  </div>
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}
