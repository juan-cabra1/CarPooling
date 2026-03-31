/**
 * ResendVerificationScreen
 * Allows users to resend email verification link
 */

import React, { useState } from 'react'
import {
  View,
  Text,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Alert,
} from 'react-native'
import { useNavigation, useRoute, type RouteProp } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Input, Button } from '@/components/ui'
import { authService, getErrorMessage } from '@/services'
import { styles } from './ResendVerificationScreen.styles'
import type { RootStackParamList } from '@/navigation/types'

type ResendVerificationScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'ResendVerification'
>
type ResendVerificationScreenRouteProp = RouteProp<
  RootStackParamList,
  'ResendVerification'
>

export default function ResendVerificationScreen() {
  const navigation = useNavigation<ResendVerificationScreenNavigationProp>()
  const route = useRoute<ResendVerificationScreenRouteProp>()

  const [email, setEmail] = useState(route.params?.email || '')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)

  const handleResendVerification = async () => {
    if (!email.trim()) {
      Alert.alert('Error', 'Por favor ingresa tu correo electrónico')
      return
    }

    // Basic email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(email)) {
      Alert.alert('Error', 'Por favor ingresa un correo electrónico válido')
      return
    }

    try {
      setLoading(true)
      await authService.resendVerification(email)
      setSuccess(true)

      Alert.alert(
        'Email Enviado',
        'Se ha enviado un nuevo enlace de verificación a tu correo electrónico. Por favor revisa tu bandeja de entrada.',
        [
          {
            text: 'OK',
            onPress: () => navigation.navigate('Login'),
          },
        ]
      )
    } catch (err) {
      Alert.alert('Error', getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={styles.container}
    >
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        keyboardShouldPersistTaps="handled"
      >
        <View style={styles.content}>
          {/* Icon */}
          <View style={styles.iconContainer}>
            <Ionicons name="mail-outline" size={80} color="#0ea5e9" />
          </View>

          {/* Title */}
          <Text style={styles.title}>Reenviar Verificación</Text>

          {/* Description */}
          <Text style={styles.description}>
            {success
              ? 'Hemos enviado un nuevo enlace de verificación a tu correo electrónico.'
              : 'Ingresa tu correo electrónico y te enviaremos un nuevo enlace de verificación.'}
          </Text>

          {!success && (
            <>
              {/* Email Input */}
              <View style={styles.inputContainer}>
                <Input
                  placeholder="Correo electrónico"
                  value={email}
                  onChangeText={setEmail}
                  keyboardType="email-address"
                  autoCapitalize="none"
                  autoCorrect={false}
                  editable={!loading}
                />
              </View>

              {/* Resend Button */}
              <Button
                title="Reenviar Email"
                onPress={handleResendVerification}
                loading={loading}
                disabled={loading}
                style={styles.resendButton}
              />
            </>
          )}

          {/* Back to Login */}
          <Button
            title="Volver a Login"
            variant="outline"
            onPress={() => navigation.navigate('Login')}
            disabled={loading}
            style={styles.backButton}
          />

          {/* Help Text */}
          <Text style={styles.helpText}>
            Si no recibes el correo en unos minutos, revisa tu carpeta de spam.
          </Text>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  )
}
