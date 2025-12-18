/**
 * RegisterScreen
 * User registration screen
 */

import React, { useState } from 'react'
import {
  View,
  Text,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  Alert,
  TouchableOpacity,
} from 'react-native'
import { useNavigation } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import DateTimePicker from '@react-native-community/datetimepicker'
import { 
  Button, 
  Input, 
  PasswordInput, 
  Card, 
  CardHeader, 
  CardTitle, 
  CardDescription, 
  CardContent, 
  CardFooter 
} from '@/components/ui'
import { getErrorMessage, authService } from '@/services'
import { styles } from './RegisterScreen.styles'
import type { RootStackParamList } from '@/navigation/types'
import type { RegisterData, UserSex } from '@/types'

type RegisterScreenNavigationProp = StackNavigationProp<RootStackParamList, 'Register'>

export default function RegisterScreen() {
  const navigation = useNavigation<RegisterScreenNavigationProp>()

  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [showDatePicker, setShowDatePicker] = useState(false)

  const [formData, setFormData] = useState<RegisterData>({
    email: '',
    password: '',
    name: '',
    lastname: '',
    phone: '',
    street: '',
    number: 0,
    sex: 'otro',
    birthdate: '',
  })

  const [confirmPassword, setConfirmPassword] = useState('')
  const [selectedDate, setSelectedDate] = useState(new Date(2000, 0, 1))

  const validateForm = (): boolean => {
    // Validar campos requeridos
    if (!formData.email || !formData.password || !formData.name || !formData.lastname) {
      Alert.alert('Error', 'Por favor completa todos los campos obligatorios')
      return false
    }

    // Validar formato de email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(formData.email)) {
      Alert.alert('Error', 'Por favor ingresa un correo electrónico válido')
      return false
    }

    // Validar contraseñas
    if (formData.password !== confirmPassword) {
      Alert.alert('Error', 'Las contraseñas no coinciden')
      return false
    }

    if (formData.password.length < 6) {
      Alert.alert('Error', 'La contraseña debe tener al menos 6 caracteres')
      return false
    }

    // Validar fecha de nacimiento
    if (!formData.birthdate) {
      Alert.alert('Error', 'Por favor selecciona tu fecha de nacimiento')
      return false
    }

    // Validar edad mínima (18 años)
    const birthDate = new Date(formData.birthdate)
    const today = new Date()
    const age = today.getFullYear() - birthDate.getFullYear()
    const monthDiff = today.getMonth() - birthDate.getMonth()
    const dayDiff = today.getDate() - birthDate.getDate()

    const actualAge = monthDiff < 0 || (monthDiff === 0 && dayDiff < 0) ? age - 1 : age

    if (actualAge < 18) {
      Alert.alert('Error', 'Debes tener al menos 18 años para registrarte')
      return false
    }

    // Validar teléfono (opcional pero si se ingresa, debe tener formato válido)
    if (formData.phone && formData.phone.length < 8) {
      Alert.alert('Error', 'Por favor ingresa un número de teléfono válido')
      return false
    }

    return true
  }

  const handleSubmit = async () => {
    setError('')

    if (!validateForm()) {
      return
    }

    setLoading(true)

    try {
      await authService.register(formData)
      setSuccess(true)

      // Redirigir después de 3 segundos
      setTimeout(() => {
        navigation.navigate('Login', {
          message: 'Registro exitoso. Por favor, verifica tu correo electrónico antes de iniciar sesión.',
        })
      }, 3000)
    } catch (err) {
      const errorMessage = getErrorMessage(err)
      setError(errorMessage)
      Alert.alert('Error', errorMessage)
    } finally {
      setLoading(false)
    }
  }

  const handleDateChange = (event: any, date?: Date) => {
    setShowDatePicker(false)
    if (date) {
      setSelectedDate(date)
      const formattedDate = date.toISOString().split('T')[0]
      setFormData({ ...formData, birthdate: formattedDate })
    }
  }

  const formatDisplayDate = (dateString: string): string => {
    if (!dateString) return 'Seleccionar fecha'
    
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('es-AR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric'
    }).format(date)
  }

  // Success Screen
  if (success) {
    return (
      <View style={{ flex: 1, backgroundColor: '#ffffff', alignItems: 'center', justifyContent: 'center', padding: 24 }}>
        <Card style={{ width: '100%' }} shadow>
          <CardHeader>
            <View style={styles.successIconContainer}>
              <Ionicons name="checkmark-circle" size={32} color="#ffffff" />
            </View>
            <CardTitle style={styles.title}>¡Registro Exitoso!</CardTitle>
            <CardDescription style={styles.description}>
              Te hemos enviado un correo de verificación a <Text style={{ fontWeight: 'bold' }}>{formData.email}</Text>
            </CardDescription>
          </CardHeader>
          <CardContent>
            <View style={styles.successInfo}>
              <Text style={styles.successInfoText}>
                Por favor, verifica tu correo electrónico antes de iniciar sesión.
                Si no recibes el correo en unos minutos, revisa tu carpeta de spam.
              </Text>
            </View>
            <Text style={styles.redirectText}>
              Serás redirigido a la página de inicio de sesión en unos segundos...
            </Text>
          </CardContent>
        </Card>
      </View>
    )
  }

  // Registration Form
  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      style={styles.container}
    >
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        keyboardShouldPersistTaps="handled"
        showsVerticalScrollIndicator={false}
      >
        <Card style={styles.card} shadow>
          <CardHeader>
            <View style={styles.iconContainer}>
              <Ionicons name="person-add" size={32} color="#ffffff" />
            </View>
            <CardTitle style={styles.title}>Crear Cuenta</CardTitle>
            <CardDescription style={styles.description}>
              Completa el formulario para unirte a CarPooling
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

            {/* Account Section */}
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Información de la Cuenta</Text>

              <Input
                label="Correo electrónico *"
                placeholder="tu@email.com"
                value={formData.email}
                onChangeText={(text) => setFormData({ ...formData, email: text })}
                keyboardType="email-address"
                autoCapitalize="none"
                autoCorrect={false}
                editable={!loading}
              />

              <PasswordInput
                label="Contraseña *"
                placeholder="Mínimo 6 caracteres"
                value={formData.password}
                onChangeText={(text) => setFormData({ ...formData, password: text })}
                editable={!loading}
              />

              <PasswordInput
                label="Confirmar Contraseña *"
                placeholder="Repite tu contraseña"
                value={confirmPassword}
                onChangeText={setConfirmPassword}
                editable={!loading}
                containerStyle={{ marginBottom: 0 }}
              />
            </View>

            {/* Personal Info Section */}
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Información Personal</Text>

              <Input
                label="Nombre *"
                placeholder="Tu nombre"
                value={formData.name}
                onChangeText={(text) => setFormData({ ...formData, name: text })}
                editable={!loading}
              />

              <Input
                label="Apellido *"
                placeholder="Tu apellido"
                value={formData.lastname}
                onChangeText={(text) => setFormData({ ...formData, lastname: text })}
                editable={!loading}
              />

              <Input
                label="Teléfono"
                placeholder="+56 9 1234 5678"
                value={formData.phone}
                onChangeText={(text) => setFormData({ ...formData, phone: text })}
                keyboardType="phone-pad"
                editable={!loading}
              />

              {/* Custom Picker for Sex */}
              <View style={{ marginBottom: 16 }}>
                <Text style={{ fontSize: 14, fontWeight: '500', color: '#1e293b', marginBottom: 8 }}>Sexo</Text>
                <View style={{ flexDirection: 'row', borderWidth: 1, borderColor: '#cbd5e1', borderRadius: 8, backgroundColor: '#ffffff' }}>
                  {[
                    { label: 'Hombre', value: 'hombre' as UserSex },
                    { label: 'Mujer', value: 'mujer' as UserSex },
                    { label: 'Otro', value: 'otro' as UserSex },
                  ].map((option) => (
                    <TouchableOpacity
                      key={option.value}
                      onPress={() => setFormData({ ...formData, sex: option.value })}
                      disabled={loading}
                      style={[
                        { flex: 1, paddingVertical: 12, alignItems: 'center' },
                        formData.sex === option.value ? { backgroundColor: '#0ea5e9' } : { backgroundColor: '#ffffff' },
                        option.value === 'hombre' ? { borderTopLeftRadius: 8, borderBottomLeftRadius: 8 } : {},
                        option.value === 'otro' ? { borderTopRightRadius: 8, borderBottomRightRadius: 8 } : {},
                      ]}
                    >
                      <Text
                        style={
                          formData.sex === option.value
                            ? { color: '#ffffff', fontWeight: '500' }
                            : { color: '#1e293b' }
                        }
                      >
                        {option.label}
                      </Text>
                    </TouchableOpacity>
                  ))}
                </View>
              </View>

              {/* Date Picker */}
              <View style={{ marginBottom: 16 }}>
                <Text style={{ fontSize: 14, fontWeight: '500', color: '#1e293b', marginBottom: 8 }}>Fecha de Nacimiento *</Text>
                <TouchableOpacity
                  onPress={() => setShowDatePicker(true)}
                  disabled={loading}
                  style={[
                    { backgroundColor: '#ffffff', borderWidth: 1, borderRadius: 8, paddingHorizontal: 16, paddingVertical: 12 },
                    loading ? { backgroundColor: '#f1f5f9' } : { borderColor: '#cbd5e1' }
                  ]}
                >
                  <Text style={formData.birthdate ? { color: '#1e293b' } : { color: '#334155' }}>
                    {formatDisplayDate(formData.birthdate)}
                  </Text>
                </TouchableOpacity>
              </View>

              {showDatePicker && (
                <DateTimePicker
                  value={selectedDate}
                  mode="date"
                  display={Platform.OS === 'ios' ? 'spinner' : 'default'}
                  onChange={handleDateChange}
                  maximumDate={new Date()}
                />
              )}
            </View>

            {/* Address Section */}
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Dirección</Text>

              <Input
                label="Calle"
                placeholder="Nombre de la calle"
                value={formData.street}
                onChangeText={(text) => setFormData({ ...formData, street: text })}
                editable={!loading}
              />

              <Input
                label="Número"
                placeholder="1234"
                value={formData.number === 0 ? '' : formData.number.toString()}
                onChangeText={(text) => {
                  const num = text === '' ? 0 : parseInt(text) || 0
                  setFormData({ ...formData, number: num })
                }}
                keyboardType="numeric"
                editable={!loading}
                containerStyle={{ marginBottom: 0 }}
              />
            </View>
          </CardContent>

          <CardFooter>
            <Button
              title="Crear Cuenta"
              onPress={handleSubmit}
              loading={loading}
              disabled={loading}
              style={styles.button}
              size="lg"
            />

            <View style={styles.loginContainer}>
              <Text style={styles.loginText}>
                ¿Ya tienes cuenta?{' '}
                <Text
                  style={styles.loginLink}
                  onPress={() => navigation.navigate('Login' as never)}
                >
                  Inicia sesión aquí
                </Text>
              </Text>
            </View>
          </CardFooter>
        </Card>
      </ScrollView>
    </KeyboardAvoidingView>
  )
}