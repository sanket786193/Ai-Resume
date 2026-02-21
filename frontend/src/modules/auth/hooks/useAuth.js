import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { authApi } from '@/services/api'

const AUTH_QUERY_KEY = ['auth', 'me']

function getStoredToken() {
  try {
    return window.localStorage.getItem('ats_token')
  } catch {
    return null
  }
}

function setStoredToken(token) {
  try {
    if (token) window.localStorage.setItem('ats_token', token)
    else window.localStorage.removeItem('ats_token')
  } catch {}
}

export function useAuth() {
  const queryClient = useQueryClient()
  const token = getStoredToken()

  const { data: user, isPending, isError, error } = useQuery({
    queryKey: AUTH_QUERY_KEY,
    queryFn: authApi.me,
    enabled: !!token,
    retry: false,
    staleTime: 5 * 60 * 1000,
  })

  const loginMutation = useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      const t = data?.access_token ?? data?.accessToken ?? data?.token
      if (t) setStoredToken(t)
      queryClient.setQueryData(AUTH_QUERY_KEY, data?.user ?? data)
    },
  })

  const registerMutation = useMutation({
    mutationFn: authApi.registerCandidate,
    onSuccess: (data) => {
      const t = data?.access_token ?? data?.accessToken ?? data?.token
      if (t) setStoredToken(t)
      queryClient.setQueryData(AUTH_QUERY_KEY, data?.user ?? data)
    },
  })

  const registerHRMutation = useMutation({
    mutationFn: authApi.registerHR,
    onSuccess: (data) => {
      const t = data?.access_token ?? data?.accessToken ?? data?.token
      if (t) setStoredToken(t)
      queryClient.setQueryData(AUTH_QUERY_KEY, data?.user ?? data)
    },
  })

  const logout = async () => {
    try {
      await authApi.logout()
    } finally {
      setStoredToken(null)
      queryClient.setQueryData(AUTH_QUERY_KEY, null)
      queryClient.removeQueries({ queryKey: AUTH_QUERY_KEY })
    }
  }

  return {
    user: user ?? null,
    isPending: !!token && isPending,
    isError,
    error,
    login: loginMutation.mutateAsync,
    loginMutation,
    register: registerMutation.mutateAsync,
    registerMutation,
    registerHR: registerHRMutation.mutateAsync,
    registerHRMutation,
    logout,
    isAuthenticated: !!user,
  }
}
