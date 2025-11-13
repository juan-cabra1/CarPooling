// User Types
export interface User {
  id: number;
  email: string;
  email_verified: boolean;
  name: string;
  lastname: string;
  role: string;
  phone: string;
  street: string;
  number: number;
  photo_url?: string;
  sex: string;
  avg_driver_rating: number;
  avg_passenger_rating: number;
  total_trips_passenger: number;
  total_trips_driver: number;
  birthdate: string;
  created_at: string;
  updated_at: string;
}

export interface AuthUser extends User {
  token: string;
}

// Trip Types
export interface Trip {
  id: string;
  driver_id: string;
  origin: Location;
  destination: Location;
  departure_time: string;
  available_seats: number;
  price_per_seat: number;
  status: TripStatus;
  created_at: string;
  updated_at: string;
}

export interface Location {
  address: string;
  latitude: number;
  longitude: number;
  city?: string;
  state?: string;
  country?: string;
}

export type TripStatus = 'active' | 'completed' | 'cancelled';

// Booking Types
export interface Booking {
  id: string;
  trip_id: string;
  passenger_id: string;
  seats_booked: number;
  total_price: number;
  status: BookingStatus;
  created_at: string;
  updated_at: string;
}

export type BookingStatus = 'pending' | 'confirmed' | 'cancelled' | 'completed';

// Search Types
export interface SearchFilters {
  origin?: string;
  destination?: string;
  departure_date?: string;
  min_seats?: number;
  max_price?: number;
}

export interface SearchResult {
  trip: Trip;
  driver: User;
  distance?: number;
}

// API Response Types
export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface ApiError {
  message: string;
  code?: string;
  errors?: Record<string, string[]>;
}
