import { useParams } from 'react-router-dom'
import { PageHeader } from '@/components/common/PageHeader'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { EmptyState } from '@/components/common/EmptyState'
import { useJob, useJobApplicants } from '@/modules/jobs/hooks/useJobs'
import { ATS_STATUS_LABEL, ROUTES } from '@/constants'

export function HrJobApplicantsPage() {
  const { jobId } = useParams()
  const validJobId = jobId != null && String(jobId) !== '' && String(jobId) !== 'undefined'
  const { data: job, isPending: jobLoading } = useJob(jobId)
  const { data: applicants, isPending: applicantsLoading } = useJobApplicants(jobId)

  const list = Array.isArray(applicants) ? applicants : applicants?.list ?? []

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
        <div className="space-y-2">
          {list.map((a) => {
            const appId = a.id ?? a.ID
            const name = a.name ?? a.Name ?? a.email ?? a.Email ?? `Candidate ${a.candidate_id ?? a.CandidateID ?? appId}`.slice(0, 40)
            const email = a.email ?? a.Email ?? ''
            const status = a.status ?? a.Status
            return (
              <Card key={appId}>
                <CardHeader className="flex flex-row items-center justify-between">
                  <div>
                    <p className="font-medium">{name}</p>
                    {email && <p className="text-sm text-muted-foreground">{email}</p>}
                  </div>
                  <Badge variant="secondary">{ATS_STATUS_LABEL[status] ?? status}</Badge>
                </CardHeader>
                {(a.currentRole ?? a.CurrentRole) && (
                  <CardContent className="pt-0">
                    <p className="text-sm text-muted-foreground">{a.currentRole ?? a.CurrentRole}</p>
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
