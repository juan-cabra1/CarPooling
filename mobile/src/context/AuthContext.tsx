/**
 * Authentication Context for CarPooling Mobile Application
 * Provides global authentication state management
 * Adapted for React Native with AsyncStorage
 */

import React, { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import AsyncStorage from '@react-native-async-storage/async-storage'
import type { User, LoginCredentials, RegisterData } from '@/types'
import { authService } from '@/services'

interface AuthContextType {
  user: User | null
  token: string | null
  loading: boolean
  isAuthenticated: boolean
  login: (credentials: LoginCredentials) => Promise<void>
  register: (data: RegisterData) => Promise<void>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

interface AuthProviderProps {
  children: ReactNode
}

/**
 * Authentication Provider Component
 * Wraps the application to provide auth state globally
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  // Initialize auth state from AsyncStorage on mount
  useEffect(() => {
    const initAuth = async () => {
      try {
        const storedUser = await authService.getStoredUser()
        if (storedUser) {
          setUser(storedUser)
          const storedToken = await AsyncStorage.getItem('token')
          setToken(storedToken)
        }
      } catch (error) {
        console.error('Error loading user from AsyncStorage:', error)
      } finally {
        setLoading(false)
      }
    }

    initAuth()
  }, [])

  /**
   * Login user with credentials
   * Delegates to authService, throws errors for component to handle
   */
  const login = async (credentials: LoginCredentials) => {
    const response = await authService.login(credentials)
    setUser(response.user)
    setToken(response.token)
  }

  /**
   * Register new user
   * Delegates to authService, throws errors for component to handle
   */
  const register = async (data: RegisterData) => {
    const user = await authService.register(data)
    setUser(user)
  }

  /**
   * Logout current user
   * Clears AsyncStorage and resets state
   */
  const logout = async () => {
    await authService.logout()
    setUser(null)
    setToken(null)
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
      await logout()
    }
  }

  const value: AuthContextType = {
    user,
    token,
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
 */
export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
