import { apiClient } from './axios';
import type { User, ApiResponse } from '@/types';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  password: string;
  name: string;
  lastname: string;
  phone: string;
  street: string;
  number: number;
  photo_url?: string;
  sex: 'hombre' | 'mujer' | 'otro';
  birthdate: string; // Format: YYYY-MM-DD
}

export interface UpdateUserData {
  name?: string;
  lastname?: string;
  phone?: string;
  street?: string;
  number?: number;
  photo_url?: string;
}

export const usersService = {
  // Authentication
  login: async (credentials: LoginCredentials): Promise<ApiResponse<{ user: User; token: string }>> => {
    const response = await apiClient.post('/login', credentials);
    return response.data;
  },

  register: async (data: RegisterData): Promise<ApiResponse<User>> => {
    const response = await apiClient.post('/users', data);
    return response.data;
  },

  logout: async (): Promise<void> => {
    // No hay endpoint de logout en el backend actual
    return Promise.resolve();
  },

  // User operations
  getCurrentUser: async (): Promise<ApiResponse<User>> => {
    const response = await apiClient.get('/users/me');
    return response.data;
  },

  getUserById: async (userId: string): Promise<ApiResponse<User>> => {
    const response = await apiClient.get(`/users/${userId}`);
    return response.data;
  },

  updateUser: async (userId: string, data: UpdateUserData): Promise<ApiResponse<User>> => {
    const response = await apiClient.patch(`/users/${userId}`, data);
    return response.data;
  },

  deleteUser: async (userId: string): Promise<ApiResponse<void>> => {
    const response = await apiClient.delete(`/users/${userId}`);
    return response.data;
  },

  // Password management
  changePassword: async (currentPassword: string, newPassword: string): Promise<ApiResponse<void>> => {
    const response = await apiClient.post('/change-password', {
      current_password: currentPassword,
      new_password: newPassword,
    });
    return response.data;
  },

  requestPasswordReset: async (email: string): Promise<ApiResponse<void>> => {
    const response = await apiClient.post('/forgot-password', { email });
    return response.data;
  },

  resetPassword: async (token: string, newPassword: string): Promise<ApiResponse<void>> => {
    const response = await apiClient.post('/reset-password', {
      token,
      new_password: newPassword,
    });
    return response.data;
  },
};
