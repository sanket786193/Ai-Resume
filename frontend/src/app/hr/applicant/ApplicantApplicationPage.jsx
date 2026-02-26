import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { useOutletContext } from 'react-router-dom'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
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
import { useUpdateCandidateStatus } from '@/modules/ats/hooks/useAtsPipeline'
import { APPLICATION_DETAIL_QUERY_KEY } from '@/modules/ats/hooks/useApplicationDetail'
import { useCreateInterview } from '@/modules/interviews/hooks/useInterviews'
import { useCreateOffer } from '@/modules/offers/hooks/useOffers'
import { ATS_STATUS_LABEL } from '@/constants'
import { JOB_APPLICANTS_QUERY_KEY } from '@/modules/jobs/hooks/useJobs'

const DetailRow = ({ label, value }) => {
  if (value == null || value === '') return null
  return (
    <div className="flex items-start gap-3 py-3 border-b border-border/50 last:border-0">
      <div className="flex-shrink-0 w-9 h-9 rounded-lg bg-muted/60 flex items-center justify-center text-muted-foreground text-xs font-medium" aria-hidden>
        {label.charAt(0)}
      </div>
      <div className="min-w-0 flex-1">
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{label}</span>
        <p className="text-sm font-medium mt-0.5 break-words">{value}</p>
      </div>
    </div>
  )
}

