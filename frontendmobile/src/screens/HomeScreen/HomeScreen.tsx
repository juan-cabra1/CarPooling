/**
 * HomeScreen
 * Main landing screen with featured trips and CTAs
 */

import React, { useState, useEffect } from 'react'
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  FlatList,
  ActivityIndicator,
  RefreshControl,
} from 'react-native'
import { useNavigation } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { LinearGradient } from 'expo-linear-gradient'
import TripCard from '@/components/trips/TripCard'
import { searchService } from '@/services'
import { styles } from './HomeScreen.styles'
import { colors } from '@/styles/colors'
import type { SearchTrip } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type HomeScreenNavigationProp = StackNavigationProp<RootStackParamList, 'MainTabs'>

export default function HomeScreen() {
  const navigation = useNavigation<HomeScreenNavigationProp>()

  const [featuredTrips, setFeaturedTrips] = useState<SearchTrip[]>([])
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)

  const fetchFeaturedTrips = async () => {
    try {
      const response = await searchService.searchTrips({
        page: 1,
        limit: 10,
      })
      setFeaturedTrips(response?.trips || [])
    } catch (error) {
      console.error('Error fetching featured trips:', error)
      setFeaturedTrips([])
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  useEffect(() => {
    fetchFeaturedTrips()
  }, [])

  const onRefresh = () => {
    setRefreshing(true)
    fetchFeaturedTrips()
  }

  return (
    <ScrollView
      style={styles.container}
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} colors={[colors.primary[500]]} />
      }
    >
      {/* Hero Section */}
      <LinearGradient
        colors={[colors.primary[50], colors.white, colors.secondary[50]]}
        start={{ x: 0, y: 0 }}
        end={{ x: 1, y: 1 }}
        style={styles.hero}
      >
        {/* Hero Title */}
        <Text style={styles.heroTitle}>
          <Text style={styles.heroTitlePrimary}>Viaja </Text>
          <Text style={styles.heroTitleSecondary}>Compartiendo</Text>
        </Text>

        {/* Subtitle */}
        <Text style={styles.heroSubtitle}>
          Conecta con conductores y pasajeros. Ahorra dinero, conoce personas y cuida el medio ambiente.
        </Text>

        {/* CTA Buttons */}
        <View style={styles.ctaContainer}>
          <TouchableOpacity
            onPress={() => navigation.navigate('Search' as never)}
            style={styles.ctaButtonSearch}
            activeOpacity={0.8}
          >
            <Ionicons name="search" size={20} color={colors.white} />
            <Text style={styles.ctaButtonText}>Buscar Viajes</Text>
          </TouchableOpacity>

          <TouchableOpacity
            onPress={() => navigation.navigate('CreateTrip' as never)}
            style={styles.ctaButtonCreate}
            activeOpacity={0.8}
          >
            <Ionicons name="add-circle" size={20} color={colors.white} />
            <Text style={styles.ctaButtonText}>Publicar Viaje</Text>
          </TouchableOpacity>
        </View>

        {/* Quick Stats */}
        <View style={styles.statsContainer}>
          <View style={styles.statCard}>
            <Text style={[styles.statNumber, styles.statNumberPrimary]}>1000+</Text>
            <Text style={styles.statLabel}>Viajes realizados</Text>
          </View>

          <View style={styles.statCard}>
            <Text style={[styles.statNumber, styles.statNumberSecondary]}>500+</Text>
            <Text style={styles.statLabel}>Usuarios activos</Text>
          </View>

          <View style={styles.statCard}>
            <Text style={[styles.statNumber, styles.statNumberAccent]}>4.8★</Text>
            <Text style={styles.statLabel}>Calificación</Text>
          </View>
        </View>
      </LinearGradient>

      {/* Featured Trips Section */}
      <View style={styles.section}>
        <View style={styles.sectionHeader}>
          <Text style={styles.sectionTitle}>Viajes Destacados</Text>
          <TouchableOpacity onPress={() => navigation.navigate('Search' as never)}>
            <Text style={styles.sectionLink}>Ver todos</Text>
          </TouchableOpacity>
        </View>

        {loading ? (
          <View style={styles.loadingContainer}>
            <ActivityIndicator size="large" color={colors.primary[500]} />
            <Text style={styles.loadingText}>Cargando viajes...</Text>
          </View>
        ) : featuredTrips.length === 0 ? (
          <View style={styles.emptyContainer}>
            <Ionicons name="car-outline" size={64} color={colors.border} />
            <Text style={styles.emptyText}>
              No hay viajes disponibles en este momento.{'\n'}
              ¡Sé el primero en publicar uno!
            </Text>
          </View>
        ) : (
          <FlatList
            data={featuredTrips}
            keyExtractor={(item) => item.trip_id.toString()}
            renderItem={({ item }) => <TripCard trip={item} />}
            scrollEnabled={false}
          />
        )}
      </View>

      {/* Features Section */}
      <View style={styles.featuresGrid}>
        <View style={styles.featureCard}>
          <View style={[styles.featureIconContainer, styles.featureIconPrimary]}>
            <Ionicons name="wallet-outline" size={24} color={colors.primary[500]} />
          </View>
          <Text style={styles.featureTitle}>Ahorra Dinero</Text>
          <Text style={styles.featureDescription}>
            Comparte los gastos del viaje y reduce tus costos de transporte significativamente.
          </Text>
        </View>

        <View style={styles.featureCard}>
          <View style={[styles.featureIconContainer, styles.featureIconSecondary]}>
            <Ionicons name="flash-outline" size={24} color={colors.secondary[500]} />
          </View>
          <Text style={styles.featureTitle}>Rápido y Fácil</Text>
          <Text style={styles.featureDescription}>
            Encuentra o publica un viaje en minutos. Sistema simple e intuitivo.
          </Text>
        </View>

        <View style={styles.featureCard}>
          <View style={[styles.featureIconContainer, styles.featureIconGreen]}>
            <Ionicons name="leaf-outline" size={24} color={colors.success} />
          </View>
          <Text style={styles.featureTitle}>Cuida el Planeta</Text>
          <Text style={styles.featureDescription}>
            Reduce tu huella de carbono compartiendo vehículo con otros pasajeros.
          </Text>
        </View>
      </View>
    </ScrollView>
  )
}
