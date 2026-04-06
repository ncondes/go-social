export function getApiBase(): string {
  const raw = import.meta.env.VITE_API_BASE_URL
  if (raw != null && String(raw).trim() !== '') {
    return String(raw).replace(/\/$/, '')
  }
  if (import.meta.env.DEV) {
    return 'http://localhost:8080'
  }
  return ''
}

export function activateAccountUrl(): string {
  return `${getApiBase()}/v1/auth/activate`
}
