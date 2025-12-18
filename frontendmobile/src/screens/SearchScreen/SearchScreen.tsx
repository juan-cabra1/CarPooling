/**
 * SearchScreen
 * Search trips with filters and location autocomplete
 */

import React, { useState, useEffect } from 'react'
import {
  View,
  Text,
  ScrollView,
  FlatList,
  ActivityIndicator,
  TouchableOpacity,
  RefreshControl,
} from 'react-native'
import { useNavigation, useRoute } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Input, Card, Label } from '@/components/ui'
import { LocationInput } from '@/components/search'
import { searchService, getErrorMessage } from '@/services'
import { styles } from './SearchScreen.styles'
import type {
  SearchTrip,
  SearchQuery,
  LocationInput as LocationInputType,
  SearchSortBy,
  SearchSortOrder,
} from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type SearchScreenNavigationProp = StackNavigationProp<RootStackParamList>

// Componente TripCard usando estilos de StyleSheet
const TripCard = ({ trip, onPress }: { trip: SearchTrip; onPress: () => void }) => (
  <TouchableOpacity onPress={onPress} style={styles.tripCard}>
    <View style={styles.tripCardHeader}>
      <View style={styles.tripCardContent}>
        <Text style={styles.tripCardRoute}>
          {trip.origin.city} → {trip.destination.city}
        </Text>
        <Text style={styles.tripCardDate}>
          {new Date(trip.departure_datetime).toLocaleDateString('es-AR', {
            weekday: 'short',
            day: 'numeric',
            month: 'short',
            hour: '2-digit',
            minute: '2-digit'
          })}
        </Text>
      </View>
      <Text style={styles.tripCardPrice}>${trip.price_per_seat}</Text>
    </View>

    <View style={styles.tripCardFooter}>
      <View style={styles.tripCardInfo}>
        <Ionicons name="person" size={16} color="#6b7280" />
        <Text style={styles.tripCardInfoText}>
          {trip.available_seats} asiento{trip.available_seats !== 1 ? 's' : ''} disponible{trip.available_seats !== 1 ? 's' : ''}
        </Text>
      </View>

      {/* Mostrar rating del conductor si existe, sino mostrar conductor nuevo */}
      <View style={styles.tripCardInfo}>
        <Ionicons name="star" size={16} color="#f59e0b" />
        <Text style={styles.tripCardInfoText}>
          {/* Usar rating del conductor si está disponible, de lo contrario mostrar "Nuevo" */}
          {(trip as any).driver_rating ? (trip as any).driver_rating.toFixed(1) :
           (trip as any).driver?.rating ? (trip as any).driver.rating.toFixed(1) : 'Nuevo'}
        </Text>
      </View>
    </View>

    {/* Información adicional del conductor si está disponible */}
    {(trip as any).driver && (
      <View style={styles.tripCardDriver}>
        <Ionicons name="person-circle" size={16} color="#6b7280" />
        <Text style={styles.tripCardDriverText}>
          {(trip as any).driver.name || 'Conductor'}
        </Text>
      </View>
    )}
  </TouchableOpacity>
)

