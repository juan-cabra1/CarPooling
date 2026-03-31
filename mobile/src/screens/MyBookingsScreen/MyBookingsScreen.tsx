/**
 * MyBookingsScreen
 * Display user's bookings as passenger with status tabs and rate driver feature
 */

import React, { useState, useCallback } from 'react'
import {
  View,
  Text,
  FlatList,
  RefreshControl,
  Alert,
  TouchableOpacity,
  Modal,
  StyleSheet,
  Linking,
} from 'react-native'
import { useNavigation, useFocusEffect } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Badge, EmptyState, TripCardSkeleton, RatingStars } from '@/components/ui'
import { bookingsService, searchService, tripsService, ratingsService, getErrorMessage } from '@/services'
import { getPaymentByBooking, createPaymentPreference } from '@/services/paymentsService'
import { useAuth } from '@/context/AuthContext'
import type { Payment } from '@/types'
import { colors } from '@/styles/colors'
import type { Booking, SearchTrip } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type MyBookingsScreenNavigationProp = StackNavigationProp<RootStackParamList>

interface BookingWithTrip extends Omit<Booking, 'trip'> {
  trip?: SearchTrip
}

type TabFilter = 'all' | 'pending' | 'confirmed' | 'completed' | 'cancelled'

const STATUS_TABS: { key: TabFilter; label: string }[] = [
  { key: 'all', label: 'Todas' },
  { key: 'confirmed', label: 'Confirmadas' },
  { key: 'pending', label: 'Pendientes' },
  { key: 'completed', label: 'Completadas' },
  { key: 'cancelled', label: 'Canceladas' },
]

const BOOKING_STATUS_CONFIG: Record<
  string,
  { label: string; variant: 'success' | 'warning' | 'destructive' | 'secondary' | 'outline' | 'default' }
> = {
  pending: { label: 'Pendiente', variant: 'warning' },
  confirmed: { label: 'Confirmada', variant: 'success' },
  cancelled: { label: 'Cancelada', variant: 'destructive' },
  completed: { label: 'Completada', variant: 'secondary' },
}

function RatingModal({
  visible,
  driverName,
  onClose,
  onSubmit,
}: {
  visible: boolean
  driverName: string
  onClose: () => void
  onSubmit: (score: number, comment: string) => void
}) {
  const [selectedScore, setSelectedScore] = useState(5)
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async () => {
    setSubmitting(true)
    onSubmit(selectedScore, '')
    setSubmitting(false)
  }

  return (
    <Modal visible={visible} transparent animationType="slide" onRequestClose={onClose}>
      <View style={styles.modalOverlay}>
        <View style={styles.modalContainer}>
          <Text style={styles.modalTitle}>Calificar conductor</Text>
          <Text style={styles.modalSubtitle}>{driverName}</Text>

          <View style={styles.starsRow}>
            {[1, 2, 3, 4, 5].map(star => (
              <TouchableOpacity key={star} onPress={() => setSelectedScore(star)} activeOpacity={0.7}>
                <Ionicons
                  name={star <= selectedScore ? 'star' : 'star-outline'}
                  size={40}
                  color={star <= selectedScore ? '#facc15' : '#d1d5db'}
                  style={{ marginHorizontal: 4 }}
                />
              </TouchableOpacity>
            ))}
          </View>
          <Text style={styles.scoreLabel}>
            {selectedScore === 1 ? 'Muy malo' : selectedScore === 2 ? 'Malo' : selectedScore === 3 ? 'Regular' : selectedScore === 4 ? 'Bueno' : 'Excelente'}
          </Text>

          <View style={styles.modalActions}>
            <Button title="Cancelar" variant="outline" onPress={onClose} style={{ flex: 1 }} />
            <Button title={submitting ? 'Enviando...' : 'Enviar'} disabled={submitting} onPress={handleSubmit} style={{ flex: 1 }} />
          </View>
        </View>
      </View>
    </Modal>
  )
}

