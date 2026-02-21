import { useQueryClient } from '@tanstack/react-query'
import { useOutletContext } from 'react-router-dom'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useUpdateCandidateStatus } from '@/modules/ats/hooks/useAtsPipeline'
import { APPLICATION_DETAIL_QUERY_KEY } from '@/modules/ats/hooks/useApplicationDetail'
import { ATS_STATUS_LABEL } from '@/constants'
import { JOB_APPLICANTS_QUERY_KEY } from '@/modules/jobs/hooks/useJobs'

const DetailRow = ({ label, value }) => {
  if (value == null || value === '') return null
  return (
    <div className="flex items-start gap-3 py-3 border-b border-border/50 last:border-0">
      <div className="flex-shrink-0 w-9 h-9 rounded-lg bg-muted/60 flex items-center justify-center text-muted-foreground text-xs font-medium" aria-hidden>
        {label.charAt(0)}
      </div>
      <div className="min-w-0 flex-1">
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{label}</span>
        <p className="text-sm font-medium mt-0.5 break-words">{value}</p>
      </div>
    </div>
  )
}

export function ApplicantApplicationPage() {
  const { detail, jobId } = useOutletContext()
  const queryClient = useQueryClient()
  const updateStatus = useUpdateCandidateStatus()
  const applicationId = detail?.id ?? detail?.ID

  const handleStatusChange = (newStatus) => {
    if (!applicationId || !newStatus) return
    updateStatus.mutate(
      { id: applicationId, status: newStatus },
      {
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: APPLICATION_DETAIL_QUERY_KEY(applicationId) })
          if (jobId) queryClient.invalidateQueries({ queryKey: JOB_APPLICANTS_QUERY_KEY(jobId) })
        },
      }
    )
  }

  if (!detail) return null

  const appliedAt = detail.created_at
    ? new Date(detail.created_at).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })
    : null

  return (
    <Card>
      <CardHeader className="pb-2">
        <h2 className="text-base font-semibold">Application</h2>
        <p className="text-sm text-muted-foreground">
          Status and job for this application
        </p>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Status</Label>
          <Select
            value={detail.status ?? ''}
            onValueChange={handleStatusChange}
            disabled={updateStatus.isPending}
          >
            <SelectTrigger className="w-full max-w-xs">
              <SelectValue placeholder="Select status" />
            </SelectTrigger>
            <SelectContent>
              {Object.entries(ATS_STATUS_LABEL).map(([value, label]) => (
                <SelectItem key={value} value={value}>
                  {label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-0 pt-2 border-t border-border/50">
          <DetailRow label="Applied" value={appliedAt} />
          <DetailRow label="Job" value={detail.job_title} />
        </div>
      </CardContent>
    </Card>
  )
}