export default function SearchScreen() {
  const navigation = useNavigation<SearchScreenNavigationProp>()
  const route = useRoute()

  const [trips, setTrips] = useState<SearchTrip[]>([])
  const [loading, setLoading] = useState(false)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState('')
  const [showFilters, setShowFilters] = useState(false)
  const [total, setTotal] = useState(0)
  const [currentPage, setCurrentPage] = useState(1)

  // Form state
  const [origin, setOrigin] = useState<LocationInputType | null>(null)
  const [destination, setDestination] = useState<LocationInputType | null>(null)
  const [departureDate, setDepartureDate] = useState('')
  const [minSeats, setMinSeats] = useState('')
  const [maxPrice, setMaxPrice] = useState('')
  const [sortOption, setSortOption] = useState('departure_time|asc')

  // Parse route params on mount
  useEffect(() => {
    const params = route.params as any

    if (params?.origin_city && params?.origin_province) {
      setOrigin({
        city: params.origin_city,
        province: params.origin_province,
        address: params.origin_address || `${params.origin_city}, ${params.origin_province}`,
        coordinates:
          params.origin_lat && params.origin_lng
            ? {
                lat: parseFloat(params.origin_lat),
                lng: parseFloat(params.origin_lng),
              }
            : undefined,
      })
    }

    if (params?.destination_city && params?.destination_province) {
      setDestination({
        city: params.destination_city,
        province: params.destination_province,
        address: params.destination_address || `${params.destination_city}, ${params.destination_province}`,
        coordinates:
          params.destination_lat && params.destination_lng
            ? {
                lat: parseFloat(params.destination_lat),
                lng: parseFloat(params.destination_lng),
              }
            : undefined,
      })
    }

    if (params?.departure_date) setDepartureDate(params.departure_date)
    if (params?.min_seats) setMinSeats(params.min_seats.toString())
    if (params?.max_price) setMaxPrice(params.max_price.toString())
  }, [route.params])

  const handleSearch = async (page = 1) => {
    try {
      setLoading(true)
      setError('')
      setCurrentPage(page)

      // Build search query
      const query: SearchQuery = {
        page,
        limit: 20,
      }

      // Add location filters
      if (origin) {
        query.origin = origin
      }

      if (destination) {
        query.destination = destination
      }

      // Add other filters
      if (departureDate) query.departure_date = departureDate
      if (minSeats) query.min_seats = parseInt(minSeats)
      if (maxPrice) query.max_price = parseInt(maxPrice)

      // Parse and add sorting from combined option
      const [sortBy, sortOrder] = sortOption.split('|') as [SearchSortBy, SearchSortOrder]
      query.sort_by = sortBy
      query.sort_order = sortOrder

      const response = await searchService.searchTrips(query)
      setTrips(response.trips)
      setTotal(response.total)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  const handleClearFilters = () => {
    setOrigin(null)
    setDestination(null)
    setDepartureDate('')
    setMinSeats('')
    setMaxPrice('')
    setSortOption('departure_time|asc')
    setTrips([])
    setTotal(0)
  }

  const onRefresh = () => {
    setRefreshing(true)
    handleSearch(currentPage)
  }

  const renderTrip = ({ item }: { item: SearchTrip }) => (
    <TripCard
      trip={item}
      onPress={() => navigation.navigate('TripDetail', { id: item.trip_id.toString() })}
    />
  )

  // Custom Picker component
  const CustomPicker = ({
    value,
    onValueChange,
    options
  }: {
    value: string
    onValueChange: (value: string) => void
    options: Array<{ label: string; value: string }>
  }) => (
    <View style={styles.pickerContainer}>
      {options.map((option, index) => (
        <TouchableOpacity
          key={option.value}
          onPress={() => onValueChange(option.value)}
          style={[
            styles.pickerOption,
            index !== options.length - 1 && styles.pickerOptionBorder,
            value === option.value && styles.pickerOptionSelected,
          ]}
        >
          <Text style={[
            styles.pickerOptionText,
            value === option.value && styles.pickerOptionTextSelected,
          ]}>
            {option.label}
          </Text>
        </TouchableOpacity>
      ))}
    </View>
  )

  const renderHeader = () => (
    <View>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Buscar Viajes</Text>
        <Text style={styles.headerSubtitle}>Encuentra el viaje perfecto para tu destino</Text>
      </View>

      {/* Search Form */}
      <Card style={styles.searchCard}>
        <View style={styles.searchCardHeader}>
          <Ionicons name="search" size={20} color="#0ea5e9" />
          <Text style={styles.searchCardTitle}>Buscar</Text>
        </View>

        <View style={styles.searchForm}>
          {/* Location Inputs */}
          <LocationInput
            label="Origen"
            value={origin}
            onChange={setOrigin}
            placeholder="Ej: Córdoba, Buenos Aires, Mendoza"
          />

          <LocationInput
            label="Destino"
            value={destination}
            onChange={setDestination}
            placeholder="Ej: Medellín, Antioquia"
          />

          {/* Departure Date */}
          <View style={styles.inputGroup}>
            <Label>Fecha de salida</Label>
            <View style={styles.dateInputContainer}>
              <Ionicons
                name="calendar"
                size={20}
                color="#9ca3af"
                style={{ position: 'absolute', left: 12, zIndex: 1 }}
              />
              <Input
                value={departureDate}
                onChangeText={setDepartureDate}
                placeholder="YYYY-MM-DD"
                style={[styles.dateInput, { paddingLeft: 40 }]}
              />
            </View>
          </View>

          {/* Filters Toggle */}
          <TouchableOpacity
            style={styles.filtersToggle}
            onPress={() => setShowFilters(!showFilters)}
          >
            <Ionicons name="filter" size={18} color="#0ea5e9" />
            <Text style={styles.filtersToggleText}>
              {showFilters ? 'Ocultar filtros' : 'Más filtros'}
            </Text>
          </TouchableOpacity>

          {/* Advanced Filters */}
          {showFilters && (
            <View style={styles.advancedFilters}>
              {/* Min Seats */}
              <View style={styles.filterItem}>
                <Label>Asientos mínimos</Label>
                <View style={styles.filterInputContainer}>
                  <Ionicons
                    name="people"
                    size={20}
                    color="#9ca3af"
                    style={{ position: 'absolute', left: 12, zIndex: 1 }}
                  />
                  <Input
                    value={minSeats}
                    onChangeText={setMinSeats}
                    placeholder="Ej: 2"
                    keyboardType="numeric"
                    style={{ paddingLeft: 40 }}
                  />
                </View>
              </View>

              {/* Max Price */}
              <View style={styles.filterItem}>
                <Label>Precio máximo</Label>
                <View style={styles.filterInputContainer}>
                  <Ionicons
                    name="cash"
                    size={20}
                    color="#9ca3af"
                    style={{ position: 'absolute', left: 12, zIndex: 1 }}
                  />
                  <Input
                    value={maxPrice}
                    onChangeText={setMaxPrice}
                    placeholder="Ej: 5000"
                    keyboardType="numeric"
                    style={{ paddingLeft: 40 }}
                  />
                </View>
              </View>

              {/* Sort Option */}
              <View style={styles.filterItem}>
                <Label>Ordenar por</Label>
                <CustomPicker
                  value={sortOption}
                  onValueChange={setSortOption}
                  options={[
                    { label: 'Salida más temprana', value: 'departure_time|asc' },
                    { label: 'Salida más tardía', value: 'departure_time|desc' },
                    { label: 'Precio más bajo', value: 'price|asc' },
                    { label: 'Precio más alto', value: 'price|desc' },
                    { label: 'Mejor calificación', value: 'rating|desc' },
                    { label: 'Más popular', value: 'popularity|desc' },
                  ]}
                />
              </View>
            </View>
          )}

          {/* Action Buttons */}
          <View style={styles.actionButtons}>
            <Button
              title="Limpiar"
              variant="outline"
              onPress={handleClearFilters}
              style={styles.clearButton}
            />
            <Button
              title="Buscar"
              onPress={() => handleSearch()}
              loading={loading}
              disabled={loading}
              style={styles.searchButton}
            />
          </View>
        </View>
      </Card>

      {/* Error Message */}
      {error && (
        <View style={styles.errorContainer}>
          <Ionicons name="alert-circle" size={20} color="#dc2626" />
          <Text style={styles.errorText}>{error}</Text>
        </View>
      )}

      {/* Results Header */}
      {trips.length > 0 && (
        <View style={styles.resultsHeader}>
          <Text style={styles.resultsTitle}>
            {origin || destination
              ? `${total} ${total === 1 ? 'viaje encontrado' : 'viajes encontrados'}`
              : `Todos los viajes disponibles (${total})`}
          </Text>
        </View>
      )}
    </View>
  )

  const renderFooter = () => {
    if (total <= 20) return null

    return (
      <View style={styles.pagination}>
        <Button
          title="Anterior"
          variant="outline"
          size="sm"
          onPress={() => handleSearch(currentPage - 1)}
          disabled={currentPage === 1 || loading}
          style={styles.paginationButton}
        />
        <Text style={styles.paginationText}>
          Página {currentPage} de {Math.ceil(total / 20)}
        </Text>
        <Button
          title="Siguiente"
          variant="outline"
          size="sm"
          onPress={() => handleSearch(currentPage + 1)}
          disabled={currentPage >= Math.ceil(total / 20) || loading}
          style={styles.paginationButton}
        />
      </View>
    )
  }

  const renderEmpty = () => {
    if (loading) return null

    return (
      <View style={styles.emptyContainer}>
        <Ionicons name="search-outline" size={80} color="#cbd5e1" />
        <Text style={styles.emptyTitle}>No se encontraron viajes</Text>
        <Text style={styles.emptyText}>
          {origin || destination
            ? 'Intenta modificar tus criterios de búsqueda o explora otros destinos'
            : 'No hay viajes disponibles en este momento'}
        </Text>
        {(origin || destination) && (
          <Button
            title="Limpiar búsqueda"
            variant="outline"
            onPress={handleClearFilters}
            style={styles.emptyButton}
          />
        )}
      </View>
    )
  }

  if (loading && trips.length === 0) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={styles.loadingText}>Buscando viajes...</Text>
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <FlatList
        data={trips}
        keyExtractor={(item) => item.trip_id.toString()}
        renderItem={renderTrip}
        ListHeaderComponent={renderHeader}
        ListFooterComponent={renderFooter}
        ListEmptyComponent={renderEmpty}
        contentContainerStyle={{ padding: 16 }}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} colors={['#0ea5e9']} />
        }
      />
    </View>
  )
}
