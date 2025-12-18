/**
 * Ratings Service for CarPooling Users API (Port 8001)
 * Handles user rating operations (driver and passenger ratings)
 */

import apiClient from './api'
import type {
  Rating,
  CreateRatingRequest,
  ApiResponse,
} from '@/types'

/**
 * Get ratings for a specific user
 * Requires authentication
 * @param userId - User ID to get ratings for
 * @param page - Page number for pagination (default: 1)
 * @param limit - Items per page (default: 10)
 * @returns Array of ratings
 * @example
 * const ratings = await ratingsService.getUserRatings(123, 1, 10)
 */
export async function getUserRatings(
  userId: number,
  page = 1,
  limit = 10
): Promise<Rating[]> {
  const params = new URLSearchParams()
  params.append('page', page.toString())
  params.append('limit', limit.toString())

  const url = `/users/${userId}/ratings?${params.toString()}`

  const response = await apiClient.get<ApiResponse<Rating[]>>(url)
  return response.data.data!
}

/**
 * Create a new rating for a user
 * Note: This uses the internal endpoint which may not be directly exposed to frontend
 * Typically called by trips-api after trip completion
 * @param data - Rating data (rater, rated user, trip, role, score, comment)
 * @returns Created rating
 * @example
 * const rating = await ratingsService.createRating({
 *   rater_id: 123,
 *   rated_user_id: 456,
 *   trip_id: '507f1f77bcf86cd799439011',
 *   role_rated: 'conductor',
 *   score: 5,
 *   comment: 'Great driver!'
 * })
 */
export async function createRating(data: CreateRatingRequest): Promise<Rating> {
  const response = await apiClient.post<ApiResponse<Rating>>(
    '/internal/ratings',
    data
  )
  return response.data.data!
}

/**
 * Rate a driver after a trip
 * Convenience method for rating drivers
 * @param raterId - ID of user giving the rating
 * @param driverId - ID of driver being rated
 * @param tripId - Trip ID (MongoDB ObjectID)
 * @param score - Rating score (1-5)
 * @param comment - Optional comment
 * @returns Created rating
 */
export async function rateDriver(
  raterId: number,
  driverId: number,
  tripId: string,
  score: number,
  comment?: string
): Promise<Rating> {
  const data: CreateRatingRequest = {
    rater_id: raterId,
    rated_user_id: driverId,
    trip_id: tripId,
    role_rated: 'conductor',
    score,
    comment,
  }

  return createRating(data)
}

/**
 * Rate a passenger after a trip
 * Convenience method for rating passengers
 * @param raterId - ID of user giving the rating (usually driver)
 * @param passengerId - ID of passenger being rated
 * @param tripId - Trip ID (MongoDB ObjectID)
 * @param score - Rating score (1-5)
 * @param comment - Optional comment
 * @returns Created rating
 */
export async function ratePassenger(
  raterId: number,
  passengerId: number,
  tripId: string,
  score: number,
  comment?: string
): Promise<Rating> {
  const data: CreateRatingRequest = {
    rater_id: raterId,
    rated_user_id: passengerId,
    trip_id: tripId,
    role_rated: 'pasajero',
    score,
    comment,
  }

  return createRating(data)
}

export default {
  getUserRatings,
  createRating,
  rateDriver,
  ratePassenger,
}
