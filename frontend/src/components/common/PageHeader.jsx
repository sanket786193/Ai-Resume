/**
 * Page title and optional actions. Stateless.
 */
export function PageHeader({ title, description, actions }) {
  return (
    <div className="flex flex-col gap-1 mb-6">
      <div className="flex items-center justify-between gap-4">
        <h1 className="text-2xl font-semibold tracking-tight">{title}</h1>
        {actions}
      </div>
      {description && <p className="text-muted-foreground text-sm">{description}</p>}
    </div>
  )
}
