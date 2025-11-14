/**
 * Search-related types for the CarPooling Search API (Port 8004)
 * Database: MongoDB + Apache Solr
 * Features: Geospatial search, full-text search, caching
 */

import type { Car, Preferences, TripStatus } from './trip'

/**
 * GeoJSON Point format for MongoDB 2dsphere indexes
 * IMPORTANT: Coordinates are [longitude, latitude] - lng first!
 */
export interface GeoJSONPoint {
  type: 'Point'
  coordinates: [number, number] // [longitude, latitude] - lng first!
}

/**
 * Location with GeoJSON coordinates for geospatial queries
 */
export interface SearchLocation {
  city: string
  province: string
  address: string
  coordinates: GeoJSONPoint
}

/**
 * Denormalized driver information (from users-api)
 * Cached in search-api for performance
 */
export interface Driver {
  id: number
  name: string
  email: string
  photo_url?: string
  rating: number // avg_driver_rating from users-api
  total_trips: number // total_trips_driver from users-api
}

/**
 * Search trip object (denormalized for performance)
 * Includes driver info and trip details in single document
 */
export interface SearchTrip {
  id: string // MongoDB ObjectID
  trip_id: string // Original trip ID
  driver_id: number
  driver: Driver // Denormalized driver data
  origin: SearchLocation
  destination: SearchLocation
  departure_datetime: string // ISO 8601
  estimated_arrival_datetime: string // ISO 8601
  price_per_seat: number
  total_seats: number
  available_seats: number
  car: Car
  preferences: Preferences
  status: TripStatus
  description: string
  search_text?: string // Concatenated text for backup search
  popularity_score?: number // For ranking/sorting
  created_at: string // ISO 8601
  updated_at: string // ISO 8601
}

/**
 * Sort options for search results
 */
export type SearchSortBy =
  | 'popularity' // Most popular/booked trips
  | 'price_asc' // Cheapest first
  | 'price_desc' // Most expensive first
  | 'date_asc' // Earliest departure first
  | 'date_desc' // Latest departure first

/**
 * Search query parameters
 * Supports city-based search, geospatial search, and filters
 */
export interface SearchQuery {
  // City-based search
  origin_city?: string
  destination_city?: string

  // Geospatial search (MongoDB 2dsphere)
  origin_lat?: number
  origin_lng?: number
  origin_radius?: number // Radius in kilometers
  destination_lat?: number
  destination_lng?: number
  destination_radius?: number // Radius in kilometers

  // Date range filters
  date_from?: string // ISO 8601 date
  date_to?: string // ISO 8601 date

  // Filters (Solr)
  min_seats?: number
  max_price?: number
  pets_allowed?: boolean
  smoking_allowed?: boolean
  music_allowed?: boolean
  min_driver_rating?: number // 0-5

  // Full-text search
  q?: string // Search text (cities, descriptions)

  // Sorting and pagination
  sort_by?: SearchSortBy
  page?: number // Default: 1
  limit?: number // Default: 20, max: 100
}

/**
 * Paginated search response
 */
export interface SearchResponse<T = SearchTrip> {
  trips: T[]
  total: number
  page: number
  limit: number
  total_pages: number
}

/**
 * Popular route analytics
 */
export interface PopularRoute {
  id: string
  origin_city: string
  destination_city: string
  search_count: number
  last_searched: string // ISO 8601
}

/**
 * Autocomplete suggestion
 */
export interface AutocompleteSuggestion {
  city: string
  province: string
  match_count?: number // Number of trips matching this location
}

/**
 * Geospatial search parameters
 */
export interface GeospatialSearchParams {
  lat: number // Latitude
  lng: number // Longitude
  radius_km: number // Search radius in kilometers
  min_seats?: number
  max_price?: number
  page?: number
  limit?: number
}
