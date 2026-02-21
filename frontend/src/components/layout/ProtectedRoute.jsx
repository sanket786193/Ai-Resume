import { Navigate, useLocation } from 'react-router-dom'
import { ROUTES, ROLES } from '@/constants'
import { useAuth } from '@/modules/auth/hooks/useAuth'

/**
 * Role-based protected route. Redirects to login or role dashboard if unauthenticated/wrong role.
 */
export function ProtectedRoute({ children, allowedRoles }) {
  const { user, isPending } = useAuth()
  const location = useLocation()

  if (isPending) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-pulse text-muted-foreground">Loading...</div>
      </div>
    )
  }

  if (!user) {
    return <Navigate to={ROUTES.LOGIN} state={{ from: location }} replace />
  }

  if (allowedRoles && allowedRoles.length > 0 && !allowedRoles.includes(user.role)) {
    const redirect = user.role === ROLES.HR ? ROUTES.HR_DASHBOARD : ROUTES.CANDIDATE_DASHBOARD
    return <Navigate to={redirect} replace />
  }

  return children
}
