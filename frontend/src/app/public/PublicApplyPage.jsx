import { useState, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useAuth } from '@/modules/auth/hooks/useAuth'
import { useJob } from '@/modules/jobs/hooks/useJobs'
import { useApplyToJob } from '@/modules/jobs/hooks/useApplications'
import { resumesApi } from '@/services/api'
import { toast } from 'sonner'
import { ROUTES, ROUTE_BUILDERS } from '@/constants'
import { Skeleton } from '@/components/ui/skeleton'

const PDF_ACCEPT = '.pdf,application/pdf'

export function PublicApplyPage() {
  const { jobId } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuth()
  const candidateId = user?.candidate_id ?? user?.candidateId
  const { data: job, isPending: jobLoading } = useJob(jobId)
  const applyMutation = useApplyToJob(jobId, candidateId)
  const { data: resumes = [] } = useQuery({
    queryKey: ['resumes', candidateId],
    queryFn: () => resumesApi.list(candidateId),
    enabled: !!candidateId,
  })
  const fileInputRef = useRef(null)

  const uploadResumeMutation = useMutation({
    mutationFn: ({ candidateId: cid, file }) => resumesApi.upload(cid, file),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['resumes', candidateId] })
      const id = data?.id ?? data?.ID
      if (id) setResume(id)
      setPendingFile(null)
      toast.success('Resume uploaded')
      if (fileInputRef.current) fileInputRef.current.value = ''
      // When user uploaded via file input (not from Submit): auto-apply and send to My Applications
      const shouldApplyAndRedirect = variables?.applyAndRedirect && id && jobId && candidateId
      if (shouldApplyAndRedirect) {
        applyMutation.mutate(
          { resumeId: id },
          {
            onSuccess: () => {
              toast.success('Application submitted')
              navigate(ROUTES.CANDIDATE_APPLICATIONS)
            },
            onError: () => toast.error('Failed to submit application'),
          }
        )
      }
    },
    onError: (err) => {
      const msg = err?.response?.data?.error?.message ?? err?.response?.data?.message ?? err?.message ?? 'Failed to upload resume.'
      toast.error(msg)
    },
  })

  const [resume, setResume] = useState(null)
  const [pendingFile, setPendingFile] = useState(null)
  const [errors, setErrors] = useState({})

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!candidateId) {
      toast.error('Please sign in to apply')
      navigate(ROUTES.LOGIN)
      return
    }
    let resumeId = typeof resume === 'string' ? resume : resume?.id ?? resume?.ID
    // If no resume selected but user has a file chosen, upload it first then apply
    if (!resumeId && pendingFile) {
      if (pendingFile.size > 100 * 1024 * 1024) {
        toast.error('File must be under 100 MB')
        return
      }
      try {
        const data = await uploadResumeMutation.mutateAsync({ candidateId, file: pendingFile })
        resumeId = data?.id ?? data?.ID
        setPendingFile(null)
        if (fileInputRef.current) fileInputRef.current.value = ''
      } catch {
        return // toast already shown in mutation onError
      }
    }
    if (!resumeId) {
      toast.error('Please select a resume or choose a PDF to upload')
      return
    }
    applyMutation.mutate(
      { resumeId },
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
        <CardTitle>Apply for {job?.title ?? job?.Title ?? 'Job'}</CardTitle>
        <CardDescription>Submit your application</CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit} noValidate>
        <CardContent className="space-y-4">
          {!candidateId && (
            <p className="text-sm text-muted-foreground">Sign in to apply. You will need a resume uploaded in your profile.</p>
          )}
          {candidateId && (
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="resume">Resume</Label>
                <Select value={resume ?? ''} onValueChange={setResume}>
                  <SelectTrigger id="resume">
                    <SelectValue placeholder="Select a resume" />
                  </SelectTrigger>
                  <SelectContent>
                    {(Array.isArray(resumes) ? resumes : []).map((r) => {
                      const id = r.id ?? r.ID
                      const name = r.file_name ?? r.fileName ?? r.FileName ?? `Resume ${(id || '').slice(0, 8)}`
                      return (
                        <SelectItem key={id} value={id}>
                          {name}
                        </SelectItem>
                      )
                    })}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Or upload a PDF</Label>
                <Input
                  ref={fileInputRef}
                  type="file"
                  accept={PDF_ACCEPT}
                  className="cursor-pointer"
                  disabled={uploadResumeMutation.isPending}
                  onChange={(e) => {
                    const file = e.target.files?.[0]
                    if (!file) {
                      setPendingFile(null)
                      return
                    }
                    if (file.type !== 'application/pdf' && !file.name.toLowerCase().endsWith('.pdf')) {
                      toast.error('Only PDF format is allowed')
                      return
                    }
                    if (file.size > 100 * 1024 * 1024) {
                      toast.error('File must be under 100 MB')
                      return
                    }
                    setPendingFile(file)
                    uploadResumeMutation.mutate({ candidateId, file, applyAndRedirect: true })
                  }}
                />
                <p className="text-xs text-muted-foreground">PDF only, max 100 MB. Stored in Supabase Storage.</p>
                {uploadResumeMutation.isPending && (
                  <p className="text-sm text-muted-foreground">Uploading…</p>
                )}
              </div>
            </div>
          )}
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
          <Button
            type="submit"
            disabled={applyMutation.isPending || (!!pendingFile && uploadResumeMutation.isPending)}
          >
            {applyMutation.isPending
              ? 'Submitting...'
              : pendingFile && uploadResumeMutation.isPending
                ? 'Uploading...'
                : 'Submit application'}
          </Button>
        </CardFooter>
      </form>
    </Card>
  )
}
