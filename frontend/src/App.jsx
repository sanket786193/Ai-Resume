import { ErrorBoundary } from '@/components/layout/ErrorBoundary'
import { AppRouter } from '@/routes'
import { Toaster } from 'sonner'

function App() {
  return (
    <ErrorBoundary>
      <AppRouter />
      <Toaster position="top-right" richColors />
    </ErrorBoundary>
  )
}

export default App
