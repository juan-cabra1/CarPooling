/**
 * VerifyEmailScreen
 * Verifies user's email address using token from email link
 */

import React, { useEffect, useState } from 'react'
import { View, Text, ActivityIndicator } from 'react-native'
import { useNavigation, useRoute, type RouteProp } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button } from '@/components/ui'
import { authService, getErrorMessage } from '@/services'
import { styles } from './VerifyEmailScreen.styles'
import type { RootStackParamList } from '@/navigation/types'

type VerifyEmailScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'VerifyEmail'
>
type VerifyEmailScreenRouteProp = RouteProp<RootStackParamList, 'VerifyEmail'>

export default function VerifyEmailScreen() {
  const navigation = useNavigation<VerifyEmailScreenNavigationProp>()
  const route = useRoute<VerifyEmailScreenRouteProp>()

  const [loading, setLoading] = useState(true)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const verifyEmail = async () => {
      const token = route.params?.token

      if (!token) {
        setError('Token de verificación no proporcionado')
        setLoading(false)
        return
      }

      try {
        await authService.verifyEmail(token)
        setSuccess(true)
      } catch (err) {
        setError(getErrorMessage(err))
      } finally {
        setLoading(false)
      }
    }

    verifyEmail()
  }, [route.params?.token])

  const handleNavigateToLogin = () => {
    navigation.navigate('Login')
  }

  if (loading) {
    return (
      <View style={styles.container}>
        <View style={styles.content}>
          <ActivityIndicator size="large" color="#0ea5e9" />
          <Text style={styles.loadingText}>Verificando tu correo electrónico...</Text>
        </View>
      </View>
    )
  }

  if (success) {
    return (
      <View style={styles.container}>
        <View style={styles.content}>
          {/* Success Icon */}
          <View style={styles.iconContainer}>
            <Ionicons name="checkmark-circle" size={80} color="#22c55e" />
          </View>

          {/* Title */}
          <Text style={styles.title}>¡Email Verificado!</Text>

          {/* Description */}
          <Text style={styles.description}>
            Tu correo electrónico ha sido verificado exitosamente. Ahora puedes iniciar sesión y
            comenzar a usar CarPooling.
          </Text>

          {/* Login Button */}
          <Button
            title="Ir a Iniciar Sesión"
            onPress={handleNavigateToLogin}
            style={styles.button}
          />
        </View>
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <View style={styles.content}>
        {/* Error Icon */}
        <View style={styles.iconContainer}>
          <Ionicons name="close-circle" size={80} color="#ef4444" />
        </View>

        {/* Title */}
        <Text style={styles.title}>Verificación Fallida</Text>

        {/* Error Message */}
        <Text style={styles.description}>{error || 'No se pudo verificar tu correo electrónico'}</Text>

        {/* Possible reasons */}
        <View style={styles.reasonsContainer}>
          <Text style={styles.reasonsTitle}>Posibles razones:</Text>
          <Text style={styles.reasonItem}>• El enlace ha expirado</Text>
          <Text style={styles.reasonItem}>• El enlace ya fue utilizado</Text>
          <Text style={styles.reasonItem}>• El enlace es inválido</Text>
        </View>

        {/* Resend Button */}
        <Button
          title="Solicitar Nuevo Enlace"
          onPress={() => navigation.navigate('ResendVerification', { email: '' })}
          style={styles.button}
        />

        {/* Back to Login */}
        <Button
          title="Volver a Login"
          variant="outline"
          onPress={handleNavigateToLogin}
          style={styles.outlineButton}
        />
      </View>
    </View>
  )
}
