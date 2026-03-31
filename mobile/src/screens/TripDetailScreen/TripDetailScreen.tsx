/**
 * TripDetailScreen
 * Display complete trip details and allow booking
 */

import React, { useState, useEffect, useRef } from 'react'
import {
  View,
  Text,
  ScrollView,
  ActivityIndicator,
  TouchableOpacity,
  Alert,
  TextInput,
  KeyboardAvoidingView,
  Platform,
  FlatList,
} from 'react-native'
import { useNavigation, useRoute } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Badge, Card, Avatar, RatingStars } from '@/components/ui'
import { BookingModal } from '@/components/modals'
import TripRouteMap from '@/components/map/TripRouteMap'
import { tripsService, searchService, getErrorMessage } from '@/services'
import { tripsApi } from '@/services/api'
import { useAuth } from '@/context/AuthContext'
import { styles } from './TripDetailScreen.styles'
import type { SearchTrip } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type TripDetailScreenNavigationProp = StackNavigationProp<RootStackParamList>

export default function TripDetailScreen() {
  const navigation = useNavigation<TripDetailScreenNavigationProp>()
  const route = useRoute()
  const { user, token } = useAuth()

  const { id } = route.params as { id: string }

  const [trip, setTrip] = useState<SearchTrip | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [deleting, setDeleting] = useState(false)
  const [isBookingModalVisible, setIsBookingModalVisible] = useState(false)

  // Chat state
  interface ChatMessage {
    id: string
    user_id: number
    user_name: string
    message: string
    created_at: string
  }
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [newMessage, setNewMessage] = useState('')
  const [sendingMessage, setSendingMessage] = useState(false)
  const [wsConnected, setWsConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (id) {
      fetchTrip()
    }
  }, [id])

  // WebSocket chat — only when authenticated
  useEffect(() => {
    if (!id || !user || !token) return

    // Load message history via HTTP
    const fetchMessages = async () => {
      try {
        const res = await tripsApi.get(`/trips/${id}/messages`)
        setMessages(res.data.messages || [])
      } catch {
        // Silently fail
      }
    }
    fetchMessages()

    // Connect WebSocket for real-time updates
    const wsUrl = `wss://trips-api-production.up.railway.app/trips/${id}/ws?token=${encodeURIComponent(token)}`
    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onopen = () => setWsConnected(true)

    ws.onmessage = (event) => {
      try {
        const newMsg: ChatMessage = JSON.parse(event.data)
        setMessages(prev =>
          prev.find(m => m.id === newMsg.id) ? prev : [...prev, newMsg]
        )
      } catch {}
    }

    ws.onerror = () => setWsConnected(false)

    ws.onclose = () => {
      setWsConnected(false)
      wsRef.current = null
    }

    return () => {
      ws.close()
    }
  }, [id, user, token])

  const handleSendMessage = async () => {
    if (!newMessage.trim() || !id) return
    setSendingMessage(true)
    try {
      await tripsApi.post(`/trips/${id}/messages`, { message: newMessage.trim() })
      setNewMessage('')
      // WebSocket will push the new message — no need to re-fetch
    } catch (err) {
      Alert.alert('Error', getErrorMessage(err))
    } finally {
      setSendingMessage(false)
    }
  }

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
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      keyboardVerticalOffset={Platform.OS === 'ios' ? 90 : 0}
    >
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

          {/* Route Map */}
          {trip.origin.coordinates && trip.destination.coordinates && (() => {
            const oc = trip.origin.coordinates as any
            const dc = trip.destination.coordinates as any
            const originCoords = oc.type === 'Point'
              ? { lat: oc.coordinates[1], lng: oc.coordinates[0] }
              : { lat: oc.lat, lng: oc.lng }
            const destCoords = dc.type === 'Point'
              ? { lat: dc.coordinates[1], lng: dc.coordinates[0] }
              : { lat: dc.lat, lng: dc.lng }
            return (
              <View style={{ height: 200, marginHorizontal: -16, marginBottom: 16, overflow: 'hidden' }}>
                <TripRouteMap origin={originCoords} destination={destCoords} />
              </View>
            )
          })()}

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
                <Avatar name={trip.driver.name || 'C'} size={48} />
                <View style={styles.driverInfo}>
                  <Text style={styles.driverName}>{trip.driver.name}</Text>
                  {trip.driver.rating > 0 ? (
                    <View style={styles.driverRating}>
                      <RatingStars rating={trip.driver.rating} size={14} showValue />
                      <Text style={styles.driverTrips}>
                        • {trip.driver.total_trips} viajes
                      </Text>
                    </View>
                  ) : (
                    <Text style={styles.driverTrips}>{trip.driver.total_trips} viajes realizados</Text>
                  )}
                  {trip.driver.email ? (
                    <Text style={styles.driverEmail}>{trip.driver.email}</Text>
                  ) : null}
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
                  <View style={styles.warningBox}>
                    <Ionicons name="alert-circle" size={20} color="#d97706" />
                    <Text style={styles.warningText}>
                      Este viaje tiene reservas activas y no puede ser editado o eliminado
                    </Text>
                  </View>
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
          </View>
          {/* Chat Section — only for authenticated users */}
          {user && (
            <View style={{ borderTopWidth: 1, borderTopColor: '#e5e7eb', paddingVertical: 24 }}>
              <View style={{ flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
                <Text style={styles.sectionTitle}>Chat del viaje</Text>
                <View style={{ flexDirection: 'row', alignItems: 'center', gap: 4 }}>
                  <View style={{ width: 8, height: 8, borderRadius: 4, backgroundColor: wsConnected ? '#22c55e' : '#d1d5db' }} />
                  <Text style={{ fontSize: 11, color: wsConnected ? '#16a34a' : '#9ca3af' }}>
                    {wsConnected ? 'En vivo' : 'Conectando...'}
                  </Text>
                </View>
              </View>

              {messages.length === 0 ? (
                <Text style={{ color: '#9ca3af', fontSize: 14, textAlign: 'center', paddingVertical: 16 }}>
                  No hay mensajes aún. ¡Sé el primero en escribir!
                </Text>
              ) : (
                <FlatList
                  data={[...messages].reverse()}
                  keyExtractor={item => item.id}
                  scrollEnabled={false}
                  renderItem={({ item }) => {
                    const isMe = item.user_id === user.id
                    return (
                      <View style={{
                        marginBottom: 12,
                        alignItems: isMe ? 'flex-end' : 'flex-start',
                      }}>
                        {!isMe && (
                          <Text style={{ fontSize: 12, color: '#6b7280', marginBottom: 2 }}>
                            {item.user_name}
                          </Text>
                        )}
                        <View style={{
                          maxWidth: '80%',
                          backgroundColor: isMe ? '#0ea5e9' : '#f3f4f6',
                          borderRadius: 12,
                          paddingHorizontal: 12,
                          paddingVertical: 8,
                        }}>
                          <Text style={{ color: isMe ? '#ffffff' : '#111827', fontSize: 14 }}>
                            {item.message}
                          </Text>
                          <Text style={{ fontSize: 10, color: isMe ? '#bae6fd' : '#9ca3af', marginTop: 2 }}>
                            {new Date(item.created_at).toLocaleTimeString('es-AR', { hour: '2-digit', minute: '2-digit' })}
                          </Text>
                        </View>
                      </View>
                    )
                  }}
                />
              )}
            </View>
          )}
        </Card>
      </ScrollView>

      {/* Message Input — only for authenticated users */}
      {user && (
        <View style={{
          flexDirection: 'row',
          alignItems: 'center',
          paddingHorizontal: 16,
          paddingVertical: 8,
          borderTopWidth: 1,
          borderTopColor: '#e5e7eb',
          backgroundColor: '#ffffff',
          gap: 8,
        }}>
          <TextInput
            value={newMessage}
            onChangeText={setNewMessage}
            placeholder="Escribe un mensaje..."
            placeholderTextColor="#9ca3af"
            style={{
              flex: 1,
              borderWidth: 1,
              borderColor: '#d1d5db',
              borderRadius: 24,
              paddingHorizontal: 16,
              paddingVertical: 10,
              fontSize: 14,
              color: '#111827',
              backgroundColor: '#f9fafb',
            }}
            multiline={false}
            returnKeyType="send"
            onSubmitEditing={handleSendMessage}
          />
          <TouchableOpacity
            onPress={handleSendMessage}
            disabled={sendingMessage || !newMessage.trim()}
            style={{
              backgroundColor: sendingMessage || !newMessage.trim() ? '#93c5fd' : '#0ea5e9',
              borderRadius: 24,
              width: 44,
              height: 44,
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <Ionicons name="send" size={20} color="#ffffff" />
          </TouchableOpacity>
        </View>
      )}

      {/* Booking Modal */}
      {trip && (
        <BookingModal
          trip={trip}
          visible={isBookingModalVisible}
          onClose={() => setIsBookingModalVisible(false)}
          onSuccess={handleBookingSuccess}
        />
      )}
    </KeyboardAvoidingView>
  )
}
