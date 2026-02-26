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
const EXPERIENCE_LABEL = { ANY: 'Any', FRESHER: 'Fresher', EXPERIENCED: 'Experienced' }

export function JobCard({ job, viewHref, onEdit, onPublish, onDelete, onClose }) {
  const isClosed = job?.status === 'CLOSED' || job?.Status === 'CLOSED'
  const jobId = job?.id ?? job?.ID
  const title = job?.title ?? job?.Title ?? 'Untitled job'
  const description = job?.description ?? job?.Description ?? ''
  const department = job?.department ?? job?.Department
  const experienceLevel = job?.experience_level ?? job?.experienceLevel
  const qualification = job?.qualification ?? job?.Qualification
  const skills = job?.skills ?? job?.Skills ?? []
  const vacancyLimits = job?.vacancy_limits ?? job?.vacancyLimits ?? job?.VacancyLimits ?? []
  const viewTo = viewHref ?? (jobId ? ROUTE_BUILDERS.publicJobDetail(jobId) : null)
  const hasRequirements = qualification || (Array.isArray(skills) && skills.length > 0) || (Array.isArray(vacancyLimits) && vacancyLimits.length > 0)
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
          {(experienceLevel && experienceLevel !== 'ANY') && (
            <p className="text-xs text-muted-foreground mt-0.5">{EXPERIENCE_LABEL[experienceLevel] ?? experienceLevel}</p>
          )}
        </div>
        {isClosed && <Badge variant="secondary">Closed</Badge>}
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground line-clamp-2">{description}</p>
        {hasRequirements && (
          <div className="mt-2 flex flex-wrap gap-1.5 text-xs text-muted-foreground">
            {qualification && <span className="rounded bg-muted px-1.5 py-0.5">{qualification}</span>}
            {Array.isArray(skills) && skills.slice(0, 4).map((s) => <span key={s} className="rounded bg-muted px-1.5 py-0.5">{s}</span>)}
            {Array.isArray(vacancyLimits) && vacancyLimits.length > 0 && (
              <span className="rounded bg-muted px-1.5 py-0.5">
                {vacancyLimits.map((v) => `${v.role ?? v.Role ?? ''}: ${v.limit ?? v.Limit ?? 0}`).filter(Boolean).join(', ')}
              </span>
            )}
          </div>
        )}
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