export function ApplicantApplicationPage() {
  const { detail, jobId } = useOutletContext()
  const queryClient = useQueryClient()
  const updateStatus = useUpdateCandidateStatus()
  const createInterview = useCreateInterview()
  const createOffer = useCreateOffer()
  const applicationId = detail?.id ?? detail?.ID
  const atsId = applicationId

  const [scheduleOpen, setScheduleOpen] = useState(false)
  const [scheduleAt, setScheduleAt] = useState('')
  const [duration, setDuration] = useState(60)
  const [location, setLocation] = useState('')
  const [notes, setNotes] = useState('')
  const [offerOpen, setOfferOpen] = useState(false)
  const [amount, setAmount] = useState('')
  const [currency, setCurrency] = useState('INR')
  const [startsAt, setStartsAt] = useState('')

  const handleStatusChange = (newStatus) => {
    if (!applicationId || !newStatus) return
    updateStatus.mutate(
      { id: applicationId, status: newStatus },
      {
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: APPLICATION_DETAIL_QUERY_KEY(applicationId) })
          if (jobId) queryClient.invalidateQueries({ queryKey: JOB_APPLICANTS_QUERY_KEY(jobId) })
        },
      }
    )
  }

  const invalidateDetail = () => {
    queryClient.invalidateQueries({ queryKey: APPLICATION_DETAIL_QUERY_KEY(applicationId) })
    if (jobId) queryClient.invalidateQueries({ queryKey: JOB_APPLICANTS_QUERY_KEY(jobId) })
  }

  const handleScheduleInterview = () => {
    if (!atsId || !scheduleAt) {
      toast.error('Select date and time')
      return
    }
    const d = new Date(scheduleAt)
    if (isNaN(d.getTime())) {
      toast.error('Invalid date/time')
      return
    }
    createInterview.mutate(
      {
        ats_id: atsId,
        scheduled_at: d.toISOString(),
        duration_minutes: duration || 60,
        location: location.trim() || undefined,
        round: 1,
        notes: notes.trim() || undefined,
      },
      {
        onSuccess: () => {
          toast.success('Interview scheduled')
          setScheduleOpen(false)
          setScheduleAt('')
          setDuration(60)
          setLocation('')
          setNotes('')
          invalidateDetail()
        },
        onError: () => toast.error('Failed to schedule interview'),
      }
    )
  }

  const handleIssueOffer = () => {
    if (!atsId) return
    createOffer.mutate(
      {
        ats_id: atsId,
        amount: amount.trim() || undefined,
        currency: currency.trim() || 'INR',
        starts_at: startsAt.trim() ? new Date(startsAt).toISOString() : undefined,
      },
      {
        onSuccess: () => {
          toast.success('Offer created')
          setOfferOpen(false)
          setAmount('')
          setCurrency('INR')
          setStartsAt('')
          invalidateDetail()
        },
        onError: () => toast.error('Failed to create offer'),
      }
    )
  }

  if (!detail) return null

  const appliedAt = detail.created_at
    ? new Date(detail.created_at).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })
    : null
  const status = detail.status ?? detail.Status
  const hasOffer = !!(detail.offer_id ?? detail.offerId)
  const isSelectedByHR = status === 'SHORTLISTED' || status === 'INTERVIEW'

  return (
    <>
    <Card>
      <CardHeader className="pb-2">
        <h2 className="text-base font-semibold">Application</h2>
        <p className="text-sm text-muted-foreground">
          Status and job for this application
        </p>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Status</Label>
          <Select
            value={detail.status ?? ''}
            onValueChange={handleStatusChange}
            disabled={updateStatus.isPending}
          >
            <SelectTrigger className="w-full max-w-xs">
              <SelectValue placeholder="Select status" />
            </SelectTrigger>
            <SelectContent>
              {Object.entries(ATS_STATUS_LABEL).map(([value, label]) => (
                <SelectItem key={value} value={value}>
                  {label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-0 pt-2 border-t border-border/50">
          <DetailRow label="Applied" value={appliedAt} />
          <DetailRow label="Job" value={detail.job_title} />
        </div>
        <div className="pt-4 border-t border-border/50 space-y-3">
          {hasOffer ? (
            <div className="rounded-lg bg-muted/40 p-3 text-sm">
              <p className="font-medium text-foreground">Offer (candidate selected)</p>
              <p className="text-muted-foreground mt-1">
                {(detail.offer_amount ?? detail.offerAmount) && (detail.offer_currency ?? detail.offerCurrency)
                  ? `${detail.offer_amount ?? detail.offerAmount} ${detail.offer_currency ?? detail.offerCurrency}`
                  : 'Offer issued'}
                {' — '}
                <span className="capitalize">{detail.offer_status ?? detail.offerStatus ?? '—'}</span>
              </p>
            </div>
          ) : isSelectedByHR ? (
            <Button variant="outline" size="sm" onClick={() => setOfferOpen(true)}>
              Issue offer
            </Button>
          ) : null}
          <Button variant="outline" size="sm" onClick={() => setScheduleOpen(true)}>
            Schedule interview
          </Button>
        </div>
      </CardContent>
    </Card>

    <Dialog open={scheduleOpen} onOpenChange={setScheduleOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Schedule interview</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="scheduled_at">Date & time</Label>
            <Input
              id="scheduled_at"
              type="datetime-local"
              value={scheduleAt}
              onChange={(e) => setScheduleAt(e.target.value)}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="duration">Duration (minutes)</Label>
            <Input
              id="duration"
              type="number"
              min={15}
              value={duration}
              onChange={(e) => setDuration(Number(e.target.value) || 60)}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="location">Location / meeting link</Label>
            <Input id="location" value={location} onChange={(e) => setLocation(e.target.value)} placeholder="e.g. Zoom link" />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="notes">Notes</Label>
            <Input id="notes" value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Optional" />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setScheduleOpen(false)}>Cancel</Button>
          <Button onClick={handleScheduleInterview} disabled={createInterview.isPending}>
            {createInterview.isPending ? 'Scheduling…' : 'Schedule'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog open={offerOpen} onOpenChange={setOfferOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Issue offer</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="amount">Amount</Label>
            <Input id="amount" value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="e.g. 15 LPA" />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="currency">Currency</Label>
            <Input id="currency" value={currency} onChange={(e) => setCurrency(e.target.value)} placeholder="INR" />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="starts_at">Start date (optional)</Label>
            <Input id="starts_at" type="date" value={startsAt} onChange={(e) => setStartsAt(e.target.value)} />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setOfferOpen(false)}>Cancel</Button>
          <Button onClick={handleIssueOffer} disabled={createOffer.isPending}>
            {createOffer.isPending ? 'Creating…' : 'Create offer'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </>
  )
}
