import { useRef } from 'react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

const TOOLBAR_ACTIONS = [
  { label: 'B', title: 'Bold', wrap: '**' },
  { label: 'I', title: 'Italic', wrap: '_' },
  { label: 'H2', title: 'Heading', prefix: '## ' },
  { label: '•', title: 'Bullet list', prefix: '\n- ' },
  { label: '1.', title: 'Numbered list', prefix: '\n1. ' },
  { label: 'Link', title: 'Link', wrap: ['[', '](url)'] },
]

/**
 * Markdown text editor with a simple toolbar. Value is raw markdown string.
 */
export function MarkdownEditor({ value, onChange, placeholder, minHeight = 160, className, id, 'aria-invalid': ariaInvalid, 'aria-describedby': ariaDescribedby }) {
  const textareaRef = useRef(null)

  const insertAtCursor = (before, after = '') => {
    const el = textareaRef.current
    if (!el) return
    const start = el.selectionStart
    const end = el.selectionEnd
    const text = value ?? ''
    const selected = text.slice(start, end)
    const newText = text.slice(0, start) + before + selected + after + text.slice(end)
    onChange(newText)
    setTimeout(() => {
      el.focus()
      const newCursor = start + before.length + selected.length
      el.setSelectionRange(newCursor, newCursor)
    }, 0)
  }

  const handleToolbar = (action) => {
    if (action.prefix) {
      insertAtCursor(action.prefix, '')
    } else if (Array.isArray(action.wrap)) {
      insertAtCursor(action.wrap[0], action.wrap[1])
    } else {
      insertAtCursor(action.wrap, action.wrap)
    }
  }

  return (
    <div className={cn('rounded-md border border-input overflow-hidden', className)}>
      <div className="flex flex-wrap gap-0.5 p-1 border-b border-input bg-muted/50">
        {TOOLBAR_ACTIONS.map((action) => (
          <Button
            key={action.label}
            type="button"
            variant="ghost"
            size="sm"
            className="h-8 w-8 p-0 font-semibold text-xs"
            title={action.title}
            onClick={() => handleToolbar(action)}
          >
            {action.label}
          </Button>
        ))}
      </div>
      <textarea
        ref={textareaRef}
        id={id}
        className="flex w-full rounded-none border-0 bg-transparent px-3 py-2 text-sm shadow-none placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-0 resize-y min-w-0"
        style={{ minHeight }}
        value={value ?? ''}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        aria-invalid={ariaInvalid}
        aria-describedby={ariaDescribedby}
      />
    </div>
  )
}
