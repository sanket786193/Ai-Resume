import { useState } from 'react'
import { PageHeader } from '@/components/common/PageHeader'
import { EmptyState } from '@/components/common/EmptyState'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useHrJobsList, useCreateJob, useUpdateJob, useDeleteJob, usePublishJob, useCloseJob } from '@/modules/jobs/hooks/useJobs'
import { JobCard } from '@/modules/jobs/components/JobCard'
import { JobForm } from '@/modules/jobs/components/JobForm'
import { toast } from 'sonner'
import { ROUTES, ROUTE_BUILDERS } from '@/constants'

export function HrJobsPage() {
  const [formOpen, setFormOpen] = useState(false)
  const [editingJob, setEditingJob] = useState(null)
  const [deleteTarget, setDeleteTarget] = useState(null)
  const { data, isPending, isError } = useHrJobsList()
  const createJob = useCreateJob()
  const updateJob = useUpdateJob(editingJob?.id ?? editingJob?.ID)
  const deleteJob = useDeleteJob()
  const publishJob = usePublishJob()
  const closeJob = useCloseJob()

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

  const handleUpdate = (payload) => {
    const id = editingJob?.id ?? editingJob?.ID
    if (!id) return
    updateJob.mutate(payload, {
      onSuccess: () => {
        toast.success('Job updated')
        setFormOpen(false)
        setEditingJob(null)
      },
      onError: () => toast.error('Failed to update job'),
    })
  }

  const handleSave = (payload) => {
    if (editingJob) {
      handleUpdate(payload)
    } else {
      handleCreate(payload)
    }
  }

  const handleDeleteConfirm = () => {
    const id = deleteTarget?.id ?? deleteTarget?.ID
    if (!id) return
    deleteJob.mutate(id, {
      onSuccess: () => {
        toast.success('Job deleted')
        setDeleteTarget(null)
      },
      onError: () => toast.error('Failed to delete job'),
    })
  }

  const handlePublishJob = (job) => {
    const id = job?.id ?? job?.ID
    if (!id) return
    publishJob.mutate(id, {
      onSuccess: () => toast.success('Job published'),
      onError: () => toast.error('Failed to publish job'),
    })
  }

  const handleCloseJob = (job) => {
    const id = job?.id ?? job?.ID
    if (!id) return
    closeJob.mutate(id, {
      onSuccess: () => toast.success('Job closed'),
      onError: () => toast.error('Failed to close job'),
    })
  }

  const isFormSaving = createJob.isPending || updateJob.isPending

  return (
    <div>
      <PageHeader
        title="Jobs"
        description="Create and manage job postings"
        actions={<Button onClick={() => { setEditingJob(null); setFormOpen(true) }}>Create job</Button>}
      />
      <JobForm
        open={formOpen}
        onOpenChange={(open) => { if (!open) setEditingJob(null); setFormOpen(open) }}
        job={editingJob}
        onSave={handleSave}
        isSaving={isFormSaving}
      />
      <Dialog open={!!deleteTarget} onOpenChange={(open) => !open && setDeleteTarget(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete job?</DialogTitle>
            <DialogDescription>
              This will remove the job. Applications for this job will remain. This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteTarget(null)}>Cancel</Button>
            <Button variant="destructive" onClick={handleDeleteConfirm}>Delete</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
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
          {list.map((job) => {
            const id = job.id ?? job.ID
            const status = job.status ?? job.Status
            return (
              <JobCard
                key={id}
                job={job}
                viewHref={id ? ROUTE_BUILDERS.hrJobApplicants(id) : undefined}
                onEdit={status === 'DRAFT' ? () => { setEditingJob(job); setFormOpen(true) } : undefined}
                onPublish={status === 'DRAFT' ? () => handlePublishJob(job) : undefined}
                onDelete={() => setDeleteTarget(job)}
                onClose={status !== 'DRAFT' && status !== 'CLOSED' ? () => handleCloseJob(job) : undefined}
              />
            )
          })}
        </div>
      )}
    </div>
  )
}
