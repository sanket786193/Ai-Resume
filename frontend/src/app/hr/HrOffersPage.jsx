import { PageHeader } from '@/components/common/PageHeader'
import { EmptyState } from '@/components/common/EmptyState'
import { useOffersList } from '@/modules/offers/hooks/useOffers'

export function HrOffersPage() {
  const { data, isPending, isError } = useOffersList()
  const list = Array.isArray(data) ? data : data?.list ?? []

  return (
    <div>
      <PageHeader title="Offers" description="Manage offers and hiring decisions" />
      {isPending && <div className="text-muted-foreground">Loading...</div>}
      {!isPending && (isError || list.length === 0) && (
        <EmptyState
          title={isError ? 'Failed to load' : 'No offers yet'}
          description={isError ? 'Please try again later.' : 'Send selections from the pipeline to create offers.'}
        />
      )}
      {!isPending && !isError && list.length > 0 && (
        <ul className="space-y-2">
          {list.map((o) => {
            const id = o.id ?? o.ID
            const amount = o.amount ?? o.Amount
            const currency = o.currency ?? o.Currency
            const status = o.status ?? o.Status
            const createdAt = o.created_at ?? o.CreatedAt
            return (
              <li key={id} className="border rounded-lg p-4">
                <p className="font-medium">
                  {amount && currency ? `${amount} ${currency}` : 'Offer'}
                  {' — '}
                  <span className="capitalize">{status ?? '—'}</span>
                </p>
                {createdAt && (
                  <p className="text-sm text-muted-foreground mt-1">
                    Created {new Date(createdAt).toLocaleDateString(undefined, { dateStyle: 'medium' })}
                  </p>
                )}
              </li>
            )
          })}
        </ul>
      )}
    </div>
  )
}
