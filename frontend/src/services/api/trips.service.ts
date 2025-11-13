import { apiClient } from './axios';
import type { Trip, ApiResponse, PaginatedResponse, Location } from '@/types';

export interface CreateTripData {
  origin: Location;
  destination: Location;
  departure_time: string;
  available_seats: number;
  price_per_seat: number;
}

export interface UpdateTripData {
  departure_time?: string;
  available_seats?: number;
  price_per_seat?: number;
  status?: 'active' | 'completed' | 'cancelled';
}

export interface TripFilters {
  driver_id?: string;
  status?: string;
  page?: number;
  per_page?: number;
}

export const tripsService = {
  // Create a new trip
  createTrip: async (data: CreateTripData): Promise<ApiResponse<Trip>> => {
    const response = await apiClient.post('/trips', data);
    return response.data;
  },

  // Get trip by ID
  getTripById: async (tripId: string): Promise<ApiResponse<Trip>> => {
    const response = await apiClient.get(`/trips/${tripId}`);
    return response.data;
  },

  // Get all trips with optional filters
  getTrips: async (filters?: TripFilters): Promise<PaginatedResponse<Trip>> => {
    const response = await apiClient.get('/trips', { params: filters });
    return response.data;
  },

  // Get trips by driver
  getTripsByDriver: async (driverId: string): Promise<ApiResponse<Trip[]>> => {
    const response = await apiClient.get(`/trips/driver/${driverId}`);
    return response.data;
  },

  // Update trip
  updateTrip: async (tripId: string, data: UpdateTripData): Promise<ApiResponse<Trip>> => {
    const response = await apiClient.patch(`/trips/${tripId}`, data);
    return response.data;
  },

  // Cancel trip
  cancelTrip: async (tripId: string): Promise<ApiResponse<Trip>> => {
    const response = await apiClient.patch(`/trips/${tripId}/cancel`);
    return response.data;
  },

  // Delete trip
  deleteTrip: async (tripId: string): Promise<ApiResponse<void>> => {
    const response = await apiClient.delete(`/trips/${tripId}`);
    return response.data;
  },

  // Get available seats for a trip
  getAvailableSeats: async (tripId: string): Promise<ApiResponse<{ available_seats: number }>> => {
    const response = await apiClient.get(`/trips/${tripId}/available-seats`);
    return response.data;
  },
};
