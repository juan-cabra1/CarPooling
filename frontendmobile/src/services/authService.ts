/**
 * Authentication Service for CarPooling Users API - Mobile Version
 * Handles user authentication, registration, profile management, and password operations
 * Adapted for React Native with AsyncStorage
 */

import AsyncStorage from '@react-native-async-storage/async-storage'
import apiClient from './api'
import type {
  User,
  LoginCredentials,
  RegisterData,
  AuthResponse,
  UpdateUserData,
  ChangePasswordRequest,
  ResetPasswordRequest,
  ForgotPasswordRequest,
  ResendVerificationRequest,
  ApiResponse,
} from '@/types'

const AUTH_BASE = '/users'

/**
 * Login user and receive JWT token
 * @param credentials - Email and password
 * @returns AuthResponse with token and user data
 * @throws Error if login fails
 */
export async function login(credentials: LoginCredentials): Promise<AuthResponse> {
  const response = await apiClient.post<ApiResponse<AuthResponse>>(
    '/login',
    credentials
  )

  const authData = response.data.data!

  // Store token and user in AsyncStorage
  await AsyncStorage.setItem('token', authData.token)
  await AsyncStorage.setItem('user', JSON.stringify(authData.user))

  return authData
}

/**
 * Register a new user
 * @param data - User registration data
 * @returns Created user object
 * @throws Error if registration fails
 */
export async function register(data: RegisterData): Promise<User> {
  const response = await apiClient.post<ApiResponse<User>>(
    `${AUTH_BASE}`,
    data
  )

  return response.data.data!
}

/**
 * Logout current user
 * Clears token and user data from AsyncStorage
 */
export async function logout(): Promise<void> {
  await AsyncStorage.removeItem('token')
  await AsyncStorage.removeItem('user')
}

/**
 * Get current authenticated user profile
 * Requires JWT token in Authorization header
 */
export async function getCurrentUser(): Promise<User> {
  const response = await apiClient.get<ApiResponse<User>>(`${AUTH_BASE}/me`)
  return response.data.data!
}

/**
 * Get user by ID
 */
export async function getUserById(id: number): Promise<User> {
  const response = await apiClient.get<ApiResponse<User>>(`${AUTH_BASE}/${id}`)
  return response.data.data!
}

/**
 * Update user profile
 */
export async function updateProfile(
  id: number,
  data: UpdateUserData
): Promise<User> {
  const response = await apiClient.put<ApiResponse<User>>(
    `${AUTH_BASE}/${id}`,
    data
  )

  const updatedUser = response.data.data!

  // Update user in AsyncStorage if updating current user
  const currentUserStr = await AsyncStorage.getItem('user')
  if (currentUserStr) {
    const currentUser = JSON.parse(currentUserStr) as User
    if (currentUser.id === id) {
      await AsyncStorage.setItem('user', JSON.stringify(updatedUser))
    }
  }

  return updatedUser
}

/**
 * Delete user account
 */
export async function deleteAccount(id: number): Promise<void> {
  await apiClient.delete(`${AUTH_BASE}/${id}`)

  // Clear storage if deleting current user
  const currentUserStr = await AsyncStorage.getItem('user')
  if (currentUserStr) {
    const currentUser = JSON.parse(currentUserStr) as User
    if (currentUser.id === id) {
      await logout()
    }
  }
}

/**
 * Change password for authenticated user
 */
export async function changePassword(data: ChangePasswordRequest): Promise<void> {
  await apiClient.post('/change-password', data)
}

/**
 * Request password reset email
 */
export async function forgotPassword(email: string): Promise<void> {
  const data: ForgotPasswordRequest = { email }
  await apiClient.post('/forgot-password', data)
}

/**
 * Reset password using token from email
 */
export async function resetPassword(data: ResetPasswordRequest): Promise<void> {
  await apiClient.post('/reset-password', data)
}

/**
 * Verify email address using token from email
 */
export async function verifyEmail(token: string): Promise<void> {
  await apiClient.get(`/verify-email?token=${token}`)
}

/**
 * Resend email verification
 */
export async function resendVerification(email: string): Promise<void> {
  const data: ResendVerificationRequest = { email }
  await apiClient.post('/resend-verification', data)
}

/**
 * Get user from AsyncStorage
 * Returns null if not logged in
 */
export async function getStoredUser(): Promise<User | null> {
  try {
    const userStr = await AsyncStorage.getItem('user')
    if (!userStr) return null

    return JSON.parse(userStr) as User
  } catch {
    return null
  }
}

/**
 * Check if user is authenticated (has valid token)
 * Note: This only checks if token exists, not if it's valid/expired
 */
export async function isAuthenticated(): Promise<boolean> {
  try {
    const token = await AsyncStorage.getItem('token')
    return !!token
  } catch {
    return false
  }
}

export default {
  login,
  register,
  logout,
  getCurrentUser,
  getUserById,
  updateProfile,
  deleteAccount,
  changePassword,
  forgotPassword,
  resetPassword,
  verifyEmail,
  resendVerification,
  getStoredUser,
  isAuthenticated,
}
