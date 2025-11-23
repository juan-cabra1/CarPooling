/**
 * LocationInput Component
 * Photon (Komoot) Autocomplete integration for structured location input
 *
 * Features:
 * - Autocompletes city/address using Photon API (OpenStreetMap data)
 * - Extracts city, province, and coordinates automatically
 * - Validates required fields (city and province)
 * - Restricts to Argentina for relevant results
 * - No API key required (free and open-source)
 *
 * User workflow:
 * 1. User types "Cordoba" in input
 * 2. Photon API shows suggestions after 300ms debounce
 * 3. User selects "Córdoba, Córdoba, Argentina"
 * 4. Component extracts:
 *    - city: "Córdoba"
 *    - province: "Córdoba"
 *    - address: "Córdoba, Córdoba, Argentina"
 *    - coordinates: {lat: -31.4166867, lng: -64.1834193}
 */

import { useEffect, useRef, useState } from 'react'
import { MapPin, X } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import type { LocationInput as LocationInputType } from '@/types/search'

interface LocationInputProps {
  label: string
  labelClassName?: string
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
  labelClassName,
  value,
  onChange,
  placeholder = 'Ej: Córdoba, Buenos Aires, Mendoza',
  required = false,
  error,
}: LocationInputProps) {
  const inputRef = useRef<HTMLInputElement>(null)
  const dropdownRef = useRef<HTMLDivElement>(null)
  const debounceTimerRef = useRef<number | null>(null)

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

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(event.target as Node)
      ) {
        setShowDropdown(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
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
    if (inputRef.current) {
      inputRef.current.value = ''
      inputRef.current.focus()
    }
  }

  return (
    <div className="space-y-2">
      <Label htmlFor={`location-${label}`}
        className={labelClassName}
      >
        {label}
        {required && <span className="text-red-500 ml-1">*</span>}
      </Label>

      <div className="relative">
        <MapPin className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />

        <Input
          ref={inputRef}
          id={`location-${label}`}
          type="text"
          placeholder={placeholder}
          value={inputValue}
          onChange={(e) => handleInputChange(e.target.value)}
          onFocus={() => suggestions.length > 0 && setShowDropdown(true)}
          className={`pl-10 pr-10 ${error ? 'border-red-500' : ''}`}
        />

        {value && (
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={handleClear}
            className="absolute right-1 top-1/2 transform -translate-y-1/2 h-7 w-7 p-0"
          >
            <X className="h-4 w-4" />
          </Button>
        )}

        {/* Dropdown with suggestions */}
        {showDropdown && suggestions.length > 0 && (
          <div
            ref={dropdownRef}
            className="absolute top-full left-0 right-0 bg-white border border-gray-300 rounded-md shadow-lg z-50 mt-1 max-h-60 overflow-y-auto"
          >
            <ul className="py-1">
              {suggestions.map((feature) => (
                <li
                  key={feature.properties.osm_id}
                  onClick={() => handleSelectSuggestion(feature)}
                  className="px-4 py-2 hover:bg-gray-100 cursor-pointer border-b last:border-b-0 transition-colors"
                >
                  <div className="font-medium text-gray-900">
                    {feature.properties.name}
                  </div>
                  {feature.properties.state && (
                    <div className="text-xs text-gray-500">
                      {feature.properties.state}, Argentina
                    </div>
                  )}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>

      {isSearching && (
        <p className="text-xs text-gray-500">Buscando ubicaciones...</p>
      )}

      {error && <p className="text-sm text-red-500">{error}</p>}

      {/* Show selected location details */}
      {value && (
        <div className="text-sm text-gray-600 bg-gray-50 p-2 rounded">
          <div>
            <span className="font-medium">Ciudad:</span> {value.city}
          </div>
          <div>
            <span className="font-medium">Provincia:</span> {value.province}
          </div>
          {value.coordinates && (
            <div className="text-xs text-gray-500 mt-1">
              Coordenadas: {value.coordinates.lat.toFixed(4)},{' '}
              {value.coordinates.lng.toFixed(4)}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
