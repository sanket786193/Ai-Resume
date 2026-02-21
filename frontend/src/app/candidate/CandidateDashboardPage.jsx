import { Link } from 'react-router-dom'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { ROUTES } from '@/constants'

export function CandidateDashboardPage() {
  return (
    <div>
      <h1 className="text-2xl font-semibold mb-6">Candidate Dashboard</h1>
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Browse jobs</CardTitle>
            <CardDescription>Find and apply to open positions</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild><Link to={ROUTES.PUBLIC_JOBS}>View jobs</Link></Button>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>My applications</CardTitle>
            <CardDescription>Track your application status</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild variant="outline"><Link to={ROUTES.CANDIDATE_APPLICATIONS}>View applications</Link></Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
