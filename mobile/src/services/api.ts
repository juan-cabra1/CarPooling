/**
 * Base Axios client for CarPooling API - Mobile Version
 * Handles authentication, error handling, and request/response interceptors
 * Adapted for React Native with AsyncStorage
 */

import axios from 'axios'
import type { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import AsyncStorage from '@react-native-async-storage/async-storage'
import type { ApiError } from '@/types'

// Navigation reference for redirects (will be set by App.tsx)
let navigationRef: any = null

export function setNavigationRef(ref: any) {
  navigationRef = ref
}

/**
 * Create Axios instance with base configuration
 * Base URLs apuntan a Railway deployment
 */
const apiClient = axios.create({
  baseURL: 'https://users-api-production-a461.up.railway.app', // Users API
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Additional API clients for other services
export const tripsApi = axios.create({
  baseURL: 'https://trips-api-production.up.railway.app', // Trips API
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

export const bookingsApi = axios.create({
  baseURL: 'https://bookings-api-production-c508.up.railway.app', // Bookings API
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

export const searchApi = axios.create({
  baseURL: 'https://search-api-production-51e4.up.railway.app', // Search API
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

export const paymentsApi = axios.create({
  baseURL: 'https://payments-api-production-e0e9.up.railway.app', // Payments API (Railway)
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

/**
 * Request interceptor function
 * Automatically adds JWT token to Authorization header if available
 */
const requestInterceptor = async (config: InternalAxiosRequestConfig) => {
  try {
    const token = await AsyncStorage.getItem('token')

    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  } catch (error) {
    console.error('Error reading token from AsyncStorage:', error)
    return config
  }
}

const requestErrorInterceptor = (error: AxiosError) => {
  return Promise.reject(error)
}

/**
 * Apply request interceptor to all API clients
 */
apiClient.interceptors.request.use(requestInterceptor, requestErrorInterceptor)
tripsApi.interceptors.request.use(requestInterceptor, requestErrorInterceptor)
bookingsApi.interceptors.request.use(requestInterceptor, requestErrorInterceptor)
searchApi.interceptors.request.use(requestInterceptor, requestErrorInterceptor)
paymentsApi.interceptors.request.use(requestInterceptor, requestErrorInterceptor)

/**
 * Response interceptor function
 * Handles global error cases like 401 Unauthorized and 403 Email not verified
 */
const responseSuccessInterceptor = (response: AxiosResponse) => {
  return response
}

const responseErrorInterceptor = async (error: AxiosError<ApiError | { success: false; error: string | ApiError }>) => {
  // Handle 401 Unauthorized - clear session and redirect to login
  if (error.response?.status === 401) {
    try {
      await AsyncStorage.removeItem('token')
      await AsyncStorage.removeItem('user')

      // Navigate to Login screen if navigation ref is available
      if (navigationRef?.isReady()) {
        const currentRoute = navigationRef.getCurrentRoute()?.name
        if (currentRoute !== 'Login') {
          navigationRef.navigate('Login')
        }
      }
    } catch (storageError) {
      console.error('Error clearing AsyncStorage:', storageError)
    }
  }

  // Handle 403 Forbidden - Email not verified
  if (error.response?.status === 403) {
    const data = error.response.data as { success?: boolean; error?: string | ApiError }
    const errorMsg = typeof data?.error === 'string' ? data.error : ''

    // Si el error es por email no verificado, redirigir a página de reenvío
    if (errorMsg.includes('verificar tu correo') || errorMsg.includes('verificar tu email')) {
      try {
        const userStr = await AsyncStorage.getItem('user')
        const user = userStr ? JSON.parse(userStr) : {}
        const email = user.email || ''

        // Navigate to ResendVerification screen if navigation ref is available
        if (navigationRef?.isReady()) {
          const currentRoute = navigationRef.getCurrentRoute()?.name
          if (currentRoute !== 'ResendVerification') {
            navigationRef.navigate('ResendVerification', { email })
          }
        }
      } catch (storageError) {
        console.error('Error reading user from AsyncStorage:', storageError)
      }
    }
  }

  return Promise.reject(error)
}

/**
 * Apply response interceptor to all API clients
 */
apiClient.interceptors.response.use(responseSuccessInterceptor, responseErrorInterceptor)
tripsApi.interceptors.response.use(responseSuccessInterceptor, responseErrorInterceptor)
bookingsApi.interceptors.response.use(responseSuccessInterceptor, responseErrorInterceptor)
searchApi.interceptors.response.use(responseSuccessInterceptor, responseErrorInterceptor)
paymentsApi.interceptors.response.use(responseSuccessInterceptor, responseErrorInterceptor)

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
