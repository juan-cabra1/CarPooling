/**
 * CarPooling Mobile App
 * Main entry point
 */

import React, { useRef, useEffect } from 'react'
import { NavigationContainer, NavigationContainerRef } from '@react-navigation/native'
import { StatusBar } from 'expo-status-bar'
import { SafeAreaProvider } from 'react-native-safe-area-context'
import { AuthProvider } from './src/context/AuthContext'
import { TrackingProvider } from './src/context/TrackingContext'
import { setNavigationRef } from './src/services/api'
import AppNavigator from './src/navigation/AppNavigator'
import type { RootStackParamList } from './src/navigation/types'

const linking = {
  prefixes: ['carpooling://', 'https://carpooling.com'],
  config: {
    screens: {
      Login: 'login',
      Register: 'register',
      ForgotPassword: 'forgot-password',
      ResetPassword: 'reset-password',
      VerifyEmail: {
        path: 'verify-email',
        parse: {
          token: (token: string) => token,
        },
      },
      ResendVerification: 'resend-verification',
      MainTabs: {
        screens: {
          Search: 'search',
          MyTrips: 'my-trips',
          MyBookings: 'my-bookings',
          Profile: 'profile',
        },
      },
      TripDetail: 'trip/:id',
      CreateTrip: 'create-trip',
      EditTrip: 'edit-trip/:id',
      Unauthorized: 'unauthorized',
    },
  },
}

export default function App() {
  const navigationRef = useRef<NavigationContainerRef<RootStackParamList>>(null)

  useEffect(() => {
    // Set navigation ref for API service to handle redirects
    if (navigationRef.current) {
      setNavigationRef(navigationRef.current)
    }
  }, [])

  return (
    <SafeAreaProvider>
      <AuthProvider>
        <TrackingProvider>
          <NavigationContainer ref={navigationRef} linking={linking}>
            <AppNavigator />
            <StatusBar style="light" />
          </NavigationContainer>
        </TrackingProvider>
      </AuthProvider>
    </SafeAreaProvider>
  )
}
