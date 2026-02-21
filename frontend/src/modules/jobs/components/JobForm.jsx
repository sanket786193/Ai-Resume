import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { MarkdownEditor } from '@/components/editor/MarkdownEditor'
import { validateRequired } from '@/utils/validation'

/**
 * Create/Edit job form in a dialog. Controlled by open/onOpenChange; submit via onSave callback.
 */
export function JobForm({ open, onOpenChange, job, onSave, isSaving }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [department, setDepartment] = useState('')
  const [location, setLocation] = useState('')
  const [errors, setErrors] = useState({})

  useEffect(() => {
    if (open) {
      setTitle(job?.title ?? job?.Title ?? '')
      setDescription(job?.description ?? job?.Description ?? '')
      setDepartment(job?.department ?? job?.Department ?? '')
      setLocation(job?.location ?? job?.Location ?? '')
      setErrors({})
    }
  }, [open, job])

  const handleSubmit = (e) => {
    e.preventDefault()
    const next = {}
    const titleErr = validateRequired(title.trim(), 'Title')
    if (titleErr) next.title = titleErr
    const descErr = validateRequired(description.trim(), 'Description')
    if (descErr) next.description = descErr
    setErrors(next)
    if (Object.keys(next).length > 0) return
    onSave({ title: title.trim(), description: description.trim(), department: department.trim() || undefined, location: location.trim() || undefined })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={handleSubmit} noValidate>
          <DialogHeader>
            <DialogTitle>{(job?.id ?? job?.ID) ? 'Edit job' : 'Create job'}</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="title">Title</Label>
              <Input
                id="title"
                value={title}
                onChange={(e) => { setTitle(e.target.value); setErrors((p) => ({ ...p, title: null })) }}
                placeholder="Job title"
                aria-invalid={!!errors.title}
                aria-describedby={errors.title ? 'title-error' : undefined}
              />
              {errors.title && <p id="title-error" className="text-sm text-destructive" role="alert">{errors.title}</p>}
            </div>
            <div className="grid gap-2">
              <Label htmlFor="department">Department</Label>
              <Input
                id="department"
                value={department}
                onChange={(e) => setDepartment(e.target.value)}
                placeholder="Department"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="location">Location</Label>
              <Input
                id="location"
                value={location}
                onChange={(e) => setLocation(e.target.value)}
                placeholder="e.g. Remote, New York"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="description">Description</Label>
              <MarkdownEditor
                id="description"
                value={description}
                onChange={(v) => { setDescription(v); setErrors((p) => ({ ...p, description: null })) }}
                placeholder="Job description (supports **bold**, _italic_, lists, links…)"
                minHeight={180}
                aria-invalid={!!errors.description}
                aria-describedby={errors.description ? 'description-error' : undefined}
              />
              {errors.description && <p id="description-error" className="text-sm text-destructive" role="alert">{errors.description}</p>}
            </div>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSaving}>{isSaving ? 'Saving...' : 'Save'}</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
