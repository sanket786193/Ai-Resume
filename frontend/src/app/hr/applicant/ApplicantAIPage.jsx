import { useOutletContext } from 'react-router-dom'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { isAIScreeningPending } from '@/modules/ats/hooks/useApplicationDetail'

export function ApplicantAIPage() {
  const { detail } = useOutletContext()
  if (!detail) return null

  const pending = isAIScreeningPending(detail)
  const hasScores =
    detail.ats_score != null ||
    detail.skill_match_score != null ||
    detail.skill_match_pct != null
  const hasMatchView =
    (Array.isArray(detail.keyword_matches) && detail.keyword_matches.length > 0) ||
    (Array.isArray(detail.semantic_matches) && detail.semantic_matches.length > 0)
  const hasAny =
    hasScores ||
    detail.qualified != null ||
    (detail.ai_summary && detail.ai_summary !== '') ||
    (detail.experience_match && detail.experience_match !== '') ||
    (Array.isArray(detail.missing_skills) && detail.missing_skills.length > 0) ||
    (Array.isArray(detail.experience_warnings) && detail.experience_warnings.length > 0) ||
    hasMatchView

  if (pending) {
    return (
      <Card>
        <CardHeader className="pb-2">
          <h2 className="text-base font-semibold">AI screening</h2>
          <p className="text-sm text-muted-foreground">
            Analyzing resume against job description…
          </p>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-3 text-muted-foreground">
            <Skeleton className="h-5 w-5 rounded-full animate-pulse" />
            <span className="text-sm">AI screening in progress. Results will appear here shortly.</span>
          </div>
          <Skeleton className="h-20 w-full rounded-lg" />
          <Skeleton className="h-24 w-full rounded-lg" />
        </CardContent>
      </Card>
    )
  }

  if (!hasAny) {
    return (
      <Card>
        <CardHeader className="pb-2">
          <h2 className="text-base font-semibold">AI feedback</h2>
          <p className="text-sm text-muted-foreground">
            Scores and insights from resume screening vs job description
          </p>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground py-4">
            No AI feedback for this application. Screening may not have run (e.g. no resume or AI disabled).
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader className="pb-2">
        <h2 className="text-base font-semibold">AI feedback</h2>
        <p className="text-sm text-muted-foreground">
          Scores and insights from resume screening vs job description
        </p>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* ATS score (0–100) & Skill match % */}
        {hasScores && (
          <div className="flex flex-wrap gap-3">
            {detail.ats_score != null && (
              <div className="rounded-lg border border-border/50 bg-muted/30 px-4 py-2 min-w-[7rem]">
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">ATS score</span>
                <p className="text-lg font-semibold mt-0.5">{detail.ats_score}/100</p>
              </div>
            )}
            {detail.skill_match_pct != null && (
              <div className="rounded-lg border border-border/50 bg-muted/30 px-4 py-2 min-w-[7rem]">
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Skill match %</span>
                <p className="text-lg font-semibold mt-0.5">{detail.skill_match_pct}%</p>
              </div>
            )}
            {detail.skill_match_score != null && detail.skill_match_pct == null && (
              <div className="rounded-lg border border-border/50 bg-muted/30 px-4 py-2 min-w-[7rem]">
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Match score</span>
                <p className="text-lg font-semibold mt-0.5">{Math.round(Number(detail.skill_match_score) * 100)}/100</p>
              </div>
            )}
          </div>
        )}

        {detail.qualified != null && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Qualified</span>
            <div className="mt-1">
              <Badge variant={detail.qualified ? 'default' : 'secondary'}>
                {detail.qualified ? 'Yes' : 'No'}
              </Badge>
            </div>
          </div>
        )}

        {/* Candidate resume summary (detailed) */}
        {detail.ai_summary && detail.ai_summary !== '' && (
          <div className="rounded-lg border border-border bg-muted/20 p-4">
            <h3 className="text-sm font-semibold text-foreground uppercase tracking-wide mb-2">Candidate resume summary</h3>
            <p className="text-base leading-relaxed whitespace-pre-wrap text-foreground">
              {detail.ai_summary}
            </p>
          </div>
        )}

        {/* Missing skills list */}
        {Array.isArray(detail.missing_skills) && detail.missing_skills.length > 0 && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Missing skills</span>
            <ul className="text-sm mt-1 list-disc list-inside space-y-0.5">
              {detail.missing_skills.map((s, i) => (
                <li key={i}>{s}</li>
              ))}
            </ul>
          </div>
        )}

        {/* Experience mismatch: summary + warnings */}
        {(detail.experience_match?.trim() || (Array.isArray(detail.experience_warnings) && detail.experience_warnings.length > 0)) && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Experience vs job</span>
            {detail.experience_match?.trim() && (
              <p className="text-sm mt-1">{detail.experience_match}</p>
            )}
            {Array.isArray(detail.experience_warnings) && detail.experience_warnings.length > 0 && (
              <ul className="text-sm mt-2 space-y-1 rounded-md bg-destructive/10 border border-destructive/20 p-3 list-disc list-inside">
                {detail.experience_warnings.map((w, i) => (
                  <li key={i} className="text-destructive/90">{w}</li>
                ))}
              </ul>
            )}
          </div>
        )}

        {/* Keyword vs semantic match view */}
        {hasMatchView && (
          <div>
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Keyword vs semantic match</span>
            <p className="text-xs text-muted-foreground mt-0.5 mb-2">
              Job requirements found in resume: exact wording vs meaning-based match.
            </p>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-2">
              <div className="rounded-lg border border-border/50 bg-muted/20 p-3">
                <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wide mb-2">Keyword match</h3>
                {Array.isArray(detail.keyword_matches) && detail.keyword_matches.length > 0 ? (
                  <div className="flex flex-wrap gap-1.5">
                    {detail.keyword_matches.map((s, i) => (
                      <Badge key={i} variant="outline" className="font-normal">{s}</Badge>
                    ))}
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground">None</p>
                )}
              </div>
              <div className="rounded-lg border border-border/50 bg-muted/20 p-3">
                <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wide mb-2">Semantic match</h3>
                {Array.isArray(detail.semantic_matches) && detail.semantic_matches.length > 0 ? (
                  <div className="flex flex-wrap gap-1.5">
                    {detail.semantic_matches.map((s, i) => (
                      <Badge key={i} variant="secondary" className="font-normal">{s}</Badge>
                    ))}
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground">None</p>
                )}
              </div>
            </div>
          </div>
        )}

      </CardContent>
    </Card>
  )
}
