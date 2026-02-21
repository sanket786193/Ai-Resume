import { useParams } from 'react-router-dom'
import { PageHeader } from '@/components/common/PageHeader'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { EmptyState } from '@/components/common/EmptyState'
import { useJob, useJobApplicants } from '@/modules/jobs/hooks/useJobs'
import { ATS_STATUS_LABEL } from '@/constants'

export function HrJobApplicantsPage() {
  const { jobId } = useParams()
  const { data: job, isPending: jobLoading } = useJob(jobId)
  const { data: applicants, isPending: applicantsLoading } = useJobApplicants(jobId)

  const list = Array.isArray(applicants) ? applicants : applicants?.list ?? []

  if (jobLoading || !job) {
    return <Skeleton className="h-64 w-full rounded-lg" />
  }

  return (
    <div>
      <PageHeader
        title={job.title}
        description={`Applicants for ${job.title}`}
      />
      {applicantsLoading && <Skeleton className="h-48 w-full rounded-lg" />}
      {!applicantsLoading && list.length === 0 && (
        <EmptyState title="No applicants yet" description="Applicants will appear here when candidates apply." />
      )}
      {!applicantsLoading && list.length > 0 && (
        <div className="space-y-2">
          {list.map((a) => (
            <Card key={a.id}>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <p className="font-medium">{a.name ?? a.email}</p>
                  <p className="text-sm text-muted-foreground">{a.email}</p>
                </div>
                <Badge variant="secondary">{ATS_STATUS_LABEL[a.status] ?? a.status}</Badge>
              </CardHeader>
              {a.currentRole && (
                <CardContent className="pt-0">
                  <p className="text-sm text-muted-foreground">{a.currentRole}</p>
                </CardContent>
              )}
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
