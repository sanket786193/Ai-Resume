import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { cn } from '@/lib/utils'

/**
 * Renders markdown content safely (no raw HTML by default). Use for job descriptions etc.
 */
export function SafeMarkdown({ content, className }) {
  if (content == null || content === '') return null
  const text = typeof content === 'string' ? content : ''
  return (
    <div className={cn('markdown-body text-sm space-y-2 [&_ul]:list-disc [&_ul]:pl-6 [&_ol]:list-decimal [&_ol]:pl-6 [&_h2]:text-lg [&_h2]:font-semibold [&_p]:leading-relaxed', className)}>
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{text}</ReactMarkdown>
    </div>
  )
}
