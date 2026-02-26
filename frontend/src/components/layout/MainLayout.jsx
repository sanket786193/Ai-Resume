import { Outlet, Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { ROUTES, ROUTE_BUILDERS } from '@/constants'
import { useAuth } from '@/modules/auth/hooks/useAuth'

function MainLayout() {
  const { user, logout, isPending } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate(ROUTES.LOGIN)
  }

  if (isPending) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-pulse text-muted-foreground">Loading...</div>
      </div>
    )
  }

  const isHr = user?.role === 'HR'
  const isCandidate = user?.role === 'CANDIDATE'

  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b border-border bg-background">
        <div className="container mx-auto flex h-14 items-center justify-between px-4">
          <Link to={isHr ? ROUTES.HR_DASHBOARD : isCandidate ? ROUTES.CANDIDATE_DASHBOARD : ROUTES.HOME} className="font-semibold">
            HR ATS
          </Link>
          <nav className="flex items-center gap-4">
            {isHr && (
              <>
                <Link to={ROUTES.HR_JOBS} className="text-sm text-muted-foreground hover:text-foreground">Jobs</Link>
                <Link to={ROUTES.HR_ATS} className="text-sm text-muted-foreground hover:text-foreground">Pipeline</Link>
                <Link to={ROUTES.HR_INTERVIEWS} className="text-sm text-muted-foreground hover:text-foreground">Interviews</Link>
                <Link to={ROUTES.HR_OFFERS} className="text-sm text-muted-foreground hover:text-foreground">Offers</Link>
                <Button variant="outline" size="sm" onClick={handleLogout}>Logout</Button>
              </>
            )}
            {isCandidate && (
              <>
                <Link to={ROUTES.PUBLIC_JOBS} className="text-sm text-muted-foreground hover:text-foreground">Jobs</Link>
                <Link to={ROUTES.CANDIDATE_APPLICATIONS} className="text-sm text-muted-foreground hover:text-foreground">My Applications</Link>
                <Button variant="outline" size="sm" onClick={handleLogout}>Logout</Button>
              </>
            )}
            {user && !isHr && !isCandidate ? (
              <Button variant="outline" size="sm" onClick={handleLogout}>Logout</Button>
            ) : null}
            {!user && (
              <>
                <Link to={ROUTES.LOGIN}><Button variant="ghost" size="sm">Login</Button></Link>
                <Link to={ROUTES.REGISTER}><Button variant="outline" size="sm">Register</Button></Link>
                <Link to={ROUTES.REGISTER_HR}><Button size="sm">HR Register</Button></Link>
              </>
            )}
          </nav>
        </div>
      </header>
      <main className="flex-1 container mx-auto px-4 py-6">
        <Outlet />
      </main>
    </div>
  )
}

export default MainLayout
