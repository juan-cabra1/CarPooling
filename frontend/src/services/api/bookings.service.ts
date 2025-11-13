import { apiClient } from './axios';
import type { Booking, ApiResponse, PaginatedResponse } from '@/types';

export interface CreateBookingData {
  trip_id: string;
  seats_booked: number;
}

export interface BookingFilters {
  trip_id?: string;
  passenger_id?: string;
  status?: string;
  page?: number;
  per_page?: number;
}

export const bookingsService = {
  // Create a new booking
  createBooking: async (data: CreateBookingData): Promise<ApiResponse<Booking>> => {
    const response = await apiClient.post('/bookings', data);
    return response.data;
  },

  // Get booking by ID
  getBookingById: async (bookingId: string): Promise<ApiResponse<Booking>> => {
    const response = await apiClient.get(`/bookings/${bookingId}`);
    return response.data;
  },

  // Get all bookings with optional filters
  getBookings: async (filters?: BookingFilters): Promise<PaginatedResponse<Booking>> => {
    const response = await apiClient.get('/bookings', { params: filters });
    return response.data;
  },

  // Get bookings by passenger
  getBookingsByPassenger: async (passengerId: string): Promise<ApiResponse<Booking[]>> => {
    const response = await apiClient.get(`/bookings/passenger/${passengerId}`);
    return response.data;
  },

  // Get bookings by trip
  getBookingsByTrip: async (tripId: string): Promise<ApiResponse<Booking[]>> => {
    const response = await apiClient.get(`/bookings/trip/${tripId}`);
    return response.data;
  },

  // Cancel booking
  cancelBooking: async (bookingId: string): Promise<ApiResponse<Booking>> => {
    const response = await apiClient.patch(`/bookings/${bookingId}/cancel`);
    return response.data;
  },

  // Confirm booking (driver action)
  confirmBooking: async (bookingId: string): Promise<ApiResponse<Booking>> => {
    const response = await apiClient.patch(`/bookings/${bookingId}/confirm`);
    return response.data;
  },

  // Complete booking
  completeBooking: async (bookingId: string): Promise<ApiResponse<Booking>> => {
    const response = await apiClient.patch(`/bookings/${bookingId}/complete`);
    return response.data;
  },
};
