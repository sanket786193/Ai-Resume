import { useMemo } from 'react'
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
import { PipelineCandidateCard } from './PipelineCandidateCard'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { EmptyState } from '@/components/common/EmptyState'
import { useState } from 'react'

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
    const list = data?.list ?? data ?? []
    const map = {}
    ATS_KANBAN_STATUSES.forEach((s) => { map[s] = [] })
    list.forEach((c) => {
      const status = c.status ?? 'APPLIED'
      if (!map[status]) map[status] = []
      map[status].push(c)
    })
    return map
  }, [data])

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
          <Skeleton key={s} className="h-[400px] w-[280px] flex-shrink-0 rounded-lg" />
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

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <div className="flex gap-4 overflow-x-auto pb-4">
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
          <Card className="w-[260px] shadow-lg">
            <CardContent className="p-3">
              <p className="font-medium text-sm">{activeCandidate.name ?? activeCandidate.email}</p>
              <p className="text-xs text-muted-foreground">{activeCandidate.email}</p>
            </CardContent>
          </Card>
        ) : null}
      </DragOverlay>
    </DndContext>
  )
}
