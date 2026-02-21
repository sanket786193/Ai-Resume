import { PageHeader } from '@/components/common/PageHeader'
import { AtsPipelineBoard } from '@/modules/ats/components/AtsPipelineBoard'

export function HrAtsPage() {
  return (
    <div>
      <PageHeader
        title="ATS Pipeline"
        description="Drag candidates between stages. Updates are saved automatically."
      />
      <AtsPipelineBoard />
    </div>
  )
}
