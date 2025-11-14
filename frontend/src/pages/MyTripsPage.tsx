import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Calendar, MapPin, Users, DollarSign, Edit, Trash2, Car, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { tripsService, getErrorMessage } from '@/services'
import { useAuth } from '@/context/AuthContext'
import type { Trip } from '@/types'

export default function MyTripsPage() {
  const { user } = useAuth()
  const [trips, setTrips] = useState<Trip[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [deletingId, setDeletingId] = useState<string | null>(null)

  useEffect(() => {
    fetchMyTrips()
  }, [])

  const fetchMyTrips = async () => {
    try {
      setLoading(true)
      setError('')

      if (!user?.id) {
        setError('No se pudo obtener la información del usuario')
        return
      }

      const response = await tripsService.getMyTrips(user.id)
      setTrips(response.trips)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteTrip = async (tripId: string) => {
    if (!confirm('¿Estás seguro de que quieres eliminar este viaje?')) {
      return
    }

    try {
      setDeletingId(tripId)
      await tripsService.deleteTrip(tripId)
      // Refrescar la lista
      await fetchMyTrips()
    } catch (err) {
      alert(getErrorMessage(err))
    } finally {
      setDeletingId(null)
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('es-AR', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
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
    const statusConfig: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }> = {
      draft: { label: 'Borrador', variant: 'outline' },
      published: { label: 'Publicado', variant: 'default' },
      full: { label: 'Completo', variant: 'secondary' },
      in_progress: { label: 'En Progreso', variant: 'secondary' },
      completed: { label: 'Completado', variant: 'secondary' },
      cancelled: { label: 'Cancelado', variant: 'destructive' },
    }

    const config = statusConfig[status] || { label: status, variant: 'outline' as const }
    return <Badge variant={config.variant}>{config.label}</Badge>
  }

  if (loading) {
    return (
      <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
        <div className="container mx-auto px-4">
          <div className="max-w-6xl mx-auto">
            <div className="flex items-center justify-center py-16">
              <div className="text-center">
                <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4" />
                <p className="text-gray-600">Cargando tus viajes...</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
      <div className="container mx-auto px-4">
        <div className="max-w-6xl mx-auto">
          {/* Header */}
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-8 gap-4">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">Mis Viajes</h1>
              <p className="text-gray-600">
                Gestiona y organiza tus viajes publicados
              </p>
            </div>
            <Link to="/create-trip">
              <Button size="lg" className="w-full sm:w-auto">
                <Plus className="w-5 h-5 mr-2" />
                Publicar Viaje
              </Button>
            </Link>
          </div>

          {/* Error Message */}
          {error && (
            <div className="mb-6 p-4 rounded-lg bg-destructive/10 border border-destructive/20 flex items-start gap-3">
              <AlertCircle className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
              <div>
                <p className="font-semibold text-destructive">Error al cargar viajes</p>
                <p className="text-sm text-destructive/80">{error}</p>
              </div>
            </div>
          )}

          {/* Trips List */}
          {trips.length === 0 ? (
            <Card className="text-center py-16">
              <CardContent>
                <div className="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Car className="w-10 h-10 text-gray-400" />
                </div>
                <CardTitle className="text-2xl mb-2">No tienes viajes publicados</CardTitle>
                <CardDescription className="text-lg mb-6">
                  Publica tu primer viaje y empieza a compartir tus trayectos
                </CardDescription>
                <Link to="/create-trip">
                  <Button size="lg">
                    <Plus className="w-5 h-5 mr-2" />
                    Publicar Mi Primer Viaje
                  </Button>
                </Link>
              </CardContent>
            </Card>
          ) : (
            <div className="grid grid-cols-1 gap-6">
              {trips.map((trip) => (
                <Card key={trip.id} className="overflow-hidden hover:shadow-lg transition-shadow">
                  <CardHeader className="bg-gradient-to-r from-primary-50 to-secondary-50 pb-4">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-1">
                          <CardTitle className="text-xl">
                            {trip.origin.city} → {trip.destination.city}
                          </CardTitle>
                          {getStatusBadge(trip.status)}
                        </div>
                        <CardDescription className="text-base">
                          {trip.origin.province} → {trip.destination.province}
                        </CardDescription>
                      </div>
                      <div className="text-right">
                        <div className="text-2xl font-bold text-primary">
                          {formatPrice(trip.price_per_seat)}
                        </div>
                        <div className="text-xs text-muted-foreground">por asiento</div>
                      </div>
                    </div>
                  </CardHeader>

                  <CardContent className="pt-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                      {/* Departure Info */}
                      <div className="space-y-3">
                        <div className="flex items-start gap-2 text-sm">
                          <Calendar className="w-4 h-4 text-muted-foreground mt-0.5 flex-shrink-0" />
                          <div>
                            <div className="font-medium text-gray-900">Salida</div>
                            <div className="text-gray-700">{formatDate(trip.departure_datetime)}</div>
                          </div>
                        </div>

                        <div className="flex items-start gap-2 text-sm">
                          <MapPin className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
                          <div>
                            <div className="font-medium text-gray-900">Origen</div>
                            <div className="text-gray-700">{trip.origin.address}</div>
                          </div>
                        </div>
                      </div>

                      {/* Arrival Info */}
                      <div className="space-y-3">
                        <div className="flex items-start gap-2 text-sm">
                          <Calendar className="w-4 h-4 text-muted-foreground mt-0.5 flex-shrink-0" />
                          <div>
                            <div className="font-medium text-gray-900">Llegada estimada</div>
                            <div className="text-gray-700">
                              {formatDate(trip.estimated_arrival_datetime)}
                            </div>
                          </div>
                        </div>

                        <div className="flex items-start gap-2 text-sm">
                          <MapPin className="w-4 h-4 text-secondary mt-0.5 flex-shrink-0" />
                          <div>
                            <div className="font-medium text-gray-900">Destino</div>
                            <div className="text-gray-700">{trip.destination.address}</div>
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Trip Details */}
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-4 pt-4 border-t">
                      <div className="flex items-center gap-2 text-sm">
                        <Car className="w-4 h-4 text-muted-foreground" />
                        <div>
                          <div className="font-medium text-gray-900">Vehículo</div>
                          <div className="text-gray-700">
                            {trip.car.brand} {trip.car.model}
                          </div>
                        </div>
                      </div>

                      <div className="flex items-center gap-2 text-sm">
                        <Users className="w-4 h-4 text-muted-foreground" />
                        <div>
                          <div className="font-medium text-gray-900">Asientos</div>
                          <div className="text-gray-700">
                            {trip.available_seats}/{trip.total_seats} disponibles
                          </div>
                        </div>
                      </div>

                      <div className="flex items-center gap-2 text-sm col-span-2 md:col-span-1">
                        <DollarSign className="w-4 h-4 text-muted-foreground" />
                        <div>
                          <div className="font-medium text-gray-900">Reservados</div>
                          <div className="text-gray-700">{trip.reserved_seats} asientos</div>
                        </div>
                      </div>
                    </div>

                    {/* Description */}
                    {trip.description && (
                      <div className="mt-4 pt-4 border-t">
                        <p className="text-sm text-gray-700">{trip.description}</p>
                      </div>
                    )}

                    {/* Preferences */}
                    {(trip.preferences.pets_allowed || trip.preferences.smoking_allowed || trip.preferences.music_allowed) && (
                      <div className="flex flex-wrap gap-2 mt-4">
                        {trip.preferences.pets_allowed && (
                          <Badge variant="secondary" className="text-xs">
                            🐕 Mascotas
                          </Badge>
                        )}
                        {trip.preferences.smoking_allowed && (
                          <Badge variant="secondary" className="text-xs">
                            🚬 Fumar
                          </Badge>
                        )}
                        {trip.preferences.music_allowed && (
                          <Badge variant="secondary" className="text-xs">
                            🎵 Música
                          </Badge>
                        )}
                      </div>
                    )}
                  </CardContent>

                  <CardFooter className="bg-gray-50 border-t">
                    <div className="w-full flex flex-col sm:flex-row gap-3">
                      <Link to={`/trips/${trip.id}`} className="flex-1">
                        <Button variant="outline" className="w-full">
                          Ver Detalles
                        </Button>
                      </Link>

                      {trip.status === 'published' && trip.reserved_seats === 0 && (
                        <>
                          <Link to={`/trips/${trip.id}/edit`} className="flex-1">
                            <Button variant="outline" className="w-full">
                              <Edit className="w-4 h-4 mr-2" />
                              Editar
                            </Button>
                          </Link>

                          <Button
                            variant="destructive"
                            onClick={() => handleDeleteTrip(trip.id)}
                            disabled={deletingId === trip.id}
                            className="flex-1"
                          >
                            {deletingId === trip.id ? (
                              <>
                                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                                Eliminando...
                              </>
                            ) : (
                              <>
                                <Trash2 className="w-4 h-4 mr-2" />
                                Eliminar
                              </>
                            )}
                          </Button>
                        </>
                      )}

                      {trip.reserved_seats > 0 && (
                        <div className="flex-1 flex items-center justify-center px-4 py-2 bg-yellow-50 border border-yellow-200 rounded-md">
                          <AlertCircle className="w-4 h-4 text-yellow-600 mr-2" />
                          <span className="text-sm text-yellow-700 font-medium">
                            No se puede editar o eliminar (tiene reservas)
                          </span>
                        </div>
                      )}
                    </div>
                  </CardFooter>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
