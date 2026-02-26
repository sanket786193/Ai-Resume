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
 * HR registration form. Delegates submit to useAuth.registerHR.
 */
export function RegisterHRForm() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [errors, setErrors] = useState({})
  const navigate = useNavigate()
  const { registerHR, registerHRMutation } = useAuth()

  const handleSubmit = async (e) => {
    e.preventDefault()
    const next = {}
    const nameErr = validateRequired(name, 'Name')
    if (nameErr) next.name = nameErr
    const emailErr = validateEmail(email)
    if (emailErr) next.email = emailErr
    const pwdErr = validateRequired(password, 'Password')
    if (pwdErr) next.password = pwdErr
    setErrors(next)
    if (Object.keys(next).length > 0) return

    try {
      await registerHR({ email, password, name })
      navigate(ROUTES.HR_DASHBOARD)
    } catch (err) {
      registerHRMutation.reset()
    }
  }

  const isSubmitting = registerHRMutation.isPending
  const error = registerHRMutation.error

  const clearError = (field) => setErrors((prev) => ({ ...prev, [field]: null }))

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Create HR account</CardTitle>
        <CardDescription>Register as a recruiter or hiring manager</CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit} noValidate>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="hr-name">Name</Label>
            <Input
              id="hr-name"
              type="text"
              placeholder="Your name"
              value={name}
              onChange={(e) => { setName(e.target.value); clearError('name') }}
              aria-invalid={!!errors.name}
              aria-describedby={errors.name ? 'hr-name-error' : undefined}
            />
            {errors.name && <p id="hr-name-error" className="text-sm text-destructive" role="alert">{errors.name}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="hr-email">Email</Label>
            <Input
              id="hr-email"
              type="email"
              placeholder="you@company.com"
              value={email}
              onChange={(e) => { setEmail(e.target.value); clearError('email') }}
              autoComplete="email"
              aria-invalid={!!errors.email}
              aria-describedby={errors.email ? 'hr-email-error' : undefined}
            />
            {errors.email && <p id="hr-email-error" className="text-sm text-destructive" role="alert">{errors.email}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="hr-password">Password</Label>
            <Input
              id="hr-password"
              type="password"
              value={password}
              onChange={(e) => { setPassword(e.target.value); clearError('password') }}
              autoComplete="new-password"
              placeholder="Min 8 characters"
              aria-invalid={!!errors.password}
              aria-describedby={errors.password ? 'hr-password-error' : undefined}
            />
            {errors.password && <p id="hr-password-error" className="text-sm text-destructive" role="alert">{errors.password}</p>}
          </div>
          {error && (
            <p className="text-sm text-destructive" role="alert">
              {error?.response?.data?.message ?? error?.message ?? 'Registration failed'}
            </p>
          )}
        </CardContent>
        <CardFooter className="flex flex-col gap-2">
          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? 'Creating account...' : 'Create HR account'}
          </Button>
          <Link to={ROUTES.LOGIN} className="text-sm text-primary underline-offset-4 hover:underline w-full text-center">
            Already have an account? Sign in
          </Link>
          <Link to={ROUTES.REGISTER} className="text-sm text-muted-foreground hover:underline w-full text-center">
            Register as candidate instead
          </Link>
        </CardFooter>
      </form>
    </Card>
  )
}
