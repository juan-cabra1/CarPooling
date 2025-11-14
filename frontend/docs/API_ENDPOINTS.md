# CarPooling Backend API Documentation

## Overview

This document provides comprehensive documentation for all backend APIs in the CarPooling microservices architecture.

### Architecture

The backend consists of 4 independent microservices:

| Service | Port | Database | Purpose |
|---------|------|----------|---------|
| **users-api** | 8001 | MySQL | User management & authentication |
| **trips-api** | 8002 | MongoDB | Trip creation & management |
| **bookings-api** | 8003 | MySQL | Booking/reservation management |
| **search-api** | 8004 | MongoDB + Solr | Advanced trip search & recommendations |

### Authentication

All protected endpoints require JWT authentication:
- Login via `POST /login` returns a JWT token
- Include token in request headers: `Authorization: Bearer <token>`
- Token contains: `user_id`, `email`, `is_verified`

### Common Response Format

All endpoints follow this response structure:

**Success Response:**
```typescript
{
  success: true,
  data: T  // Type varies by endpoint
}
```

**Error Response:**
```typescript
{
  success: false,
  error: {
    code: string,        // e.g., "USER_NOT_FOUND"
    message: string,     // Human-readable error message
    details?: any        // Optional additional error details
  }
}
```

---

## 1. Users API (Port 8001)

**Base URL:** `http://localhost:8001`

### 1.1 Public Endpoints

#### POST /users
Register a new user account.

**Authentication:** Not required

**Request Body:**
```typescript
{
  email: string,              // Valid email format
  password: string,           // Min 6 characters
  first_name: string,         // Required
  last_name: string,          // Required
  phone_number: string,       // Required
  date_of_birth: string,      // ISO date format
  gender?: "male" | "female" | "other",
  profile_picture_url?: string,
  bio?: string,
  preferences?: {
    smoking_allowed: boolean,
    pets_allowed: boolean,
    music_allowed: boolean,
    conversation_level: "quiet" | "moderate" | "chatty"
  }
}
```

**Response (201):**
```typescript
{
  success: true,
  data: {
    id: string,
    email: string,
    first_name: string,
    last_name: string,
    phone_number: string,
    date_of_birth: string,
    gender: string,
    is_verified: boolean,
    is_active: boolean,
    profile_picture_url: string,
    bio: string,
    average_rating: number,
    total_trips_as_driver: number,
    total_trips_as_passenger: number,
    preferences: {
      smoking_allowed: boolean,
      pets_allowed: boolean,
      music_allowed: boolean,
      conversation_level: string
    },
    created_at: string,
    updated_at: string
  }
}
```

**Status Codes:**
- `201`: User created successfully
- `400`: Validation error or email already exists
- `500`: Internal server error

---

#### POST /login
Authenticate user and receive JWT token.

**Authentication:** Not required

**Request Body:**
```typescript
{
  email: string,
  password: string
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    token: string,  // JWT token
    user: {
      id: string,
      email: string,
      first_name: string,
      last_name: string,
      is_verified: boolean,
      profile_picture_url: string
    }
  }
}
```

**Status Codes:**
- `200`: Login successful
- `400`: Invalid credentials
- `500`: Internal server error

---

#### GET /verify-email
Verify user email with token.

**Authentication:** Not required

**Query Parameters:**
```typescript
{
  token: string  // Email verification token
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    message: "Email verified successfully"
  }
}
```

**Status Codes:**
- `200`: Email verified
- `400`: Invalid or expired token
- `500`: Internal server error

---

#### POST /resend-verification
Resend email verification link.

**Authentication:** Not required

**Request Body:**
```typescript
{
  email: string
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    message: "Verification email sent"
  }
}
```

**Status Codes:**
- `200`: Email sent
- `400`: User not found or already verified
- `500`: Internal server error

---

#### POST /forgot-password
Request password reset email.

**Authentication:** Not required

**Request Body:**
```typescript
{
  email: string
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    message: "Password reset email sent"
  }
}
```

**Status Codes:**
- `200`: Email sent
- `400`: User not found
- `500`: Internal server error

---

#### POST /reset-password
Reset password using token from email.

**Authentication:** Not required

