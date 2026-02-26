import { Link } from 'react-router-dom'
import { PageHeader } from '@/components/common/PageHeader'
import { EmptyState } from '@/components/common/EmptyState'
import { Skeleton } from '@/components/ui/skeleton'
import { useJobsList } from '@/modules/jobs/hooks/useJobs'
import { JobCard } from '@/modules/jobs/components/JobCard'

export function PublicJobsPage() {
  const { data, isPending, isError } = useJobsList({ status: 'PUBLISHED' })
  const list = Array.isArray(data) ? data : data?.list ?? []

  return (
    <div>
      <PageHeader title="Open positions" description="Browse and apply to open roles" />
      {isPending && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-[180px] rounded-lg" />
          ))}
        </div>
      )}
      {!isPending && isError && (
        <EmptyState title="Failed to load jobs" description="Please try again later." />
      )}
      {!isPending && !isError && list.length === 0 && (
        <EmptyState title="No open positions" description="Check back later for new roles." />
      )}
      {!isPending && !isError && list.length > 0 && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {list.map((job) => (
            <JobCard key={job.id} job={job} />
          ))}
        </div>
      )}
    </div>
  )
}
