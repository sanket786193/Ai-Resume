import { Link } from 'react-router-dom'
import { PageHeader } from '@/components/common/PageHeader'
import { EmptyState } from '@/components/common/EmptyState'
import { useInterviewsList } from '@/modules/interviews/hooks/useInterviews'
import { ROUTE_BUILDERS } from '@/constants'

export function HrInterviewsPage() {
  const { data, isPending, isError } = useInterviewsList()
  const list = Array.isArray(data) ? data : data?.list ?? []

  return (
    <div className="space-y-6">
      <PageHeader title="Interviews" description="Schedule and manage interviews" />
      {isPending && <div className="text-muted-foreground">Loading...</div>}
      {!isPending && (isError || list.length === 0) && (
        <EmptyState
          title={isError ? 'Failed to load' : 'No interviews scheduled'}
          description={isError ? 'Please try again later.' : 'Schedule an interview from the pipeline.'}
        />
      )}
      {!isPending && !isError && list.length > 0 && (
        <ul className="space-y-3">
          {list.map((i) => {
            const id = i.id ?? i.ID
            const candidateName = i.candidate_name ?? i.candidateName ?? '—'
            const jobTitle = i.job_title ?? i.jobTitle ?? '—'
            const scheduledAt = i.scheduled_at ?? i.scheduledAt ?? i.ScheduledAt
            const duration = i.duration_minutes ?? i.durationMinutes ?? i.DurationMin ?? i.Duration
            const location = i.location ?? i.Location
            const round = i.round ?? i.Round
            const status = i.status ?? i.Status
            const atsId = i.ats_id ?? i.atsId ?? i.ATSID
            const jobId = i.job_id ?? i.jobId ?? i.JobID
            const whenStr = scheduledAt
              ? new Date(scheduledAt).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })
              : '—'
            return (
              <li key={id} className="border rounded-lg p-4 bg-card">
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div className="min-w-0 flex-1 space-y-1">
                    <p className="font-semibold text-foreground">{candidateName}</p>
                    <p className="text-sm text-muted-foreground">{jobTitle}</p>
                    <dl className="text-sm text-muted-foreground mt-2 space-y-0.5">
                      <div>
                        <span className="font-medium text-foreground">When: </span>
                        {whenStr}
                      </div>
                      {duration != null && duration > 0 && (
                        <div>
                          <span className="font-medium text-foreground">Duration: </span>
                          {duration} min
                        </div>
                      )}
                      {location && (
                        <div>
                          <span className="font-medium text-foreground">Where: </span>
                          {location}
                        </div>
                      )}
                      {round != null && round > 0 && (
                        <div>
                          <span className="font-medium text-foreground">Round: </span>
                          {round}
                        </div>
                      )}
                    </dl>
                  </div>
                  <div className="flex items-center gap-2 shrink-0">
                    <span className="text-xs font-medium uppercase tracking-wide text-muted-foreground capitalize">
                      {status ?? '—'}
                    </span>
                    {jobId && (
                      <Link
                        to={ROUTE_BUILDERS.hrJobApplicants(jobId)}
                        className="text-xs text-primary hover:underline"
                      >
                        View in pipeline
                      </Link>
                    )}
                  </div>
                </div>
              </li>
            )
          })}
        </ul>
      )}
    </div>
  )
}
