import { useDroppable } from '@dnd-kit/core'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ATS_STATUS_LABEL, ATS_STATUS_BADGE_VARIANT } from '@/constants'
import { PipelineCandidateCard } from './PipelineCandidateCard'
import { cn } from '@/lib/utils'

/**
 * Single Kanban column. Droppable; renders candidate cards.
 */
export function PipelineColumn({ status, candidates, onCardClick }) {
  const { isOver, setNodeRef } = useDroppable({ id: status })
  const label = ATS_STATUS_LABEL[status] ?? status
  const count = candidates?.length ?? 0
  const variant = ATS_STATUS_BADGE_VARIANT[status] ?? 'secondary'

  return (
    <div ref={setNodeRef} className="min-w-[288px] flex-shrink-0">
      <Card
        className={cn(
          'h-full min-h-[320px] flex flex-col transition-all duration-200',
          isOver && 'ring-2 ring-primary bg-primary/5'
        )}
      >
        <CardHeader className="py-4 pb-2">
          <div className="flex items-center justify-between gap-2">
            <span className="font-semibold text-foreground">{label}</span>
            <Badge variant={variant} className="shrink-0">
              {count}
            </Badge>
          </div>
        </CardHeader>
        <CardContent className="space-y-2 pt-0 flex-1 overflow-y-auto">
          {count > 0 ? (
            candidates.map((c) => (
              <PipelineCandidateCard key={c.id} candidate={c} onClick={() => onCardClick?.(c)} />
            ))
          ) : (
            <div
              className="rounded-lg border border-dashed border-muted-foreground/25 bg-muted/30 py-8 px-4 text-center text-sm text-muted-foreground"
              aria-hidden
            >
              Drop candidates here
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
