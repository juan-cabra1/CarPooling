/**
 * Per-service Axios clients for CarPooling API
 *
 * Dev:  no env vars set → all clients use '/api' base, Vite proxy routes by path prefix
 * Prod: VITE_*_API_URL set → each client calls its Railway URL directly
 *
 * Env vars (set as build args in Railway):
 *   VITE_USERS_API_URL     = https://users-api-xxx.up.railway.app
 *   VITE_TRIPS_API_URL     = https://trips-api-xxx.up.railway.app
 *   VITE_BOOKINGS_API_URL  = https://bookings-api-xxx.up.railway.app
 *   VITE_SEARCH_API_URL    = https://search-api-xxx.up.railway.app
 *   VITE_PAYMENTS_API_URL  = https://payments-api-xxx.up.railway.app
 */

import axios from 'axios'
import type { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import type { ApiError } from '@/types'

function createApiClient(baseURL: string) {
  const client = axios.create({
    baseURL,
    timeout: 10000,
    headers: { 'Content-Type': 'application/json' },
  })

  client.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      const token = localStorage.getItem('token')
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    },
    (error: AxiosError) => Promise.reject(error)
  )

  client.interceptors.response.use(
    (response: AxiosResponse) => response,
    (error: AxiosError<ApiError | { success: false; error: string | ApiError }>) => {
      if (error.response?.status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        if (window.location.pathname !== '/login') {
          window.location.href = '/login'
        }
      }

      if (error.response?.status === 403) {
        const data = error.response.data as { success?: boolean; error?: string | ApiError }
        const errorMsg = typeof data?.error === 'string' ? data.error : ''
        if (errorMsg.includes('verificar tu correo') || errorMsg.includes('verificar tu email')) {
          const user = JSON.parse(localStorage.getItem('user') || '{}')
          const email = user.email || ''
          if (window.location.pathname !== '/resend-verification') {
            window.location.href = `/resend-verification?email=${encodeURIComponent(email)}`
          }
        }
      }

      return Promise.reject(error)
    }
  )

  return client
}

// ─── Per-service base URLs ────────────────────────────────────────────────────
//
// In dev (no env vars): all fall back to '/api' — the Vite proxy routes each
// path prefix (/api/trips, /api/bookings, etc.) to the correct local port.
//
// In prod: use Railway URL directly. Bookings/Search/Payments APIs have an
// /api/v1 prefix, so we append it here once, keeping service files clean.

export const usersApi = createApiClient(
  import.meta.env.VITE_USERS_API_URL || '/api'
)

export const tripsApi = createApiClient(
  import.meta.env.VITE_TRIPS_API_URL || '/api'
)

export const bookingsApi = createApiClient(
  import.meta.env.VITE_BOOKINGS_API_URL
    ? `${import.meta.env.VITE_BOOKINGS_API_URL}/api/v1`
    : '/api'
)

export const searchApi = createApiClient(
  import.meta.env.VITE_SEARCH_API_URL
    ? `${import.meta.env.VITE_SEARCH_API_URL}/api/v1`
    : '/api'
)

export const paymentsApi = createApiClient(
  import.meta.env.VITE_PAYMENTS_API_URL
    ? `${import.meta.env.VITE_PAYMENTS_API_URL}/api/v1`
    : '/api'
)

// ─── Helpers ─────────────────────────────────────────────────────────────────

export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as { success?: boolean; error?: string | ApiError }

    if (data?.error) {
      if (typeof data.error === 'object' && 'message' in data.error) {
        return data.error.message
      }
      if (typeof data.error === 'string') {
        return data.error
      }
    }

    return error.message || 'An unexpected error occurred'
  }

  return 'An unexpected error occurred'
}

// Default export kept for backward compatibility (points to usersApi)
export default usersApi