**Request Body:**
```typescript
{
  token: string,      // Password reset token
  password: string    // New password (min 6 chars)
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    message: "Password reset successfully"
  }
}
```

**Status Codes:**
- `200`: Password reset successful
- `400`: Invalid or expired token, or validation error
- `500`: Internal server error

---

#### GET /health
Health check endpoint.

**Authentication:** Not required

**Response (200):**
```typescript
{
  status: "healthy",
  timestamp: string
}
```

**Status Codes:**
- `200`: Service is healthy

---

### 1.2 Protected Endpoints (JWT Required)

#### GET /users/me
Get authenticated user's profile.

**Authentication:** Required (JWT)

**Response (200):**
```typescript
{
  success: true,
  data: {
    id: string,
    email: string,
    first_name: string,
    last_name: string,
    phone_number: string,
    date_of_birth: string,
    gender: string,
    is_verified: boolean,
    is_active: boolean,
    profile_picture_url: string,
    bio: string,
    average_rating: number,
    total_trips_as_driver: number,
    total_trips_as_passenger: number,
    preferences: {
      smoking_allowed: boolean,
      pets_allowed: boolean,
      music_allowed: boolean,
      conversation_level: string
    },
    created_at: string,
    updated_at: string
  }
}
```

**Status Codes:**
- `200`: Success
- `401`: Unauthorized (invalid/missing token)
- `404`: User not found
- `500`: Internal server error

---

#### GET /users/:id
Get user profile by ID.

**Authentication:** Required (JWT)

**Path Parameters:**
- `id`: User UUID

**Response (200):**
```typescript
{
  success: true,
  data: {
    id: string,
    email: string,
    first_name: string,
    last_name: string,
    phone_number: string,
    profile_picture_url: string,
    bio: string,
    average_rating: number,
    total_trips_as_driver: number,
    total_trips_as_passenger: number,
    preferences: {
      smoking_allowed: boolean,
      pets_allowed: boolean,
      music_allowed: boolean,
      conversation_level: string
    },
    created_at: string
  }
}
```

**Status Codes:**
- `200`: Success
- `401`: Unauthorized
- `404`: User not found
- `500`: Internal server error

---

#### PUT /users/:id
Update user profile (owner only).

**Authentication:** Required (JWT, must be profile owner)

**Path Parameters:**
- `id`: User UUID

**Request Body:**
```typescript
{
  first_name?: string,
  last_name?: string,
  phone_number?: string,
  date_of_birth?: string,
  gender?: "male" | "female" | "other",
  profile_picture_url?: string,
  bio?: string,
  preferences?: {
    smoking_allowed?: boolean,
    pets_allowed?: boolean,
    music_allowed?: boolean,
    conversation_level?: "quiet" | "moderate" | "chatty"
  }
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    // Updated user object (same as GET /users/me)
  }
}
```

**Status Codes:**
- `200`: Update successful
- `400`: Validation error
- `401`: Unauthorized
- `403`: Forbidden (not profile owner)
- `404`: User not found
- `500`: Internal server error

---

#### DELETE /users/:id
Delete user account (owner only).

**Authentication:** Required (JWT, must be profile owner)

**Path Parameters:**
- `id`: User UUID

**Response (200):**
```typescript
{
  success: true,
  data: {
    message: "User deleted successfully"
  }
}
```

**Status Codes:**
- `200`: Deletion successful
- `401`: Unauthorized
- `403`: Forbidden (not profile owner)
- `404`: User not found
- `500`: Internal server error

---

#### POST /change-password
Change user password (authenticated user).

**Authentication:** Required (JWT)

**Request Body:**
```typescript
{
  current_password: string,
  new_password: string  // Min 6 characters
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    message: "Password changed successfully"
  }
}
```

**Status Codes:**
- `200`: Password changed
- `400`: Invalid current password or validation error
- `401`: Unauthorized
- `500`: Internal server error

---

#### GET /users/:id/ratings
Get ratings for a user (paginated).

**Authentication:** Required (JWT)

**Path Parameters:**
- `id`: User UUID

