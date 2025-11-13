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
export interface Coordinates {
  lat: number;
  lng: number;
}

export interface Location {
  city: string;
  province: string;
  address: string;
  coordinates: Coordinates;
}

export interface Car {
  brand: string;
  model: string;
  year: number;
  color: string;
  plate: string;
}

export interface Preferences {
  pets_allowed: boolean;
  smoking_allowed: boolean;
  music_allowed: boolean;
}

export interface Trip {
  id: string;
  driver_id: number;
  origin: Location;
  destination: Location;
  departure_datetime: string;
  estimated_arrival_datetime: string;
  price_per_seat: number;
  total_seats: number;
  reserved_seats: number;
  available_seats: number;
  availability_version: number;
  car: Car;
  preferences: Preferences;
  status: TripStatus;
  description: string;
  cancelled_at?: string | null;
  cancelled_by?: number | null;
  cancellation_reason?: string;
  created_at: string;
  updated_at: string;
}

export type TripStatus = 'draft' | 'published' | 'full' | 'in_progress' | 'completed' | 'cancelled';

export interface CreateTripData {
  origin: Location;
  destination: Location;
  departure_datetime: string;
  estimated_arrival_datetime: string;
  price_per_seat: number;
  total_seats: number;
  car: Car;
  preferences?: Preferences;
  description?: string;
}

export interface UpdateTripData {
  origin?: Location;
  destination?: Location;
  departure_datetime?: string;
  estimated_arrival_datetime?: string;
  price_per_seat?: number;
  total_seats?: number;
  car?: Car;
  preferences?: Preferences;
  description?: string;
}

export interface TripFilters {
  driver_id?: number;
  status?: string;
  origin_city?: string;
  destination_city?: string;
  page?: number;
  limit?: number;
}

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
