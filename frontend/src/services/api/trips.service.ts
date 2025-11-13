import { apiClient } from './axios';
import type { Trip, CreateTripData, UpdateTripData, TripFilters } from '@/types';

interface TripsListResponse {
  success: boolean;
  data: {
    trips: Trip[];
    total: number;
    page: number;
    limit: number;
  };
}

interface TripResponse {
  success: boolean;
  data: Trip;
}

interface DeleteResponse {
  success: boolean;
  message: string;
}

export const tripsService = {
  // Create a new trip (requires authentication)
  createTrip: async (data: CreateTripData): Promise<Trip> => {
    const response = await apiClient.post<TripResponse>('/api/trips', data);
    return response.data.data;
  },

  // Get trip by ID (public)
  getTripById: async (tripId: string): Promise<Trip> => {
    const response = await apiClient.get<TripResponse>(`/api/trips/${tripId}`);
    return response.data.data;
  },

  // Get all trips with optional filters (public)
  getTrips: async (filters?: TripFilters): Promise<TripsListResponse['data']> => {
    const response = await apiClient.get<TripsListResponse>('/api/trips', { params: filters });
    return response.data.data;
  },

  // Update trip (requires authentication)
  updateTrip: async (tripId: string, data: UpdateTripData): Promise<Trip> => {
    const response = await apiClient.put<TripResponse>(`/api/trips/${tripId}`, data);
    return response.data.data;
  },

  // Partial update trip (requires authentication)
  partialUpdateTrip: async (tripId: string, data: Partial<UpdateTripData>): Promise<Trip> => {
    const response = await apiClient.patch<TripResponse>(`/api/trips/${tripId}`, data);
    return response.data.data;
  },

  // Cancel trip (requires authentication)
  cancelTrip: async (tripId: string, reason: string): Promise<void> => {
    await apiClient.patch(`/api/trips/${tripId}/cancel`, { reason });
  },

  // Delete trip (requires authentication)
  deleteTrip: async (tripId: string): Promise<void> => {
    await apiClient.delete<DeleteResponse>(`/api/trips/${tripId}`);
  },
};
