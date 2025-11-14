/**
 * Bookings Service for CarPooling Bookings API (Port 8003)
 * Handles reservation/booking operations
 * Note: All endpoints require authentication
 */

import apiClient from './api'
import type {
  Booking,
  CreateBookingRequest,
  CancelBookingRequest,
  BookingListResponse,
  BookingFilters,
  ApiResponse,
} from '@/types'

const BOOKINGS_BASE = '/bookings' // Note: Proxy adds /api/v1 prefix automatically

/**
 * Create a new booking/reservation
 * Requires authentication
 * @param data - Booking creation data (trip_id, passenger_id, seats_reserved)
 * @returns Created booking object
 * @throws Error if trip not found (404), insufficient seats (400), or duplicate booking (409)
 * @example
 * const booking = await bookingsService.createBooking({
 *   trip_id: '507f1f77bcf86cd799439011',
 *   passenger_id: 123,
 *   seats_reserved: 2
 * })
 */
export async function createBooking(data: CreateBookingRequest): Promise<Booking> {
  const response = await apiClient.post<ApiResponse<Booking>>(BOOKINGS_BASE, data)
  return response.data.data!
}

/**
 * Get bookings for the authenticated user
 * Requires authentication
 * @param filters - Optional pagination parameters
 * @returns Paginated list of user's bookings
 * @example
 * const { bookings, total, page, limit, total_pages } = await bookingsService.getMyBookings({
 *   page: 1,
 *   limit: 10
 * })
 */
export async function getMyBookings(
  filters?: BookingFilters
): Promise<BookingListResponse> {
  const params = new URLSearchParams()

  if (filters?.page) params.append('page', filters.page.toString())
  if (filters?.limit) params.append('limit', filters.limit.toString())

  const queryString = params.toString()
  const url = queryString ? `${BOOKINGS_BASE}?${queryString}` : BOOKINGS_BASE

  const response = await apiClient.get<ApiResponse<BookingListResponse>>(url)
  return response.data.data!
}

/**
 * Get a specific booking by ID
 * Requires authentication - only booking owner can view
 * @param id - Booking UUID (external identifier)
 * @returns Booking object
 * @throws Error if booking not found (404) or not authorized (403)
 * @example
 * const booking = await bookingsService.getBookingById('550e8400-e29b-41d4-a716-446655440000')
 */
export async function getBookingById(id: string): Promise<Booking> {
  const response = await apiClient.get<ApiResponse<Booking>>(`${BOOKINGS_BASE}/${id}`)
  return response.data.data!
}

/**
 * Cancel a booking
 * Requires authentication - booking owner or trip driver can cancel
 * @param id - Booking UUID
 * @param reason - Optional cancellation reason
 * @throws Error if booking not found (404), already cancelled (400), or not authorized (403)
 * @example
 * await bookingsService.cancelBooking('550e8400-e29b-41d4-a716-446655440000', 'Changed plans')
 */
export async function cancelBooking(id: string, reason?: string): Promise<void> {
  const data: CancelBookingRequest = reason ? { reason } : {}
  await apiClient.patch(`${BOOKINGS_BASE}/${id}/cancel`, data)
}

export default {
  createBooking,
  getMyBookings,
  getBookingById,
  cancelBooking,
}
