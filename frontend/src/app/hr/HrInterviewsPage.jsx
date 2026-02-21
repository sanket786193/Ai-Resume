import { PageHeader } from '@/components/common/PageHeader'
import { EmptyState } from '@/components/common/EmptyState'
import { useInterviewsList } from '@/modules/interviews/hooks/useInterviews'

export function HrInterviewsPage() {
  const { data, isPending, isError } = useInterviewsList()
  const list = Array.isArray(data) ? data : data?.list ?? []

  return (
    <div>
      <PageHeader title="Interviews" description="Schedule and manage interviews" />
      {isPending && <div className="text-muted-foreground">Loading...</div>}
      {!isPending && (isError || list.length === 0) && (
        <EmptyState
          title={isError ? 'Failed to load' : 'No interviews scheduled'}
          description={isError ? 'Please try again later.' : 'Schedule an interview from the pipeline.'}
        />
      )}
      {!isPending && !isError && list.length > 0 && (
        <ul className="space-y-2">
          {list.map((i) => (
            <li key={i.id} className="border rounded-lg p-4">
              <p className="font-medium">{i.candidateName ?? i.title}</p>
              <p className="text-sm text-muted-foreground">{i.scheduledAt ?? i.date}</p>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
