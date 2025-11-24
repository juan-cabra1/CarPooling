/**
 * Admin service for CarPooling API
 * Handles admin-only operations like user management, trip management, and booking management
 */

import apiClient from './api'
import type { User } from '@/types/user'
import type { Trip } from '@/types/trip'
import type { Booking } from '@/types/booking'

/**
 * Pagination metadata
 */
export interface PaginationMeta {
  page: number
  limit: number
  total: number
  totalPages: number
}

/**
 * Admin users list response
 */
export interface AdminUsersResponse {
  users: User[]
  pagination: PaginationMeta
}

/**
 * Admin bookings list response
 */
export interface AdminBookingsResponse {
  bookings: Booking[]
  pagination: PaginationMeta
}

/**
 * Get all users with pagination and filters (admin only)
 */
export async function getAllUsers(
  page = 1,
  limit = 10,
  role?: string,
  search?: string
): Promise<AdminUsersResponse> {
  const params = new URLSearchParams()
  params.append('page', page.toString())
  params.append('limit', limit.toString())
  if (role) params.append('role', role)
  if (search) params.append('search', search)

  const response = await apiClient.get(`/admin/users?${params.toString()}`)
  return response.data.data
}

/**
 * Get all bookings with pagination and filters (admin only)
 */
export async function getAllBookings(
  page = 1,
  limit = 10,
  status?: string,
  tripId?: string,
  passengerId?: number
): Promise<AdminBookingsResponse> {
  const params = new URLSearchParams()
  params.append('page', page.toString())
  params.append('limit', limit.toString())
  if (status) params.append('status', status)
  if (tripId) params.append('trip_id', tripId)
  if (passengerId) params.append('passenger_id', passengerId.toString())

  const response = await apiClient.get(`/admin/bookings?${params.toString()}`)
  return response.data.data
}

/**
 * Get all trips with pagination (admin only - uses existing trips endpoint)
 */
export async function getAllTrips(
  page = 1,
  limit = 10,
  status?: string
): Promise<{ trips: Trip[]; total: number }> {
  const params = new URLSearchParams()
  params.append('page', page.toString())
  params.append('limit', limit.toString())
  if (status) params.append('status', status)

  const response = await apiClient.get(`/trips?${params.toString()}`)
  return response.data.data
}

/**
 * Force user to re-authenticate by un-verifying and resending email verification
 * Admin-only operation
 * @param userId - User ID to force re-auth
 * @param _email - User email (unused but kept for function signature compatibility)
 */
export async function forceReauthentication(userId: number, _email: string): Promise<void> {
  // Call admin endpoint that will unverify and resend email
  const response = await apiClient.post(`/admin/users/${userId}/force-reauth`)
  return response.data
}

const adminService = {
  getAllUsers,
  getAllBookings,
  getAllTrips,
  forceReauthentication,
}

export default adminService
