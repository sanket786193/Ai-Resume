import { Link } from 'react-router-dom'
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ROUTE_BUILDERS } from '@/constants'

/**
 * Single job card for list views.
 * viewHref: optional (e.g. HR applicants page); default is public job detail.
 * onEdit, onPublish, onDelete, onClose: optional HR actions.
 */
export function JobCard({ job, viewHref, onEdit, onPublish, onDelete, onClose }) {
  const isClosed = job?.status === 'CLOSED' || job?.Status === 'CLOSED'
  const jobId = job?.id ?? job?.ID
  const title = job?.title ?? job?.Title ?? 'Untitled job'
  const description = job?.description ?? job?.Description ?? ''
  const department = job?.department ?? job?.Department
  const viewTo = viewHref ?? (jobId ? ROUTE_BUILDERS.publicJobDetail(jobId) : null)
  return (
    <Card>
      <CardHeader className="flex flex-row items-start justify-between gap-2">
        <div>
          <CardContent className="p-0">
            {viewTo ? (
              <Link to={viewTo} className="font-semibold hover:underline">
                {title}
              </Link>
            ) : (
              <span className="font-semibold">{title}</span>
            )}
          </CardContent>
          {department && <p className="text-sm text-muted-foreground mt-1">{department}</p>}
        </div>
        {isClosed && <Badge variant="secondary">Closed</Badge>}
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground line-clamp-2">{description}</p>
      </CardContent>
      <CardFooter className="flex flex-wrap gap-2">
        {viewTo && (
          <Button variant="outline" size="sm" asChild>
            <Link to={viewTo}>View</Link>
          </Button>
        )}
        {onEdit && (
          <Button variant="outline" size="sm" onClick={onEdit}>Edit</Button>
        )}
        {onPublish && (
          <Button variant="default" size="sm" onClick={onPublish}>Publish</Button>
        )}
        {onClose && (
          <Button variant="outline" size="sm" onClick={onClose}>Close</Button>
        )}
        {onDelete && (
          <Button variant="outline" size="sm" className="text-destructive hover:text-destructive" onClick={onDelete}>Delete</Button>
        )}
      </CardFooter>
    </Card>
  )
}
