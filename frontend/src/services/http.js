import { config } from '@/config/env'

/**
 * Normalize API errors into a consistent shape. Single place for error handling.
 */
export function normalizeApiError(error) {
  if (!error) return { message: 'An unexpected error occurred' }
  if (typeof error === 'string') return { message: error }
  if (error.response) {
    const data = error.response.data
    const message = data?.error?.message ?? data?.message ?? error.message ?? 'Request failed'
    const status = error.response.status
    return { message, status, data }
  }
  if (error.message) return { message: error.message }
  return { message: 'Network or unknown error' }
}

/**
 * HTTP client. All API calls go through this. No business logic.
 */
async function request(path, options = {}) {
  const url = path.startsWith('http') ? path : `${config.apiBaseUrl}${path}`
  const headers = { ...options.headers }
  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json'
  }
  const token = getStoredToken()
  if (token) headers.Authorization = `Bearer ${token}`

  const requestBody = options.body instanceof FormData ? options.body : options.body
  const response = await fetch(url, {
    ...options,
    headers,
    body: requestBody ? (typeof requestBody === 'string' ? requestBody : requestBody instanceof FormData ? requestBody : JSON.stringify(requestBody)) : undefined,
  })

  const contentType = response.headers.get('content-type')
  const isJson = contentType && contentType.includes('application/json')
  const responseData = isJson ? await response.json().catch(() => null) : await response.text().catch(() => null)

  if (!response.ok) {
    // Go API error shape: { success: false, error: { code, message } }
    const message = responseData?.error?.message ?? responseData?.message ?? response.statusText
    const error = new Error(message)
    error.response = { status: response.status, data: responseData }
    if (response.status === 401) {
      try { window.localStorage.removeItem('ats_token') } catch {}
    }
    throw error
  }

  // Go API success shape: { success: true, data: T } — unwrap so callers get T
  if (responseData && responseData.success === true && 'data' in responseData) {
    return responseData.data
  }
  return responseData
}

function getStoredToken() {
  try {
    return window.localStorage.getItem('ats_token')
  } catch {
    return null
  }
}

export const http = {
  get: (path, opts) => request(path, { ...opts, method: 'GET' }),
  post: (path, body, opts) => request(path, { ...opts, method: 'POST', body: body ?? undefined }),
  postForm: (path, formData, opts) => request(path, { ...opts, method: 'POST', body: formData }),
  put: (path, body, opts) => request(path, { ...opts, method: 'PUT', body: body ? JSON.stringify(body) : undefined }),
  patch: (path, body, opts) => request(path, { ...opts, method: 'PATCH', body: body ? JSON.stringify(body) : undefined }),
  delete: (path, opts) => request(path, { ...opts, method: 'DELETE' }),
}
