import { Outlet } from 'react-router-dom'

/**
 * Minimal layout for login/register. Centered card.
 */
function AuthLayout() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-muted/30 p-4">
      <Outlet />
    </div>
  )
}

export default AuthLayout
