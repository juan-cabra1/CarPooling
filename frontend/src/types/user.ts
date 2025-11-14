/**
 * User-related types for the CarPooling Users API (Port 8001)
 * Database: MySQL
 */

/**
 * User role enum
 */
export type UserRole = 'user' | 'admin'

/**
 * User sex/gender options
 */
export type UserSex = 'hombre' | 'mujer' | 'otro'

/**
 * Role being rated in a trip (driver or passenger)
 */
export type RoleRated = 'conductor' | 'pasajero'

/**
 * Complete user object as returned by the API
 */
export interface User {
  id: number
  email: string
  email_verified: boolean
  name: string
  lastname: string
  role: UserRole
  phone: string
  street: string
  number: number
  photo_url?: string
  sex: UserSex
  avg_driver_rating: number
  avg_passenger_rating: number
  total_trips_passenger: number
  total_trips_driver: number
  birthdate: string // ISO 8601 date string
  created_at: string // ISO 8601 datetime string
  updated_at: string // ISO 8601 datetime string
}

/**
 * Login credentials for authentication
 */
export interface LoginCredentials {
  email: string
  password: string
}

/**
 * Registration data for new user signup
 */
export interface RegisterData {
  email: string
  password: string
  name: string
  lastname: string
  phone: string
  street: string
  number: number
  photo_url?: string
  sex: UserSex
  birthdate: string // Format: YYYY-MM-DD
}

/**
 * Update user profile data (all fields optional)
 */
export interface UpdateUserData {
  name?: string
  lastname?: string
  phone?: string
  street?: string
  number?: number
  photo_url?: string
}

/**
 * Authentication response with JWT token and user data
 */
export interface AuthResponse {
  token: string // JWT token
  user: User
}

/**
 * Change password request
 */
export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

/**
 * Reset password with token (from email)
 */
export interface ResetPasswordRequest {
  token: string
  new_password: string
}

/**
 * Request password reset email
 */
export interface ForgotPasswordRequest {
  email: string
}

/**
 * Resend email verification
 */
export interface ResendVerificationRequest {
  email: string
}

/**
 * User rating object
 */
export interface Rating {
  id: number
  rater_id: number // User who gave the rating
  rated_user_id: number // User receiving the rating
  trip_id: string // MongoDB ObjectID
  role_rated: RoleRated // Whether rating the driver or passenger
  score: number // 1-5 stars
  comment?: string
  created_at: string // ISO 8601 datetime
}

/**
 * Create a new rating for a user
 */
export interface CreateRatingRequest {
  rater_id: number
  rated_user_id: number
  trip_id: string
  role_rated: RoleRated
  score: number // 1-5
  comment?: string
}