**Query Parameters:**
```typescript
{
  page?: number,    // Default: 1
  limit?: number    // Default: 10
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    ratings: Array<{
      id: string,
      rated_user_id: string,
      rating_user_id: string,
      trip_id: string,
      rating: number,        // 1-5
      comment: string,
      created_at: string,
      rating_user: {
        id: string,
        first_name: string,
        last_name: string,
        profile_picture_url: string
      }
    }>,
    pagination: {
      current_page: number,
      total_pages: number,
      total_items: number,
      items_per_page: number
    }
  }
}
```

**Status Codes:**
- `200`: Success
- `401`: Unauthorized
- `404`: User not found
- `500`: Internal server error

---

### 1.3 Internal Endpoints

#### POST /internal/ratings
Create a rating (called by trips-api after trip completion).

**Authentication:** Internal service call

**Request Body:**
```typescript
{
  rated_user_id: string,
  rating_user_id: string,
  trip_id: string,
  rating: number,        // 1-5
  comment?: string
}
```

**Response (201):**
```typescript
{
  success: true,
  data: {
    id: string,
    rated_user_id: string,
    rating_user_id: string,
    trip_id: string,
    rating: number,
    comment: string,
    created_at: string
  }
}
```

**Status Codes:**
- `201`: Rating created
- `400`: Validation error
- `500`: Internal server error

---

## 2. Trips API (Port 8002)

**Base URL:** `http://localhost:8002`

### 2.1 Public Endpoints

#### GET /trips
List trips with optional filters.

**Authentication:** Not required

**Query Parameters:**
```typescript
{
  driver_id?: string,           // Filter by driver UUID
  status?: "scheduled" | "in_progress" | "completed" | "cancelled",
  origin_city?: string,         // Filter by origin city name
  destination_city?: string,    // Filter by destination city name
  page?: number,                // Default: 1
  limit?: number                // Default: 10
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    trips: Array<{
      id: string,
      driver_id: string,
      origin: {
        address: string,
        city: string,
        state: string,
        country: string,
        postal_code: string,
        coordinates: {
          latitude: number,
          longitude: number
        }
      },
      destination: {
        address: string,
        city: string,
        state: string,
        country: string,
        postal_code: string,
        coordinates: {
          latitude: number,
          longitude: number
        }
      },
      departure_time: string,     // ISO datetime
      arrival_time: string,       // ISO datetime
      available_seats: number,
      price_per_seat: number,
      status: "scheduled" | "in_progress" | "completed" | "cancelled",
      car: {
        make: string,
        model: string,
        year: number,
        color: string,
        license_plate: string
      },
      preferences: {
        smoking_allowed: boolean,
        pets_allowed: boolean,
        music_allowed: boolean,
        max_two_in_back: boolean
      },
      stops: Array<{
        location: {
          address: string,
          city: string,
          coordinates: {
            latitude: number,
            longitude: number
          }
        },
        estimated_arrival: string
      }>,
      notes?: string,
      created_at: string,
      updated_at: string
    }>,
    pagination: {
      current_page: number,
      total_pages: number,
      total_items: number,
      items_per_page: number
    }
  }
}
```

**Status Codes:**
- `200`: Success
- `400`: Invalid query parameters
- `500`: Internal server error

---

#### GET /trips/:id
Get trip details by ID.

**Authentication:** Not required

**Path Parameters:**
- `id`: Trip UUID

**Response (200):**
```typescript
{
  success: true,
  data: {
    // Same trip object as in GET /trips array
  }
}
```

**Status Codes:**
- `200`: Success
- `404`: Trip not found
- `500`: Internal server error

---

#### GET /health
Health check endpoint.

**Authentication:** Not required

**Response (200):**
```typescript
{
  status: "healthy",
  timestamp: string
}
```

**Status Codes:**
- `200`: Service is healthy

---

### 2.2 Protected Endpoints (JWT Required)

#### POST /trips
Create a new trip.

**Authentication:** Required (JWT)