function BookingCard({ booking, onRefresh }: { booking: BookingWithTrip; onRefresh: () => void }) {
  const navigation = useNavigation<MyBookingsScreenNavigationProp>()
  const { user } = useAuth()
  const [ratingModalVisible, setRatingModalVisible] = useState(false)
  const [cancelling, setCancelling] = useState(false)
  const [payment, setPayment] = useState<Payment | null | 'loading'>('loading')
  const [paying, setPaying] = useState(false)

  React.useEffect(() => {
    getPaymentByBooking(booking.id)
      .then(p => setPayment(p))
      .catch(() => setPayment(null))
  }, [booking.id])

  const handlePay = async () => {
    if (!user || !booking.trip) return
    try {
      setPaying(true)
      const pref = await createPaymentPreference({
        booking_id: booking.id,
        trip_id: booking.trip_id,
        driver_id: String(booking.trip.driver_id),
        seats_count: booking.seats_requested,
        price_per_seat: Math.round((booking.total_price / booking.seats_requested) * 100),
        origin: `${booking.trip.origin.city}, ${booking.trip.origin.province}`,
        destination: `${booking.trip.destination.city}, ${booking.trip.destination.province}`,
        departure_at: booking.trip.departure_datetime,
        payer_email: user.email,
        return_base_url: 'carpooling:/',
      })
      await Linking.openURL(pref.init_point)
    } catch (err) {
      Alert.alert('Error', getErrorMessage(err))
    } finally {
      setPaying(false)
    }
  }

  const paymentStatus = payment === 'loading' ? null : payment?.status
  const isPendingPayment = booking.status !== 'cancelled' && (payment === null || paymentStatus === 'pending')

  const formatDate = (dateString: string) =>
    new Intl.DateTimeFormat('es-AR', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    }).format(new Date(dateString))

  const formatPrice = (price: number) =>
    new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS',
      minimumFractionDigits: 0,
    }).format(price)

  const statusConfig = BOOKING_STATUS_CONFIG[booking.status] ?? { label: booking.status, variant: 'outline' as const }

  const handleCancel = () => {
    Alert.alert(
      'Cancelar Reserva',
      '¿Estás seguro de que quieres cancelar esta reserva?',
      [
        { text: 'No', style: 'cancel' },
        {
          text: 'Sí, cancelar',
          style: 'destructive',
          onPress: async () => {
            try {
              setCancelling(true)
              await bookingsService.cancelBooking(booking.id, 'Cancelado por el pasajero')
              onRefresh()
            } catch (err) {
              Alert.alert('Error', getErrorMessage(err))
            } finally {
              setCancelling(false)
            }
          },
        },
      ]
    )
  }

  const handleRatingSubmit = async (score: number, _comment: string) => {
    if (!user || !booking.trip) return
    try {
      await ratingsService.rateDriver(user.id, booking.trip.driver_id, booking.trip_id, score)
      setRatingModalVisible(false)
      Alert.alert('¡Gracias!', 'Tu calificación fue enviada correctamente.')
      onRefresh()
    } catch (err) {
      setRatingModalVisible(false)
      Alert.alert('Error', getErrorMessage(err))
    }
  }

  const isCompleted = booking.status === 'completed' || booking.trip?.status === 'completed'
  const canCancel = booking.status === 'pending' || booking.status === 'confirmed'

  return (
    <>
      <View style={styles.card}>
        {/* Trip info */}
        <TouchableOpacity
          onPress={() => navigation.navigate('TripDetail', { id: booking.trip_id })}
          activeOpacity={0.8}
        >
          {booking.trip ? (
            <View style={styles.cardHeader}>
              <View style={{ flex: 1 }}>
                <Text style={styles.route} numberOfLines={1}>
                  {booking.trip.origin.city} → {booking.trip.destination.city}
                </Text>
                <Text style={styles.date}>{formatDate(booking.trip.departure_datetime)}</Text>
                {booking.trip.driver?.name ? (
                  <View style={styles.driverRow}>
                    <Ionicons name="person-circle-outline" size={14} color={colors.mutedForeground} />
                    <Text style={styles.driverName}>{booking.trip.driver.name}</Text>
                    {booking.trip.driver.rating > 0 && (
                      <RatingStars rating={booking.trip.driver.rating} size={12} showValue />
                    )}
                  </View>
                ) : null}
              </View>
              <View style={{ alignItems: 'flex-end', gap: 6 }}>
                <Text style={styles.price}>{formatPrice(booking.total_price)}</Text>
                <Text style={styles.seats}>
                  {booking.seats_requested} asiento{booking.seats_requested > 1 ? 's' : ''}
                </Text>
                <Badge variant={statusConfig.variant}>{statusConfig.label}</Badge>
                {payment !== 'loading' && paymentStatus === 'approved' && (
                  <Badge variant="success">Pagado</Badge>
                )}
                {payment !== 'loading' && isPendingPayment && (
                  <Badge variant="warning">Sin pagar</Badge>
                )}
              </View>
            </View>
          ) : (
            <View style={styles.cardHeader}>
              <Text style={styles.route}>Viaje {booking.trip_id.slice(-6)}</Text>
              <Badge variant={statusConfig.variant}>{statusConfig.label}</Badge>
            </View>
          )}
        </TouchableOpacity>

        {/* Actions */}
        <View style={styles.actionsRow}>
          <Button
            title="Ver viaje"
            variant="outline"
            size="sm"
            onPress={() => navigation.navigate('TripDetail', { id: booking.trip_id })}
            style={{ flex: 1 }}
          />
          {booking.trip?.status === 'in_progress' && (
            <Button
              title="Seguir viaje"
              size="sm"
              onPress={() => navigation.navigate('TripTracking', { tripId: booking.trip_id })}
              style={{ flex: 1 }}
            />
          )}
          {isCompleted && booking.trip?.driver?.name && (
            <Button
              title="Calificar"
              size="sm"
              onPress={() => setRatingModalVisible(true)}
              style={{ flex: 1 }}
            />
          )}
          {isPendingPayment && booking.status !== 'cancelled' && booking.trip && (
            <Button
              title={paying ? '...' : 'Pagar ahora'}
              size="sm"
              disabled={paying}
              onPress={handlePay}
              style={{ flex: 1, backgroundColor: '#22c55e' }}
            />
          )}
          {canCancel && (
            <Button
              title={cancelling ? '...' : 'Cancelar'}
              variant="destructive"
              size="sm"
              disabled={cancelling}
              onPress={handleCancel}
              style={{ flex: 1 }}
            />
          )}
        </View>
      </View>

      {booking.trip?.driver?.name && (
        <RatingModal
          visible={ratingModalVisible}
          driverName={booking.trip.driver.name}
          onClose={() => setRatingModalVisible(false)}
          onSubmit={handleRatingSubmit}
        />
      )}
    </>
  )
}

