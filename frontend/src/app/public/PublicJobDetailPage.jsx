import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { SafeMarkdown } from '@/components/editor/SafeMarkdown'
import { useJob } from '@/modules/jobs/hooks/useJobs'
import { ROUTES, ROUTE_BUILDERS } from '@/constants'

export function PublicJobDetailPage() {
  const { jobId } = useParams()
  const [descriptionDialogOpen, setDescriptionDialogOpen] = useState(false)
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

  const id = job.id ?? job.ID
  const title = job.title ?? job.Title ?? 'Job'
  const department = job.department ?? job.Department
  const description = job.description ?? job.Description ?? ''
  const status = job.status ?? job.Status
  const experienceLevel = job.experience_level ?? job.experienceLevel
  const qualification = job.qualification ?? job.Qualification
  const skills = job.skills ?? job.Skills ?? []
  const vacancyLimits = job.vacancy_limits ?? job.vacancyLimits ?? job.VacancyLimits ?? []
  const isClosed = status === 'CLOSED'
  const hasLongDescription = description.length > 400
  const experienceLabel = { ANY: 'Any', FRESHER: 'Fresher', EXPERIENCED: 'Experienced' }[experienceLevel] ?? experienceLevel

  return (
    <div>
      <Card>
        <CardHeader className="flex flex-row items-start justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold">{title}</h1>
            {department && <p className="text-muted-foreground">{department}</p>}
            {(experienceLevel && experienceLevel !== 'ANY') && <p className="text-sm text-muted-foreground">{experienceLabel}</p>}
          </div>
          {isClosed && <Badge variant="secondary">Closed</Badge>}
        </CardHeader>
        <CardContent className="space-y-4">
          {(qualification || (Array.isArray(skills) && skills.length > 0) || (Array.isArray(vacancyLimits) && vacancyLimits.length > 0)) && (
            <div className="rounded-lg border bg-muted/30 p-4 space-y-2">
              <h3 className="text-sm font-medium">Requirements</h3>
              {qualification && <p className="text-sm"><span className="text-muted-foreground">Qualification:</span> {qualification}</p>}
              {Array.isArray(skills) && skills.length > 0 && (
                <p className="text-sm"><span className="text-muted-foreground">Skills:</span> {skills.join(', ')}</p>
              )}
              {Array.isArray(vacancyLimits) && vacancyLimits.length > 0 && (
                <p className="text-sm"><span className="text-muted-foreground">Vacancies:</span> {vacancyLimits.map((v) => `${v.role ?? v.Role ?? ''} (${v.limit ?? v.Limit ?? 0})`).filter(Boolean).join(', ')}</p>
              )}
            </div>
          )}
          <div>
            <SafeMarkdown content={description} />
          </div>
          {description.trim() && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setDescriptionDialogOpen(true)}
            >
              {hasLongDescription ? 'View full description' : 'View description'}
            </Button>
          )}
        </CardContent>
        {!isClosed && id && (
          <CardContent className="pt-0">
            <Button asChild>
              <Link to={ROUTE_BUILDERS.publicApply(id)}>Apply for this job</Link>
            </Button>
          </CardContent>
        )}
      </Card>

      <Dialog open={descriptionDialogOpen} onOpenChange={setDescriptionDialogOpen}>
        <DialogContent className="max-w-2xl max-h-[85vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Job description — {title}</DialogTitle>
          </DialogHeader>
          <SafeMarkdown content={description} className="pt-2" />
        </DialogContent>
      </Dialog>
    </div>
  )
}