**Request Body:**
```typescript
{
  origin: {
    address: string,
    city: string,
    state: string,
    country: string,
    postal_code: string,
    coordinates: {
      latitude: number,
      longitude: number
    }
  },
  destination: {
    address: string,
    city: string,
    state: string,
    country: string,
    postal_code: string,
    coordinates: {
      latitude: number,
      longitude: number
    }
  },
  departure_time: string,     // ISO datetime (future)
  available_seats: number,    // Min: 1
  price_per_seat: number,     // Min: 0
  car: {
    make: string,
    model: string,
    year: number,
    color: string,
    license_plate: string
  },
  preferences?: {
    smoking_allowed?: boolean,
    pets_allowed?: boolean,
    music_allowed?: boolean,
    max_two_in_back?: boolean
  },
  stops?: Array<{
    location: {
      address: string,
      city: string,
      coordinates: {
        latitude: number,
        longitude: number
      }
    },
    estimated_arrival: string  // ISO datetime
  }>,
  notes?: string
}
```

**Response (201):**
```typescript
{
  success: true,
  data: {
    // Complete trip object (same as GET /trips/:id)
  }
}
```

**Status Codes:**
- `201`: Trip created successfully
- `400`: Validation error (e.g., departure_time in past, invalid coordinates)
- `401`: Unauthorized
- `500`: Internal server error

---

#### PUT /trips/:id
Update trip (owner only).

**Authentication:** Required (JWT, must be trip creator)

**Path Parameters:**
- `id`: Trip UUID

**Request Body:**
```typescript
{
  // All fields from POST /trips are optional here
  available_seats?: number,
  price_per_seat?: number,
  departure_time?: string,
  preferences?: {
    smoking_allowed?: boolean,
    pets_allowed?: boolean,
    music_allowed?: boolean,
    max_two_in_back?: boolean
  },
  notes?: string
  // Note: Cannot update if status is "completed" or "cancelled"
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    // Updated trip object
  }
}
```

**Status Codes:**
- `200`: Update successful
- `400`: Validation error or trip cannot be modified
- `401`: Unauthorized
- `403`: Forbidden (not trip owner)
- `404`: Trip not found
- `500`: Internal server error

---

#### PATCH /trips/:id
Partial update trip (owner only).

**Authentication:** Required (JWT, must be trip creator)

**Path Parameters:**
- `id`: Trip UUID

**Request Body:**
```typescript
{
  // Same as PUT, accepts partial updates
  status?: "scheduled" | "in_progress" | "completed" | "cancelled"
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    // Updated trip object
  }
}
```

**Status Codes:**
- `200`: Update successful
- `400`: Validation error
- `401`: Unauthorized
- `403`: Forbidden (not trip owner)
- `404`: Trip not found
- `500`: Internal server error

---

#### DELETE /trips/:id
Delete trip (owner only, only if no bookings).

**Authentication:** Required (JWT, must be trip creator)

**Path Parameters:**
- `id`: Trip UUID

**Response (200):**
```typescript
{
  success: true,
  data: {
    message: "Trip deleted successfully"
  }
}
```

**Status Codes:**
- `200`: Deletion successful
- `400`: Cannot delete (has active bookings)
- `401`: Unauthorized
- `403`: Forbidden (not trip owner)
- `404`: Trip not found
- `500`: Internal server error

---

## 3. Bookings API (Port 8003)

**Base URL:** `http://localhost:8003`

### 3.1 Public Endpoints

#### GET /health
Health check endpoint.

**Authentication:** Not required

**Response (200):**
```typescript
{
  status: "healthy",
  timestamp: string
}
```

**Status Codes:**
- `200`: Service is healthy

---

### 3.2 Protected Endpoints (JWT Required)

#### POST /api/v1/bookings
Create a new booking.

**Authentication:** Required (JWT)

**Request Body:**
```typescript
{
  trip_id: string,           // UUID of the trip
  seats_reserved: number,    // Number of seats to book (min: 1)
  pickup_location?: {
    address: string,
    city: string,
    coordinates: {
      latitude: number,
      longitude: number
    }
  },
  dropoff_location?: {
    address: string,
    city: string,
    coordinates: {
      latitude: number,
      longitude: number
    }
  },
  notes?: string
}
```

**Response (201):**
```typescript
{
  success: true,
  data: {
    id: string,
    trip_id: string,
    passenger_id: string,
    seats_reserved: number,
    total_price: number,
    status: "pending" | "confirmed" | "cancelled" | "completed",
    pickup_location?: {
      address: string,
      city: string,
      coordinates: {
        latitude: number,
        longitude: number
      }
    },
    dropoff_location?: {
      address: string,
      city: string,
      coordinates: {
        latitude: number,
        longitude: number
      }
    },
    notes?: string,
    created_at: string,
    updated_at: string
  }
}
```

