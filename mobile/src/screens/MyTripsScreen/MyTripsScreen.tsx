/**
 * MyTripsScreen
 * Display user's trips as driver with status tabs and lifecycle actions
 */

import React, { useState, useCallback } from 'react'
import {
  View,
  Text,
  FlatList,
  RefreshControl,
  Alert,
  TouchableOpacity,
  StyleSheet,
} from 'react-native'
import { useNavigation, useFocusEffect } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Badge, EmptyState, TripCardSkeleton } from '@/components/ui'
import { tripsService, getErrorMessage } from '@/services'
import { useAuth } from '@/context/AuthContext'
import { colors } from '@/styles/colors'
import type { Trip } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type MyTripsScreenNavigationProp = StackNavigationProp<RootStackParamList>

type TabFilter = 'all' | 'published' | 'in_progress' | 'completed' | 'cancelled'

const STATUS_TABS: { key: TabFilter; label: string }[] = [
  { key: 'all', label: 'Todos' },
  { key: 'published', label: 'Publicados' },
  { key: 'in_progress', label: 'En curso' },
  { key: 'completed', label: 'Completados' },
  { key: 'cancelled', label: 'Cancelados' },
]

const STATUS_CONFIG: Record<
  string,
  { label: string; variant: 'success' | 'secondary' | 'destructive' | 'warning' | 'outline' | 'default' }
> = {
  draft: { label: 'Borrador', variant: 'outline' },
  published: { label: 'Publicado', variant: 'success' },
  full: { label: 'Completo', variant: 'warning' },
  in_progress: { label: 'En curso', variant: 'secondary' },
  completed: { label: 'Completado', variant: 'secondary' },
  cancelled: { label: 'Cancelado', variant: 'destructive' },
}

function TripCard({ trip, onRefresh }: { trip: Trip; onRefresh: () => void }) {
  const navigation = useNavigation<MyTripsScreenNavigationProp>()
  const [actionLoading, setActionLoading] = useState(false)

  const formatDate = (dateString: string) =>
    new Intl.DateTimeFormat('es-AR', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    }).format(new Date(dateString))

  const statusConfig = STATUS_CONFIG[trip.status] ?? { label: trip.status, variant: 'outline' as const }

  const handleStatusUpdate = async (newStatus: string, confirmMsg: string) => {
    Alert.alert('Confirmar', confirmMsg, [
      { text: 'Cancelar', style: 'cancel' },
      {
        text: 'Confirmar',
        onPress: async () => {
          try {
            setActionLoading(true)
            await tripsService.updateTripStatus(trip.id.toString(), newStatus)
            onRefresh()
            if (newStatus === 'in_progress') {
              navigation.navigate('DriverTracking', { tripId: trip.id.toString() })
            }
          } catch (err) {
            Alert.alert('Error', getErrorMessage(err))
          } finally {
            setActionLoading(false)
          }
        },
      },
    ])
  }

  return (
    <View style={styles.card}>
      {/* Header */}
      <TouchableOpacity
        onPress={() => navigation.navigate('TripDetail', { id: trip.id.toString() })}
        activeOpacity={0.8}
      >
        <View style={styles.cardHeader}>
          <View style={{ flex: 1 }}>
            <Text style={styles.route} numberOfLines={1}>
              {trip.origin.city} → {trip.destination.city}
            </Text>
            <Text style={styles.date}>{formatDate(trip.departure_datetime)}</Text>
          </View>
          <View style={{ alignItems: 'flex-end', gap: 4 }}>
            <Text style={styles.price}>${trip.price_per_seat}</Text>
            <Badge variant={statusConfig.variant}>{statusConfig.label}</Badge>
          </View>
        </View>

        {/* Meta row */}
        <View style={styles.metaRow}>
          <View style={styles.metaItem}>
            <Ionicons name="people-outline" size={14} color={colors.mutedForeground} />
            <Text style={styles.metaText}>
              {trip.available_seats}/{trip.total_seats} disponibles
            </Text>
          </View>
          <View style={styles.metaItem}>
            <Ionicons name="location-outline" size={14} color={colors.mutedForeground} />
            <Text style={styles.metaText} numberOfLines={1}>
              {trip.origin.province} → {trip.destination.province}
            </Text>
          </View>
        </View>
      </TouchableOpacity>

      {/* Actions */}
      <View style={styles.actionsRow}>
        <Button
          title="Detalles"
          variant="outline"
          size="sm"
          onPress={() => navigation.navigate('TripDetail', { id: trip.id.toString() })}
          style={{ flex: 1 }}
        />
        {trip.status === 'published' && (
          <Button
            title={actionLoading ? '...' : 'Iniciar'}
            size="sm"
            disabled={actionLoading}
            onPress={() => handleStatusUpdate('in_progress', '¿Iniciar este viaje ahora?')}
            style={{ flex: 1 }}
          />
        )}
        {trip.status === 'in_progress' && (
          <>
            <Button
              title="Navegar"
              size="sm"
              onPress={() => navigation.navigate('DriverTracking', { tripId: trip.id.toString() })}
              style={{ flex: 1 }}
            />
            <Button
              title={actionLoading ? '...' : 'Completar'}
              size="sm"
              disabled={actionLoading}
              onPress={() => handleStatusUpdate('completed', '¿Marcar este viaje como completado?')}
              style={{ flex: 1 }}
            />
          </>
        )}
        {(trip.status === 'published' || trip.status === 'draft') && (
          <Button
            title={actionLoading ? '...' : 'Cancelar'}
            variant="destructive"
            size="sm"
            disabled={actionLoading}
            onPress={() => handleStatusUpdate('cancelled', '¿Cancelar este viaje? Esta acción no se puede deshacer.')}
            style={{ flex: 1 }}
          />
        )}
      </View>
    </View>
  )
}

