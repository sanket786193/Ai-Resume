import { useState, useCallback } from 'react'
import { useOutletContext } from 'react-router-dom'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { IconUser, IconMail, IconPhone, IconLinkedin, IconCopy, IconCheck, IconExternalLink } from '@/components/icons/ApplicantIcons'

function CopyButton({ value, label }) {
  const [copied, setCopied] = useState(false)
  const copy = useCallback(() => {
    if (!value) return
    navigator.clipboard.writeText(value).then(
      () => {
        setCopied(true)
        toast.success('Copied to clipboard')
        setTimeout(() => setCopied(false), 2000)
      },
      () => toast.error('Could not copy')
    )
  }, [value])
  if (!value) return null
  return (
    <Button
      type="button"
      variant="ghost"
      size="icon"
      className="h-8 w-8 shrink-0 text-muted-foreground hover:text-foreground"
      onClick={copy}
      aria-label={`Copy ${label}`}
    >
      {copied ? <IconCheck className="w-4 h-4" /> : <IconCopy className="w-4 h-4" />}
    </Button>
  )
}

const ContactRow = ({ icon: Icon, label, value, href, mailto, tel, copyable }) => {
  if (value == null || value === '') return null
  const linkHref = mailto ? `mailto:${value}` : tel ? `tel:${value}` : href
  const isLink = Boolean(linkHref)
  return (
    <div className="flex items-start gap-3 py-3 first:pt-0 last:pb-0 border-b border-border/50 last:border-0 group">
      <div className="flex-shrink-0 w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center text-primary" aria-hidden>
        <Icon className="w-5 h-5" />
      </div>
      <div className="min-w-0 flex-1">
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{label}</span>
        <div className="flex items-center gap-2 mt-1">
          {isLink ? (
            <a
              href={linkHref}
              target={href ? '_blank' : undefined}
              rel={href ? 'noopener noreferrer' : undefined}
              className="text-sm font-medium text-primary hover:underline break-all inline-flex items-center gap-1"
            >
              {value}
              {href && <IconExternalLink className="w-3.5 h-3.5 shrink-0 opacity-70" />}
            </a>
          ) : (
            <p className="text-sm font-medium break-words">{value}</p>
          )}
          {copyable && <CopyButton value={value} label={label} />}
        </div>
      </div>
    </div>
  )
}

export function ApplicantContactPage() {
  const { detail } = useOutletContext()
  if (!detail) return null

  const hasAny =
    detail.candidate_name ||
    detail.candidate_email ||
    detail.candidate_phone ||
    detail.candidate_linkedin

  const socialLinks = []
  if (detail.candidate_linkedin?.trim()) {
    socialLinks.push({ label: 'LinkedIn', value: detail.candidate_linkedin.trim(), icon: IconLinkedin })
  }

  return (
    <Card className="shadow-sm border-border/80">
      <CardHeader className="pb-2">
        <h2 className="text-base font-semibold">Contact</h2>
        <p className="text-sm text-muted-foreground">
          How to reach this candidate
        </p>
      </CardHeader>
      <CardContent className="space-y-1">
        {hasAny ? (
          <div className="space-y-0">
            <ContactRow icon={IconUser} label="Name" value={detail.candidate_name} />
            <ContactRow
              icon={IconMail}
              label="Email"
              value={detail.candidate_email}
              mailto
              copyable
            />
            <ContactRow
              icon={IconPhone}
              label="Phone"
              value={detail.candidate_phone}
              tel
              copyable
            />
            {socialLinks.length > 0 && (
              <>
                <div className="pt-2 pb-1">
                  <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Social media</span>
                </div>
                {socialLinks.map(({ label, value, icon: Icon }) => (
                  <ContactRow
                    key={label}
                    icon={Icon}
                    label={label}
                    value={value}
                    href={value.startsWith('http') ? value : `https://${value}`}
                    copyable
                  />
                ))}
              </>
            )}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-6">
            No contact information yet. It may appear here after the resume is parsed.
          </p>
        )}
      </CardContent>
    </Card>
  )
}
