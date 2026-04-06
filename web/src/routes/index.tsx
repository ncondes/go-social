import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: HomePage,
})

function HomePage() {
  return (
    <main className="gs-panel">
      <p className="gs-brand">Gopher Social</p>
      <h1 className="gs-title">Account activation</h1>
      <p className="gs-muted">
        Open the link in your registration email to confirm your account. It
        opens this site with a secure token (for example{' '}
        <code className="gs-code">/activate?token=…</code>).
      </p>
    </main>
  )
}
