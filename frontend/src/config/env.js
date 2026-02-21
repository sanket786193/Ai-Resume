/**
 * App configuration. In production, replace with env vars.
 */
export const config = {
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  appName: 'HR ATS',
}
