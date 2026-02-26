import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { offersApi } from '@/services/api'

export const OFFERS_QUERY_KEY = ['offers']

export function useOffersList(params) {
  return useQuery({
    queryKey: [...OFFERS_QUERY_KEY, params],
    queryFn: () => offersApi.list(params),
  })
}

export function useCreateOffer() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: offersApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: OFFERS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: ['ats', 'pipeline'] })
    },
  })
}

export function useSendSelection() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: offersApi.accept,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: OFFERS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: ['ats', 'pipeline'] })
    },
  })
}

export function useSendRejection() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: offersApi.reject,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: OFFERS_QUERY_KEY })
      queryClient.invalidateQueries({ queryKey: ['ats', 'pipeline'] })
    },
  })
}
