import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Search, MapPin, Calendar, Users, DollarSign, Star, Filter, X, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { searchService, getErrorMessage } from '@/services'
import type { SearchTrip, SearchQuery } from '@/types'

interface SearchFilters {
  origin_city?: string
  destination_city?: string
  date_from?: string
  min_seats?: number
  max_price?: number
  sort_by?: 'date_asc' | 'date_desc' | 'price_asc' | 'price_desc'
}

export default function SearchPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [trips, setTrips] = useState<SearchTrip[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [showFilters, setShowFilters] = useState(false)
  const [total, setTotal] = useState(0)
  const [currentPage, setCurrentPage] = useState(1)

  // Form state
  const [searchQuery, setSearchQuery] = useState<SearchFilters>({
    origin_city: searchParams.get('origin_city') || '',
    destination_city: searchParams.get('destination_city') || '',
    date_from: searchParams.get('date_from') || '',
    min_seats: searchParams.get('min_seats') ? parseInt(searchParams.get('min_seats')!) : undefined,
    max_price: searchParams.get('max_price') ? parseInt(searchParams.get('max_price')!) : undefined,
    sort_by: (searchParams.get('sort_by') as any) || 'date_asc',
  })

  useEffect(() => {
    // Auto-search on page load to show all trips by default
    handleSearch()
  }, [])

  const handleSearch = async (page = 1) => {
    try {
      setLoading(true)
      setError('')
      setCurrentPage(page)

      const query: SearchQuery = {
        ...searchQuery,
        page,
        limit: 20,
      }

      // Remove empty values
      Object.keys(query).forEach(key => {
        if (query[key as keyof SearchQuery] === '' || query[key as keyof SearchQuery] === undefined) {
          delete query[key as keyof SearchQuery]
        }
      })

      const response = await searchService.searchTrips(query)
      setTrips(response.trips)
      setTotal(response.total)

      // Update URL params
      const params = new URLSearchParams()
      if (query.origin_city) params.set('origin_city', query.origin_city)
      if (query.destination_city) params.set('destination_city', query.destination_city)
      if (query.date_from) params.set('date_from', query.date_from)
      if (query.min_seats) params.set('min_seats', query.min_seats.toString())
      if (query.max_price) params.set('max_price', query.max_price.toString())
      // Only add sort_by to URL if there are other filters active
      if (query.sort_by && (query.origin_city || query.destination_city || query.date_from || query.min_seats || query.max_price)) {
        params.set('sort_by', query.sort_by)
      }
      setSearchParams(params)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleClearFilters = () => {
    setSearchQuery({
      origin_city: '',
      destination_city: '',
      date_from: '',
      min_seats: undefined,
      max_price: undefined,
      sort_by: 'date_asc',
      page: 1,
      limit: 20,
    })
    setTrips([])
    setTotal(0)
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
                  handleSearch()
                }}
                className="space-y-6"
              >
                {/* Main Search Fields */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <Label htmlFor="origin_city">Origen</Label>
                    <div className="relative">
                      <MapPin className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                      <Input
                        id="origin_city"
                        placeholder="Ciudad de origen"
                        value={searchQuery.origin_city || ''}
                        onChange={(e) => setSearchQuery({ ...searchQuery, origin_city: e.target.value })}
                        className="pl-10"
                      />
                    </div>
                  </div>

                  <div>
                    <Label htmlFor="destination_city">Destino</Label>
                    <div className="relative">
                      <MapPin className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                      <Input
                        id="destination_city"
                        placeholder="Ciudad de destino"
                        value={searchQuery.destination_city || ''}
                        onChange={(e) => setSearchQuery({ ...searchQuery, destination_city: e.target.value })}
                        className="pl-10"
                      />
                    </div>
                  </div>

                  <div>
                    <Label htmlFor="date_from">Fecha desde</Label>
                    <div className="relative">
                      <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
                      <Input
                        id="date_from"
                        type="date"
                        value={searchQuery.date_from || ''}
                        onChange={(e) => setSearchQuery({ ...searchQuery, date_from: e.target.value })}
                        className="pl-10"
                      />
                    </div>
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
                          value={searchQuery.min_seats || ''}
                          onChange={(e) => setSearchQuery({ ...searchQuery, min_seats: e.target.value ? parseInt(e.target.value) : undefined })}
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
                          value={searchQuery.max_price || ''}
                          onChange={(e) => setSearchQuery({ ...searchQuery, max_price: e.target.value ? parseInt(e.target.value) : undefined })}
                          className="pl-10"
                        />
                      </div>
                    </div>

                    <div>
                      <Label htmlFor="sort_by">Ordenar por</Label>
                      <select
                        id="sort_by"
                        value={searchQuery.sort_by || 'date_asc'}
                        onChange={(e) => setSearchQuery({ ...searchQuery, sort_by: e.target.value as any })}
                        className="w-full h-10 px-3 rounded-md border border-input bg-background text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                      >
                        <option value="date_asc">Fecha más cercana</option>
                        <option value="date_desc">Fecha más lejana</option>
                        <option value="price_asc">Precio más bajo</option>
                        <option value="price_desc">Precio más alto</option>
                        <option value="popularity">Más popular</option>
                      </select>
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
                {searchQuery.origin_city || searchQuery.destination_city
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
                    {searchQuery.origin_city || searchQuery.destination_city
                      ? 'Intenta modificar tus criterios de búsqueda o explora otros destinos'
                      : 'No hay viajes disponibles en este momento'}
                  </p>
                  {(searchQuery.origin_city || searchQuery.destination_city) && (
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
