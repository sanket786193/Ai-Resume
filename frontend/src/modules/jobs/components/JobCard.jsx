import { Link } from 'react-router-dom'
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ROUTE_BUILDERS } from '@/constants'

/**
 * Single job card for list views. Stateless.
 * viewHref: optional (e.g. HR applicants page); default is public job detail.
 */
export function JobCard({ job, viewHref }) {
  const isClosed = job?.status === 'CLOSED'
  const viewTo = viewHref ?? ROUTE_BUILDERS.publicJobDetail(job.id)
  return (
    <Card>
      <CardHeader className="flex flex-row items-start justify-between gap-2">
        <div>
          <CardContent className="p-0">
            <Link to={viewTo} className="font-semibold hover:underline">
              {job.title}
            </Link>
          </CardContent>
          {job.department && <p className="text-sm text-muted-foreground mt-1">{job.department}</p>}
        </div>
        {isClosed && <Badge variant="secondary">Closed</Badge>}
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground line-clamp-2">{job.description}</p>
      </CardContent>
      <CardFooter className="gap-2">
        <Button variant="outline" size="sm" asChild>
          <Link to={viewTo}>View</Link>
        </Button>
      </CardFooter>
    </Card>
  )
}
