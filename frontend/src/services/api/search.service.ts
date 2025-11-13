import { apiClient } from './axios';
import type { SearchFilters, SearchResult, ApiResponse } from '@/types';

export const searchService = {
  // Search trips
  searchTrips: async (filters: SearchFilters): Promise<ApiResponse<SearchResult[]>> => {
    const response = await apiClient.get('/search/trips', { params: filters });
    return response.data;
  },

  // Advanced search with geolocation
  searchNearbyTrips: async (
    latitude: number,
    longitude: number,
    radius: number,
    filters?: SearchFilters
  ): Promise<ApiResponse<SearchResult[]>> => {
    const response = await apiClient.get('/search/nearby', {
      params: {
        lat: latitude,
        lng: longitude,
        radius,
        ...filters,
      },
    });
    return response.data;
  },

  // Get popular routes
  getPopularRoutes: async (): Promise<ApiResponse<Array<{ origin: string; destination: string; count: number }>>> => {
    const response = await apiClient.get('/search/popular-routes');
    return response.data;
  },

  // Autocomplete locations
  autocompleteLocation: async (query: string): Promise<ApiResponse<string[]>> => {
    const response = await apiClient.get('/search/autocomplete', { params: { q: query } });
    return response.data;
  },
};
