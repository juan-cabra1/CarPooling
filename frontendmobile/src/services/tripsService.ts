/**
 * Trips Service for CarPooling Trips API (Port 8002)
 * Handles trip CRUD operations and queries
 */

import { tripsApi } from './api'
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
  console.log('getTrips called with filters:', filters);
  const params = new URLSearchParams()

  if (filters?.driver_id !== null && filters?.driver_id !== undefined) {
    params.append('driver_id', filters.driver_id.toString());
  }
  if (filters?.status) params.append('status', filters.status)
  if (filters?.origin_city) params.append('origin_city', filters.origin_city)
  if (filters?.destination_city) params.append('destination_city', filters.destination_city)
  if (filters?.page !== null && filters?.page !== undefined) {
    params.append('page', filters.page.toString());
  }
  if (filters?.limit !== null && filters?.limit !== undefined) {
    params.append('limit', filters.limit.toString());
  }

  const queryString = params.toString()
  const url = queryString ? `${TRIPS_BASE}?${queryString}` : TRIPS_BASE

  const response = await tripsApi.get<ApiResponse<TripListResponse>>(url)
  const data = response.data.data!

  // FIX: Filter trips that are null/undefined or don't have an ID.
  // The API should return 'id', not '_id'. No mapping is needed.
  return {
    ...data,
    trips: (data.trips || []).filter((trip) => trip && trip.id),
  };
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
  const response = await tripsApi.get<ApiResponse<Trip>>(`${TRIPS_BASE}/${id}`)
  return response.data.data!;
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
  const response = await tripsApi.post<ApiResponse<Trip>>(TRIPS_BASE, data)
  return response.data.data!;
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
  const response = await tripsApi.put<ApiResponse<Trip>>(
    `${TRIPS_BASE}/${id}`,
    data
  )
  return response.data.data!;
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
  await tripsApi.delete(`${TRIPS_BASE}/${id}`)
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

/**
 * Start a trip (change status to in_progress)
 * Requires authentication - only trip owner (driver) can start
 * @param id - Trip ID
 * @returns Updated trip object
 * @throws Error if not owner (403), trip not found (404), or trip not in valid status (400)
 * @example
 * const trip = await tripsService.startTrip('507f1f77bcf86cd799439011')
 */
export async function startTrip(id: string): Promise<Trip> {
  const response = await tripsApi.post<ApiResponse<Trip>>(`${TRIPS_BASE}/${id}/start`)
  return response.data.data!;
}

/**
 * Complete a trip (change status to completed)
 * Requires authentication - only trip owner (driver) can complete
 * @param id - Trip ID
 * @returns Updated trip object
 * @throws Error if not owner (403), trip not found (404), or trip not in progress (400)
 * @example
 * const trip = await tripsService.completeTrip('507f1f77bcf86cd799439011')
 */
export async function completeTrip(id: string): Promise<Trip> {
  const response = await tripsApi.post<ApiResponse<Trip>>(`${TRIPS_BASE}/${id}/complete`)
  return response.data.data!;
}

export default {
  getTrips,
  getTripById,
  createTrip,
  updateTrip,
  deleteTrip,
  getMyTrips,
  startTrip,
  completeTrip,
}