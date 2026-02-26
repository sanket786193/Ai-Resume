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
 * Registration form. Delegates submit to useAuth.
 */
export function RegisterForm() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [errors, setErrors] = useState({})
  const navigate = useNavigate()
  const { register, registerMutation } = useAuth()

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
      await register({ email, password, name })
      navigate(ROUTES.CANDIDATE_DASHBOARD)
    } catch (err) {
      registerMutation.reset()
    }
  }

  const isSubmitting = registerMutation.isPending
  const error = registerMutation.error

  const clearError = (field) => setErrors((prev) => ({ ...prev, [field]: null }))

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Create account</CardTitle>
        <CardDescription>Register as a candidate</CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit} noValidate>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              type="text"
              placeholder="Your name"
              value={name}
              onChange={(e) => { setName(e.target.value); clearError('name') }}
              aria-invalid={!!errors.name}
              aria-describedby={errors.name ? 'name-error' : undefined}
            />
            {errors.name && <p id="name-error" className="text-sm text-destructive" role="alert">{errors.name}</p>}
          </div>
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => { setEmail(e.target.value); clearError('email') }}
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
              onChange={(e) => { setPassword(e.target.value); clearError('password') }}
              autoComplete="new-password"
              aria-invalid={!!errors.password}
              aria-describedby={errors.password ? 'password-error' : undefined}
            />
            {errors.password && <p id="password-error" className="text-sm text-destructive" role="alert">{errors.password}</p>}
          </div>
          {error && (
            <p className="text-sm text-destructive" role="alert">
              {error?.response?.data?.message ?? error?.message ?? 'Registration failed'}
            </p>
          )}
        </CardContent>
        <CardFooter className="flex flex-col gap-2">
          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? 'Creating account...' : 'Create account'}
          </Button>
          <Link to={ROUTES.LOGIN} className="text-sm text-primary underline-offset-4 hover:underline w-full text-center">
            Already have an account? Sign in
          </Link>
          <Link to={ROUTES.REGISTER_HR} className="text-sm text-muted-foreground hover:underline w-full text-center">
            Register as HR instead
          </Link>
        </CardFooter>
      </form>
    </Card>
  )
}
