import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Search, Calendar, Users, DollarSign, Star, Filter, X, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { LocationInput, RadiusSlider } from '@/components/search'
import { searchService, getErrorMessage } from '@/services'
import type { SearchTrip, SearchQuery, LocationInput as LocationInputType, SearchSortBy, SearchSortOrder } from '@/types'

export default function SearchPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [trips, setTrips] = useState<SearchTrip[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [validationError, setValidationError] = useState('')
  const [showFilters, setShowFilters] = useState(false)
  const [total, setTotal] = useState(0)
  const [currentPage, setCurrentPage] = useState(1)

  // Form state - individual states for better control
  const [origin, setOrigin] = useState<LocationInputType | null>(null)
  const [destination, setDestination] = useState<LocationInputType | null>(null)
  const [useOriginRadius, setUseOriginRadius] = useState(false)
  const [useDestinationRadius, setUseDestinationRadius] = useState(false)
  const [originRadius, setOriginRadius] = useState(10)
  const [destinationRadius, setDestinationRadius] = useState(10)
  const [departureDate, setDepartureDate] = useState('')
  const [minSeats, setMinSeats] = useState<number | undefined>(undefined)
  const [maxPrice, setMaxPrice] = useState<number | undefined>(undefined)
  const [sortOption, setSortOption] = useState('departure_time|asc')
  const [petsAllowed, setPetsAllowed] = useState<boolean | undefined>(undefined)
  const [smokingAllowed, setSmokingAllowed] = useState<boolean | undefined>(undefined)
  const [musicAllowed, setMusicAllowed] = useState<boolean | undefined>(undefined)

  // Initialize from URL params on mount
  useEffect(() => {
    // Parse location from URL params
    const originCity = searchParams.get('origin_city')
    const originProvince = searchParams.get('origin_province')
    const originAddress = searchParams.get('origin_address')
    const originLat = searchParams.get('origin_lat')
    const originLng = searchParams.get('origin_lng')

    if (originCity && originProvince) {
      setOrigin({
        city: originCity,
        province: originProvince,
        address: originAddress || `${originCity}, ${originProvince}`,
        coordinates: originLat && originLng ? {
          lat: parseFloat(originLat),
          lng: parseFloat(originLng)
        } : undefined
      })
    }

    const destCity = searchParams.get('destination_city')
    const destProvince = searchParams.get('destination_province')
    const destAddress = searchParams.get('destination_address')
    const destLat = searchParams.get('destination_lat')
    const destLng = searchParams.get('destination_lng')

    if (destCity && destProvince) {
      setDestination({
        city: destCity,
        province: destProvince,
        address: destAddress || `${destCity}, ${destProvince}`,
        coordinates: destLat && destLng ? {
          lat: parseFloat(destLat),
          lng: parseFloat(destLng)
        } : undefined
      })
    }

    // Parse other filters
    const date = searchParams.get('departure_date')
    if (date) setDepartureDate(date)

    const seats = searchParams.get('min_seats')
    if (seats) setMinSeats(parseInt(seats))

    const price = searchParams.get('max_price')
    if (price) setMaxPrice(parseInt(price))

    const originRad = searchParams.get('origin_radius')
    if (originRad) {
      setOriginRadius(parseInt(originRad))
      setUseOriginRadius(true)
    }

    const destRad = searchParams.get('destination_radius')
    if (destRad) {
      setDestinationRadius(parseInt(destRad))
      setUseDestinationRadius(true)
    }

    const sortByParam = searchParams.get('sort_by') as SearchSortBy
    const sortOrderParam = searchParams.get('sort_order') as SearchSortOrder
    if (sortByParam && sortOrderParam) {
      setSortOption(`${sortByParam}|${sortOrderParam}`)
    }

    // Only auto-search if there are URL params (bookmarked search)
    if (searchParams.toString()) {
      handleSearch()
    }
  }, [])

  const validateSearch = (): boolean => {
    // Clear any previous validation errors
    setValidationError('')
    return true
  }

  const handleSearch = async (page = 1) => {
    try {
      setLoading(true)
      setError('')
      setValidationError('')
      setCurrentPage(page)

      // Build search query
      const query: SearchQuery = {
        page,
        limit: 20,
      }

      // Add location filters
      if (origin) {
        query.origin = origin
        // Only add radius if user enabled it AND coordinates are available
        if (origin.coordinates && useOriginRadius) {
          query.origin_radius = originRadius
        }
      }

      if (destination) {
        query.destination = destination
        // Only add radius if user enabled it AND coordinates are available
        if (destination.coordinates && useDestinationRadius) {
          query.destination_radius = destinationRadius
        }
      }

      // Add other filters
      if (departureDate) query.departure_date = departureDate
      if (minSeats) query.min_seats = minSeats
      if (maxPrice) query.max_price = maxPrice
      if (petsAllowed !== undefined) query.pets_allowed = petsAllowed
      if (smokingAllowed !== undefined) query.smoking_allowed = smokingAllowed
      if (musicAllowed !== undefined) query.music_allowed = musicAllowed

      // Parse and add sorting from combined option
      const [sortBy, sortOrder] = sortOption.split('|') as [SearchSortBy, SearchSortOrder]
      query.sort_by = sortBy
      query.sort_order = sortOrder

      const response = await searchService.searchTrips(query)
      setTrips(response.trips)
      setTotal(response.total)

      // Update URL params for bookmarking
      const params = new URLSearchParams()

      // Origin location
      if (origin) {
        params.set('origin_city', origin.city)
        params.set('origin_province', origin.province)
        if (origin.address) params.set('origin_address', origin.address)
        if (origin.coordinates) {
          params.set('origin_lat', origin.coordinates.lat.toString())
          params.set('origin_lng', origin.coordinates.lng.toString())
          // Only add radius to URL if user enabled it
          if (useOriginRadius) {
            params.set('origin_radius', originRadius.toString())
          }
        }
      }

      // Destination location
      if (destination) {
        params.set('destination_city', destination.city)
        params.set('destination_province', destination.province)
        if (destination.address) params.set('destination_address', destination.address)
        if (destination.coordinates) {
          params.set('destination_lat', destination.coordinates.lat.toString())
          params.set('destination_lng', destination.coordinates.lng.toString())
          // Only add radius to URL if user enabled it
          if (useDestinationRadius) {
            params.set('destination_radius', destinationRadius.toString())
          }
        }
      }

      // Other filters
      if (departureDate) params.set('departure_date', departureDate)
      if (minSeats) params.set('min_seats', minSeats.toString())
      if (maxPrice) params.set('max_price', maxPrice.toString())

      // Sorting
      params.set('sort_by', sortBy)
      params.set('sort_order', sortOrder)

      setSearchParams(params)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleClearFilters = () => {
    setOrigin(null)
    setDestination(null)
    setUseOriginRadius(false)
    setUseDestinationRadius(false)
    setOriginRadius(10)
    setDestinationRadius(10)
    setDepartureDate('')
    setMinSeats(undefined)
    setMaxPrice(undefined)
    setSortOption('departure_time|asc')
    setPetsAllowed(undefined)
    setSmokingAllowed(undefined)
    setMusicAllowed(undefined)
    setTrips([])
    setTotal(0)
    setValidationError('')
    setSearchParams({})
  }

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
    <div className="min-h-screen bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
      <div className="container mx-auto px-4">
        <div className="max-w-7xl mx-auto">
          {/* Header */}
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-2">Buscar Viajes</h1>
            <p className="text-gray-600">Encuentra el viaje perfecto para tu destino</p>
          </div>

          {/* Search Form */}
          <Card className="mb-8">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Search className="w-5 h-5" />
                Buscar
              </CardTitle>
            </CardHeader>
            <CardContent>
              <form
                onSubmit={(e) => {
                  e.preventDefault()
                  if (validateSearch()) {
                    handleSearch()
                  }
                }}
                className="space-y-6"
              >
                {/* Validation Error */}
                {validationError && (
                  <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
                    {validationError}
                  </div>
                )}

                {/* Main Search Fields */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* Origin Location */}
                  <div className="space-y-4">
                    <LocationInput
                      label="Origen"
                      value={origin}
                      onChange={setOrigin}
                      placeholder="Ej: Córdoba, Buenos Aires, Mendoza"
                    />

                    {/* Origin Radius Toggle - only shown if coordinates available */}
                    {origin?.coordinates && (
                      <div className="space-y-3">
                        <div className="flex items-center space-x-2">
                          <Checkbox
                            id="use-origin-radius"
                            checked={useOriginRadius}
                            onCheckedChange={(checked) => setUseOriginRadius(!!checked)}
                          />
                          <label
                            htmlFor="use-origin-radius"
                            className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                          >
                            Buscar por radio desde este punto
                          </label>
                        </div>

                        {useOriginRadius && (
                          <RadiusSlider
                            value={originRadius}
                            onChange={setOriginRadius}
                            label="Radio de búsqueda"
                          />
                        )}
                      </div>
                    )}
                  </div>

                  {/* Destination Location */}
                  <div className="space-y-4">
                    <LocationInput
                      label="Destino"
                      value={destination}
                      onChange={setDestination}
                      placeholder="Ej: Medellín, Antioquia"
                    />

                    {/* Destination Radius - optional with checkbox */}
                    {destination?.coordinates && (
                      <div className="space-y-3">
                        <div className="flex items-center space-x-2">
                          <Checkbox
                            id="use-destination-radius"
                            checked={useDestinationRadius}
                            onCheckedChange={(checked) =>
                              setUseDestinationRadius(!!checked)
                            }
                          />
                          <label
                            htmlFor="use-destination-radius"
                            className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                          >
                            Buscar por radio hasta este punto
                          </label>
                        </div>

                        {useDestinationRadius && (
                          <RadiusSlider
                            value={destinationRadius}
                            onChange={setDestinationRadius}
                            label="Radio de búsqueda hasta destino"
                          />
                        )}
                      </div>
                    )}
                  </div>
                </div>

                {/* Departure Date */}
                <div>
                  <Label htmlFor="departure_date">Fecha de salida</Label>
                  <div className="relative">
                    <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <Input
                      id="departure_date"
                      type="date"
                      value={departureDate}
                      onChange={(e) => setDepartureDate(e.target.value)}
                      className="pl-10"
                    />
                  </div>
                </div>

                {/* Advanced Filters Toggle */}
                <div className="flex items-center justify-between pt-4 border-t">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setShowFilters(!showFilters)}
                  >
                    <Filter className="w-4 h-4 mr-2" />
                    {showFilters ? 'Ocultar filtros' : 'Más filtros'}
                  </Button>

                  <div className="flex gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleClearFilters}
                    >
                      <X className="w-4 h-4 mr-2" />
                      Limpiar
                    </Button>
                    <Button type="submit" disabled={loading}>
                      {loading ? (
                        <>
                          <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                          Buscando...
                        </>
                      ) : (
                        <>
                          <Search className="w-4 h-4 mr-2" />
                          Buscar
                        </>
                      )}
                    </Button>
                  </div>
                </div>

                {/* Advanced Filters */}
                {showFilters && (
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-4 border-t">
                    <div>
                      <Label htmlFor="min_seats">Asientos mínimos</Label>
                      <div className="relative">
                        <Users className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <Input
                          id="min_seats"
                          type="number"
                          min="1"
                          placeholder="Ej: 2"
                          value={minSeats || ''}
                          onChange={(e) => setMinSeats(e.target.value ? parseInt(e.target.value) : undefined)}
                          className="pl-10"
                        />
                      </div>
                    </div>

                    <div>
                      <Label htmlFor="max_price">Precio máximo</Label>
                      <div className="relative">
                        <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <Input
                          id="max_price"
                          type="number"
                          min="0"
                          placeholder="Ej: 5000"
                          value={maxPrice || ''}
                          onChange={(e) => setMaxPrice(e.target.value ? parseInt(e.target.value) : undefined)}
                          className="pl-10"
                        />
                      </div>
                    </div>

                    <div>
                      <Label htmlFor="sort_option">Ordenar por</Label>
                      <select
                        id="sort_option"
                        value={sortOption}
                        onChange={(e) => setSortOption(e.target.value)}
                        className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                      >
                        <option value="departure_time|asc">Salida más temprana</option>
                        <option value="departure_time|desc">Salida más tardía</option>
                        <option value="price|asc">Precio más bajo</option>
                        <option value="price|desc">Precio más alto</option>
                        <option value="rating|desc">Mejor calificación</option>
                        <option value="popularity|desc">Más popular</option>
                      </select>
                    </div>

                    {/* Preference Filters */}
                    <div className="md:col-span-3 pt-2 border-t">
                      <Label className="block mb-3">Preferencias del conductor</Label>
                      <div className="flex flex-wrap gap-6">
                        <div className="flex items-center gap-2">
                          <Checkbox
                            id="pets_allowed"
                            checked={petsAllowed === true}
                            onCheckedChange={(checked) => setPetsAllowed(checked ? true : undefined)}
                          />
                          <label htmlFor="pets_allowed" className="text-sm font-medium cursor-pointer">
                            🐕 Permite mascotas
                          </label>
                        </div>
                        <div className="flex items-center gap-2">
                          <Checkbox
                            id="smoking_allowed"
                            checked={smokingAllowed === true}
                            onCheckedChange={(checked) => setSmokingAllowed(checked ? true : undefined)}
                          />
                          <label htmlFor="smoking_allowed" className="text-sm font-medium cursor-pointer">
                            🚬 Permite fumar
                          </label>
                        </div>
                        <div className="flex items-center gap-2">
                          <Checkbox
                            id="music_allowed"
                            checked={musicAllowed === true}
                            onCheckedChange={(checked) => setMusicAllowed(checked ? true : undefined)}
                          />
                          <label htmlFor="music_allowed" className="text-sm font-medium cursor-pointer">
                            🎵 Permite música
                          </label>
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </form>
            </CardContent>
          </Card>

          {/* Error Message */}
          {error && (
            <Card className="mb-6 border-red-200 bg-red-50">
              <CardContent className="pt-6">
                <p className="text-red-700">{error}</p>
              </CardContent>
            </Card>
          )}

          {/* Results */}
          {trips.length > 0 && (
            <div className="mb-6">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">
                {origin || destination
                  ? `${total} ${total === 1 ? 'viaje encontrado' : 'viajes encontrados'}`
                  : `Todos los viajes disponibles (${total})`}
              </h2>

              <div className="grid gap-6">
                {trips.map((trip) => (
                  <Card key={trip.id} className="overflow-hidden hover:shadow-lg transition-shadow">
                    <CardHeader className="bg-gradient-to-r from-primary-50 to-secondary-50">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <CardTitle className="text-2xl mb-2 flex items-center gap-3 text-gray-900">
                            <span>{trip.origin.city}</span>
                            <ArrowRight className="w-5 h-5 text-gray-400" />
                            <span>{trip.destination.city}</span>
                          </CardTitle>
                          <CardDescription className="text-base text-gray-700">
                            {trip.origin.province} → {trip.destination.province}
                          </CardDescription>
                        </div>
                        <div className="text-right">
                          <div className="text-3xl font-bold text-primary">
                            {formatPrice(trip.price_per_seat)}
                          </div>
                          <div className="text-sm text-muted-foreground">por asiento</div>
                        </div>
                      </div>
                    </CardHeader>

                    <CardContent className="pt-6">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        {/* Trip Info */}
                        <div className="space-y-3">
                          <div className="flex items-center gap-2 text-gray-700">
                            <Calendar className="w-5 h-5 text-primary" />
                            <div>
                              <div className="font-semibold">Salida</div>
                              <div className="text-sm">{formatDate(trip.departure_datetime)}</div>
                            </div>
                          </div>

                          <div className="flex items-center gap-2 text-gray-700">
                            <Users className="w-5 h-5 text-muted-foreground" />
                            <div>
                              <div className="font-semibold">Asientos disponibles</div>
                              <div className="text-sm">
                                {trip.available_seats} de {trip.total_seats}
                              </div>
                            </div>
                          </div>
                        </div>

                        {/* Driver Info */}
                        <div className="space-y-3">
                          <div className="flex items-center gap-2 text-gray-700">
                            <div className="w-10 h-10 rounded-full bg-primary-100 flex items-center justify-center text-primary font-bold">
                              {trip.driver.name.charAt(0).toUpperCase()}
                            </div>
                            <div>
                              <div className="font-semibold">{trip.driver.name}</div>
                              <div className="text-sm flex items-center gap-1">
                                <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                                <span>{trip.driver.rating.toFixed(1)}</span>
                                <span className="text-gray-500">
                                  ({trip.driver.total_trips} viajes)
                                </span>
                              </div>
                            </div>
                          </div>

                          <div className="flex flex-wrap gap-2">
                            {trip.preferences.pets_allowed && (
                              <Badge variant="outline">🐕 Mascotas</Badge>
                            )}
                            {trip.preferences.music_allowed && (
                              <Badge variant="outline">🎵 Música</Badge>
                            )}
                            {!trip.preferences.smoking_allowed && (
                              <Badge variant="outline">🚭 No fumar</Badge>
                            )}
                          </div>
                        </div>
                      </div>

                      {trip.description && (
                        <div className="mt-4 pt-4 border-t">
                          <p className="text-sm text-gray-600 line-clamp-2">{trip.description}</p>
                        </div>
                      )}
                    </CardContent>

                    <CardFooter className="bg-gray-50 border-t p-0">
                      <Link to={`/trips/${trip.trip_id}`} className="w-full">
                        <Button className="w-full rounded-t-none" size="lg">
                          Ver Detalles y Reservar
                          <ArrowRight className="w-4 h-4 ml-2" />
                        </Button>
                      </Link>
                    </CardFooter>
                  </Card>
                ))}
              </div>

              {/* Pagination */}
              {total > 20 && (
                <div className="flex justify-center gap-2 mt-8">
                  <Button
                    variant="outline"
                    onClick={() => handleSearch(currentPage - 1)}
                    disabled={currentPage === 1 || loading}
                  >
                    Anterior
                  </Button>
                  <span className="flex items-center px-4 text-gray-700">
                    Página {currentPage} de {Math.ceil(total / 20)}
                  </span>
                  <Button
                    variant="outline"
                    onClick={() => handleSearch(currentPage + 1)}
                    disabled={currentPage >= Math.ceil(total / 20) || loading}
                  >
                    Siguiente
                  </Button>
                </div>
              )}
            </div>
          )}

          {/* Empty State - No results found */}
          {!loading && trips.length === 0 && !error && (
            <Card>
              <CardContent className="pt-6">
                <div className="text-center py-12">
                  <Search className="w-16 h-16 text-gray-400 mx-auto mb-4" />
                  <h2 className="text-2xl font-bold text-gray-900 mb-2">No se encontraron viajes</h2>
                  <p className="text-gray-600 mb-6">
                    {origin || destination
                      ? 'Intenta modificar tus criterios de búsqueda o explora otros destinos'
                      : 'No hay viajes disponibles en este momento'}
                  </p>
                  {(origin || destination) && (
                    <Button onClick={handleClearFilters} variant="outline">
                      Limpiar búsqueda
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}
