/**
 * LocationInput Component (Mobile)
 * Photon (Komoot) Autocomplete integration for structured location input
 *
 * Features:
 * - Autocompletes city/address using Photon API (OpenStreetMap data)
 * - Extracts city, province, and coordinates automatically
 * - Validates required fields (city and province)
 * - Restricts to Argentina for relevant results
 * - No API key required (free and open-source)
 */

import React, { useState, useEffect, useRef } from 'react'
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  FlatList,
  ActivityIndicator,
  Modal,
} from 'react-native'
import { Ionicons } from '@expo/vector-icons'
import { Input, Label } from '@/components/ui'
import { styles } from './LocationInput.styles'
import type { LocationInput as LocationInputType } from '@/types'

interface LocationInputProps {
  label: string
  value: LocationInputType | null
  onChange: (location: LocationInputType | null) => void
  placeholder?: string
  required?: boolean
  error?: string
}

interface PhotonFeature {
  properties: {
    name: string
    state?: string
    country: string
    osm_id: number
    osm_type: string
  }
  geometry: {
    type: 'Point'
    coordinates: [number, number] // [lng, lat] - GeoJSON format
  }
}

interface PhotonResponse {
  features: PhotonFeature[]
}

export function LocationInput({
  label,
  value,
  onChange,
  placeholder = 'Ej: Córdoba, Buenos Aires, Mendoza',
  required = false,
  error,
}: LocationInputProps) {
  const debounceTimerRef = useRef<NodeJS.Timeout | null>(null)

  const [inputValue, setInputValue] = useState('')
  const [suggestions, setSuggestions] = useState<PhotonFeature[]>([])
  const [isSearching, setIsSearching] = useState(false)
  const [showDropdown, setShowDropdown] = useState(false)

  // Sync input value with prop value
  useEffect(() => {
    if (value) {
      setInputValue(value.address || `${value.city}, ${value.province}`)
    } else {
      setInputValue('')
    }
  }, [value])

  // Cleanup debounce timer on unmount
  useEffect(() => {
    return () => {
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current)
      }
    }
  }, [])

  // Fetch suggestions from Photon API
  const fetchPhotonSuggestions = async (searchText: string) => {
    try {
      setIsSearching(true)

      const response = await fetch(
        `https://photon.komoot.io/api/?q=${encodeURIComponent(searchText)}&limit=5`
      )

      if (!response.ok) {
        throw new Error('Failed to fetch from Photon API')
      }

      const data: PhotonResponse = await response.json()

      // Filter only Argentina results
      const argentineFeatures = (data.features || []).filter(
        (feature) => feature.properties.country === 'Argentina'
      )

      setSuggestions(argentineFeatures)
      setShowDropdown(argentineFeatures.length > 0)
    } catch (err) {
      console.error('Error fetching from Photon API:', err)
      setSuggestions([])
      setShowDropdown(false)
    } finally {
      setIsSearching(false)
    }
  }

  // Handle input change with debounce
  const handleInputChange = (newValue: string) => {
    setInputValue(newValue)

    // Clear previous timer
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current)
    }

    // Clear suggestions if input is too short
    if (newValue.length < 2) {
      setSuggestions([])
      setShowDropdown(false)
      return
    }

    // Debounce: wait 300ms before making request
    debounceTimerRef.current = setTimeout(() => {
      fetchPhotonSuggestions(newValue)
    }, 300)
  }

  // Handle suggestion selection
  const handleSelectSuggestion = (feature: PhotonFeature) => {
    const city = feature.properties.name
    const province = feature.properties.state || ''
    const address = province
      ? `${city}, ${province}, Argentina`
      : `${city}, Argentina`

    // CRITICAL: Convert GeoJSON [lng, lat] to {lat, lng}
    const coordinates = {
      lat: feature.geometry.coordinates[1], // Index 1 is latitude
      lng: feature.geometry.coordinates[0], // Index 0 is longitude
    }

    const location: LocationInputType = {
      city,
      province,
      address,
      coordinates,
    }

    setInputValue(address)
    onChange(location)
    setShowDropdown(false)
    setSuggestions([])
  }

  // Clear location
  const handleClear = () => {
    setInputValue('')
    onChange(null)
    setSuggestions([])
    setShowDropdown(false)
  }

  const renderSuggestion = ({ item }: { item: PhotonFeature }) => (
    <TouchableOpacity
      style={styles.suggestionItem}
      onPress={() => handleSelectSuggestion(item)}
    >
      <View>
        <Text style={styles.suggestionCity}>{item.properties.name}</Text>
        {item.properties.state && (
          <Text style={styles.suggestionProvince}>
            {item.properties.state}, Argentina
          </Text>
        )}
      </View>
    </TouchableOpacity>
  )

  return (
    <View style={styles.container}>
      <Label>
        {label}
        {required && <Text style={styles.required}> *</Text>}
      </Label>

      <View style={styles.inputContainer}>
        <Ionicons
          name="location"
          size={20}
          color="#9ca3af"
          style={{ position: 'absolute', left: 12, zIndex: 1 }}
        />

        <TextInput
          placeholder={placeholder}
          value={inputValue}
          onChangeText={handleInputChange}
          onFocus={() => suggestions.length > 0 && setShowDropdown(true)}
          style={[styles.input, { paddingLeft: 40, paddingRight: value ? 40 : 16 }]}
        />

        {value && (
          <TouchableOpacity
            style={styles.clearButton}
            onPress={handleClear}
          >
            <Ionicons name="close-circle" size={20} color="#9ca3af" />
          </TouchableOpacity>
        )}
      </View>

      {isSearching && (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="small" color="#0ea5e9" />
          <Text style={styles.loadingText}>Buscando ubicaciones...</Text>
        </View>
      )}

      {error && <Text style={styles.errorText}>{error}</Text>}

      {/* Show selected location details */}
      {value && (
        <View style={styles.selectedLocationCard}>
          <View style={styles.selectedLocationRow}>
            <Text style={styles.selectedLocationLabel}>Ciudad:</Text>
            <Text style={styles.selectedLocationValue}>{value.city}</Text>
          </View>
          <View style={styles.selectedLocationRow}>
            <Text style={styles.selectedLocationLabel}>Provincia:</Text>
            <Text style={styles.selectedLocationValue}>{value.province}</Text>
          </View>
          {value.coordinates && (
            <Text style={styles.coordinatesText}>
              Coordenadas: {value.coordinates.lat.toFixed(4)}, {value.coordinates.lng.toFixed(4)}
            </Text>
          )}
        </View>
      )}

      {/* Suggestions Modal */}
      <Modal
        visible={showDropdown}
        transparent
        animationType="fade"
        onRequestClose={() => setShowDropdown(false)}
      >
        <TouchableOpacity
          style={styles.modalOverlay}
          activeOpacity={1}
          onPress={() => setShowDropdown(false)}
        >
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Selecciona una ubicación</Text>
              <TouchableOpacity onPress={() => setShowDropdown(false)}>
                <Ionicons name="close" size={24} color="#6b7280" />
              </TouchableOpacity>
            </View>

            <FlatList
              data={suggestions}
              keyExtractor={(item) => item.properties.osm_id.toString()}
              renderItem={renderSuggestion}
              style={styles.suggestionsList}
            />
          </View>
        </TouchableOpacity>
      </Modal>
    </View>
  )
}
