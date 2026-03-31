/**
 * TripCard Component
 * Displays trip information in a card format
 */

import React from 'react'
import { View, Text, TouchableOpacity } from 'react-native'
import { useNavigation } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Card, CardHeader, CardContent, CardFooter, Badge } from '@/components/ui'
import { styles } from './TripCard.styles'
import type { SearchTrip, Trip } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

interface TripCardProps {
  trip: SearchTrip | Trip
}

type NavigationProp = StackNavigationProp<RootStackParamList>

export default function TripCard({ trip }: TripCardProps) {
  const navigation = useNavigation<NavigationProp>()

  // Check if trip is SearchTrip or Trip
  const isSearchTrip = 'driver' in trip && 'trip_id' in trip
  const tripId = isSearchTrip ? (trip as SearchTrip).trip_id : trip.id

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
    <Card style={styles.card}>
      {/* Header with Route */}
      <CardHeader style={styles.header}>
        <View style={styles.headerContent}>
          <View style={styles.routeContainer}>
            {/* Origin */}
            <View style={styles.routeRow}>
              <Ionicons name="location" size={16} color="#0ea5e9" />
              <Text style={styles.cityText}>
                {trip.origin.city}, {trip.origin.province}
              </Text>
            </View>

            <View style={styles.routeLine} />

            {/* Destination */}
            <View style={styles.routeRow}>
              <Ionicons name="location" size={16} color="#14b8a6" />
              <Text style={styles.cityText}>
                {trip.destination.city}, {trip.destination.province}
              </Text>
            </View>
          </View>

          {/* Price */}
          <View style={styles.priceContainer}>
            <Text style={styles.price}>{formatPrice(trip.price_per_seat)}</Text>
            <Text style={styles.priceLabel}>por persona</Text>
          </View>
        </View>
      </CardHeader>

      <CardContent style={styles.content}>
        {/* Driver Info - Only for SearchTrip */}
        {isSearchTrip && (
          <View style={styles.driverSection}>
            <View style={styles.driverAvatar}>
              <Text style={styles.driverAvatarText}>
                {(trip as SearchTrip).driver.name.charAt(0).toUpperCase()}
              </Text>
            </View>

            <View style={styles.driverInfo}>
              <Text style={styles.driverName}>{(trip as SearchTrip).driver.name}</Text>
              <View style={styles.driverRating}>
                <Ionicons name="star" size={16} color="#eab308" />
                <Text style={styles.driverRatingText}>
                  {(trip as SearchTrip).driver.rating.toFixed(1)}
                </Text>
                <Text style={styles.driverTrips}>
                  • {(trip as SearchTrip).driver.total_trips} viajes
                </Text>
              </View>
            </View>
          </View>
        )}

        {/* Trip Details */}
        <View style={styles.detailsGrid}>
          {/* Date & Time */}
          <View style={styles.detailRow}>
            <Ionicons name="calendar-outline" size={20} color="#64748b" />
            <Text style={styles.detailText}>{formatDate(trip.departure_datetime)}</Text>
          </View>

          {/* Available Seats */}
          <View style={styles.detailRow}>
            <Ionicons name="people-outline" size={20} color="#64748b" />
            <Text style={styles.detailText}>
              {trip.available_seats} asiento{trip.available_seats !== 1 ? 's' : ''} disponible
              {trip.available_seats !== 1 ? 's' : ''}
            </Text>
          </View>

          {/* Car */}
          {trip.car && (
            <View style={styles.detailRow}>
              <Ionicons name="car-outline" size={20} color="#64748b" />
              <Text style={styles.detailText}>
                {trip.car.brand} {trip.car.model}
              </Text>
            </View>
          )}
        </View>

        {/* Preferences */}
        {trip.preferences && Object.keys(trip.preferences).some((key) => trip.preferences[key as keyof typeof trip.preferences]) && (
          <View style={styles.preferencesSection}>
            <Text style={styles.preferencesTitle}>PREFERENCIAS</Text>
            <View style={styles.preferencesBadges}>
              {'no_smoking' in trip.preferences && trip.preferences.no_smoking === true && <Badge variant="secondary">Sin fumar</Badge>}
              {'pets_allowed' in trip.preferences && trip.preferences.pets_allowed && <Badge variant="secondary">Mascotas</Badge>}
              {'music_allowed' in trip.preferences && trip.preferences.music_allowed && <Badge variant="secondary">Música</Badge>}
            </View>
          </View>
        )}
      </CardContent>

      {/* Footer */}
      <CardFooter style={styles.footer}>
        <TouchableOpacity
          onPress={() => navigation.navigate('TripDetail', { id: tripId.toString() })}
          style={styles.button}
          activeOpacity={0.8}
        >
          <Text style={styles.buttonText}>Ver Detalles</Text>
        </TouchableOpacity>
      </CardFooter>
    </Card>
  )
}
