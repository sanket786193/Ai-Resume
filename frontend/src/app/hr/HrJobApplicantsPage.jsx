import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { PageHeader } from '@/components/common/PageHeader'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Skeleton } from '@/components/ui/skeleton'
import { EmptyState } from '@/components/common/EmptyState'
import { useJob, useJobApplicants, useBulkApply } from '@/modules/jobs/hooks/useJobs'
import { ATS_STATUS_LABEL, ATS_STATUS_BADGE_VARIANT, ROUTES, ROUTE_BUILDERS } from '@/constants'
import { cn } from '@/lib/utils'
import { toast } from 'sonner'

function getDisplayName(a) {
  return a?.candidate_name ?? a?.name ?? a?.candidate_email ?? a?.email ?? 'Applicant'
}

function getEmail(a) {
  return a?.candidate_email ?? a?.email ?? ''
}

function getInitial(a) {
  const name = getDisplayName(a)
  if (name && name !== 'Applicant') return name.charAt(0).toUpperCase()
  const email = getEmail(a)
  return email ? email.charAt(0).toUpperCase() : '?'
}

export function HrJobApplicantsPage() {
  const { jobId } = useParams()
  const navigate = useNavigate()
  const [bulkOpen, setBulkOpen] = useState(false)
  const [bulkRows, setBulkRows] = useState([{ candidate_id: '', resume_id: '' }])

  const validJobId = jobId != null && String(jobId) !== '' && String(jobId) !== 'undefined'
  const { data: job, isPending: jobLoading } = useJob(jobId)
  const { data: applicants, isPending: applicantsLoading } = useJobApplicants(jobId)
  const bulkApply = useBulkApply(jobId)

  const list = Array.isArray(applicants) ? applicants : applicants?.list ?? []

  const openDetail = (appId) => {
    navigate(ROUTE_BUILDERS.hrApplicantDetail(jobId, appId, 'contact'))
  }

  const addBulkRow = () => setBulkRows((r) => [...r, { candidate_id: '', resume_id: '' }])
  const updateBulkRow = (index, field, value) => {
    setBulkRows((r) => {
      const next = [...r]
      next[index] = { ...next[index], [field]: value }
      return next
    })
  }
  const removeBulkRow = (index) => setBulkRows((r) => r.filter((_, i) => i !== index))
  const submitBulk = () => {
    const applications = bulkRows
      .filter((row) => (row.candidate_id ?? '').trim() && (row.resume_id ?? '').trim())
      .map((row) => ({ candidate_id: String(row.candidate_id).trim(), resume_id: String(row.resume_id).trim() }))
    if (applications.length === 0) {
      toast.error('Add at least one row with candidate_id and resume_id')
      return
    }
    bulkApply.mutate(
      { applications },
      {
        onSuccess: (data) => {
          toast.success(`Created ${data?.created ?? 0} application(s), skipped ${data?.skipped ?? 0}`)
          if (data?.errors?.length) toast.error(data.errors.slice(0, 3).join('; '))
          setBulkOpen(false)
          setBulkRows([{ candidate_id: '', resume_id: '' }])
        },
        onError: () => toast.error('Bulk apply failed'),
      }
    )
  }

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
        actions={
          <Button variant="outline" size="sm" onClick={() => setBulkOpen(true)}>
            Bulk add applicants
          </Button>
        }
      />
      <Dialog open={bulkOpen} onOpenChange={setBulkOpen}>
        <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Bulk add applicants</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            Add one row per application: candidate_id and resume_id (must belong to that candidate). AI screening will run for each new application.
          </p>
          <div className="space-y-2">
            {bulkRows.map((row, index) => (
              <div key={index} className="flex gap-2 items-center">
                <Input
                  placeholder="Candidate ID"
                  value={row.candidate_id ?? ''}
                  onChange={(e) => updateBulkRow(index, 'candidate_id', e.target.value)}
                  className="flex-1"
                />
                <Input
                  placeholder="Resume ID"
                  value={row.resume_id ?? ''}
                  onChange={(e) => updateBulkRow(index, 'resume_id', e.target.value)}
                  className="flex-1"
                />
                <Button type="button" variant="ghost" size="sm" onClick={() => removeBulkRow(index)} aria-label="Remove row">
                  ×
                </Button>
              </div>
            ))}
            <Button type="button" variant="outline" size="sm" onClick={addBulkRow}>
              Add row
            </Button>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setBulkOpen(false)}>Cancel</Button>
            <Button onClick={submitBulk} disabled={bulkApply.isPending}>
              {bulkApply.isPending ? 'Adding…' : 'Add applicants'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      {applicantsLoading && <Skeleton className="h-48 w-full rounded-lg" />}
      {!applicantsLoading && list.length === 0 && (
        <EmptyState title="No applicants yet" description="Applicants will appear here when candidates apply." />
      )}
      {!applicantsLoading && list.length > 0 && (
        <div className="space-y-3">
          {list.map((a) => {
            const appId = a.id ?? a.ID
            const status = a.status ?? a.Status
            const variant = ATS_STATUS_BADGE_VARIANT[status] ?? 'secondary'
            const appliedAt = a.created_at ? new Date(a.created_at).toLocaleDateString(undefined, { dateStyle: 'medium' }) : null
            return (
              <Card
                key={appId}
                className={cn(
                  'transition-colors cursor-pointer hover:bg-muted/50 hover:border-primary/20',
                  'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
                )}
                tabIndex={0}
                onClick={() => openDetail(appId)}
                onKeyDown={(e) => (e.key === 'Enter' || e.key === ' ') && openDetail(appId)}
              >
                <CardContent className="p-4 flex items-center gap-4">
                  <div
                    className="flex-shrink-0 w-11 h-11 rounded-full bg-primary/10 text-primary flex items-center justify-center text-base font-semibold"
                    aria-hidden
                  >
                    {getInitial(a)}
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="font-medium truncate">{getDisplayName(a)}</p>
                    {getEmail(a) && (
                      <p className="text-sm text-muted-foreground truncate">{getEmail(a)}</p>
                    )}
                    {appliedAt && (
                      <p className="text-xs text-muted-foreground mt-0.5">Applied {appliedAt}</p>
                    )}
                  </div>
                  <div className="flex items-center gap-2 shrink-0">
                    {a.ats_score != null && (
                      <span className="text-xs text-muted-foreground" title="ATS score">
                        {a.ats_score}/100
                      </span>
                    )}
                    <Badge variant={variant}>{ATS_STATUS_LABEL[status] ?? status}</Badge>
                  </div>
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}
