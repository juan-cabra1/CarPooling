/**
 * ProfileScreen
 * Display and manage user profile
 */

import React, { useState, useEffect } from 'react'
import { View, Text, ScrollView, TouchableOpacity, Alert, Switch } from 'react-native'
import { useNavigation } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Card } from '@/components/ui'
import { useAuth } from '@/context/AuthContext'
import { DevConfig } from '@/config/dev'
import { styles } from './ProfileScreen.styles'
import type { RootStackParamList } from '@/navigation/types'

type ProfileScreenNavigationProp = StackNavigationProp<RootStackParamList>

// Extender el tipo User para incluir propiedades opcionales
type ExtendedUser = {
  gender?: string
  address?: {
    street?: string
    city?: string
    province?: string
    postal_code?: string
  }
  is_verified?: boolean
  driver_stats?: {
    total_trips?: number
    rating?: number
  }
  phone?: string
  birthdate?: string
  created_at?: string
} & any // Usamos any para las propiedades base del User

export default function ProfileScreen() {
  const navigation = useNavigation<ProfileScreenNavigationProp>()
  const { user, logout } = useAuth()
  const [isDevMode, setIsDevMode] = useState(false)
  const [loadingDevMode, setLoadingDevMode] = useState(false)

  // Hacer type assertion para el usuario
  const extendedUser = user as ExtendedUser

  // Load dev mode state on mount
  useEffect(() => {
    const loadDevMode = async () => {
      const devMode = await DevConfig.isDevModeEnabled()
      setIsDevMode(devMode)
    }
    loadDevMode()
  }, [])

  const handleLogout = () => {
    Alert.alert('Cerrar Sesión', '¿Estás seguro de que quieres cerrar sesión?', [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Cerrar Sesión',
        style: 'destructive',
        onPress: async () => {
          await logout()
        },
      },
    ])
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('es-AR', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    }).format(date)
  }

  const handleDevModeToggle = async () => {
    setLoadingDevMode(true)
    try {
      const newValue = await DevConfig.toggleDevMode()
      setIsDevMode(newValue)

      Alert.alert(
        newValue ? 'Modo Desarrollo Activado' : 'Modo Desarrollo Desactivado',
        newValue
          ? 'Las reservas se crearán sin procesar pagos. Esto es solo para testing.'
          : 'Las reservas ahora requerirán pago a través de MercadoPago.',
        [{ text: 'OK' }]
      )
    } catch (error) {
      Alert.alert('Error', 'No se pudo cambiar el modo de desarrollo')
    } finally {
      setLoadingDevMode(false)
    }
  }

  if (!user) {
    return (
      <View style={styles.loadingContainer}>
        <Text style={styles.loadingText}>Cargando perfil...</Text>
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <ScrollView>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.avatarContainer}>
            <View style={styles.avatar}>
              <Text style={styles.avatarText}>{user.name.charAt(0).toUpperCase()}</Text>
            </View>
          </View>
          <Text style={styles.userName}>{user.name}</Text>
          <Text style={styles.userEmail}>{user.email}</Text>
        </View>

        {/* Personal Information */}
        <Card style={styles.card}>
          <View style={styles.cardHeader}>
            <Ionicons name="person" size={20} color="#0ea5e9" />
            <Text style={styles.cardTitle}>Información Personal</Text>
          </View>

          <View style={styles.infoSection}>
            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Nombre completo</Text>
              <Text style={styles.infoValue}>{user.name}</Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Email</Text>
              <Text style={styles.infoValue}>{user.email}</Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Teléfono</Text>
              <Text style={styles.infoValue}>{extendedUser.phone || 'No especificado'}</Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Fecha de nacimiento</Text>
              <Text style={styles.infoValue}>
                {extendedUser.birthdate ? formatDate(extendedUser.birthdate) : 'No especificado'}
              </Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Género</Text>
              <Text style={styles.infoValue}>
                {extendedUser.gender === 'male'
                  ? 'Masculino'
                  : extendedUser.gender === 'female'
                  ? 'Femenino'
                  : extendedUser.gender === 'other'
                  ? 'Otro'
                  : 'No especificado'}
              </Text>
            </View>
          </View>
        </Card>

        {/* Address Information */}
        <Card style={styles.card}>
          <View style={styles.cardHeader}>
            <Ionicons name="location" size={20} color="#0ea5e9" />
            <Text style={styles.cardTitle}>Dirección</Text>
          </View>

          <View style={styles.infoSection}>
            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Calle</Text>
              <Text style={styles.infoValue}>{extendedUser.address?.street || 'No especificado'}</Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Ciudad</Text>
              <Text style={styles.infoValue}>{extendedUser.address?.city || 'No especificado'}</Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Provincia</Text>
              <Text style={styles.infoValue}>
                {extendedUser.address?.province || 'No especificado'}
              </Text>
            </View>

            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Código Postal</Text>
              <Text style={styles.infoValue}>
                {extendedUser.address?.postal_code || 'No especificado'}
              </Text>
            </View>
          </View>
        </Card>

        {/* Account Information */}
        <Card style={styles.card}>
          <View style={styles.cardHeader}>
            <Ionicons name="shield-checkmark" size={20} color="#0ea5e9" />
            <Text style={styles.cardTitle}>Cuenta</Text>
          </View>

          <View style={styles.infoSection}>
            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Estado de verificación</Text>
              <View style={styles.verificationBadge}>
                <Ionicons
                  name={extendedUser.is_verified ? 'checkmark-circle' : 'close-circle'}
                  size={16}
                  color={extendedUser.is_verified ? '#22c55e' : '#dc2626'}
                />
                <Text
                  style={
                    extendedUser.is_verified ? styles.verificationTextSuccess : styles.verificationTextError
                  }
                >
                  {extendedUser.is_verified ? 'Verificado' : 'No verificado'}
                </Text>
              </View>
            </View>

            <View style={styles.infoRow}>
              <Text style={styles.infoLabel}>Miembro desde</Text>
              <Text style={styles.infoValue}>
                {extendedUser.created_at ? formatDate(extendedUser.created_at) : 'Desconocido'}
              </Text>
            </View>
          </View>
        </Card>

        {/* Statistics (if available) */}
        {extendedUser.driver_stats && (
          <Card style={styles.card}>
            <View style={styles.cardHeader}>
              <Ionicons name="stats-chart" size={20} color="#0ea5e9" />
              <Text style={styles.cardTitle}>Estadísticas como Conductor</Text>
            </View>

            <View style={styles.statsGrid}>
              <View style={styles.statItem}>
                <Text style={styles.statValue}>{extendedUser.driver_stats.total_trips || 0}</Text>
                <Text style={styles.statLabel}>Viajes</Text>
              </View>

              <View style={styles.statItem}>
                <View style={styles.ratingContainer}>
                  <Ionicons name="star" size={20} color="#facc15" />
                  <Text style={styles.statValue}>
                    {extendedUser.driver_stats.rating?.toFixed(1) || '0.0'}
                  </Text>
                </View>
                <Text style={styles.statLabel}>Calificación</Text>
              </View>
            </View>
          </Card>
        )}

        {/* Quick Actions */}
        <Card style={styles.card}>
          <View style={styles.cardHeader}>
            <Ionicons name="settings" size={20} color="#0ea5e9" />
            <Text style={styles.cardTitle}>Acciones</Text>
          </View>

          <View style={styles.actionsSection}>
            <TouchableOpacity
              style={styles.actionButton}
              onPress={() => navigation.navigate('Wallet')}
            >
              <Ionicons name="wallet" size={20} color="#22c55e" />
              <Text style={styles.actionButtonText}>Mi Billetera</Text>
              <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
            </TouchableOpacity>

            {/* Dev Mode Toggle */}
            <View style={styles.actionButton}>
              <Ionicons name="construct" size={20} color="#f59e0b" />
              <Text style={styles.actionButtonText}>Modo Desarrollo</Text>
              <Switch
                value={isDevMode}
                onValueChange={handleDevModeToggle}
                disabled={loadingDevMode}
                trackColor={{ false: '#d1d5db', true: '#fbbf24' }}
                thumbColor={isDevMode ? '#f59e0b' : '#f3f4f6'}
              />
            </View>

            <TouchableOpacity style={styles.actionButton}>
              <Ionicons name="create" size={20} color="#6b7280" />
              <Text style={styles.actionButtonText}>Editar Perfil</Text>
              <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
            </TouchableOpacity>

            <TouchableOpacity style={styles.actionButton}>
              <Ionicons name="lock-closed" size={20} color="#6b7280" />
              <Text style={styles.actionButtonText}>Cambiar Contraseña</Text>
              <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
            </TouchableOpacity>

            <TouchableOpacity style={styles.actionButton}>
              <Ionicons name="notifications" size={20} color="#6b7280" />
              <Text style={styles.actionButtonText}>Notificaciones</Text>
              <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
            </TouchableOpacity>

            <TouchableOpacity style={styles.actionButton}>
              <Ionicons name="help-circle" size={20} color="#6b7280" />
              <Text style={styles.actionButtonText}>Ayuda y Soporte</Text>
              <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
            </TouchableOpacity>
          </View>
        </Card>

        {/* Logout Button */}
        <View style={styles.logoutContainer}>
          <Button
            title="Cerrar Sesión"
            variant="destructive"
            onPress={handleLogout}
            style={styles.logoutButton}
          />
        </View>

        {/* Version Info */}
        <View style={styles.versionContainer}>
          <Text style={styles.versionText}>Carpooling App v1.0.0</Text>
          <Text style={styles.versionText}>React Native + Expo</Text>
        </View>
      </ScrollView>
    </View>
  )
}
