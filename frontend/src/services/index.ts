/**
 * Central export point for all CarPooling API services
 * Provides convenient barrel exports for importing services throughout the application
 */

// Export base API client and utilities
export { default as apiClient, getErrorMessage } from './api'

// Export authentication service
export { default as authService } from './authService'
export * from './authService'

// Export trips service
export { default as tripsService } from './tripsService'
export * from './tripsService'

// Export bookings service
export { default as bookingsService } from './bookingsService'
export * from './bookingsService'

// Export search service
export { default as searchService } from './searchService'
export * from './searchService'

// Export ratings service
export { default as ratingsService } from './ratingsService'
export * from './ratingsService'