export default function MyTripsScreen() {
  const { user } = useAuth()
  const [trips, setTrips] = useState<Trip[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState('')
  const [activeTab, setActiveTab] = useState<TabFilter>('all')
  const navigation = useNavigation<MyTripsScreenNavigationProp>()

  const fetchTrips = useCallback(async () => {
    try {
      setError('')
      if (!user?.id) return
      const data = await tripsService.getMyTrips(user.id, 1, 50)
      setTrips(data.trips ?? [])
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }, [user?.id])

  useFocusEffect(
    useCallback(() => {
      fetchTrips()
    }, [fetchTrips])
  )

  const filteredTrips = activeTab === 'all' ? trips : trips.filter(t => t.status === activeTab)

  if (loading) {
    return (
      <View style={{ flex: 1, backgroundColor: colors.background, padding: 16 }}>
        {[1, 2, 3].map(i => <TripCardSkeleton key={i} />)}
      </View>
    )
  }

  if (error) {
    return (
      <EmptyState
        icon="alert-circle-outline"
        title="Error al cargar viajes"
        subtitle={error}
        actionLabel="Reintentar"
        onAction={() => { setLoading(true); fetchTrips() }}
      />
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
            const count = item.key === 'all' ? trips.length : trips.filter(t => t.status === item.key).length
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

      {/* Trip list */}
      {filteredTrips.length === 0 ? (
        <EmptyState
          icon="car-outline"
          title={activeTab === 'all' ? 'No tienes viajes publicados' : `No hay viajes ${STATUS_TABS.find(t => t.key === activeTab)?.label.toLowerCase()}`}
          subtitle={activeTab === 'all' ? 'Publica tu primer viaje y comparte tus rutas' : undefined}
          actionLabel={activeTab === 'all' ? 'Publicar Viaje' : undefined}
          onAction={activeTab === 'all' ? () => navigation.navigate('CreateTrip' as never) : undefined}
        />
      ) : (
        <FlatList
          data={filteredTrips}
          keyExtractor={item => item.id.toString()}
          renderItem={({ item }) => <TripCard trip={item} onRefresh={() => { setLoading(true); fetchTrips() }} />}
          contentContainerStyle={{ padding: 16, gap: 12 }}
          refreshControl={
            <RefreshControl
              refreshing={refreshing}
              onRefresh={() => { setRefreshing(true); fetchTrips() }}
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
    marginBottom: 10,
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
  },
  price: {
    fontSize: 18,
    fontWeight: '700',
    color: colors.primary[500],
    marginBottom: 4,
  },
  metaRow: {
    flexDirection: 'row',
    gap: 16,
    marginBottom: 12,
  },
  metaItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
    flex: 1,
  },
  metaText: {
    fontSize: 12,
    color: colors.mutedForeground,
    flex: 1,
  },
  actionsRow: {
    flexDirection: 'row',
    gap: 8,
    paddingTop: 12,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
})
