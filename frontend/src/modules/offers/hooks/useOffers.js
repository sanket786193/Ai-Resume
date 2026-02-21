import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { offersApi } from '@/services/api'

export const OFFERS_QUERY_KEY = ['offers']

export function useOffersList(params) {
  return useQuery({
    queryKey: [...OFFERS_QUERY_KEY, params],
    queryFn: () => offersApi.list(params),
  })
}

export function useSendSelection() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: offersApi.sendSelection,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: OFFERS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: ['ats', 'pipeline'] })
    },
  })
}

export function useSendRejection() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: offersApi.sendRejection,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: OFFERS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: ['ats', 'pipeline'] })
    },
  })
}
