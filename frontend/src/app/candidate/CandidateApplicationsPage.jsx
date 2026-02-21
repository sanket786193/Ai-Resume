import { Link, useNavigate } from 'react-router-dom'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { EmptyState } from '@/components/common/EmptyState'
import { PageHeader } from '@/components/common/PageHeader'
import { useMyApplications } from '@/modules/jobs/hooks/useApplications'
import { ATS_STATUS_LABEL, ROUTES, ROUTE_BUILDERS } from '@/constants'

export function CandidateApplicationsPage() {
  const navigate = useNavigate()
  const { data, isPending, isError } = useMyApplications()
  const list = Array.isArray(data) ? data : data?.list ?? []
  const goToJobs = () => navigate(ROUTES.PUBLIC_JOBS)

  return (
    <div>
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
        <div className="space-y-2">
          {list.map((app) => (
            <Card key={app.id}>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <Link to={ROUTE_BUILDERS.publicJobDetail(app.jobId)} className="font-medium hover:underline">
                    {app.jobTitle ?? `Job #${app.jobId}`}
                  </Link>
                  <p className="text-sm text-muted-foreground">{app.appliedAt && new Date(app.appliedAt).toLocaleDateString()}</p>
                </div>
                <Badge variant="secondary">{ATS_STATUS_LABEL[app.status] ?? app.status}</Badge>
              </CardHeader>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
