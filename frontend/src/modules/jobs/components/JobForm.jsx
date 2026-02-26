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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { MarkdownEditor } from '@/components/editor/MarkdownEditor'
import { validateRequired } from '@/utils/validation'

const EXPERIENCE_LEVELS = [
  { value: 'ANY', label: 'Any' },
  { value: 'FRESHER', label: 'Fresher' },
  { value: 'EXPERIENCED', label: 'Experienced' },
]

/**
 * Parse skills from comma-separated or newline-separated string; trim and dedupe.
 */
function parseSkillsText(text) {
  if (!text || typeof text !== 'string') return []
  return [...new Set(text.split(/[\n,]/).map((s) => s.trim()).filter(Boolean))]
}

/**
 * Create/Edit job form in a dialog. Controlled by open/onOpenChange; submit via onSave callback.
 */
export function JobForm({ open, onOpenChange, job, onSave, isSaving }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [department, setDepartment] = useState('')
  const [location, setLocation] = useState('')
  const [experienceLevel, setExperienceLevel] = useState('ANY')
  const [qualification, setQualification] = useState('')
  const [skillsText, setSkillsText] = useState('')
  const [vacancyLimits, setVacancyLimits] = useState([])
  const [errors, setErrors] = useState({})

  useEffect(() => {
    if (open) {
      setTitle(job?.title ?? job?.Title ?? '')
      setDescription(job?.description ?? job?.Description ?? '')
      setDepartment(job?.department ?? job?.Department ?? '')
      setLocation(job?.location ?? job?.Location ?? '')
      setExperienceLevel(job?.experience_level ?? job?.experienceLevel ?? 'ANY')
      setQualification(job?.qualification ?? job?.Qualification ?? '')
      const skills = job?.skills ?? job?.Skills ?? []
      setSkillsText(Array.isArray(skills) ? skills.join(', ') : '')
      const limits = job?.vacancy_limits ?? job?.vacancyLimits ?? job?.VacancyLimits ?? []
      setVacancyLimits(Array.isArray(limits) ? limits.map((l) => ({ role: l.role ?? l.Role ?? '', limit: Number(l.limit ?? l.Limit ?? 0) || 1 })) : [])
      setErrors({})
    }
  }, [open, job])

  const addVacancyRow = () => {
    setVacancyLimits((prev) => [...prev, { role: '', limit: 1 }])
  }

  const updateVacancyRow = (index, field, value) => {
    setVacancyLimits((prev) => {
      const next = [...prev]
      next[index] = { ...next[index], [field]: field === 'limit' ? (Number(value) || 0) : value }
      return next
    })
  }

  const removeVacancyRow = (index) => {
    setVacancyLimits((prev) => prev.filter((_, i) => i !== index))
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    const next = {}
    const titleErr = validateRequired(title.trim(), 'Title')
    if (titleErr) next.title = titleErr
    const descErr = validateRequired(description.trim(), 'Description')
    if (descErr) next.description = descErr
    setErrors(next)
    if (Object.keys(next).length > 0) return
    const skills = parseSkillsText(skillsText)
    const limits = vacancyLimits.filter((l) => (l.role ?? '').trim()).map((l) => ({ role: String(l.role).trim(), limit: Math.max(1, Number(l.limit) || 1) }))
    onSave({
      title: title.trim(),
      description: description.trim(),
      department: department.trim() || undefined,
      location: location.trim() || undefined,
      experience_level: experienceLevel,
      qualification: qualification.trim() || undefined,
      skills,
      vacancy_limits: limits,
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto">
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
              <Label htmlFor="experience_level">Experience level</Label>
              <Select value={experienceLevel} onValueChange={setExperienceLevel}>
                <SelectTrigger id="experience_level">
                  <SelectValue placeholder="Select" />
                </SelectTrigger>
                <SelectContent>
                  {EXPERIENCE_LEVELS.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>{opt.label}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="qualification">Qualification</Label>
              <Input
                id="qualification"
                value={qualification}
                onChange={(e) => setQualification(e.target.value)}
                placeholder="e.g. B.Tech, MBA"
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="skills">Skills (comma-separated)</Label>
              <Input
                id="skills"
                value={skillsText}
                onChange={(e) => setSkillsText(e.target.value)}
                placeholder="e.g. Golang, Python, DevOps"
              />
            </div>
            <div className="grid gap-2">
              <div className="flex items-center justify-between">
                <Label>Vacancy limits (role / count)</Label>
                <Button type="button" variant="outline" size="sm" onClick={addVacancyRow}>Add row</Button>
              </div>
              {vacancyLimits.length === 0 && (
                <p className="text-sm text-muted-foreground">No vacancy limits. Add rows like &quot;Golang — 4&quot;, &quot;Python — 3&quot;.</p>
              )}
              {vacancyLimits.map((row, index) => (
                <div key={index} className="flex gap-2 items-center">
                  <Input
                    value={row.role}
                    onChange={(e) => updateVacancyRow(index, 'role', e.target.value)}
                    placeholder="Role / skill"
                    className="flex-1"
                  />
                  <Input
                    type="number"
                    min={1}
                    value={row.limit}
                    onChange={(e) => updateVacancyRow(index, 'limit', e.target.value)}
                    placeholder="Limit"
                    className="w-20"
                  />
                  <Button type="button" variant="ghost" size="sm" onClick={() => removeVacancyRow(index)} aria-label="Remove row">×</Button>
                </div>
              ))}
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
