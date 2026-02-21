import { useParams, Link } from 'react-router-dom'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { useJob } from '@/modules/jobs/hooks/useJobs'
import { ROUTES, ROUTE_BUILDERS } from '@/constants'

export function PublicJobDetailPage() {
  const { jobId } = useParams()
  const { data: job, isPending, isError } = useJob(jobId)

  if (isPending || !job) {
    return <Skeleton className="h-64 w-full rounded-lg" />
  }

  if (isError) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">Job not found.</p>
        <Button variant="link" asChild><Link to={ROUTES.PUBLIC_JOBS}>Back to jobs</Link></Button>
      </div>
    )
  }

  const isClosed = job.status === 'CLOSED'

  return (
    <div>
      <Card>
        <CardHeader className="flex flex-row items-start justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold">{job.title}</h1>
            {job.department && <p className="text-muted-foreground">{job.department}</p>}
          </div>
          {isClosed && <Badge variant="secondary">Closed</Badge>}
        </CardHeader>
        <CardContent className="prose prose-sm max-w-none">
          <p className="whitespace-pre-wrap">{job.description}</p>
        </CardContent>
        {!isClosed && (
          <CardContent className="pt-0">
            <Button asChild>
              <Link to={ROUTE_BUILDERS.publicApply(job.id)}>Apply for this job</Link>
            </Button>
          </CardContent>
        )}
      </Card>
    </div>
  )
}
