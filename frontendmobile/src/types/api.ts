/**
 * Generic API types used across all CarPooling microservices
 * Defines standard response formats, error handling, and HTTP status codes
 */

/**
 * Standard API error object
 */
export interface ApiError {
  code: string // Error code (e.g., "NOT_FOUND", "UNAUTHORIZED")
  message: string // Human-readable error message
  details?: Record<string, any> // Optional additional error details
}

/**
 * Standard API response wrapper
 * All API responses follow this format
 */
export interface ApiResponse<T = any> {
  success: boolean
  data?: T // Present when success is true
  error?: ApiError | string // Present when success is false
}

/**
 * Paginated response format
 * Used for list endpoints across all services
 */
export interface PaginatedResponse<T> {
  items: T[] // Array of items (could be trips, bookings, etc.)
  total: number // Total count across all pages
  page: number // Current page number (1-indexed)
  limit: number // Items per page
  total_pages: number // Total pages (calculated: ceil(total / limit))
}

/**
 * HTTP status codes used by the API
 */
export const HttpStatus = {
  OK: 200,
  CREATED: 201,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,
  INTERNAL_SERVER_ERROR: 500,
} as const

/**
 * Common error codes across all microservices
 */
export type ErrorCode =
  // Authentication/Authorization
  | 'UNAUTHORIZED'
  | 'FORBIDDEN'
  | 'INVALID_TOKEN'
  | 'TOKEN_EXPIRED'

  // Resource errors
  | 'NOT_FOUND'
  | 'USER_NOT_FOUND'
  | 'TRIP_NOT_FOUND'
  | 'BOOKING_NOT_FOUND'
  | 'DRIVER_NOT_FOUND'

  // Validation errors
  | 'VALIDATION_ERROR'
  | 'INVALID_EMAIL'
  | 'INVALID_PASSWORD'
  | 'INVALID_CREDENTIALS'

  // Conflict errors
  | 'DUPLICATE_EMAIL'
  | 'DUPLICATE_BOOKING'
  | 'OPTIMISTIC_LOCK_FAILED'

  // Business logic errors
  | 'PAST_DEPARTURE'
  | 'NO_SEATS_AVAILABLE'
  | 'INSUFFICIENT_SEATS'
  | 'HAS_RESERVATIONS'
  | 'CANNOT_BOOK_OWN_TRIP'
  | 'CANNOT_CANCEL_COMPLETED'
  | 'BOOKING_ALREADY_CANCELLED'
  | 'TRIP_NOT_PUBLISHED'

  // Service availability
  | 'TRIPS_API_UNAVAILABLE'
  | 'USERS_API_UNAVAILABLE'
  | 'SERVICE_UNAVAILABLE'

  // Generic errors
  | 'INTERNAL_ERROR'
  | 'UNKNOWN_ERROR'

/**
 * Health check response format
 */
export interface HealthCheckResponse {
  status: 'ok' | 'degraded' | 'error'
  service: string
  port?: string
  timestamp: string
  dependencies?: {
    mongodb?: 'connected' | 'disconnected'
    mysql?: 'connected' | 'disconnected'
    solr?: 'connected' | 'disconnected'
    memcached?: 'connected' | 'disconnected'
    rabbitmq?: 'connected' | 'disconnected'
  }
}

/**
 * JWT token payload (decoded)
 */
export interface JwtPayload {
  user_id: number
  email: string
  role: string
  exp: number // Expiration timestamp
  iat?: number // Issued at timestamp
}

/**
 * Request configuration for API calls
 */
export interface RequestConfig {
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  headers?: Record<string, string>
  body?: any
  params?: Record<string, string | number | boolean>
}

/**
 * Pagination parameters for list requests
 */
export interface PaginationParams {
  page?: number // Default: 1
  limit?: number // Default: 10-20 depending on service
}
