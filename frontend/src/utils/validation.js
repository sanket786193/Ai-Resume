/**
 * Simple validation helpers. Returns error message or null.
 */
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export function validateRequired(value, fieldName = 'This field') {
  if (value == null || String(value).trim() === '') {
    return `${fieldName} is required`
  }
  return null
}

export function validateEmail(value) {
  const required = validateRequired(value, 'Email')
  if (required) return required
  if (!EMAIL_REGEX.test(String(value).trim())) {
    return 'Please enter a valid email address'
  }
  return null
}