**Status Codes:**
- `201`: Booking created successfully
- `400`: Validation error (not enough seats, trip not available, etc.)
- `401`: Unauthorized
- `404`: Trip not found
- `500`: Internal server error

---

#### GET /api/v1/bookings
List user's bookings (paginated).

**Authentication:** Required (JWT)

**Query Parameters:**
```typescript
{
  status?: "pending" | "confirmed" | "cancelled" | "completed",
  page?: number,    // Default: 1
  limit?: number    // Default: 10
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    bookings: Array<{
      id: string,
      trip_id: string,
      passenger_id: string,
      seats_reserved: number,
      total_price: number,
      status: "pending" | "confirmed" | "cancelled" | "completed",
      pickup_location?: {
        address: string,
        city: string,
        coordinates: {
          latitude: number,
          longitude: number
        }
      },
      dropoff_location?: {
        address: string,
        city: string,
        coordinates: {
          latitude: number,
          longitude: number
        }
      },
      notes?: string,
      created_at: string,
      updated_at: string
    }>,
    pagination: {
      current_page: number,
      total_pages: number,
      total_items: number,
      items_per_page: number
    }
  }
}
```

**Status Codes:**
- `200`: Success
- `401`: Unauthorized
- `500`: Internal server error

---

#### GET /api/v1/bookings/:id
Get booking details by ID (owner only).

**Authentication:** Required (JWT, must be booking owner)

**Path Parameters:**
- `id`: Booking UUID

**Response (200):**
```typescript
{
  success: true,
  data: {
    // Same booking object as in GET /api/v1/bookings array
  }
}
```

**Status Codes:**
- `200`: Success
- `401`: Unauthorized
- `403`: Forbidden (not booking owner)
- `404`: Booking not found
- `500`: Internal server error

---

#### PATCH /api/v1/bookings/:id/cancel
Cancel a booking (owner only).

**Authentication:** Required (JWT, must be booking owner)

**Path Parameters:**
- `id`: Booking UUID

**Request Body:** None required

**Response (200):**
```typescript
{
  success: true,
  data: {
    id: string,
    status: "cancelled",
    // ... rest of booking object
  }
}
```

**Status Codes:**
- `200`: Cancellation successful
- `400`: Booking cannot be cancelled (already completed/cancelled)
- `401`: Unauthorized
- `403`: Forbidden (not booking owner)
- `404`: Booking not found
- `500`: Internal server error

---

## 4. Search API (Port 8004)

**Base URL:** `http://localhost:8004`

### 4.1 Public Endpoints

#### GET /health
Health check endpoint.

**Authentication:** Not required

**Response (200):**
```typescript
{
  status: "healthy",
  timestamp: string
}
```

**Status Codes:**
- `200`: Service is healthy

---

#### GET /api/v1/search/trips
Advanced search for trips with multiple filters.

**Authentication:** Not required

**Query Parameters:**
```typescript
{
  origin_city?: string,
  destination_city?: string,
  departure_date?: string,       // ISO date (YYYY-MM-DD)
  min_seats?: number,
  max_price?: number,
  smoking_allowed?: boolean,
  pets_allowed?: boolean,
  music_allowed?: boolean,
  page?: number,                 // Default: 1
  limit?: number,                // Default: 10
  sort_by?: "price" | "departure_time" | "rating",
  sort_order?: "asc" | "desc"
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    trips: Array<{
      // Same trip object structure as trips-api
      id: string,
      driver_id: string,
      origin: { /* Location object */ },
      destination: { /* Location object */ },
      departure_time: string,
      available_seats: number,
      price_per_seat: number,
      status: string,
      car: { /* Car object */ },
      preferences: { /* Preferences object */ },
      // Additional search metadata
      relevance_score?: number
    }>,
    pagination: {
      current_page: number,
      total_pages: number,
      total_items: number,
      items_per_page: number
    },
    filters_applied: {
      origin_city?: string,
      destination_city?: string,
      departure_date?: string,
      // ... echoes applied filters
    }
  }
}
```

