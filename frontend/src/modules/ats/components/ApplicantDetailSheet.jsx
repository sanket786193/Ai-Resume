import { useQueryClient } from '@tanstack/react-query'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetClose,
} from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { useApplicationDetail } from '../hooks/useApplicationDetail'
import { useUpdateCandidateStatus } from '../hooks/useAtsPipeline'
import { ATS_STATUS_LABEL } from '@/constants'
import { JOB_APPLICANTS_QUERY_KEY } from '@/modules/jobs/hooks/useJobs'
import { APPLICATION_DETAIL_QUERY_KEY } from '../hooks/useApplicationDetail'
import { cn } from '@/lib/utils'

const DetailSection = ({ title, children, className }) => (
  <div className={cn('space-y-2', className)}>
    <h3 className="text-sm font-semibold text-foreground border-b pb-1">{title}</h3>
    {children}
  </div>
)

const DetailRow = ({ label, value, href }) => {
  if (value == null || value === '') return null
  return (
    <div className="flex flex-col gap-0.5">
      <span className="text-xs text-muted-foreground">{label}</span>
      {href ? (
        <a href={href} target="_blank" rel="noopener noreferrer" className="text-sm text-primary hover:underline truncate">
          {value}
        </a>
      ) : (
        <span className="text-sm">{value}</span>
      )}
    </div>
  )
}

export function ApplicantDetailSheet({ open, onOpenChange, applicationId, jobId }) {
  const queryClient = useQueryClient()
  const { data: detail, isPending, isError } = useApplicationDetail(applicationId, { enabled: open && Boolean(applicationId) })
  const updateStatus = useUpdateCandidateStatus()

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

  const name = detail?.candidate_name ?? detail?.candidate_email ?? 'Applicant'
  const hasContact = name || detail?.candidate_email || detail?.candidate_phone || detail?.candidate_linkedin
  const hasApplication = detail?.status != null || detail?.created_at != null || detail?.job_title
  const hasResume = detail?.resume_file_name
  const hasAI =
    detail?.ats_score != null ||
    detail?.skill_match_score != null ||
    detail?.skill_match_pct != null ||
    detail?.qualified != null ||
    (detail?.ai_summary && detail.ai_summary !== '') ||
    (detail?.experience_match && detail.experience_match !== '') ||
    (Array.isArray(detail?.missing_skills) && detail.missing_skills.length > 0)

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-md overflow-y-auto">
        <SheetHeader className="flex flex-row items-start justify-between space-y-0">
          <SheetTitle className="pr-8">{name}</SheetTitle>
          <SheetClose asChild>
            <button
              type="button"
              className="rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none absolute right-6 top-6 p-1 text-muted-foreground hover:text-foreground"
              aria-label="Close"
            >
              ✕
            </button>
          </SheetClose>
        </SheetHeader>

        <div className="mt-6 space-y-6">
          {isPending && (
            <div className="space-y-4">
              <Skeleton className="h-20 w-full" />
              <Skeleton className="h-24 w-full" />
              <Skeleton className="h-16 w-full" />
            </div>
          )}

          {isError && (
            <p className="text-sm text-destructive">Failed to load applicant details. Please try again.</p>
          )}

          {!isPending && !isError && detail && (
            <>
              {hasContact && (
                <DetailSection title="Contact">
                  <div className="grid gap-3">
                    <DetailRow label="Name" value={detail.candidate_name} />
                    <DetailRow label="Email" value={detail.candidate_email} />
                    <DetailRow label="Phone" value={detail.candidate_phone} />
                    <DetailRow
                      label="LinkedIn"
                      value={detail.candidate_linkedin ? 'Profile' : null}
                      href={detail.candidate_linkedin || undefined}
                    />
                  </div>
                </DetailSection>
              )}

              {hasApplication && (
                <DetailSection title="Application">
                  <div className="space-y-3">
                    <div className="flex flex-col gap-1.5">
                      <Label className="text-xs">Status</Label>
                      <Select
                        value={detail.status ?? ''}
                        onValueChange={handleStatusChange}
                        disabled={updateStatus.isPending}
                      >
                        <SelectTrigger className="w-full">
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
                    <DetailRow
                      label="Applied"
                      value={detail.created_at ? new Date(detail.created_at).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' }) : null}
                    />
                    <DetailRow label="Job" value={detail.job_title} />
                  </div>
                </DetailSection>
              )}

              {hasResume && (
                <DetailSection title="Resume">
                  <DetailRow label="File" value={detail.resume_file_name} />
                </DetailSection>
              )}

              {hasAI && (
                <DetailSection title="AI feedback">
                  <div className="space-y-3">
                    {(detail.ats_score != null || detail.skill_match_score != null || detail.skill_match_pct != null) && (
                      <div className="flex flex-wrap gap-3">
                        {detail.ats_score != null && (
                          <div>
                            <span className="text-xs text-muted-foreground">ATS score</span>
                            <p className="text-sm font-medium">{detail.ats_score}/100</p>
                          </div>
                        )}
                        {detail.skill_match_pct != null && (
                          <div>
                            <span className="text-xs text-muted-foreground">Skill match</span>
                            <p className="text-sm font-medium">{detail.skill_match_pct}%</p>
                          </div>
                        )}
                        {detail.skill_match_score != null && detail.skill_match_pct == null && (
                          <div>
                            <span className="text-xs text-muted-foreground">Match score</span>
                            <p className="text-sm font-medium">{Number(detail.skill_match_score).toFixed(2)}</p>
                          </div>
                        )}
                      </div>
                    )}
                    {detail.qualified != null && (
                      <DetailRow label="Qualified" value={detail.qualified ? 'Yes' : 'No'} />
                    )}
                    {detail.ai_summary && detail.ai_summary !== '' && (
                      <div className="flex flex-col gap-0.5">
                        <span className="text-xs text-muted-foreground">Summary</span>
                        <p className="text-sm whitespace-pre-wrap">{detail.ai_summary}</p>
                      </div>
                    )}
                    {detail.experience_match && detail.experience_match !== '' && (
                      <div className="flex flex-col gap-0.5">
                        <span className="text-xs text-muted-foreground">Experience match</span>
                        <p className="text-sm">{detail.experience_match}</p>
                      </div>
                    )}
                    {Array.isArray(detail.missing_skills) && detail.missing_skills.length > 0 && (
                      <div className="flex flex-col gap-1">
                        <span className="text-xs text-muted-foreground">Missing skills</span>
                        <ul className="text-sm list-disc list-inside space-y-0.5">
                          {detail.missing_skills.map((s, i) => (
                            <li key={i}>{s}</li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                </DetailSection>
              )}
            </>
          )}
        </div>
      </SheetContent>
    </Sheet>
  )
}
