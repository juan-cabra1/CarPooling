import axios, { AxiosError } from 'axios';
import type { AxiosInstance, AxiosRequestConfig } from 'axios';
import type { ApiError } from '@/types';

// Create axios instance with default config
const createAxiosInstance = (baseURL?: string): AxiosInstance => {
  const instance = axios.create({
    baseURL: baseURL || import.meta.env.VITE_API_BASE_URL,
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Request interceptor for adding auth token
  instance.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('auth_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Response interceptor for handling errors
  instance.interceptors.response.use(
    (response) => response,
    (error: AxiosError<ApiError>) => {
      if (error.response) {
        // Log detailed error information
        console.error('API Error Response:', {
          status: error.response.status,
          data: error.response.data,
          url: error.config?.url,
          method: error.config?.method,
        });

        // Handle specific error codes
        switch (error.response.status) {
          case 401:
            // Unauthorized - clear token and redirect to login
            localStorage.removeItem('auth_token');
            localStorage.removeItem('user');
            window.location.href = '/login';
            break;
          case 403:
            // Forbidden
            console.error('Access forbidden');
            break;
          case 404:
            // Not found
            console.error('Resource not found');
            break;
          case 500:
            // Server error
            console.error('Server error:', error.response.data);
            break;
          default:
            console.error('API error:', error.response.data);
        }
      } else if (error.request) {
        // Request made but no response received
        console.error('Network error - no response received');
      } else {
        // Something else happened
        console.error('Error:', error.message);
      }
      return Promise.reject(error);
    }
  );

  return instance;
};

// Main API instance (uses proxy)
export const apiClient = createAxiosInstance();

// Direct API instances for each microservice
export const usersApi = createAxiosInstance(import.meta.env.VITE_USERS_API_URL);
export const tripsApi = createAxiosInstance(import.meta.env.VITE_TRIPS_API_URL);
export const bookingsApi = createAxiosInstance(import.meta.env.VITE_BOOKINGS_API_URL);
export const searchApi = createAxiosInstance(import.meta.env.VITE_SEARCH_API_URL);

// Generic request wrapper with error handling
export async function apiRequest<T>(
  config: AxiosRequestConfig,
  client: AxiosInstance = apiClient
): Promise<T> {
  try {
    const response = await client.request<T>(config);
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      throw error.response.data;
    }
    throw error;
  }
}
