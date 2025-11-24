/**
 * Authentication Context for CarPooling Application
 * Provides global authentication state management
 * Delegates all logic to authService (no duplication)
 */

import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import type { User, LoginCredentials, RegisterData } from '@/types'
import { authService } from '@/services'

interface AuthContextType {
  user: User | null
  loading: boolean
  isAuthenticated: boolean
  login: (credentials: LoginCredentials) => Promise<void>
  register: (data: RegisterData) => Promise<void>
  logout: () => void
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

interface AuthProviderProps {
  children: ReactNode
}

/**
 * Authentication Provider Component
 * Wraps the application to provide auth state globally
 * @example
 * <AuthProvider>
 *   <App />
 * </AuthProvider>
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  // Initialize auth state from localStorage on mount
  useEffect(() => {
    const storedUser = authService.getStoredUser()
    if (storedUser) {
      setUser(storedUser)
    }
    setLoading(false)
  }, [])

  /**
   * Login user with credentials
   * Delegates to authService, throws errors for component to handle
   * @throws Error if login fails
   */
  const login = async (credentials: LoginCredentials) => {
    const response = await authService.login(credentials)
    setUser(response.user)
  }

  /**
   * Register new user
   * Delegates to authService, throws errors for component to handle
   * @throws Error if registration fails
   */
  const register = async (data: RegisterData) => {
    const user = await authService.register(data)
    setUser(user)
  }

  /**
   * Logout current user
   * Clears localStorage and resets state
   */
  const logout = () => {
    authService.logout()
    setUser(null)
  }

  /**
   * Refresh current user data from API
   * Has try-catch because it's called automatically
   * Logs out user if refresh fails (invalid/expired token)
   */
  const refreshUser = async () => {
    try {
      const freshUser = await authService.getCurrentUser()
      setUser(freshUser)
    } catch (error) {
      // Token invalid or expired, logout
      logout()
    }
  }

  const value: AuthContextType = {
    user,
    loading,
    isAuthenticated: !!user,
    login,
    register,
    logout,
    refreshUser,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

/**
 * Custom hook to access authentication context
 * Must be used within AuthProvider
 * @throws Error if used outside AuthProvider
 * @example
 * const { user, login, logout, isAuthenticated } = useAuth()
 */
export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