export default function MyBookingsScreen() {
  const navigation = useNavigation<MyBookingsScreenNavigationProp>()
  const [bookings, setBookings] = useState<BookingWithTrip[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [activeTab, setActiveTab] = useState<TabFilter>('all')

  const fetchBookings = useCallback(async () => {
    try {
      const data = await bookingsService.getMyBookings()
      const bookingsData = data.bookings || []

      const enriched = await Promise.all(
        bookingsData.map(async (booking): Promise<BookingWithTrip> => {
          try {
            const trip = await searchService.getTripDetails(booking.trip_id)
            return { ...booking, trip }
          } catch (err: any) {
            if (err?.response?.status === 404) {
              try {
                const tripData = await tripsService.getTripById(booking.trip_id)
                return {
                  ...booking,
                  trip: {
                    ...tripData,
                    trip_id: tripData.id,
                    driver: { id: tripData.driver_id, name: '', email: '', rating: 0, total_trips: 0 },
                  } as SearchTrip,
                }
              } catch {
                return booking as BookingWithTrip
              }
            }
            return booking as BookingWithTrip
          }
        })
      )

      setBookings(enriched)
    } catch {
      setBookings([])
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }, [])

  useFocusEffect(useCallback(() => { fetchBookings() }, [fetchBookings]))

  const filteredBookings = activeTab === 'all' ? bookings : bookings.filter(b => b.status === activeTab)

  if (loading) {
    return (
      <View style={{ flex: 1, backgroundColor: colors.background, padding: 16 }}>
        {[1, 2, 3].map(i => <TripCardSkeleton key={i} />)}
      </View>
    )
  }

  return (
    <View style={{ flex: 1, backgroundColor: colors.background }}>
      {/* Status tabs */}
      <View style={styles.tabsContainer}>
        <FlatList
          horizontal
          data={STATUS_TABS}
          keyExtractor={item => item.key}
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={{ paddingHorizontal: 16, gap: 8 }}
          renderItem={({ item }) => {
            const count = item.key === 'all' ? bookings.length : bookings.filter(b => b.status === item.key).length
            const isActive = activeTab === item.key
            return (
              <TouchableOpacity
                onPress={() => setActiveTab(item.key)}
                style={[styles.tab, isActive && styles.tabActive]}
                activeOpacity={0.7}
              >
                <Text style={[styles.tabText, isActive && styles.tabTextActive]}>
                  {item.label}{count > 0 ? ` (${count})` : ''}
                </Text>
              </TouchableOpacity>
            )
          }}
        />
      </View>

      {filteredBookings.length === 0 ? (
        <EmptyState
          icon="ticket-outline"
          title={activeTab === 'all' ? 'No tienes reservas' : `No hay reservas ${STATUS_TABS.find(t => t.key === activeTab)?.label.toLowerCase()}`}
          subtitle={activeTab === 'all' ? 'Busca un viaje y reserva tu asiento' : undefined}
          actionLabel={activeTab === 'all' ? 'Buscar Viajes' : undefined}
          onAction={activeTab === 'all' ? () => navigation.navigate('MainTabs', { screen: 'Search' }) : undefined}
        />
      ) : (
        <FlatList
          data={filteredBookings}
          keyExtractor={item => item.id}
          renderItem={({ item }) => (
            <BookingCard booking={item} onRefresh={() => { setLoading(true); fetchBookings() }} />
          )}
          contentContainerStyle={{ padding: 16, gap: 12 }}
          refreshControl={
            <RefreshControl
              refreshing={refreshing}
              onRefresh={() => { setRefreshing(true); fetchBookings() }}
              colors={[colors.primary[500]]}
            />
          }
        />
      )}
    </View>
  )
}

const styles = StyleSheet.create({
  tabsContainer: {
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    backgroundColor: colors.background,
  },
  tab: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 9999,
    borderWidth: 1,
    borderColor: colors.border,
    backgroundColor: colors.white,
  },
  tabActive: {
    backgroundColor: colors.primary[500],
    borderColor: colors.primary[500],
  },
  tabText: {
    fontSize: 13,
    fontWeight: '500',
    color: colors.mutedForeground,
  },
  tabTextActive: {
    color: colors.white,
    fontWeight: '600',
  },
  card: {
    backgroundColor: colors.card,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: colors.border,
    padding: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 4,
    elevation: 2,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 12,
    gap: 8,
  },
  route: {
    fontSize: 16,
    fontWeight: '700',
    color: colors.foreground,
    marginBottom: 4,
  },
  date: {
    fontSize: 13,
    color: colors.mutedForeground,
    marginBottom: 4,
  },
  driverRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    marginTop: 2,
  },
  driverName: {
    fontSize: 12,
    color: colors.mutedForeground,
  },
  price: {
    fontSize: 17,
    fontWeight: '700',
    color: colors.primary[500],
  },
  seats: {
    fontSize: 12,
    color: colors.mutedForeground,
  },
  actionsRow: {
    flexDirection: 'row',
    gap: 8,
    paddingTop: 12,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.5)',
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  modalContainer: {
    backgroundColor: colors.background,
    borderRadius: 16,
    padding: 24,
    width: '100%',
    alignItems: 'center',
  },
  modalTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: colors.foreground,
    marginBottom: 4,
  },
  modalSubtitle: {
    fontSize: 14,
    color: colors.mutedForeground,
    marginBottom: 24,
  },
  starsRow: {
    flexDirection: 'row',
    marginBottom: 12,
  },
  scoreLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: colors.foreground,
    marginBottom: 24,
  },
  modalActions: {
    flexDirection: 'row',
    gap: 12,
    width: '100%',
  },
})
