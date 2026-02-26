import { Link, useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { EmptyState } from '@/components/common/EmptyState'
import { PageHeader } from '@/components/common/PageHeader'
import { useAuth } from '@/modules/auth/hooks/useAuth'
import { useMyApplications } from '@/modules/jobs/hooks/useApplications'
import { useAcceptOffer, useRejectOffer } from '@/modules/offers/hooks/useOffers'
import { ATS_STATUS, ATS_STATUS_LABEL, ATS_STATUS_BADGE_VARIANT, ROUTES, ROUTE_BUILDERS } from '@/constants'

function OfferActions({ candidateId, offerId, onSuccess }) {
  const acceptOffer = useAcceptOffer()
  const rejectOffer = useRejectOffer()
  const handleAccept = () => {
    if (!candidateId || !offerId) return
    acceptOffer.mutate(
      { candidateId, offerId },
      {
        onSuccess: () => {
          toast.success('Offer accepted')
          onSuccess?.()
        },
        onError: () => toast.error('Failed to accept offer'),
      }
    )
  }
  const handleReject = () => {
    if (!candidateId || !offerId) return
    rejectOffer.mutate(
      { candidateId, offerId },
      {
        onSuccess: () => {
          toast.success('Offer declined')
          onSuccess?.()
        },
        onError: () => toast.error('Failed to decline offer'),
      }
    )
  }
  return (
    <>
      <Button size="sm" onClick={handleAccept} disabled={acceptOffer.isPending || rejectOffer.isPending}>
        {acceptOffer.isPending ? 'Accepting…' : 'Accept offer'}
      </Button>
      <Button size="sm" variant="outline" onClick={handleReject} disabled={acceptOffer.isPending || rejectOffer.isPending}>
        {rejectOffer.isPending ? 'Declining…' : 'Decline'}
      </Button>
    </>
  )
}

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
                : `Job …${shortId(jobId ?? app.id ?? app.ID)}`
            const appliedAt = app.created_at ?? app.appliedAt ?? app.CreatedAt
            const status = app.status ?? app.Status
            const statusLabel = ATS_STATUS_LABEL[status] ?? status ?? '—'
            const badgeVariant = ATS_STATUS_BADGE_VARIANT[status] ?? 'secondary'
            const interview = app.interview
            const isInterviewScheduled = status === ATS_STATUS.INTERVIEW && interview
            return (
              <Card key={app.id ?? app.ID ?? jobId} className="overflow-hidden transition-shadow hover:shadow-md">
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
                {isInterviewScheduled && (
                  <CardContent className="pt-0 border-t border-border/50">
                    <div className="rounded-lg bg-muted/40 p-4 space-y-2">
                      <p className="text-sm font-medium">Interview details</p>
                      <dl className="text-sm text-muted-foreground space-y-1">
                        <div className="flex flex-wrap gap-x-4 gap-y-1">
                          <span>
                            <span className="font-medium text-foreground">When: </span>
                            {new Date(interview.scheduled_at ?? interview.ScheduledAt).toLocaleString(undefined, {
                              dateStyle: 'medium',
                              timeStyle: 'short',
                            })}
                          </span>
                          {(interview.duration_minutes ?? interview.DurationMin ?? interview.duration_minutes) != null && (
                            <span>
                              <span className="font-medium text-foreground">Duration: </span>
                              {interview.duration_minutes ?? interview.DurationMin ?? interview.duration_minutes} min
                            </span>
                          )}
                        </div>
                        {(interview.location ?? interview.Location) && (
                          <div>
                            <span className="font-medium text-foreground">Where: </span>
                            {interview.location ?? interview.Location}
                          </div>
                        )}
                        {(interview.round ?? interview.Round) != null && (interview.round ?? interview.Round) > 0 && (
                          <div>
                            <span className="font-medium text-foreground">Round: </span>
                            {interview.round ?? interview.Round}
                          </div>
                        )}
                      </dl>
                    </div>
                  </CardContent>
                )}
                {app.offer && (
                  <CardContent className="pt-0 border-t border-border/50">
                    <div className="rounded-lg bg-muted/40 p-4 space-y-2">
                      <p className="text-sm font-medium">Offer (you were selected)</p>
                      <p className="text-sm text-muted-foreground">
                        {(app.offer.amount ?? app.offer.Amount) && (app.offer.currency ?? app.offer.Currency)
                          ? `${app.offer.amount ?? app.offer.Amount} ${app.offer.currency ?? app.offer.Currency}`
                          : 'Offer issued'}
                        {' — '}
                        <span className="capitalize">{app.offer.status ?? app.offer.Status ?? '—'}</span>
                      </p>
                      {(app.offer.status ?? app.offer.Status) === 'PENDING' && (
                        <div className="flex flex-wrap gap-2 pt-2">
                          <OfferActions
                            candidateId={candidateId}
                            offerId={app.offer.id ?? app.offer.ID}
                            onSuccess={() => {}}
                          />
                        </div>
                      )}
                    </div>
                  </CardContent>
                )}
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}
