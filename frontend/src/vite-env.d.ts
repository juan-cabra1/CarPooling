/// <reference types="vite/client" />

/**
 * Environment variable type definitions for Vite
 * Provides type safety for import.meta.env in the CarPooling frontend
 */

interface ImportMetaEnv {
  /**
   * Users API URL (Port 8001)
   * Handles authentication, user management, and ratings
   * @example 'http://localhost:8001'
   */
  readonly VITE_USERS_API_URL: string

  /**
   * Trips API URL (Port 8002)
   * Handles trip CRUD operations
   * @example 'http://localhost:8002'
   */
  readonly VITE_TRIPS_API_URL: string

  /**
   * Bookings API URL (Port 8003)
   * Handles reservation/booking management
   * @example 'http://localhost:8003'
   */
  readonly VITE_BOOKINGS_API_URL: string

  /**
   * Search API URL (Port 8004)
   * Handles trip search with filters and geospatial queries
   * @example 'http://localhost:8004'
   */
  readonly VITE_SEARCH_API_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
