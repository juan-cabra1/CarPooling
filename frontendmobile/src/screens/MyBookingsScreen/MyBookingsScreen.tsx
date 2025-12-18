/**
 * MyBookingsScreen
 * Display user's bookings as passenger
 */

import React, { useState } from 'react'
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
import { Button, Badge } from '@/components/ui'
import { bookingsService, searchService, tripsService, getErrorMessage } from '@/services'
import { styles } from './MyBookingsScreen.styles'
import type { Booking, SearchTrip } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type MyBookingsScreenNavigationProp = StackNavigationProp<RootStackParamList>

// Extended booking type with populated trip details
interface BookingWithTrip extends Omit<Booking, 'trip'> {
  trip?: SearchTrip
}

export default function MyBookingsScreen() {
  const navigation = useNavigation<MyBookingsScreenNavigationProp>()

  const [bookings, setBookings] = useState<BookingWithTrip[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)

  const fetchBookings = async () => {
    try {
      const data = await bookingsService.getMyBookings()
      const bookingsData = data.bookings || []

      // Fetch trip details for each booking
      const bookingsWithTrips = await Promise.all(
        bookingsData.map(async (booking): Promise<BookingWithTrip> => {
          try {
            // Try search-api first (has denormalized driver data)
            const tripDetails = await searchService.getTripDetails(booking.trip_id)
            return { ...booking, trip: tripDetails }
          } catch (searchErr: any) {
            // If search-api returns 404 (trip not indexed yet), fallback to trips-api
            if (searchErr?.response?.status === 404) {
              try {
                console.log(`Trip ${booking.trip_id} not in search-api, falling back to trips-api`)
                const tripData = await tripsService.getTripById(booking.trip_id)

                // Convert Trip to SearchTrip format (without driver data)
                const searchTrip: SearchTrip = {
                  ...tripData,
                  trip_id: tripData.id, // Add trip_id field for SearchTrip compatibility
                  driver: {
                    id: tripData.driver_id,
                    name: 'Cargando...',
                    email: '',
                    rating: 0,
                    total_trips: 0,
                  },
                }
                return { ...booking, trip: searchTrip }
              } catch (tripsErr) {
                console.error(`Error fetching trip ${booking.trip_id} from trips-api:`, tripsErr)
                return booking as BookingWithTrip
              }
            }
            console.error(`Error fetching trip ${booking.trip_id}:`, searchErr)
            return booking as BookingWithTrip
          }
        })
      )

      setBookings(bookingsWithTrips)
    } catch (err) {
      console.error('Error fetching bookings:', err)
      setBookings([])
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  useFocusEffect(
    React.useCallback(() => {
      fetchBookings()
    }, [])
  )

  const onRefresh = () => {
    setRefreshing(true)
    fetchBookings()
  }

  const handleCancelBooking = async (bookingId: string) => {
    Alert.prompt(
      'Cancelar Reserva',
      'Por favor indica el motivo de la cancelación:',
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Confirmar',
          onPress: async (reason: string | undefined) => {
            try {
              await bookingsService.cancelBooking(bookingId, reason || 'Sin motivo especificado')
              fetchBookings()
              Alert.alert('Éxito', 'Reserva cancelada correctamente')
            } catch (err) {
              Alert.alert('Error', getErrorMessage(err))
            }
          },
        },
      ],
      'plain-text'
    )
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'confirmed':
        return <Badge variant="success">Confirmada</Badge>
      case 'pending':
        return <Badge variant="warning">Pendiente</Badge>
      case 'cancelled':
        return <Badge variant="destructive">Cancelada</Badge>
      case 'completed':
        return <Badge variant="secondary">Completada</Badge>
      default:
        return <Badge>{status}</Badge>
    }
  }

  const renderBooking = ({ item }: { item: Booking }) => {
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

    const formatPrice = (price: number) => {
      return new Intl.NumberFormat('es-AR', {
        style: 'currency',
        currency: 'ARS',
        minimumFractionDigits: 0,
      }).format(price)
    }

    return (
      <View style={styles.bookingCard}>
        <View style={styles.bookingHeader}>
          <View style={styles.tripInfo}>
            {item.trip ? (
              <>
                <View style={styles.route}>
                  <Ionicons name="location" size={16} color="#0ea5e9" />
                  <Text style={styles.cityText}>{item.trip.origin.city}</Text>
                </View>
                <View style={styles.route}>
                  <Ionicons name="location" size={16} color="#14b8a6" />
                  <Text style={styles.cityText}>{item.trip.destination.city}</Text>
                </View>
                <Text style={styles.dateText}>{formatDate(item.trip.departure_datetime)}</Text>
              </>
            ) : (
              <Text style={styles.tripIdText}>Viaje ID: {item.trip_id}</Text>
            )}
          </View>

          <View style={styles.priceContainer}>
            <Text style={styles.totalPrice}>{formatPrice(item.total_price)}</Text>
            <Text style={styles.seatsText}>
              {item.seats_requested} asiento{item.seats_requested > 1 ? 's' : ''}
            </Text>
          </View>
        </View>

        <View style={styles.bookingStatus}>{getStatusBadge(item.status)}</View>

        <View style={styles.actionsRow}>
          <Button
            title="Ver Viaje"
            variant="outline"
            size="sm"
            style={styles.actionButton}
            onPress={() => navigation.navigate('TripDetail', { id: item.trip_id })}
          />
          {(item.status === 'pending' || item.status === 'confirmed') && (
            <Button
              title="Cancelar"
              variant="destructive"
              size="sm"
              style={styles.actionButton}
              onPress={() => handleCancelBooking(item.id)}
            />
          )}
        </View>
      </View>
    )
  }

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={styles.loadingText}>Cargando tus reservas...</Text>
      </View>
    )
  }

  if (bookings.length === 0) {
    return (
      <View style={styles.emptyContainer}>
        <Ionicons name="list-outline" size={80} color="#cbd5e1" style={styles.emptyIcon} />
        <Text style={styles.emptyTitle}>No tienes reservas</Text>
        <Text style={styles.emptyText}>
          Busca un viaje y reserva tu asiento para empezar a viajar
        </Text>
        <Button
          title="Buscar Viajes"
          onPress={() => navigation.navigate('MainTabs', { screen: 'Search' })}
          style={styles.emptyButton}
        />
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <FlatList
        data={bookings}
        keyExtractor={(item) => item.id}
        renderItem={renderBooking}
        contentContainerStyle={{ padding: 16 }}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} colors={['#0ea5e9']} />
        }
      />
    </View>
  )
}