**Status Codes:**
- `200`: Success
- `400`: Invalid query parameters
- `500`: Internal server error

---

#### GET /api/v1/search/location
Geospatial search for trips near coordinates.

**Authentication:** Not required

**Query Parameters:**
```typescript
{
  latitude: number,          // Required
  longitude: number,         // Required
  radius_km?: number,        // Default: 50km
  search_type?: "origin" | "destination" | "both",  // Default: "both"
  page?: number,             // Default: 1
  limit?: number             // Default: 10
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    trips: Array<{
      // Same trip object as search/trips
      distance_km: number  // Distance from search point
    }>,
    search_center: {
      latitude: number,
      longitude: number
    },
    radius_km: number,
    pagination: {
      current_page: number,
      total_pages: number,
      total_items: number,
      items_per_page: number
    }
  }
}
```

**Status Codes:**
- `200`: Success
- `400`: Invalid coordinates or parameters
- `500`: Internal server error

---

#### GET /api/v1/search/autocomplete
City autocomplete suggestions for search input.

**Authentication:** Not required

**Query Parameters:**
```typescript
{
  query: string,    // Search term (min 2 characters)
  limit?: number    // Default: 10, max: 20
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    suggestions: Array<{
      city: string,
      state: string,
      country: string,
      trip_count: number  // Number of trips with this city
    }>
  }
}
```

**Status Codes:**
- `200`: Success
- `400`: Invalid query (too short)
- `500`: Internal server error

---

#### GET /api/v1/search/popular-routes
Get trending/popular routes.

**Authentication:** Not required

**Query Parameters:**
```typescript
{
  limit?: number    // Default: 10, max: 50
}
```

**Response (200):**
```typescript
{
  success: true,
  data: {
    routes: Array<{
      origin_city: string,
      destination_city: string,
      trip_count: number,
      avg_price: number,
      route_id: string
    }>
  }
}
```

**Status Codes:**
- `200`: Success
- `500`: Internal server error

---

#### GET /api/v1/trips/:id
Get trip details (mirrors trips-api with caching).

**Authentication:** Not required

**Path Parameters:**
- `id`: Trip UUID

**Response (200):**
```typescript
{
  success: true,
  data: {
    // Same trip object as trips-api GET /trips/:id
  }
}
```

**Status Codes:**
- `200`: Success
- `404`: Trip not found
- `500`: Internal server error

---

## 5. Common Models

### Location
```typescript
interface Location {
  address: string;
  city: string;
  state: string;
  country: string;
  postal_code: string;
  coordinates: {
    latitude: number;
    longitude: number;
  };
}
```

### Car
```typescript
interface Car {
  make: string;
  model: string;
  year: number;
  color: string;
  license_plate: string;
}
```

### Preferences
```typescript
interface Preferences {
  smoking_allowed: boolean;
  pets_allowed: boolean;
  music_allowed: boolean;
  conversation_level?: "quiet" | "moderate" | "chatty";
  max_two_in_back?: boolean;
}
```

### Pagination
```typescript
interface Pagination {
  current_page: number;
  total_pages: number;
  total_items: number;
  items_per_page: number;
}
```

---

## 6. Error Handling

### Standard Error Codes

**Users API:**
- `USER_NOT_FOUND`: User does not exist
- `USER_ALREADY_EXISTS`: Email already registered
- `INVALID_CREDENTIALS`: Wrong email/password
- `EMAIL_NOT_VERIFIED`: User hasn't verified email
- `UNAUTHORIZED`: Missing or invalid JWT token
- `FORBIDDEN`: User doesn't have permission
- `VALIDATION_ERROR`: Request validation failed

**Trips API:**
- `TRIP_NOT_FOUND`: Trip does not exist
- `INVALID_TRIP_DATA`: Validation error
- `TRIP_NOT_MODIFIABLE`: Trip status prevents modification
- `UNAUTHORIZED`: Missing or invalid JWT token
- `FORBIDDEN`: User is not trip owner

**Bookings API:**
- `BOOKING_NOT_FOUND`: Booking does not exist
- `INSUFFICIENT_SEATS`: Not enough seats available
- `TRIP_NOT_AVAILABLE`: Trip cannot be booked
- `BOOKING_NOT_MODIFIABLE`: Booking status prevents changes
- `UNAUTHORIZED`: Missing or invalid JWT token
- `FORBIDDEN`: User is not booking owner

