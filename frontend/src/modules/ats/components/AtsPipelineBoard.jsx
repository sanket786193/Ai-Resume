import { useMemo, useState } from 'react'
import {
  DndContext,
  DragOverlay,
  closestCenter,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import { ATS_KANBAN_STATUSES } from '@/constants'
import { useAtsPipeline, useUpdateCandidateStatus } from '../hooks/useAtsPipeline'
import { PipelineColumn } from './PipelineColumn'
import { PipelineCandidateCard, getDisplayName, getEmail, getJobTitle, getInitial } from './PipelineCandidateCard'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { EmptyState } from '@/components/common/EmptyState'

/**
 * Kanban board for ATS pipeline. Groups candidates by status, drag-and-drop with optimistic update.
 */
export function AtsPipelineBoard() {
  const { data, isPending, isError, error } = useAtsPipeline()
  const updateStatus = useUpdateCandidateStatus()
  const [activeCandidate, setActiveCandidate] = useState(null)

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } })
  )

  const candidatesByStatus = useMemo(() => {
    const list = Array.isArray(data) ? data : data?.list ?? []
    const map = {}
    ATS_KANBAN_STATUSES.forEach((s) => { map[s] = [] })
    list.forEach((c) => {
      const status = c.status ?? 'APPLIED'
      if (!map[status]) map[status] = []
      map[status].push(c)
    })
    return map
  }, [data])

  const totalCandidates = useMemo(
    () => Object.values(candidatesByStatus).reduce((sum, arr) => sum + (arr?.length ?? 0), 0),
    [candidatesByStatus]
  )

  const handleDragStart = (event) => {
    setActiveCandidate(event.active?.data?.current ?? null)
  }

  const handleDragEnd = (event) => {
    setActiveCandidate(null)
    const { active, over } = event
    if (!over?.id) return
    const candidate = active?.data?.current
    const newStatus = String(over.id)
    if (!candidate?.id || !ATS_KANBAN_STATUSES.includes(newStatus)) return
    if (candidate.status === newStatus) return
    updateStatus.mutate({ id: candidate.id, status: newStatus })
  }

  if (isPending) {
    return (
      <div className="flex gap-4 overflow-x-auto pb-4">
        {ATS_KANBAN_STATUSES.map((s) => (
          <Skeleton key={s} className="h-[420px] w-[288px] flex-shrink-0 rounded-xl" />
        ))}
      </div>
    )
  }

  if (isError) {
    return (
      <EmptyState
        title="Failed to load pipeline"
        description={error?.message ?? 'Please try again later.'}
      />
    )
  }

  if (totalCandidates === 0) {
    return (
      <EmptyState
        title="No applications yet"
        description="When candidates apply to your jobs, they will appear here. You can drag them between stages to update their status."
      />
    )
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <div className="flex gap-4 overflow-x-auto pb-4 -mx-1 px-1" style={{ minHeight: 360 }}>
        {ATS_KANBAN_STATUSES.map((status) => (
          <PipelineColumn
            key={status}
            status={status}
            candidates={candidatesByStatus[status]}
            onCardClick={() => {}}
          />
        ))}
      </div>
      <DragOverlay>
        {activeCandidate ? (
          <Card className="w-[272px] shadow-xl border-2 bg-background">
            <CardContent className="p-3 flex gap-3">
              <div className="flex-shrink-0 w-9 h-9 rounded-full bg-primary/15 text-primary flex items-center justify-center text-sm font-semibold">
                {getInitial(activeCandidate)}
              </div>
              <div className="min-w-0 flex-1">
                <p className="font-medium text-sm truncate">{getDisplayName(activeCandidate)}</p>
                {getEmail(activeCandidate) && (
                  <p className="text-xs text-muted-foreground truncate">{getEmail(activeCandidate)}</p>
                )}
                {getJobTitle(activeCandidate) && (
                  <p className="text-xs text-muted-foreground mt-0.5 truncate">{getJobTitle(activeCandidate)}</p>
                )}
              </div>
            </CardContent>
          </Card>
        ) : null}
      </DragOverlay>
    </DndContext>
  )
}
