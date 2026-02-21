import { useOutletContext } from 'react-router-dom'
import { Card, CardContent, CardHeader } from '@/components/ui/card'

export function ApplicantResumePage() {
  const { detail } = useOutletContext()
  if (!detail) return null

  const fileName = detail.resume_file_name

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
          <div className="flex items-start gap-3 py-2">
            <div className="flex-shrink-0 w-10 h-10 rounded-lg bg-muted/60 flex items-center justify-center text-muted-foreground text-sm font-medium" aria-hidden>
              PDF
            </div>
            <div className="min-w-0 flex-1">
              <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">File</span>
              <p className="text-sm font-medium mt-0.5 break-all">{fileName}</p>
            </div>
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
