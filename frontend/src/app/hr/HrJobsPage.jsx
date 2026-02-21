import { useState } from 'react'
import { Link } from 'react-router-dom'
import { PageHeader } from '@/components/common/PageHeader'
import { EmptyState } from '@/components/common/EmptyState'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { useJobsList, useCreateJob } from '@/modules/jobs/hooks/useJobs'
import { JobCard } from '@/modules/jobs/components/JobCard'
import { JobForm } from '@/modules/jobs/components/JobForm'
import { toast } from 'sonner'
import { ROUTES, ROUTE_BUILDERS } from '@/constants'

export function HrJobsPage() {
  const [formOpen, setFormOpen] = useState(false)
  const { data, isPending, isError } = useJobsList()
  const createJob = useCreateJob()

  const list = Array.isArray(data) ? data : data?.list ?? []

  const handleCreate = (payload) => {
    createJob.mutate(payload, {
      onSuccess: () => {
        toast.success('Job created')
        setFormOpen(false)
      },
      onError: () => toast.error('Failed to create job'),
    })
  }

  return (
    <div>
      <PageHeader
        title="Jobs"
        description="Create and manage job postings"
        actions={<Button onClick={() => setFormOpen(true)}>Create job</Button>}
      />
      <JobForm
        open={formOpen}
        onOpenChange={setFormOpen}
        onSave={handleCreate}
        isSaving={createJob.isPending}
      />
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
        <EmptyState
          title="No jobs yet"
          description="Create your first job posting."
          actionLabel="Create job"
          onAction={() => setFormOpen(true)}
        />
      )}
      {!isPending && !isError && list.length > 0 && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {list.map((job) => (
            <JobCard key={job.id} job={job} viewHref={ROUTE_BUILDERS.hrJobApplicants(job.id)} />
          ))}
        </div>
      )}
    </div>
  )
}
