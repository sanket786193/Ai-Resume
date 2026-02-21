import { useDraggable } from '@dnd-kit/core'
import { Card, CardContent } from '@/components/ui/card'

/**
 * Draggable candidate card for pipeline Kanban.
 */
export function PipelineCandidateCard({ candidate, onClick }) {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: String(candidate.id),
    data: candidate,
  })

  return (
    <Card
      ref={setNodeRef}
      className={`cursor-grab active:cursor-grabbing ${isDragging ? 'opacity-50' : ''}`}
      {...listeners}
      {...attributes}
      onClick={onClick}
    >
      <CardContent className="p-3">
        <p className="font-medium text-sm truncate">{candidate.name ?? candidate.email ?? 'Unknown'}</p>
        <p className="text-xs text-muted-foreground truncate">{candidate.email}</p>
        {candidate.currentRole && (
          <p className="text-xs text-muted-foreground mt-1 truncate">{candidate.currentRole}</p>
        )}
      </CardContent>
    </Card>
  )
}
