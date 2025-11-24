/**
 * Search Service for CarPooling Search API (Port 8004)
 * Handles trip search, geospatial queries, autocomplete, and popular routes
 * Note: All endpoints are public (no authentication required)
 */

import apiClient from './api'
import type {
  SearchQuery,
  SearchResponse,
  SearchTrip,
  GeospatialSearchParams,
  PopularRoute,
  ApiResponse,
} from '@/types'

const SEARCH_BASE = '/search' // Note: Proxy adds /api/v1 prefix automatically

/**
 * Search trips with filters
 * Public endpoint - no authentication required
 * Supports structured location search with Google Places integration
 * @param query - Search filters and pagination
 * @returns Paginated search results with denormalized trip and driver data
 * @example
 * const { trips, total, page, limit } = await searchService.searchTrips({
 *   origin: { city: 'Bogotá', province: 'Cundinamarca', address: '...' },
 *   destination: { city: 'Medellín', province: 'Antioquia', address: '...' },
 *   origin_radius: 10,
 *   min_seats: 2,
 *   max_price: 50000,
 *   sort_by: 'price',
 *   sort_order: 'asc',
 *   page: 1,
 *   limit: 20
 * })
 */
export async function searchTrips(query: SearchQuery): Promise<SearchResponse<SearchTrip>> {
  const params = new URLSearchParams()

  // Origin location (structured)
  if (query.origin) {
    if (query.origin.city) params.append('origin_city', query.origin.city)
    if (query.origin.province) params.append('origin_province', query.origin.province)
    if (query.origin.address) params.append('origin_address', query.origin.address)

    // Coordinates from Google Places (if available)
    if (query.origin.coordinates) {
      params.append('origin_lat', query.origin.coordinates.lat.toString())
      params.append('origin_lng', query.origin.coordinates.lng.toString())
    }
  }

  // Destination location (structured)
  if (query.destination) {
    if (query.destination.city) params.append('destination_city', query.destination.city)
    if (query.destination.province) params.append('destination_province', query.destination.province)
    if (query.destination.address) params.append('destination_address', query.destination.address)

    // Coordinates from Google Places (if available)
    if (query.destination.coordinates) {
      params.append('destination_lat', query.destination.coordinates.lat.toString())
      params.append('destination_lng', query.destination.coordinates.lng.toString())
    }
  }

  // Geospatial search radius
  if (query.origin_radius) params.append('origin_radius', query.origin_radius.toString())
  if (query.destination_radius) params.append('destination_radius', query.destination_radius.toString())

  // Single departure date (NOT a range)
  if (query.departure_date) params.append('departure_date', query.departure_date)

  // Numeric filters
  if (query.min_seats) params.append('min_seats', query.min_seats.toString())
  if (query.max_price) params.append('max_price', query.max_price.toString())
  if (query.min_driver_rating) params.append('min_driver_rating', query.min_driver_rating.toString())

  // Boolean filters
  if (query.pets_allowed !== undefined) params.append('pets_allowed', query.pets_allowed.toString())
  if (query.smoking_allowed !== undefined) params.append('smoking_allowed', query.smoking_allowed.toString())
  if (query.music_allowed !== undefined) params.append('music_allowed', query.music_allowed.toString())

  // Full-text search
  if (query.q) params.append('q', query.q)

  // NEW: Flexible sorting
  if (query.sort_by) params.append('sort_by', query.sort_by)
  if (query.sort_order) params.append('sort_order', query.sort_order)

  // Pagination
  if (query.page) params.append('page', query.page.toString())
  if (query.limit) params.append('limit', query.limit.toString())

  const queryString = params.toString()
  const url = `${SEARCH_BASE}/trips${queryString ? `?${queryString}` : ''}`

  const response = await apiClient.get<ApiResponse<SearchResponse<SearchTrip>>>(url)
  return response.data.data!
}

/**
 * Search trips by geospatial coordinates
 * Public endpoint - finds trips within radius of given coordinates
 * @param params - Latitude, longitude, radius, and optional filters
 * @returns Paginated search results
 * @example
 * const results = await searchService.searchByLocation({
 *   lat: 4.7110,
 *   lng: -74.0721,
 *   radius_km: 10,
 *   min_seats: 2,
 *   max_price: 50000
 * })
 */
export async function searchByLocation(
  params: GeospatialSearchParams
): Promise<SearchResponse<SearchTrip>> {
  const queryParams = new URLSearchParams()

  queryParams.append('lat', params.lat.toString())
  queryParams.append('lng', params.lng.toString())
  queryParams.append('radius_km', params.radius_km.toString())

  if (params.min_seats) queryParams.append('min_seats', params.min_seats.toString())
  if (params.max_price) queryParams.append('max_price', params.max_price.toString())
  if (params.page) queryParams.append('page', params.page.toString())
  if (params.limit) queryParams.append('limit', params.limit.toString())

  const url = `${SEARCH_BASE}/location?${queryParams.toString()}`

  const response = await apiClient.get<ApiResponse<SearchResponse<SearchTrip>>>(url)
  return response.data.data!
}

/**
 * Autocomplete city/location search
 * Public endpoint - returns suggestions for location input
 * @param searchQuery - Partial city/location name (min 2 characters)
 * @param limit - Max number of suggestions (default: 10, max: 50)
 * @returns Array of autocomplete suggestions
 * @example
 * const suggestions = await searchService.autocomplete('bog', 10)
 * // Returns: ['Bogotá', 'Bogotá D.C.', ...]
 */
export async function autocomplete(
  searchQuery: string,
  limit = 10
): Promise<string[]> {
  const params = new URLSearchParams()
  params.append('q', searchQuery)
  params.append('limit', limit.toString())

  const url = `${SEARCH_BASE}/autocomplete?${params.toString()}`

  const response = await apiClient.get<ApiResponse<string[]>>(url)
  return response.data.data!
}

/**
 * Get popular/trending routes
 * Public endpoint - returns most searched origin-destination pairs
 * @param limit - Max number of routes to return (default: 10, max: 50)
 * @returns Array of popular routes with search counts
 * @example
 * const routes = await searchService.getPopularRoutes(10)
 */
export async function getPopularRoutes(limit = 10): Promise<PopularRoute[]> {
  const params = new URLSearchParams()
  params.append('limit', limit.toString())

  const url = `${SEARCH_BASE}/popular-routes?${params.toString()}`

  const response = await apiClient.get<ApiResponse<PopularRoute[]>>(url)
  return response.data.data!
}

/**
 * Get trip details from search API
 * Public endpoint - returns denormalized trip with driver info
 * Note: This is different from trips-api GET /trips/:id
 * @param id - Trip ID (MongoDB ObjectID)
 * @returns SearchTrip object with denormalized driver data
 * @example
 * const trip = await searchService.getTripDetails('507f1f77bcf86cd799439011')
 */
export async function getTripDetails(id: string): Promise<SearchTrip> {
  // Note: This uses /api/v1/trips/:id which routes to search-api
  const response = await apiClient.get<ApiResponse<SearchTrip>>(`/trips/${id}`)
  return response.data.data!
}

export default {
  searchTrips,
  searchByLocation,
  autocomplete,
  getPopularRoutes,
  getTripDetails,
}
