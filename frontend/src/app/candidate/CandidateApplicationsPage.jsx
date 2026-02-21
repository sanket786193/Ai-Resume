import { Link, useNavigate } from 'react-router-dom'
import { Card, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { EmptyState } from '@/components/common/EmptyState'
import { PageHeader } from '@/components/common/PageHeader'
import { useAuth } from '@/modules/auth/hooks/useAuth'
import { useMyApplications } from '@/modules/jobs/hooks/useApplications'
import { ATS_STATUS_LABEL, ATS_STATUS_BADGE_VARIANT, ROUTES, ROUTE_BUILDERS } from '@/constants'

function shortId(id) {
  if (!id || typeof id !== 'string') return ''
  return id.slice(0, 8)
}

export function CandidateApplicationsPage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const candidateId = user?.candidate_id ?? user?.candidateId
  const { data, isPending, isError } = useMyApplications(candidateId)
  const list = Array.isArray(data) ? data : data?.list ?? []
  const goToJobs = () => navigate(ROUTES.PUBLIC_JOBS)

  return (
    <div className="space-y-6">
      <PageHeader title="My applications" description="Track your application status" />
      {isPending && <Skeleton className="h-48 w-full rounded-lg" />}
      {!isPending && isError && (
        <EmptyState title="Failed to load applications" description="Please try again later." />
      )}
      {!isPending && !isError && list.length === 0 && (
        <EmptyState
          title="No applications yet"
          description="Apply to a job to see your applications here."
          actionLabel="Browse jobs"
          onAction={goToJobs}
        />
      )}
      {!isPending && !isError && list.length > 0 && (
        <div className="space-y-3">
          {list.map((app) => {
            const jobId = app.job_id ?? app.jobId ?? app.JobID
            const jobTitle =
              app.job_title ?? app.jobTitle ?? app.JobTitle
                ? (app.job_title ?? app.jobTitle ?? app.JobTitle)
                : `Job …${shortId(jobId ?? app.id)}`
            const appliedAt = app.created_at ?? app.appliedAt ?? app.CreatedAt
            const statusLabel = ATS_STATUS_LABEL[app.status] ?? app.status
            const badgeVariant = ATS_STATUS_BADGE_VARIANT[app.status] ?? 'secondary'
            return (
              <Card key={app.id ?? jobId} className="overflow-hidden transition-shadow hover:shadow-md">
                <CardHeader className="flex flex-row items-start justify-between gap-4 pb-2">
                  <div className="min-w-0 flex-1 space-y-1">
                    {jobId ? (
                      <Link
                        to={ROUTE_BUILDERS.publicJobDetail(jobId)}
                        className="font-semibold text-base hover:underline focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 rounded"
                      >
                        {jobTitle}
                      </Link>
                    ) : (
                      <span className="font-semibold text-base">{jobTitle}</span>
                    )}
                    {appliedAt && (
                      <p className="text-sm text-muted-foreground">
                        Applied {new Date(appliedAt).toLocaleDateString(undefined, { dateStyle: 'medium' })}
                      </p>
                    )}
                  </div>
                  <Badge variant={badgeVariant} className="shrink-0 capitalize">
                    {statusLabel}
                  </Badge>
                </CardHeader>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}
