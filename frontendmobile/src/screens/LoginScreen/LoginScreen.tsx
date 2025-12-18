/**
 * LoginScreen
 * User authentication screen
 */

import React, { useState, useEffect } from 'react'
import {
  View,
  Text,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  TouchableOpacity,
  Alert,
} from 'react-native'
import { useNavigation, useRoute, type RouteProp } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { useAuth } from '@/context/AuthContext'
import { Button, Input, PasswordInput, Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui'
import { getErrorMessage } from '@/services'
import { styles } from './LoginScreen.styles'
import type { RootStackParamList } from '@/navigation/types'

type LoginScreenNavigationProp = StackNavigationProp<RootStackParamList, 'Login'>
type LoginScreenRouteProp = RouteProp<RootStackParamList, 'Login'>

export default function LoginScreen() {
  const navigation = useNavigation<LoginScreenNavigationProp>()
  const route = useRoute<LoginScreenRouteProp>()
  const { login } = useAuth()

  const [formData, setFormData] = useState({
    email: '',
    password: '',
  })

  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [successMessage, setSuccessMessage] = useState('')
  const [showResendLink, setShowResendLink] = useState(false)

  // Check for success message from navigation params
  useEffect(() => {
    if (route.params?.message) {
      setSuccessMessage(route.params.message)
    }
  }, [route.params])

  const handleSubmit = async () => {
    if (!formData.email || !formData.password) {
      Alert.alert('Error', 'Por favor completa todos los campos')
      return
    }

    setError('')
    setSuccessMessage('')
    setLoading(true)

    try {
      await login({
        email: formData.email,
        password: formData.password,
      })

      // Login exitoso, la navegación se manejará automáticamente por AppNavigator
    } catch (err) {
      const errorMsg = getErrorMessage(err)

      // Verificar si es error de email no verificado
      if (errorMsg.includes('verificar tu correo') || errorMsg.includes('verificar tu email')) {
        setError(errorMsg)
        setShowResendLink(true)
      } else {
        setError(errorMsg)
        setShowResendLink(false)
      }
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
        <Card style={styles.card} shadow>
          <CardHeader>
            <View style={styles.iconContainer}>
              <Ionicons name="log-in-outline" size={32} color="#ffffff" />
            </View>
            <CardTitle style={styles.title}>
              Bienvenido de nuevo
            </CardTitle>
            <CardDescription style={styles.description}>
              Ingresa tus credenciales para acceder a tu cuenta
            </CardDescription>
          </CardHeader>

          <CardContent>
            {/* Success Message */}
            {successMessage ? (
              <View style={styles.successContainer}>
                <Ionicons name="checkmark-circle" size={20} color="#16a34a" />
                <Text style={styles.successText}>{successMessage}</Text>
              </View>
            ) : null}

            {/* Error Message */}
            {error ? (
              <View style={styles.errorContainer}>
                <View style={styles.errorInner}>
                  <Ionicons name="alert-circle" size={20} color="#dc2626" />
                  <Text style={styles.errorText}>{error}</Text>
                </View>
                {showResendLink && (
                  <View style={styles.resendLink}>
                    <TouchableOpacity
                      onPress={() => navigation.navigate('ResendVerification', { email: formData.email })}
                    >
                      <Text style={styles.resendLinkText}>
                        Reenviar email de verificación
                      </Text>
                    </TouchableOpacity>
                  </View>
                )}
              </View>
            ) : null}

            {/* Form */}
            <View style={styles.formContainer}>
              {/* Email Input */}
              <View style={styles.inputGroup}>
                <View style={styles.labelContainer}>
                  <Ionicons name="mail-outline" size={16} color="#1e293b" />
                  <Text style={styles.label}>Correo electrónico</Text>
                </View>
                <Input
                  placeholder="tu@email.com"
                  value={formData.email}
                  onChangeText={(text) => setFormData({ ...formData, email: text })}
                  keyboardType="email-address"
                  autoCapitalize="none"
                  autoCorrect={false}
                  editable={!loading}
                  containerStyle={{ marginBottom: 0 }}
                />
              </View>

              {/* Password Input */}
              <View style={styles.inputGroup}>
                <View style={styles.labelContainer}>
                  <Ionicons name="lock-closed-outline" size={16} color="#1e293b" />
                  <Text style={styles.label}>Contraseña</Text>
                </View>
                <PasswordInput
                  placeholder="••••••••"
                  value={formData.password}
                  onChangeText={(text) => setFormData({ ...formData, password: text })}
                  editable={!loading}
                  containerStyle={{ marginBottom: 0 }}
                />
              </View>

              {/* Forgot Password Link */}
              <TouchableOpacity onPress={() => navigation.navigate('ForgotPassword')}>
                <Text style={styles.forgotPassword}>
                  ¿Olvidaste tu contraseña?
                </Text>
              </TouchableOpacity>
            </View>
          </CardContent>

          <CardFooter>
            {/* Login Button */}
            <Button
              title="Iniciar Sesión"
              onPress={handleSubmit}
              loading={loading}
              disabled={loading}
              style={{ width: '100%' }}
              size="lg"
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
