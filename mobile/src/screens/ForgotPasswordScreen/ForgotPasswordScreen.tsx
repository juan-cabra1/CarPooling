/**
 * ForgotPasswordScreen
 * Password recovery request screen
 */

import React, { useState } from 'react'
import { View, Text, ScrollView, KeyboardAvoidingView, Platform } from 'react-native'
import { useNavigation } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Input, Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui'
import { authService, getErrorMessage } from '@/services'
import { styles } from './ForgotPasswordScreen.styles'
import type { RootStackParamList } from '@/navigation/types'

type ForgotPasswordScreenNavigationProp = StackNavigationProp<RootStackParamList, 'ForgotPassword'>

export default function ForgotPasswordScreen() {
  const navigation = useNavigation<ForgotPasswordScreenNavigationProp>()

  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  const handleSubmit = async () => {
    setError('')
    setLoading(true)

    try {
      await authService.forgotPassword(email)
      setSuccess(true)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  // Success Screen
  if (success) {
    return (
      <View style={styles.container}>
        <Card style={styles.card} shadow>
          <CardHeader>
            <View style={styles.successIconContainer}>
              <Ionicons name="checkmark-circle" size={32} color="#ffffff" />
            </View>
            <CardTitle style={styles.title}>¡Email Enviado!</CardTitle>
            <CardDescription style={styles.description}>
              Hemos enviado las instrucciones de recuperación a <Text style={{ fontWeight: 'bold' }}>{email}</Text>
            </CardDescription>
          </CardHeader>

          <CardContent>
            <View style={styles.infoBox}>
              <Text style={styles.infoText}>
                Por favor, revisa tu bandeja de entrada y sigue las instrucciones para restablecer tu contraseña.
                El enlace expirará en 1 hora.
              </Text>
            </View>

            <View style={styles.warningBox}>
              <Text style={styles.warningText}>
                Si no recibes el correo en unos minutos, revisa tu carpeta de spam o correo no deseado.
              </Text>
            </View>
          </CardContent>

          <CardFooter>
            <Button
              title="Volver al Inicio de Sesión"
              variant="outline"
              size="lg"
              style={styles.button}
              onPress={() => navigation.navigate('Login', { message: undefined })}
            />
          </CardFooter>
        </Card>
      </View>
    )
  }

  // Request Form
  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={styles.container}
    >
      <ScrollView
        contentContainerStyle={{ flexGrow: 1, justifyContent: 'center', padding: 24 }}
        keyboardShouldPersistTaps="handled"
      >
        <Card style={styles.card} shadow>
          <CardHeader>
            <View style={styles.iconContainer}>
              <Ionicons name="key" size={32} color="#ffffff" />
            </View>
            <CardTitle style={styles.title}>¿Olvidaste tu contraseña?</CardTitle>
            <CardDescription style={styles.description}>
              Ingresa tu correo electrónico y te enviaremos instrucciones para recuperarla
            </CardDescription>
          </CardHeader>

          <CardContent>
            {/* Error Message */}
            {error ? (
              <View style={styles.errorContainer}>
                <Ionicons name="alert-circle" size={20} color="#dc2626" />
                <Text style={styles.errorText}>{error}</Text>
              </View>
            ) : null}

            {/* Email Input */}
            <Input
              label="Correo electrónico"
              placeholder="tu@email.com"
              value={email}
              onChangeText={setEmail}
              keyboardType="email-address"
              autoCapitalize="none"
              autoCorrect={false}
              editable={!loading}
            />

            {/* Info Box */}
            <View style={styles.infoBox}>
              <Text style={styles.infoText}>
                Te enviaremos un correo con un enlace para restablecer tu contraseña.
                El enlace será válido durante 1 hora.
              </Text>
            </View>
          </CardContent>

          <CardFooter>
            {/* Submit Button */}
            <Button
              title="Enviar Instrucciones"
              onPress={handleSubmit}
              loading={loading}
              disabled={loading}
              size="lg"
              style={styles.button}
            />

            {/* Back to Login Button */}
            <Button
              title="Volver al Inicio de Sesión"
              variant="outline"
              size="lg"
              style={styles.button}
              onPress={() => navigation.navigate('Login', { message: undefined })}
            />

            {/* Register Link */}
            <View style={styles.registerContainer}>
              <Text style={styles.registerText}>
                ¿No tienes cuenta?{' '}
                <Text
                  style={styles.registerLink}
                  onPress={() => navigation.navigate('Register')}
                >
                  Regístrate aquí
                </Text>
              </Text>
            </View>
          </CardFooter>
        </Card>
      </ScrollView>
    </KeyboardAvoidingView>
  )
}