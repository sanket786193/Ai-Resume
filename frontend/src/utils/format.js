/**
 * Formatting utilities. No business logic.
 */
export function formatDate(value) {
  if (value == null) return ''
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? '' : d.toLocaleDateString()
}

export function formatDateTime(value) {
  if (value == null) return ''
  const d = new Date(value)
  return Number.isNaN(d.getTime()) ? '' : d.toLocaleString()
}
