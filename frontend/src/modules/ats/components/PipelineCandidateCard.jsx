import { useDraggable } from '@dnd-kit/core'
import { Card, CardContent } from '@/components/ui/card'
import { cn } from '@/lib/utils'

/** Display name from API (candidate_name/name) or email, or fallback */
function getDisplayName(c) {
  return c?.candidate_name ?? c?.name ?? c?.candidate_email ?? c?.email ?? 'Applicant'
}

/** Display email from API */
function getEmail(c) {
  return c?.candidate_email ?? c?.email ?? ''
}

/** Job title applied to */
function getJobTitle(c) {
  return c?.job_title ?? c?.currentRole ?? ''
}

/** Initial for avatar (first letter of name or email) */
function getInitial(c) {
  const name = getDisplayName(c)
  if (name && name !== 'Applicant') return name.charAt(0).toUpperCase()
  const email = getEmail(c)
  return email ? email.charAt(0).toUpperCase() : '?'
}

/**
 * Draggable candidate card for pipeline Kanban.
 */
export function PipelineCandidateCard({ candidate, onClick }) {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: String(candidate.id),
    data: candidate,
  })
  const displayName = getDisplayName(candidate)
  const email = getEmail(candidate)
  const jobTitle = getJobTitle(candidate)
  const initial = getInitial(candidate)

  return (
    <Card
      ref={setNodeRef}
      className={cn(
        'cursor-grab active:cursor-grabbing transition-shadow hover:shadow-md',
        isDragging && 'opacity-50 shadow-lg ring-2 ring-primary/20'
      )}
      {...listeners}
      {...attributes}
      onClick={onClick}
    >
      <CardContent className="p-3 flex gap-3">
        <div
          className="flex-shrink-0 w-9 h-9 rounded-full bg-primary/10 text-primary flex items-center justify-center text-sm font-semibold"
          aria-hidden
        >
          {initial}
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-medium text-sm truncate">{displayName}</p>
          {email && (
            <p className="text-xs text-muted-foreground truncate">{email}</p>
          )}
          {jobTitle && (
            <p className="text-xs text-muted-foreground mt-0.5 truncate" title={jobTitle}>
              {jobTitle}
            </p>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

export { getDisplayName, getEmail, getJobTitle, getInitial }
