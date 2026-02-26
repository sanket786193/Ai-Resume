import { useMemo } from 'react'
import { useOutletContext } from 'react-router-dom'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { config } from '@/config/env'

export function ApplicantResumePage() {
  const { detail } = useOutletContext()
  if (!detail) return null

  const fileName = detail.resume_file_name
  const resumeUrl = detail.resume_url
  const applicationId = detail.id ?? detail.ID

  // Proxy URL (same-origin + token) so the iframe can load the PDF; direct Cloudinary URL often fails in iframes (embed/CORS)
  const iframeSrc = useMemo(() => {
    if (!applicationId || !resumeUrl) return null
    try {
      const token = typeof window !== 'undefined' ? window.localStorage?.getItem('ats_token') : null
      if (!token) return null
      const base = config.apiBaseUrl.replace(/\/$/, '')
      return `${base}/api/hr/applications/${applicationId}/resume?access_token=${encodeURIComponent(token)}`
    } catch {
      return null
    }
  }, [applicationId, resumeUrl])

  return (
    <Card>
      <CardHeader className="pb-2">
        <h2 className="text-base font-semibold">Resume</h2>
        <p className="text-sm text-muted-foreground">
          Uploaded resume for this application
        </p>
      </CardHeader>
      <CardContent>
        {fileName ? (
          <div className="space-y-4">
            <div className="flex items-start gap-3 py-2">
              <div className="flex-shrink-0 w-10 h-10 rounded-lg bg-muted/60 flex items-center justify-center text-muted-foreground text-sm font-medium" aria-hidden>
                PDF
              </div>
              <div className="min-w-0 flex-1">
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">File</span>
                <p className="text-sm font-medium mt-0.5 break-all">{fileName}</p>
              </div>
            </div>
            {(iframeSrc || resumeUrl) && (
              <div className="space-y-2">
                {resumeUrl && (
                  <a
                    href={resumeUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-primary hover:underline"
                  >
                    Open in new tab
                  </a>
                )}
                {iframeSrc && (
                  <div className="rounded-lg border border-border overflow-hidden bg-muted/30">
                    <iframe
                      src={iframeSrc}
                      title="Resume PDF"
                      className="w-full min-h-[70vh] h-[800px] border-0"
                    />
                  </div>
                )}
              </div>
            )}
          </div>
        ) : (
          <p className="text-sm text-muted-foreground py-4">
            No resume file linked to this application.
          </p>
        )}
      </CardContent>
    </Card>
  )
}
