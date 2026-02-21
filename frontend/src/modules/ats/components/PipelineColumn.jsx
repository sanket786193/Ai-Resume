import { useDroppable } from '@dnd-kit/core'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ATS_STATUS_LABEL } from '@/constants'
import { PipelineCandidateCard } from './PipelineCandidateCard'

/**
 * Single Kanban column. Droppable; renders candidate cards.
 */
export function PipelineColumn({ status, candidates, onCardClick }) {
  const { isOver, setNodeRef } = useDroppable({ id: status })
  const label = ATS_STATUS_LABEL[status] ?? status

  return (
    <div ref={setNodeRef} className="min-w-[280px] flex-shrink-0">
      <Card className={`h-full transition-colors ${isOver ? 'ring-2 ring-primary' : ''}`}>
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <span className="font-medium">{label}</span>
            <Badge variant="secondary">{candidates?.length ?? 0}</Badge>
          </div>
        </CardHeader>
        <CardContent className="space-y-2 pt-0">
          {candidates?.map((c) => (
            <PipelineCandidateCard key={c.id} candidate={c} onClick={() => onCardClick?.(c)} />
          ))}
        </CardContent>
      </Card>
    </div>
  )
}
