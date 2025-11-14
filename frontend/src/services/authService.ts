/**
 * Authentication Service for CarPooling Users API (Port 8001)
 * Handles user authentication, registration, profile management, and password operations
 */

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
 * @example
 * const { token, user } = await authService.login({ email: 'user@example.com', password: 'password123' })
 */
export async function login(credentials: LoginCredentials): Promise<AuthResponse> {
  const response = await apiClient.post<ApiResponse<AuthResponse>>(
    '/login',
    credentials
  )

  const authData = response.data.data!

  // Store token and user in localStorage
  localStorage.setItem('token', authData.token)
  localStorage.setItem('user', JSON.stringify(authData.user))

  return authData
}

/**
 * Register a new user
 * @param data - User registration data
 * @returns Created user object
 * @throws Error if registration fails (e.g., email already exists)
 * @example
 * const user = await authService.register({ email: '...', password: '...', ... })
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
 * Clears token and user data from localStorage
 */
export function logout(): void {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
}

/**
 * Get current authenticated user profile
 * Requires JWT token in Authorization header (handled by interceptor)
 * @returns Current user data
 * @throws Error if not authenticated (401)
 * @example
 * const user = await authService.getCurrentUser()
 */
export async function getCurrentUser(): Promise<User> {
  const response = await apiClient.get<ApiResponse<User>>(`${AUTH_BASE}/me`)
  return response.data.data!
}

/**
 * Get user by ID
 * Requires authentication
 * @param id - User ID
 * @returns User data
 * @throws Error if user not found (404) or unauthorized (401)
 * @example
 * const user = await authService.getUserById(123)
 */
export async function getUserById(id: number): Promise<User> {
  const response = await apiClient.get<ApiResponse<User>>(`${AUTH_BASE}/${id}`)
  return response.data.data!
}

/**
 * Update user profile
 * Only the user themselves can update their profile
 * @param id - User ID
 * @param data - Fields to update (all optional)
 * @returns Updated user data
 * @throws Error if unauthorized (403) or user not found (404)
 * @example
 * const updatedUser = await authService.updateProfile(123, { name: 'New Name' })
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

  // Update user in localStorage if updating current user
  const currentUser = JSON.parse(localStorage.getItem('user') || 'null') as User | null
  if (currentUser && currentUser.id === id) {
    localStorage.setItem('user', JSON.stringify(updatedUser))
  }

  return updatedUser
}

/**
 * Delete user account
 * Only the user themselves can delete their account
 * @param id - User ID
 * @throws Error if unauthorized (403) or user not found (404)
 * @example
 * await authService.deleteAccount(123)
 */
export async function deleteAccount(id: number): Promise<void> {
  await apiClient.delete(`${AUTH_BASE}/${id}`)

  // Clear local storage if deleting current user
  const currentUser = JSON.parse(localStorage.getItem('user') || 'null') as User | null
  if (currentUser && currentUser.id === id) {
    logout()
  }
}

/**
 * Change password for authenticated user
 * Requires current password and new password
 * @param data - Current password and new password
 * @throws Error if current password is incorrect (400)
 * @example
 * await authService.changePassword({ current_password: 'old', new_password: 'new' })
 */
export async function changePassword(data: ChangePasswordRequest): Promise<void> {
  await apiClient.post('/change-password', data)
}

/**
 * Request password reset email
 * Public endpoint - doesn't reveal if email exists
 * @param email - User email address
 * @example
 * await authService.forgotPassword('user@example.com')
 */
export async function forgotPassword(email: string): Promise<void> {
  const data: ForgotPasswordRequest = { email }
  await apiClient.post('/forgot-password', data)
}

/**
 * Reset password using token from email
 * Public endpoint
 * @param data - Reset token and new password
 * @throws Error if token is invalid or expired (400)
 * @example
 * await authService.resetPassword({ token: 'xxx', new_password: 'newpass123' })
 */
export async function resetPassword(data: ResetPasswordRequest): Promise<void> {
  await apiClient.post('/reset-password', data)
}

/**
 * Verify email address using token from email
 * Public endpoint
 * @param token - Email verification token
 * @throws Error if token is invalid or expired (400)
 * @example
 * await authService.verifyEmail('verification-token-from-email')
 */
export async function verifyEmail(token: string): Promise<void> {
  await apiClient.get(`/verify-email?token=${token}`)
}

/**
 * Resend email verification
 * Public endpoint
 * @param email - User email address
 * @example
 * await authService.resendVerification('user@example.com')
 */
export async function resendVerification(email: string): Promise<void> {
  const data: ResendVerificationRequest = { email }
  await apiClient.post('/resend-verification', data)
}

/**
 * Get user from localStorage
 * Returns null if not logged in
 */
export function getStoredUser(): User | null {
  const userStr = localStorage.getItem('user')
  if (!userStr) return null

  try {
    return JSON.parse(userStr) as User
  } catch {
    return null
  }
}

/**
 * Check if user is authenticated (has valid token)
 * Note: This only checks if token exists, not if it's valid/expired
 */
export function isAuthenticated(): boolean {
  return !!localStorage.getItem('token')
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
