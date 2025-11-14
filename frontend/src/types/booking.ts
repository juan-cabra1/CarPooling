/**
 * Booking/Reservation-related types for the CarPooling Bookings API (Port 8003)
 * Database: MySQL
 */

/**
 * Booking status values
 */
export type BookingStatus =
  | 'pending' // Initial state, waiting for confirmation
  | 'confirmed' // Trip confirmed reservation, seats reserved
  | 'cancelled' // Cancelled by passenger or driver
  | 'completed' // Trip completed successfully
  | 'failed' // Reservation failed (no seats, etc.)

/**
 * Complete booking object
 * Note: Uses BookingUUID as the external ID (not internal MySQL ID)
 */
export interface Booking {
  id: string // UUID (external identifier)
  trip_id: string // MongoDB ObjectID from trips-api
  passenger_id: number // User ID from users-api
  driver_id?: number // Set after confirmation
  seats_requested: number
  total_price: number
  status: BookingStatus
  cancelled_at?: string // ISO 8601 datetime, only if cancelled
  cancellation_reason?: string
  created_at: string // ISO 8601 datetime
  updated_at: string // ISO 8601 datetime
}

/**
 * Create a new booking/reservation
 */
export interface CreateBookingRequest {
  trip_id: string // MongoDB ObjectID
  passenger_id: number
  seats_reserved: number // Must be >= 1
}

/**
 * Cancel a booking
 */
export interface CancelBookingRequest {
  reason?: string // Optional cancellation reason
}

/**
 * Paginated booking list response
 */
export interface BookingListResponse {
  bookings: Booking[]
  total: number
  page: number
  limit: number
  total_pages: number
}

/**
 * Query parameters for listing bookings
 */
export interface BookingFilters {
  page?: number // Default: 1
  limit?: number // Default: 10
}
