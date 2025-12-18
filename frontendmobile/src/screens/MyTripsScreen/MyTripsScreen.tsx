/**
 * MyTripsScreen
 * Display user's trips as driver
 */

import React, { useState, useEffect } from 'react'
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  RefreshControl,
  Alert,
} from 'react-native'
import { useNavigation, useFocusEffect } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button } from '@/components/ui'
import { tripsService, getErrorMessage } from '@/services'
import { useAuth } from '@/context/AuthContext'
import { styles } from './MyTripsScreen.styles'
import type { Trip } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type MyTripsScreenNavigationProp = StackNavigationProp<RootStackParamList>

// Definir los posibles estados del viaje localmente
type LocalTripStatus = 'active' | 'completed' | 'cancelled' | 'pending' | 'full'

// Componente TripCard temporal
const TripCard = ({ trip }: { trip: Trip }) => {
  const navigation = useNavigation<MyTripsScreenNavigationProp>()

  // El filtro en fetchTrips debería evitar que lleguen aquí viajes sin ID, 
  // pero esta verificación asegura que la lógica interna de la tarjeta es robusta.
  if (!trip || trip.id === undefined || trip.id === null) {
    // Ya no deberías ver esta advertencia si el filtro funciona correctamente
    console.warn('TripCard recibió un objeto de viaje inválido o sin ID:', trip);
    return null; 
  }
  
  // Variable segura para el ID del viaje
  const tripIdString = trip.id.toString();


  // Debug: verificar estructura del objeto trip
  console.log('Trip object in TripCard:', JSON.stringify(trip, null, 2))
  console.log('Trip ID:', trip.id)
  console.log('Trip _id:', (trip as any)._id)

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('es-AR', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date)
  }

  const getStatusText = (status: string): string => {
    switch (status) {
      case 'active':
        return 'Activo'
      case 'completed':
        return 'Completado'
      case 'cancelled':
        return 'Cancelado'
      case 'pending':
        return 'Pendiente'
      case 'full':
        return 'Completo'
      case 'in_progress':
        return 'En Progreso'
      default:
        return status
    }
  }

  const isTripActive = (status: string): boolean => {
    return status === 'active' || status === 'pending'
  }

  const isTripInProgress = (status: string): boolean => {
    return status === 'in_progress'
  }

  return (
    <View style={styles.tripCard}>
      <View style={styles.tripCardHeader}>
        <View style={styles.tripCardContent}>
          <Text style={styles.tripRoute}>
            {trip.origin.city} → {trip.destination.city}
          </Text>
          <Text style={styles.tripDate}>
            {formatDate(trip.departure_datetime)}
          </Text>
        </View>
        <Text style={styles.tripPrice}>${trip.price_per_seat}</Text>
      </View>

      <View style={styles.tripMetaRow}>
        <View style={styles.tripSeatsContainer}>
          <Ionicons name="person" size={16} color="#6b7280" />
          <Text style={styles.tripSeatsText}>
            {trip.available_seats} asiento{trip.available_seats !== 1 ? 's' : ''} disponible{trip.available_seats !== 1 ? 's' : ''}
          </Text>
        </View>
        <View style={styles.statusContainer}>
          {isTripInProgress(trip.status) && <View style={styles.liveBadge}><Text style={styles.liveBadgeText}>EN VIVO</Text></View>}
          <Text style={styles.tripStatus}>
            {getStatusText(trip.status)}
          </Text>
        </View>
      </View>

      <View style={styles.tripActionsRow}>
        {isTripInProgress(trip.status) ? (
          <Button
            title="Rastrear"
            variant="primary"
            size="sm"
            style={styles.tripActionButton}
            onPress={() => navigation.navigate('DriverTracking', { tripId: tripIdString })}
          />
        ) : (
          <Button
            title="Ver Detalles"
            variant="outline"
            size="sm"
            style={styles.tripActionButton}
            onPress={() => navigation.navigate('TripDetail', { id: tripIdString })}
          />
        )}
        {isTripActive(trip.status) && (
          <Button
            title="Cancelar"
            variant="destructive"
            size="sm"
            style={styles.tripActionButton}
            onPress={() => {
              Alert.alert(
                'Cancelar Viaje',
                '¿Estás seguro de que quieres cancelar este viaje?',
                [
                  { text: 'No', style: 'cancel' },
                  {
                    text: 'Sí, cancelar',
                    style: 'destructive',
                    onPress: () => {
                      // Aquí puedes agregar la lógica para cancelar el viaje
                      console.log('Cancelar viaje:', tripIdString)
                    },
                  },
                ]
              )
            }}
          />
        )}
      </View>
    </View>
  )
}

export default function MyTripsScreen() {
  const navigation = useNavigation<MyTripsScreenNavigationProp>()
  const { user } = useAuth()

  const [trips, setTrips] = useState<Trip[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState('')

  const fetchTrips = async () => {
    try {
      setError('')
      if (!user?.id) {
        throw new Error('Usuario no autenticado')
      }

      const data = await tripsService.getMyTrips(user.id)
      // El servicio ahora devuelve viajes limpios y validados.
      setTrips(data.trips || [])

    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  useFocusEffect(
    React.useCallback(() => {
      if (user?.id) {
        fetchTrips()
      }
    }, [user?.id])
  )

  const onRefresh = () => {
    setRefreshing(true)
    fetchTrips()
  }

  const handleDeleteTrip = async (tripId: string) => {
    Alert.alert(
      'Eliminar Viaje',
      '¿Estás seguro de que quieres eliminar este viaje?',
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Eliminar',
          style: 'destructive',
          onPress: async () => {
            try {
              await tripsService.deleteTrip(tripId)
              setTrips(trips.filter((t) => t.id.toString() !== tripId))
              Alert.alert('Éxito', 'Viaje eliminado correctamente')
            } catch (err) {
              Alert.alert('Error', getErrorMessage(err))
            }
          },
        },
      ]
    )
  }

  // Loading State
  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={styles.loadingText}>Cargando tus viajes...</Text>
      </View>
    )
  }

  // Error State
  if (error) {
    return (
      <View style={styles.errorContainer}>
        <Ionicons name="alert-circle" size={64} color="#dc2626" />
        <Text style={styles.errorText}>{error}</Text>
        <Button title="Reintentar" onPress={fetchTrips} />
      </View>
    )
  }

  // Empty State
  if (trips.length === 0) {
    return (
      <View style={styles.emptyContainer}>
        <Ionicons name="car-outline" size={80} color="#cbd5e1" style={styles.emptyIcon} />
        <Text style={styles.emptyTitle}>No tienes viajes publicados</Text>
        <Text style={styles.emptyText}>
          Publica tu primer viaje y empieza a compartir tus rutas con otros usuarios
        </Text>
        <Button
          title="Publicar Viaje"
          onPress={() => navigation.navigate('CreateTrip' as never)}
          style={styles.emptyButton}
        />
      </View>
    )
  }

  // List State
  return (
    <View style={styles.container}>
      <FlatList
        data={trips}
        // keyExtractor robusto (ya corregido)
        keyExtractor={(item, index) => (item?.id?.toString() || index.toString())}
        renderItem={({ item }) => <TripCard trip={item} />}
        contentContainerStyle={{ padding: 16 }}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} colors={['#0ea5e9']} />
        }
      />
    </View>
  )
}