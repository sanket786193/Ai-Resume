import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { useJob } from '@/modules/jobs/hooks/useJobs'
import { useApplyToJob } from '@/modules/jobs/hooks/useApplications'
import { toast } from 'sonner'
import { ROUTES, ROUTE_BUILDERS } from '@/constants'
import { validateEmail, validateRequired } from '@/utils/validation'
import { Skeleton } from '@/components/ui/skeleton'

export function PublicApplyPage() {
  const { jobId } = useParams()
  const navigate = useNavigate()
  const { data: job, isPending: jobLoading } = useJob(jobId)
  const applyMutation = useApplyToJob(jobId)

  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [resume, setResume] = useState(null)
  const [errors, setErrors] = useState({})

  const handleSubmit = (e) => {
    e.preventDefault()
    const next = {}
    const nameErr = validateRequired(name.trim(), 'Full name')
    if (nameErr) next.name = nameErr
    const emailErr = validateEmail(email)
    if (emailErr) next.email = emailErr
    setErrors(next)
    if (Object.keys(next).length > 0) return

    applyMutation.mutate(
      { name: name.trim(), email, resume },
      {
        onSuccess: () => {
          toast.success('Application submitted')
          navigate(ROUTES.CANDIDATE_APPLICATIONS)
        },
        onError: () => toast.error('Failed to submit application'),
      }
    )
  }

  if (jobLoading || !job) {
    return <Skeleton className="h-64 w-full rounded-lg" />
  }

  return (
    <Card className="max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Apply for {job.title}</CardTitle>
        <CardDescription>Submit your application</CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit} noValidate>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Full name</Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => { setName(e.target.value); setErrors((p) => ({ ...p, name: null })) }}
              aria-invalid={!!errors.name}
              aria-describedby={errors.name ? 'name-error' : undefined}
            />
            {errors.name && <p id="name-error" className="text-sm text-destructive" role="alert">{errors.name}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => { setEmail(e.target.value); setErrors((p) => ({ ...p, email: null })) }}
              aria-invalid={!!errors.email}
              aria-describedby={errors.email ? 'email-error' : undefined}
            />
            {errors.email && <p id="email-error" className="text-sm text-destructive" role="alert">{errors.email}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="resume">Resume (file)</Label>
            <Input
              id="resume"
              type="file"
              accept=".pdf,.doc,.docx"
              onChange={(e) => setResume(e.target.files?.[0] ?? null)}
            />
          </div>
          {applyMutation.error && (
            <p className="text-sm text-destructive" role="alert">
              {applyMutation.error?.response?.data?.message ?? applyMutation.error?.message ?? 'Application failed'}
            </p>
          )}
        </CardContent>
        <CardFooter className="flex gap-2">
          <Button type="button" variant="outline" onClick={() => navigate(ROUTE_BUILDERS.publicJobDetail(jobId))}>
            Cancel
          </Button>
          <Button type="submit" disabled={applyMutation.isPending}>
            {applyMutation.isPending ? 'Submitting...' : 'Submit application'}
          </Button>
        </CardFooter>
      </form>
    </Card>
  )
}
