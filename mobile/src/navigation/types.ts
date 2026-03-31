/**
 * Navigation types for type-safe navigation
 */

import type { NavigatorScreenParams } from '@react-navigation/native'
import { Coordinates } from '../types'

// Root Stack Navigator params
export type RootStackParamList = {
  // Auth screens
  Login: { message?: string }
  Register: undefined
  ForgotPassword: undefined
  ResetPassword: { token: string }
  VerifyEmail: { token: string }
  ResendVerification: { email?: string }

  // Main app (with tabs)
  MainTabs: NavigatorScreenParams<TabParamList>

  // Detail screens (modals/stack)
  TripDetail: { id: string }
  CreateTrip: undefined
  EditTrip: { id: string }
  EditProfile: undefined
  ChangePassword: undefined
  TripRouteMap: { origin: Coordinates, destination: Coordinates }

  // Tracking
  DriverTracking: { tripId: string }
  TripTracking: { tripId: string }

  // Wallet & Payments
  Wallet: undefined
  Withdraw: undefined
  Transactions: undefined
  SellerOnboarding: undefined

  // Other
  Unauthorized: undefined

  // Admin
  AdminDashboard: undefined
  AdminUsers: undefined
  AdminTrips: undefined
  AdminBookings: undefined
}

// Tab Navigator params
export type TabParamList = {
  Home: undefined
  Search: undefined
  MyTrips: undefined
  MyBookings: undefined
  Profile: undefined
}

declare global {
  namespace ReactNavigation {
    interface RootParamList extends RootStackParamList {}
  }
}
