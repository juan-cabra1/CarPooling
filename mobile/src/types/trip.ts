/**
 * Trip-related types for the CarPooling Trips API (Port 8002)
 * Database: MongoDB
 */

/**
 * Trip status values
 */
export type TripStatus =
  | 'draft' // Trip created but not published
  | 'published' // Trip visible and accepting reservations
  | 'full' // All seats reserved
  | 'in_progress' // Trip has started
  | 'completed' // Trip finished successfully
  | 'cancelled' // Trip cancelled

/**
 * Geographic coordinates (latitude and longitude)
 */
export interface Coordinates {
  lat: number // Latitude
  lng: number // Longitude
}

/**
 * Location with city, province, address and coordinates
 */
export interface Location {
  city: string
  province: string
  address: string
  coordinates: Coordinates
}

/**
 * Vehicle information
 */
export interface Car {
  brand: string
  model: string
  year: number
  color: string
  plate: string
}

/**
 * Trip preferences for passengers
 */
export interface Preferences {
  pets_allowed: boolean
  smoking_allowed: boolean
  music_allowed: boolean
}

/**
 * Complete trip object
 */
export interface Trip {
  id: string // MongoDB ObjectID as hex string
  driver_id: number // User ID from users-api
  origin: Location
  destination: Location
  departure_datetime: string // ISO 8601 datetime
  estimated_arrival_datetime: string // ISO 8601 datetime
  price_per_seat: number
  total_seats: number
  reserved_seats: number
  available_seats: number
  availability_version: number // For optimistic locking
  car: Car
  preferences: Preferences
  status: TripStatus
  description: string
  cancelled_at?: string // ISO 8601 datetime, only if status is 'cancelled'
  cancelled_by?: number // User ID who cancelled
  cancellation_reason?: string
  created_at: string // ISO 8601 datetime
  updated_at: string // ISO 8601 datetime
}

/**
 * Data for creating a new trip
 */
export interface CreateTripRequest {
  origin: Location
  destination: Location
  departure_datetime: string // ISO 8601 datetime (RFC3339)
  estimated_arrival_datetime: string // ISO 8601 datetime (RFC3339)
  price_per_seat: number // Must be >= 0
  total_seats: number // 1-8 seats
  car: Car
  preferences: Preferences
  description?: string
}

/**
 * Data for updating an existing trip (all fields optional)
 * Note: Cannot update if trip has reservations
 */
export interface UpdateTripRequest {
  origin?: Location
  destination?: Location
  departure_datetime?: string // ISO 8601 datetime
  estimated_arrival_datetime?: string // ISO 8601 datetime
  price_per_seat?: number
  total_seats?: number
  car?: Car
  preferences?: Preferences
  description?: string
}

/**
 * Cancel trip request
 */
export interface CancelTripRequest {
  reason: string
}

/**
 * Query filters for listing trips
 */
export interface TripFilters {
  driver_id?: number
  status?: TripStatus
  origin_city?: string
  destination_city?: string
  page?: number // Default: 1
  limit?: number // Default: 10, max: 100
}

/**
 * Paginated trip list response
 */
export interface TripListResponse {
  trips: Trip[]
  total: number
  page: number
  limit: number
}
