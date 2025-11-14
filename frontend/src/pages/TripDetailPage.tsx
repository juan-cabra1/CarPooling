import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { MapPin, Calendar, Users, DollarSign, Car, ArrowLeft, Star, Edit, Trash2, AlertCircle } from 'lucide-react'
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

export default function TripDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()
  const [trip, setTrip] = useState<Trip | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    if (id) {
      fetchTrip()
    }
  }, [id])

  const fetchTrip = async () => {
    if (!id) return

    try {
      setLoading(true)
      setError('')
      const data = await tripsService.getTripById(id)
      setTrip(data)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!trip || !id) return

    if (!confirm('¿Estás seguro de que quieres eliminar este viaje?')) {
      return
    }

    try {
      setDeleting(true)
      await tripsService.deleteTrip(id)
      navigate('/my-trips', { state: { message: 'Viaje eliminado exitosamente' } })
    } catch (err) {
      alert(getErrorMessage(err))
    } finally {
      setDeleting(false)
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('es-AR', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
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

  const isOwner = user && trip && trip.driver_id === user.id

  if (loading) {
    return (
      <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <div className="flex items-center justify-center py-16">
              <div className="text-center">
                <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4" />
                <p className="text-gray-600">Cargando detalles del viaje...</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !trip) {
    return (
      <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <Card>
              <CardContent className="pt-6">
                <div className="text-center py-8">
                  <AlertCircle className="w-16 h-16 text-destructive mx-auto mb-4" />
                  <h2 className="text-2xl font-bold text-gray-900 mb-2">Error al cargar viaje</h2>
                  <p className="text-gray-600 mb-6">{error || 'No se pudo encontrar el viaje'}</p>
                  <Link to="/search">
                    <Button>
                      <ArrowLeft className="w-4 h-4 mr-2" />
                      Volver a búsqueda
                    </Button>
                  </Link>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
      <div className="container mx-auto px-4">
        <div className="max-w-4xl mx-auto">
          {/* Back Button */}
          <div className="mb-6">
            <Button
              variant="outline"
              onClick={() => navigate(-1)}
              className="gap-2"
            >
              <ArrowLeft className="w-4 h-4" />
              Volver
            </Button>
          </div>

          {/* Main Card */}
          <Card className="overflow-hidden shadow-lg">
            <CardHeader className="bg-gradient-to-r from-primary-50 to-secondary-50 pb-6">
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <CardTitle className="text-3xl">
                      {trip.origin.city} → {trip.destination.city}
                    </CardTitle>
                    {getStatusBadge(trip.status)}
                  </div>
                  <CardDescription className="text-lg">
                    {trip.origin.province} → {trip.destination.province}
                  </CardDescription>
                </div>
                <div className="text-right">
                  <div className="text-4xl font-bold text-primary">
                    {formatPrice(trip.price_per_seat)}
                  </div>
                  <div className="text-sm text-muted-foreground">por asiento</div>
                </div>
              </div>

              {/* Route Details */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6 pt-6 border-t border-gray-200">
                <div className="space-y-2">
                  <div className="flex items-center gap-2 text-primary font-semibold">
                    <MapPin className="w-5 h-5" />
                    <span>Origen</span>
                  </div>
                  <div className="ml-7">
                    <p className="font-medium">{trip.origin.city}, {trip.origin.province}</p>
                    <p className="text-sm text-gray-600">{trip.origin.address}</p>
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="flex items-center gap-2 text-secondary font-semibold">
                    <MapPin className="w-5 h-5" />
                    <span>Destino</span>
                  </div>
                  <div className="ml-7">
                    <p className="font-medium">{trip.destination.city}, {trip.destination.province}</p>
                    <p className="text-sm text-gray-600">{trip.destination.address}</p>
                  </div>
                </div>
              </div>
            </CardHeader>

            <CardContent className="pt-6 space-y-6">
              {/* Date & Time */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                    <Calendar className="w-5 h-5 text-primary" />
                    <span>Salida</span>
                  </div>
                  <p className="ml-7 text-gray-900">{formatDate(trip.departure_datetime)}</p>
                </div>

                <div>
                  <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                    <Calendar className="w-5 h-5 text-secondary" />
                    <span>Llegada estimada</span>
                  </div>
                  <p className="ml-7 text-gray-900">{formatDate(trip.estimated_arrival_datetime)}</p>
                </div>
              </div>

              {/* Trip Info */}
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 pt-6 border-t">
                <div>
                  <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                    <Car className="w-5 h-5 text-muted-foreground" />
                    <span>Vehículo</span>
                  </div>
                  <p className="ml-7 text-gray-900">
                    {trip.car.brand} {trip.car.model}
                  </p>
                  <p className="ml-7 text-sm text-gray-600">
                    {trip.car.year} - {trip.car.color}
                  </p>
                  <p className="ml-7 text-sm text-gray-600">
                    Patente: {trip.car.plate}
                  </p>
                </div>

                <div>
                  <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                    <Users className="w-5 h-5 text-muted-foreground" />
                    <span>Asientos</span>
                  </div>
                  <p className="ml-7 text-gray-900">
                    {trip.available_seats} de {trip.total_seats} disponibles
                  </p>
                  <p className="ml-7 text-sm text-gray-600">
                    {trip.reserved_seats} reservados
                  </p>
                </div>

                <div>
                  <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                    <DollarSign className="w-5 h-5 text-muted-foreground" />
                    <span>Precio total</span>
                  </div>
                  <p className="ml-7 text-gray-900">
                    {formatPrice(trip.price_per_seat * trip.total_seats)}
                  </p>
                  <p className="ml-7 text-sm text-gray-600">
                    ({formatPrice(trip.price_per_seat)} × {trip.total_seats})
                  </p>
                </div>
              </div>

              {/* Description */}
              {trip.description && (
                <div className="pt-6 border-t">
                  <h3 className="font-semibold text-gray-900 mb-3 text-lg">Descripción</h3>
                  <p className="text-gray-700 whitespace-pre-wrap">{trip.description}</p>
                </div>
              )}

              {/* Preferences */}
              <div className="pt-6 border-t">
                <h3 className="font-semibold text-gray-900 mb-3 text-lg">Preferencias del viaje</h3>
                <div className="flex flex-wrap gap-3">
                  <Badge variant={trip.preferences.pets_allowed ? 'default' : 'outline'} className="text-sm">
                    🐕 Mascotas {trip.preferences.pets_allowed ? 'permitidas' : 'no permitidas'}
                  </Badge>
                  <Badge variant={trip.preferences.smoking_allowed ? 'default' : 'outline'} className="text-sm">
                    🚬 Fumar {trip.preferences.smoking_allowed ? 'permitido' : 'no permitido'}
                  </Badge>
                  <Badge variant={trip.preferences.music_allowed ? 'default' : 'outline'} className="text-sm">
                    🎵 Música {trip.preferences.music_allowed ? 'permitida' : 'no permitida'}
                  </Badge>
                </div>
              </div>
            </CardContent>

            <CardFooter className="bg-gray-50 border-t flex flex-col sm:flex-row gap-3">
              {isOwner ? (
                <>
                  {trip.status === 'published' && trip.reserved_seats === 0 && (
                    <>
                      <Link to={`/trips/${trip.id}/edit`} className="flex-1">
                        <Button variant="outline" className="w-full">
                          <Edit className="w-4 h-4 mr-2" />
                          Editar Viaje
                        </Button>
                      </Link>

                      <Button
                        variant="destructive"
                        onClick={handleDelete}
                        disabled={deleting}
                        className="flex-1"
                      >
                        {deleting ? (
                          <>
                            <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                            Eliminando...
                          </>
                        ) : (
                          <>
                            <Trash2 className="w-4 h-4 mr-2" />
                            Eliminar Viaje
                          </>
                        )}
                      </Button>
                    </>
                  )}
                  {trip.reserved_seats > 0 && (
                    <div className="w-full flex items-center justify-center px-4 py-3 bg-yellow-50 border border-yellow-200 rounded-md">
                      <AlertCircle className="w-5 h-5 text-yellow-600 mr-2" />
                      <span className="text-sm text-yellow-700 font-medium">
                        Este viaje tiene reservas activas y no puede ser editado o eliminado
                      </span>
                    </div>
                  )}
                </>
              ) : (
                <Button className="w-full" size="lg">
                  <Users className="w-5 h-5 mr-2" />
                  Reservar Asientos
                </Button>
              )}
            </CardFooter>
          </Card>
        </div>
      </div>
    </div>
  )
}
