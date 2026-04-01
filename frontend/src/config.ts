/** Set `VITE_API_BASE` in `frontend/.env` (see `.env.example`). Empty uses Vite dev proxy for `/passkey`. */
export const API_BASE = (import.meta.env.VITE_API_BASE ?? '').replace(/\/$/, '')

export const JSON_HEADERS = { 'Content-Type': 'application/json' } as const
