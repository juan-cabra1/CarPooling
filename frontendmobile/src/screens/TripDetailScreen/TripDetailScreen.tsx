/**
 * TripDetailScreen
 * Display complete trip details and allow booking
 */

import React, { useState, useEffect } from 'react'
import {
  View,
  Text,
  ScrollView,
  ActivityIndicator,
  TouchableOpacity,
  Alert,
} from 'react-native'
import { useNavigation, useRoute } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Badge, Card } from '@/components/ui'
import { BookingModal } from '@/components/modals'
import { tripsService, searchService, getErrorMessage } from '@/services'
import { useAuth } from '@/context/AuthContext'
import { styles } from './TripDetailScreen.styles'
import type { SearchTrip } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type TripDetailScreenNavigationProp = StackNavigationProp<RootStackParamList>

export default function TripDetailScreen() {
  const navigation = useNavigation<TripDetailScreenNavigationProp>()
  const route = useRoute()
  const { user } = useAuth()

  const { id } = route.params as { id: string }

  const [trip, setTrip] = useState<SearchTrip | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [deleting, setDeleting] = useState(false)
  const [starting, setStarting] = useState(false)
  const [isBookingModalVisible, setIsBookingModalVisible] = useState(false)

  useEffect(() => {
    if (id) {
      fetchTrip()
    }
  }, [id])

  const fetchTrip = async () => {
    if (!id) return

    try {
      setLoading(true)
      setError('')

      // Try search-api first (has denormalized driver data)
      try {
        const data = await searchService.getTripDetails(id)
        setTrip(data)
      } catch (searchErr: any) {
        // If search-api returns 404 (trip not indexed yet), fallback to trips-api
        if (searchErr?.response?.status === 404) {
          console.log('Trip not in search-api yet, falling back to trips-api')
          const tripData = await tripsService.getTripById(id)

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
          setTrip(searchTrip)
        } else {
          throw searchErr
        }
      }
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = () => {
    if (!trip || !id) return

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
              setDeleting(true)
              await tripsService.deleteTrip(id)
              Alert.alert('Éxito', 'Viaje eliminado exitosamente')
              navigation.navigate('MainTabs', { screen: 'MyTrips' })
            } catch (err) {
              Alert.alert('Error', getErrorMessage(err))
            } finally {
              setDeleting(false)
            }
          },
        },
      ]
    )
  }

  const handleBookingSuccess = () => {
    fetchTrip()
    navigation.navigate('MainTabs', { screen: 'MyBookings' })
  }

  const handleStartTrip = async () => {
    if (!trip || !id) return

    Alert.alert(
      'Iniciar Viaje',
      '¿Estás listo para iniciar el viaje? Los pasajeros podrán rastrear tu ubicación en tiempo real.',
      [
        { text: 'Cancelar', style: 'cancel' },
        {
          text: 'Iniciar',
          onPress: async () => {
            try {
              setStarting(true)
              await tripsService.startTrip(id)
              Alert.alert('Viaje Iniciado', 'El viaje ha comenzado. Los pasajeros fueron notificados.')
              // Refrescar para obtener el nuevo status
              await fetchTrip()
              // Navegar directamente al tracking
              navigation.navigate('DriverTracking', { tripId: id })
            } catch (err) {
              Alert.alert('Error', getErrorMessage(err))
            } finally {
              setStarting(false)
            }
          },
        },
      ]
    )
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('es-AR', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric',
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

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<
      string,
      { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' | 'success' | 'warning' }
    > = {
      draft: { label: 'Borrador', variant: 'outline' },
      published: { label: 'Publicado', variant: 'success' },
      full: { label: 'Completo', variant: 'secondary' },
      in_progress: { label: 'En Progreso', variant: 'secondary' },
      completed: { label: 'Completado', variant: 'secondary' },
      cancelled: { label: 'Cancelado', variant: 'destructive' },
    }

    const config = statusConfig[status] || { label: status, variant: 'outline' as const }
    return <Badge variant={config.variant}>{config.label}</Badge>
  }

  const isOwner = user && trip && trip.driver_id === user.id

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={styles.loadingText}>Cargando detalles del viaje...</Text>
      </View>
    )
  }

  if (error || !trip) {
    return (
      <View style={styles.errorContainer}>
        <Ionicons name="alert-circle-outline" size={64} color="#dc2626" />
        <Text style={styles.errorTitle}>Error al cargar viaje</Text>
        <Text style={styles.errorText}>{error || 'No se pudo encontrar el viaje'}</Text>
        <Button
          title="Volver a búsqueda"
          onPress={() => navigation.navigate('MainTabs', { screen: 'Search' })}
          style={styles.errorButton}
        />
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <ScrollView>
        {/* Back Button */}
        <View style={styles.backButtonContainer}>
          <TouchableOpacity style={styles.backButton} onPress={() => navigation.goBack()}>
            <Ionicons name="arrow-back" size={24} color="#0ea5e9" />
            <Text style={styles.backButtonText}>Volver</Text>
          </TouchableOpacity>
        </View>

        {/* Main Card */}
        <Card style={styles.mainCard}>
          {/* Header */}
          <View style={styles.header}>
            <View style={styles.headerContent}>
              <View style={styles.headerTitleRow}>
                <Text style={styles.headerTitle}>
                  {trip.origin.city} → {trip.destination.city}
                </Text>
                {getStatusBadge(trip.status)}
              </View>
              <Text style={styles.headerSubtitle}>
                {trip.origin.province} → {trip.destination.province}
              </Text>
            </View>
            <View style={styles.priceContainer}>
              <Text style={styles.price}>{formatPrice(trip.price_per_seat)}</Text>
              <Text style={styles.priceLabel}>por asiento</Text>
            </View>
          </View>

          {/* Route Details */}
          <View style={styles.routeDetails}>
            <View style={styles.routeItem}>
              <View style={styles.routeIconRow}>
                <Ionicons name="location" size={20} color="#0ea5e9" />
                <Text style={styles.routeLabel}>Origen</Text>
              </View>
              <View style={styles.routeInfo}>
                <Text style={styles.routeCity}>
                  {trip.origin.city}, {trip.origin.province}
                </Text>
                <Text style={styles.routeAddress}>{trip.origin.address}</Text>
              </View>
            </View>

            <View style={styles.routeItem}>
              <View style={styles.routeIconRow}>
                <Ionicons name="location" size={20} color="#14b8a6" />
                <Text style={styles.routeLabel}>Destino</Text>
              </View>
              <View style={styles.routeInfo}>
                <Text style={styles.routeCity}>
                  {trip.destination.city}, {trip.destination.province}
                </Text>
                <Text style={styles.routeAddress}>{trip.destination.address}</Text>
              </View>
            </View>
          </View>

          {/* Date & Time */}
          <View style={styles.dateSection}>
            <View style={styles.dateItem}>
              <View style={styles.dateIconRow}>
                <Ionicons name="calendar" size={20} color="#0ea5e9" />
                <Text style={styles.dateLabel}>Salida</Text>
              </View>
              <Text style={styles.dateValue}>{formatDate(trip.departure_datetime)}</Text>
            </View>

            <View style={styles.dateItem}>
              <View style={styles.dateIconRow}>
                <Ionicons name="calendar" size={20} color="#14b8a6" />
                <Text style={styles.dateLabel}>Llegada estimada</Text>
              </View>
              <Text style={styles.dateValue}>{formatDate(trip.estimated_arrival_datetime)}</Text>
            </View>
          </View>

          {/* Trip Info */}
          <View style={styles.tripInfoSection}>
            <View style={styles.infoItem}>
              <View style={styles.infoIconRow}>
                <Ionicons name="car" size={20} color="#6b7280" />
                <Text style={styles.infoLabel}>Vehículo</Text>
              </View>
              <Text style={styles.infoValue}>
                {trip.car.brand} {trip.car.model}
              </Text>
              <Text style={styles.infoSubValue}>
                {trip.car.year} - {trip.car.color}
              </Text>
              <Text style={styles.infoSubValue}>Patente: {trip.car.plate}</Text>
            </View>

            <View style={styles.infoItem}>
              <View style={styles.infoIconRow}>
                <Ionicons name="people" size={20} color="#6b7280" />
                <Text style={styles.infoLabel}>Asientos</Text>
              </View>
              <Text style={styles.infoValue}>
                {trip.available_seats} de {trip.total_seats} disponibles
              </Text>
              <Text style={styles.infoSubValue}>
                {trip.total_seats - trip.available_seats} reservados
              </Text>
            </View>

            <View style={styles.infoItem}>
              <View style={styles.infoIconRow}>
                <Ionicons name="cash" size={20} color="#6b7280" />
                <Text style={styles.infoLabel}>Precio total</Text>
              </View>
              <Text style={styles.infoValue}>
                {formatPrice(trip.price_per_seat * trip.total_seats)}
              </Text>
              <Text style={styles.infoSubValue}>
                ({formatPrice(trip.price_per_seat)} × {trip.total_seats})
              </Text>
            </View>
          </View>

          {/* Driver Information */}
          {trip.driver && (
            <View style={styles.driverSection}>
              <View style={styles.driverIconRow}>
                <Ionicons name="person" size={20} color="#0ea5e9" />
                <Text style={styles.driverSectionTitle}>Conductor</Text>
              </View>
              <View style={styles.driverCard}>
                <View style={styles.driverAvatar}>
                  <Text style={styles.driverAvatarText}>
                    {trip.driver.name.charAt(0).toUpperCase()}
                  </Text>
                </View>
                <View style={styles.driverInfo}>
                  <Text style={styles.driverName}>{trip.driver.name}</Text>
                  <View style={styles.driverRating}>
                    <Ionicons name="star" size={16} color="#facc15" />
                    <Text style={styles.driverRatingText}>{trip.driver.rating.toFixed(1)}</Text>
                    <Text style={styles.driverTrips}>
                      • {trip.driver.total_trips} viajes realizados
                    </Text>
                  </View>
                  <Text style={styles.driverEmail}>{trip.driver.email}</Text>
                </View>
              </View>
            </View>
          )}

          {/* Description */}
          {trip.description && (
            <View style={styles.descriptionSection}>
              <Text style={styles.sectionTitle}>Descripción</Text>
              <Text style={styles.descriptionText}>{trip.description}</Text>
            </View>
          )}

          {/* Preferences */}
          <View style={styles.preferencesSection}>
            <Text style={styles.sectionTitle}>Preferencias del viaje</Text>
            <View style={styles.preferencesBadges}>
              <Badge variant={trip.preferences.pets_allowed ? 'default' : 'outline'}>
                🐕 Mascotas {trip.preferences.pets_allowed ? 'permitidas' : 'no permitidas'}
              </Badge>
              <Badge variant={trip.preferences.smoking_allowed ? 'default' : 'outline'}>
                🚬 Fumar {trip.preferences.smoking_allowed ? 'permitido' : 'no permitido'}
              </Badge>
              <Badge variant={trip.preferences.music_allowed ? 'default' : 'outline'}>
                🎵 Música {trip.preferences.music_allowed ? 'permitida' : 'no permitida'}
              </Badge>
            </View>
          </View>

          {/* Action Buttons */}
          <View style={styles.actionSection}>
            {trip.status === 'in_progress' ? (
              isOwner ? (
                <Button
                  title="Continuar Viaje"
                  size="lg"
                  onPress={() => navigation.navigate('DriverTracking', { tripId: trip.trip_id.toString() })}
                  style={styles.bookButton}
                />
              ) : (
                <Button
                  title="Rastrear Viaje"
                  size="lg"
                  onPress={() => navigation.navigate('TripTracking', { tripId: trip.trip_id.toString() })}
                  style={styles.bookButton}
                />
              )
            ) : (
              <>
                {isOwner ? (
                  <>
                    {trip.status === 'published' && trip.total_seats - trip.available_seats === 0 ? (
                      <View style={styles.ownerActions}>
                        <Button
                          title="Editar Viaje"
                          variant="outline"
                          onPress={() => (navigation.navigate as any)('EditTrip', { id: trip.trip_id.toString() })}
                          style={styles.editButton}
                        />
                        <Button
                          title={deleting ? 'Eliminando...' : 'Eliminar Viaje'}
                          variant="destructive"
                          onPress={handleDelete}
                          loading={deleting}
                          disabled={deleting}
                          style={styles.deleteButton}
                        />
                      </View>
                    ) : trip.total_seats - trip.available_seats > 0 ? (
                      <>
                        <View style={styles.infoBox}>
                          <Ionicons name="checkmark-circle" size={20} color="#22c55e" />
                          <Text style={styles.infoText}>
                            {trip.total_seats - trip.available_seats} pasajero{trip.total_seats - trip.available_seats > 1 ? 's' : ''} confirmado{trip.total_seats - trip.available_seats > 1 ? 's' : ''}
                          </Text>
                        </View>
                        <Button
                          title={starting ? 'Iniciando...' : 'Iniciar Viaje'}
                          size="lg"
                          onPress={handleStartTrip}
                          loading={starting}
                          disabled={starting}
                          style={styles.startButton}
                        />
                      </>
                    ) : null}
                  </>
                ) : (
                  <>
                    {trip.status === 'published' && trip.available_seats > 0 ? (
                      <Button
                        title="Reservar Asientos"
                        size="lg"
                        onPress={() => setIsBookingModalVisible(true)}
                        style={styles.bookButton}
                      />
                    ) : (
                      <View style={styles.unavailableBox}>
                        <Ionicons name="alert-circle" size={20} color="#6b7280" />
                        <Text style={styles.unavailableText}>
                          {trip.available_seats === 0
                            ? 'No hay asientos disponibles'
                            : 'Viaje no disponible para reservas'}
                        </Text>
                      </View>
                    )}
                  </>
                )}
              </>
            )}
            <Button
              title="Ver en Mapa"
              variant="outline"
              onPress={() => navigation.navigate('TripRouteMap', { origin: trip.origin.coordinates, destination: trip.destination.coordinates })}
              style={styles.mapButton}
            />
          </View>
        </Card>
      </ScrollView>

      {/* Booking Modal */}
      {trip && (
        <BookingModal
          trip={trip}
          visible={isBookingModalVisible}
          onClose={() => setIsBookingModalVisible(false)}
          onSuccess={handleBookingSuccess}
        />
      )}
    </View>
  )
}
