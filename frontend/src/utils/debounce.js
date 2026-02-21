/**
 * Debounce for search and other high-frequency inputs.
 */
export function debounce(fn, ms) {
  let timeoutId
  return function debounced(...args) {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn.apply(this, args), ms)
  }
}