**Search API:**
- `INVALID_SEARCH_PARAMS`: Invalid search parameters
- `SEARCH_SERVICE_ERROR`: Solr or cache error

### Error Response Example
```typescript
{
  success: false,
  error: {
    code: "USER_NOT_FOUND",
    message: "No user found with the provided ID",
    details: {
      user_id: "123e4567-e89b-12d3-a456-426614174000"
    }
  }
}
```

---

## 7. Frontend Integration Guide

### 7.1 Authentication Flow

1. **Register:** `POST /users` → Store user data
2. **Login:** `POST /login` → Store JWT token
3. **Store Token:** Save to localStorage/sessionStorage
4. **Include in Requests:** Add `Authorization: Bearer <token>` header
5. **Handle 401:** Redirect to login page

```typescript
// Example: Store token
localStorage.setItem('auth_token', data.token);

// Example: Include in requests
headers: {
  'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
}
```

### 7.2 Pagination Handling

All paginated endpoints accept `page` and `limit` parameters and return consistent pagination metadata.

```typescript
// Example: Fetch next page
const nextPage = pagination.current_page + 1;
if (nextPage <= pagination.total_pages) {
  fetchTrips({ page: nextPage, limit: 10 });
}
```

### 7.3 Error Handling Best Practices

```typescript
try {
  const response = await fetch(url, options);
  const data = await response.json();

  if (!data.success) {
    // Handle API error
    switch (data.error.code) {
      case 'USER_NOT_FOUND':
        showNotification('User not found');
        break;
      case 'UNAUTHORIZED':
        redirectToLogin();
        break;
      default:
        showNotification(data.error.message);
    }
  }

  return data.data;
} catch (error) {
  // Handle network error
  showNotification('Network error. Please try again.');
}
```

### 7.4 Search Query Building

```typescript
// Example: Build search query
const searchParams = new URLSearchParams({
  origin_city: 'Bogotá',
  destination_city: 'Medellín',
  departure_date: '2024-03-15',
  min_seats: '2',
  max_price: '50000',
  page: '1',
  limit: '10'
});

const url = `http://localhost:8004/api/v1/search/trips?${searchParams}`;
```

### 7.5 Date/Time Handling

All datetime fields use ISO 8601 format:
- Dates: `YYYY-MM-DD`
- DateTimes: `YYYY-MM-DDTHH:mm:ssZ`

```typescript
// Example: Format for API
const departureTime = new Date('2024-03-15T14:30:00').toISOString();

// Example: Parse from API
const date = new Date(trip.departure_time);
```

### 7.6 Recommended API Client Structure

```
src/
  services/
    api/
      users.api.ts       // Users API calls
      trips.api.ts       // Trips API calls
      bookings.api.ts    // Bookings API calls
      search.api.ts      // Search API calls
    auth.service.ts      // Authentication helpers
    http.service.ts      // Base HTTP client
  types/
    user.types.ts        // User interfaces
    trip.types.ts        // Trip interfaces
    booking.types.ts     // Booking interfaces
    common.types.ts      // Shared interfaces
```

---

## 8. Event-Driven Communication (RabbitMQ)

The backend services communicate via RabbitMQ events. While the frontend doesn't directly interact with these, understanding the flow helps with debugging.

### Events Published:

**trips-api:**
- `trip.created` → Consumed by search-api
- `trip.updated` → Consumed by search-api
- `trip.deleted` → Consumed by search-api
- `trip.completed` → Triggers rating flow

**bookings-api:**
- `booking.created` → Updates trip available_seats
- `booking.cancelled` → Updates trip available_seats

**users-api:**
- `user.verified` → May trigger welcome emails
- `rating.created` → Updates user average_rating

---

## Notes

- All IDs are UUIDs (UUID v4)
- All timestamps are in UTC
- Prices are in COP (Colombian Pesos) as integers
- Ratings are 1-5 scale (integers)
- All services implement idempotency checking
- CORS is enabled for frontend development
- Rate limiting may be implemented (check response headers)

---

**Last Updated:** 2024-03-13
**Version:** 1.0.0
