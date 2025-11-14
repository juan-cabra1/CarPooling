/**
 * Base Axios client for CarPooling API
 * Handles authentication, error handling, and request/response interceptors
 */

import axios from 'axios'
import type { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import type { ApiError } from '@/types'

/**
 * Create Axios instance with base configuration
 * Base URL is '/api' which Vite proxy will route to appropriate microservice
 */
const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000, // 10 seconds
  headers: {
    'Content-Type': 'application/json',
  },
})

/**
 * Request interceptor
 * Automatically adds JWT token to Authorization header if available
 */
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')

    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

/**
 * Response interceptor
 * Handles global error cases like 401 Unauthorized and 403 Email not verified
 */
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  (error: AxiosError<ApiError | { success: false; error: string | ApiError }>) => {
    // Handle 401 Unauthorized - clear session and redirect to login
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')

      // Only redirect if not already on login page
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }

    // Handle 403 Forbidden - Email not verified
    if (error.response?.status === 403) {
      const data = error.response.data as { success?: boolean; error?: string | ApiError }
      const errorMsg = typeof data?.error === 'string' ? data.error : ''

      // Si el error es por email no verificado, redirigir a página de reenvío
      if (errorMsg.includes('verificar tu correo') || errorMsg.includes('verificar tu email')) {
        const user = JSON.parse(localStorage.getItem('user') || '{}')
        const email = user.email || ''

        // Solo redirigir si no está ya en la página de reenvío
        if (window.location.pathname !== '/resend-verification') {
          window.location.href = `/resend-verification?email=${encodeURIComponent(email)}`
        }
      }
    }

    return Promise.reject(error)
  }
)

/**
 * Helper function to extract error message from API error response
 * Handles both simple string errors and structured error objects
 */
export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as { success?: boolean; error?: string | ApiError }

    if (data?.error) {
      // Handle structured error object (bookings-api, search-api)
      if (typeof data.error === 'object' && 'message' in data.error) {
        return data.error.message
      }
      // Handle simple string error (users-api, trips-api)
      if (typeof data.error === 'string') {
        return data.error
      }
    }

    // Fallback to Axios error message
    return error.message || 'An unexpected error occurred'
  }

  return 'An unexpected error occurred'
}

export default apiClient
