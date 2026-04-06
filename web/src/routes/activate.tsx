import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { activateAccountUrl } from '../lib/api'

type ActivationResult = { ok: true } | { ok: false; message: string }

const activationByToken = new Map<string, Promise<ActivationResult>>()

function activateOnce(token: string): Promise<ActivationResult> {
  const existing = activationByToken.get(token)
  if (existing) {
    return existing
  }
  const p = (async (): Promise<ActivationResult> => {
    try {
      const res = await fetch(activateAccountUrl(), {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token }),
      })
      if (res.status === 204) {
        return { ok: true }
      }
      let message = `Request failed (${res.status})`
      try {
        const body = (await res.json()) as { error?: string; errors?: string[] }
        if (body.error) {
          message = body.error
        } else if (body.errors?.length) {
          message = body.errors.join(' ')
        }
      } catch {
        void 0
      }
      return { ok: false, message }
    } catch (e) {
      return {
        ok: false,
        message: e instanceof Error ? e.message : 'Network error',
      }
    }
  })()
  activationByToken.set(token, p)
  void p.then((result) => {
    if (!result.ok) {
      activationByToken.delete(token)
    }
  })
  return p
}

export const Route = createFileRoute('/activate')({
  validateSearch: (search: Record<string, unknown>) => ({
    token: typeof search.token === 'string' ? search.token : '',
  }),
  component: ActivatePage,
})

function ActivatePage() {
  const { token } = Route.useSearch()
  if (!token) {
    return <ActivateMissing />
  }
  return <ActivateWithToken token={token} />
}

function ActivateMissing() {
  return (
    <main className="gs-panel gs-panel--narrow">
      <p className="gs-brand">Gopher Social</p>
      <h1 className="gs-title">Missing link</h1>
      <p className="gs-muted">
        This page needs an activation token from your email. Use the full link
        from your invitation message.
      </p>
    </main>
  )
}

function ActivateWithToken({ token }: { token: string }) {
  const [phase, setPhase] = useState<'loading' | 'success' | 'error'>('loading')
  const [detail, setDetail] = useState<string>('')

  useEffect(() => {
    let ignore = false
    activateOnce(token).then((result) => {
      if (ignore) {
        return
      }
      if (result.ok) {
        setPhase('success')
      } else {
        setDetail(result.message)
        setPhase('error')
      }
    })
    return () => {
      ignore = true
    }
  }, [token])

  return (
    <main className="gs-panel gs-panel--narrow">
      <p className="gs-brand">Gopher Social</p>

      {phase === 'loading' && (
        <>
          <div className="gs-spinner" aria-hidden />
          <h1 className="gs-title gs-title--center">Activating your account…</h1>
          <p className="gs-muted gs-muted--center">One moment.</p>
        </>
      )}

      {phase === 'success' && (
        <>
          <div className="gs-success-icon" aria-hidden>
            <svg viewBox="0 0 24 24" width="48" height="48" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="1.5" />
              <path
                d="M8 12.5l2.5 2.5 5-5"
                stroke="currentColor"
                strokeWidth="1.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </div>
          <h1 className="gs-title gs-title--center">You’re all set</h1>
          <p className="gs-muted gs-muted--center">
            Your account is activated. You can close this tab and sign in from the
            app.
          </p>
        </>
      )}

      {phase === 'error' && (
        <>
          <div className="gs-error-icon" aria-hidden>
            <svg viewBox="0 0 24 24" width="44" height="44" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="1.5" />
              <path
                d="M12 8v5M12 16h.01"
                stroke="currentColor"
                strokeWidth="1.5"
                strokeLinecap="round"
              />
            </svg>
          </div>
          <h1 className="gs-title gs-title--center">Activation didn’t work</h1>
          <p className="gs-muted gs-muted--center">{detail}</p>
          <p className="gs-muted gs-muted--center gs-small">
            The link may have expired or already been used. Request a new
            invitation from the app if you still need access.
          </p>
        </>
      )}
    </main>
  )
}
