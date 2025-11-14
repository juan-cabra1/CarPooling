/**
 * Trips Service for CarPooling Trips API (Port 8002)
 * Handles trip CRUD operations and queries
 */

import apiClient from './api'
import type {
  Trip,
  CreateTripRequest,
  UpdateTripRequest,
  TripFilters,
  TripListResponse,
  ApiResponse,
} from '@/types'

const TRIPS_BASE = '/trips'

/**
 * Get list of trips with optional filters and pagination
 * Public endpoint - no authentication required
 * @param filters - Optional filters (driver_id, status, cities, pagination)
 * @returns Paginated list of trips
 * @example
 * const { trips, total, page, limit } = await tripsService.getTrips({
 *   origin_city: 'Bogotá',
 *   destination_city: 'Medellín',
 *   status: 'published',
 *   page: 1,
 *   limit: 10
 * })
 */
export async function getTrips(filters?: TripFilters): Promise<TripListResponse> {
  const params = new URLSearchParams()

  if (filters?.driver_id) params.append('driver_id', filters.driver_id.toString())
  if (filters?.status) params.append('status', filters.status)
  if (filters?.origin_city) params.append('origin_city', filters.origin_city)
  if (filters?.destination_city) params.append('destination_city', filters.destination_city)
  if (filters?.page) params.append('page', filters.page.toString())
  if (filters?.limit) params.append('limit', filters.limit.toString())

  const queryString = params.toString()
  const url = queryString ? `${TRIPS_BASE}?${queryString}` : TRIPS_BASE

  const response = await apiClient.get<ApiResponse<TripListResponse>>(url)
  return response.data.data!
}

/**
 * Get a specific trip by ID
 * Public endpoint - no authentication required
 * @param id - Trip ID (MongoDB ObjectID as hex string)
 * @returns Trip object
 * @throws Error if trip not found (404)
 * @example
 * const trip = await tripsService.getTripById('507f1f77bcf86cd799439011')
 */
export async function getTripById(id: string): Promise<Trip> {
  const response = await apiClient.get<ApiResponse<Trip>>(`${TRIPS_BASE}/${id}`)
  return response.data.data!
}

/**
 * Create a new trip
 * Requires authentication (JWT token)
 * @param data - Trip creation data
 * @returns Created trip object
 * @throws Error if validation fails (400) or driver not found (404)
 * @example
 * const trip = await tripsService.createTrip({
 *   origin: { city: 'Bogotá', province: 'Cundinamarca', ... },
 *   destination: { city: 'Medellín', province: 'Antioquia', ... },
 *   departure_datetime: '2025-12-01T10:00:00Z',
 *   price_per_seat: 50000,
 *   total_seats: 4,
 *   car: { brand: 'Toyota', model: 'Corolla', ... },
 *   preferences: { pets_allowed: false, ... }
 * })
 */
export async function createTrip(data: CreateTripRequest): Promise<Trip> {
  const response = await apiClient.post<ApiResponse<Trip>>(TRIPS_BASE, data)
  return response.data.data!
}

/**
 * Update an existing trip
 * Requires authentication - only trip owner can update
 * Cannot update if trip has reservations
 * @param id - Trip ID
 * @param data - Fields to update (all optional)
 * @returns Updated trip object
 * @throws Error if not owner (403), trip not found (404), or has reservations (400)
 * @example
 * const updated = await tripsService.updateTrip('507f1f77bcf86cd799439011', {
 *   price_per_seat: 45000,
 *   description: 'Updated description'
 * })
 */
export async function updateTrip(id: string, data: UpdateTripRequest): Promise<Trip> {
  const response = await apiClient.put<ApiResponse<Trip>>(
    `${TRIPS_BASE}/${id}`,
    data
  )
  return response.data.data!
}

/**
 * Delete a trip
 * Requires authentication - only trip owner can delete
 * Cannot delete if trip has active reservations
 * @param id - Trip ID
 * @throws Error if not owner (403), trip not found (404), or has reservations (400)
 * @example
 * await tripsService.deleteTrip('507f1f77bcf86cd799439011')
 */
export async function deleteTrip(id: string): Promise<void> {
  await apiClient.delete(`${TRIPS_BASE}/${id}`)
}

/**
 * Get trips for a specific driver
 * Convenience method that wraps getTrips with driver_id filter
 * @param driverId - Driver user ID
 * @param page - Page number (default: 1)
 * @param limit - Items per page (default: 10)
 * @returns Paginated list of trips for the driver
 * @example
 * const myTrips = await tripsService.getMyTrips(123)
 */
export async function getMyTrips(
  driverId: number,
  page = 1,
  limit = 10
): Promise<TripListResponse> {
  return getTrips({ driver_id: driverId, page, limit })
}

export default {
  getTrips,
  getTripById,
  createTrip,
  updateTrip,
  deleteTrip,
  getMyTrips,
}
