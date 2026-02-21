import { useState, useEffect, useCallback } from 'react'
import { debounce } from '@/utils/debounce'

/**
 * Debounced value for search inputs. Updates after delay ms of no changes.
 */
export function useDebounce(value, delay) {
  const [debouncedValue, setDebouncedValue] = useState(value)

  useEffect(() => {
    const id = setTimeout(() => setDebouncedValue(value), delay)
    return () => clearTimeout(id)
  }, [value, delay])

  return debouncedValue
}
