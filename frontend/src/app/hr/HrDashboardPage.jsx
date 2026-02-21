import { Link } from 'react-router-dom'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { ROUTES } from '@/constants'

export function HrDashboardPage() {
  return (
    <div>
      <h1 className="text-2xl font-semibold mb-6">HR Dashboard</h1>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle>Jobs</CardTitle>
            <CardDescription>Create and manage job postings</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild><Link to={ROUTES.HR_JOBS}>View jobs</Link></Button>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Pipeline</CardTitle>
            <CardDescription>ATS Kanban board</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild><Link to={ROUTES.HR_ATS}>Open pipeline</Link></Button>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Interviews</CardTitle>
            <CardDescription>Schedule and manage interviews</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild><Link to={ROUTES.HR_INTERVIEWS}>View interviews</Link></Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
