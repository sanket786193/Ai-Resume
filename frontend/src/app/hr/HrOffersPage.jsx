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
          {list.map((o) => (
            <li key={o.id} className="border rounded-lg p-4">
              <p className="font-medium">{o.candidateName ?? o.title}</p>
              <p className="text-sm text-muted-foreground">{o.status ?? o.createdAt}</p>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
