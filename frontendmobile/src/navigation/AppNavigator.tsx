/**
 * Root Stack Navigator
 * Handles all app navigation including auth flow
 */

import React from 'react'
import { createStackNavigator } from '@react-navigation/stack'
import { useAuth } from '@/context/AuthContext'
import type { RootStackParamList } from './types'

// Tab Navigator
import TabNavigator from './TabNavigator'

// Auth Screens
import LoginScreen from '@/screens/LoginScreen/LoginScreen'
import RegisterScreen from '@/screens/RegisterScreen/RegisterScreen'
import ForgotPasswordScreen from '@/screens/ForgotPasswordScreen/ForgotPasswordScreen'
import ResetPasswordScreen from '@/screens/ResetPasswordScreen/ResetPasswordScreen'
import VerifyEmailScreen from '@/screens/VerifyEmailScreen/VerifyEmailScreen'
import ResendVerificationScreen from '@/screens/ResendVerificationScreen/ResendVerificationScreen'

// Detail Screens
import TripDetailScreen from '@/screens/TripDetailScreen/TripDetailScreen'
import CreateTripScreen from '@/screens/CreateTripScreen/CreateTripScreen'
import EditTripScreen from '@/screens/EditTripScreen/EditTripScreen'

// Wallet & Payment Screens
import WalletScreen from '@/screens/WalletScreen'
import WithdrawScreen from '@/screens/WithdrawScreen'
import TransactionsScreen from '@/screens/TransactionsScreen'
import SellerOnboardingScreen from '@/screens/SellerOnboardingScreen'

import TripTrackingScreen from '@/screens/TripTrackingScreen/TripTrackingScreen'
import DriverTrackingScreen from '@/screens/DriverTrackingScreen/DriverTrackingScreen'
import TripRouteMapScreen from '@/screens/TripRouteMapScreen/TripRouteMapScreen'

// Other Screens
import UnauthorizedScreen from '@/screens/UnauthorizedScreen/UnauthorizedScreen'

const Stack = createStackNavigator<RootStackParamList>()

export default function AppNavigator() {
  const { isAuthenticated, loading, user } = useAuth()

  if (loading) {
    // You can return a loading screen here
    return null
  }

  return (
    <Stack.Navigator
      screenOptions={{
        headerStyle: {
          backgroundColor: '#0ea5e9', // primary-500
        },
        headerTintColor: '#ffffff',
        headerTitleStyle: {
          fontWeight: 'bold',
        },
      }}
    >
      {!isAuthenticated ? (
        // Auth Stack
        <>
          <Stack.Screen
            name="Login"
            component={LoginScreen}
            options={{ headerShown: false }}
          />
          <Stack.Screen
            name="Register"
            component={RegisterScreen}
            options={{ title: 'Crear Cuenta', headerShown: false }}
          />
          <Stack.Screen
            name="ForgotPassword"
            component={ForgotPasswordScreen}
            options={{ title: 'Recuperar Contraseña' }}
          />
          <Stack.Screen
            name="ResetPassword"
            component={ResetPasswordScreen}
            options={{ title: 'Restablecer Contraseña' }}
          />
          <Stack.Screen
            name="VerifyEmail"
            component={VerifyEmailScreen}
            options={{ title: 'Verificar Email' }}
          />
          <Stack.Screen
            name="ResendVerification"
            component={ResendVerificationScreen}
            options={{ title: 'Reenviar Verificación' }}
          />
        </>
      ) : (
        // Main App Stack
        <>
          <Stack.Screen
            name="MainTabs"
            component={TabNavigator}
            options={{ headerShown: false }}
          />
          <Stack.Screen
            name="TripDetail"
            component={TripDetailScreen}
            options={{ title: 'Detalle del Viaje' }}
          />
          <Stack.Screen
            name="CreateTrip"
            component={CreateTripScreen}
            options={{ title: 'Publicar Viaje' }}
          />
          <Stack.Screen
            name="EditTrip"
            component={EditTripScreen}
            options={{ title: 'Editar Viaje' }}
          />
          <Stack.Screen
            name="TripRouteMap"
            component={TripRouteMapScreen}
            options={{ title: 'Ruta del Viaje' }}
          />
          <Stack.Screen
            name="DriverTracking"
            component={DriverTrackingScreen}
            options={{ title: 'Siguiendo Viaje' }}
          />
          <Stack.Screen
            name="TripTracking"
            component={TripTrackingScreen}
            options={{ title: 'Tu Viaje' }}
          />
          <Stack.Screen
            name="Wallet"
            component={WalletScreen}
            options={{ title: 'Mi Billetera' }}
          />
          <Stack.Screen
            name="Withdraw"
            component={WithdrawScreen}
            options={{ title: 'Retirar Fondos' }}
          />
          <Stack.Screen
            name="Transactions"
            component={TransactionsScreen}
            options={{ title: 'Historial de Transacciones' }}
          />
          <Stack.Screen
            name="SellerOnboarding"
            component={SellerOnboardingScreen}
            options={{ title: 'Vincular MercadoPago' }}
          />
          <Stack.Screen
            name="Unauthorized"
            component={UnauthorizedScreen}
            options={{ title: 'No Autorizado' }}
          />
        </>
      )}
    </Stack.Navigator>
  )
}
