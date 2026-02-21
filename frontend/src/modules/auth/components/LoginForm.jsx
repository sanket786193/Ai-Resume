import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { ROUTES } from '@/constants'
import { validateEmail, validateRequired } from '@/utils/validation'
import { useAuth } from '../hooks/useAuth'

/**
 * Login form. Delegates submit to useAuth; no API calls in component.
 */
export function LoginForm() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [errors, setErrors] = useState({})
  const navigate = useNavigate()
  const { login, loginMutation } = useAuth()

  const handleSubmit = async (e) => {
    e.preventDefault()
    const next = {}
    const emailErr = validateEmail(email)
    if (emailErr) next.email = emailErr
    const pwdErr = validateRequired(password, 'Password')
    if (pwdErr) next.password = pwdErr
    setErrors(next)
    if (Object.keys(next).length > 0) return

    try {
      const data = await login({ email, password })
      const role = data?.user?.role ?? data?.role
      navigate(role === 'HR' ? ROUTES.HR_DASHBOARD : ROUTES.CANDIDATE_DASHBOARD)
    } catch (err) {
      loginMutation.reset()
    }
  }

  const isSubmitting = loginMutation.isPending
  const error = loginMutation.error

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Login</CardTitle>
        <CardDescription>Sign in to your account</CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit} noValidate>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="you@company.com"
              value={email}
              onChange={(e) => { setEmail(e.target.value); setErrors((prev) => ({ ...prev, email: null })) }}
              autoComplete="email"
              aria-invalid={!!errors.email}
              aria-describedby={errors.email ? 'email-error' : undefined}
            />
            {errors.email && <p id="email-error" className="text-sm text-destructive" role="alert">{errors.email}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => { setPassword(e.target.value); setErrors((prev) => ({ ...prev, password: null })) }}
              autoComplete="current-password"
              aria-invalid={!!errors.password}
              aria-describedby={errors.password ? 'password-error' : undefined}
            />
            {errors.password && <p id="password-error" className="text-sm text-destructive" role="alert">{errors.password}</p>}
          </div>
          {error && (
            <p className="text-sm text-destructive" role="alert">
              {error?.response?.data?.message ?? error?.message ?? 'Login failed'}
            </p>
          )}
        </CardContent>
        <CardFooter className="flex flex-col gap-2">
          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? 'Signing in...' : 'Sign in'}
          </Button>
          <Link to={ROUTES.REGISTER} className="text-sm text-primary underline-offset-4 hover:underline w-full text-center">
            Create account
          </Link>
        </CardFooter>
      </form>
    </Card>
  )
}
