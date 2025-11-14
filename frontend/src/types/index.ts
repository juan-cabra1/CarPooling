/**
 * Central export point for all CarPooling TypeScript types
 * Provides convenient barrel exports for importing types throughout the application
 */

// API generic types
export type {
  ApiError,
  ApiResponse,
  PaginatedResponse,
  HealthCheckResponse,
  JwtPayload,
  RequestConfig,
  PaginationParams,
} from './api'

export { HttpStatus } from './api'

export type { ErrorCode } from './api'

// User types (users-api)
export type {
  UserRole,
  UserSex,
  RoleRated,
  User,
  LoginCredentials,
  RegisterData,
  UpdateUserData,
  AuthResponse,
  ChangePasswordRequest,
  ResetPasswordRequest,
  ForgotPasswordRequest,
  ResendVerificationRequest,
  Rating,
  CreateRatingRequest,
} from './user'

// Trip types (trips-api)
export type {
  TripStatus,
  Coordinates,
  Location,
  Car,
  Preferences,
  Trip,
  CreateTripRequest,
  UpdateTripRequest,
  CancelTripRequest,
  TripFilters,
  TripListResponse,
} from './trip'

// Booking types (bookings-api)
export type {
  BookingStatus,
  Booking,
  CreateBookingRequest,
  CancelBookingRequest,
  BookingListResponse,
  BookingFilters,
} from './booking'

// Search types (search-api)
export type {
  GeoJSONPoint,
  SearchLocation,
  Driver,
  SearchTrip,
  SearchSortBy,
  SearchQuery,
  SearchResponse,
  PopularRoute,
  AutocompleteSuggestion,
  GeospatialSearchParams,
} from './search'
