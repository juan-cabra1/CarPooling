import { useState, useEffect } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { MapPin, Calendar, Users, DollarSign, X, AlertCircle, CheckCircle, Clock } from 'lucide-react'
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
import { bookingsService, getErrorMessage } from '@/services'
import { useAuth } from '@/context/AuthContext'
import type { Booking } from '@/types'

export default function MyBookingsPage() {
  const { } = useAuth()
  const location = useLocation()
  const [bookings, setBookings] = useState<Booking[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [cancellingId, setCancellingId] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState('')

  useEffect(() => {
    fetchBookings()

    // Show success message if redirected after booking
    if (location.state?.message) {
      setSuccessMessage(location.state.message)
      setTimeout(() => setSuccessMessage(''), 5000)
      // Clear the state
      window.history.replaceState({}, document.title)
    }
  }, [location])

  const fetchBookings = async () => {
    try {
      setLoading(true)
      setError('')
      const response = await bookingsService.getMyBookings({ page: 1, limit: 50 })
      setBookings(response.bookings)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleCancelBooking = async (bookingId: string) => {
    if (!confirm('¿Estás seguro de que quieres cancelar esta reserva?')) {
      return
    }

    const reason = prompt('Motivo de cancelación (opcional):')

    try {
      setCancellingId(bookingId)
      await bookingsService.cancelBooking(bookingId, reason || undefined)
      // Refresh bookings list
      await fetchBookings()
      setSuccessMessage('Reserva cancelada exitosamente')
      setTimeout(() => setSuccessMessage(''), 5000)
    } catch (err) {
      alert(getErrorMessage(err))
    } finally {
      setCancellingId(null)
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
    const statusConfig: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline'; icon: React.ReactNode }> = {
      pending: {
        label: 'Pendiente',
        variant: 'outline',
        icon: <Clock className="w-3 h-3" />
      },
      confirmed: {
        label: 'Confirmada',
        variant: 'default',
        icon: <CheckCircle className="w-3 h-3" />
      },
      cancelled: {
        label: 'Cancelada',
        variant: 'destructive',
        icon: <X className="w-3 h-3" />
      },
      completed: {
        label: 'Completada',
        variant: 'secondary',
        icon: <CheckCircle className="w-3 h-3" />
      },
      failed: {
        label: 'Fallida',
        variant: 'destructive',
        icon: <AlertCircle className="w-3 h-3" />
      },
    }

    const config = statusConfig[status] || {
      label: status,
      variant: 'outline' as const,
      icon: <AlertCircle className="w-3 h-3" />
    }

    return (
      <Badge variant={config.variant} className="gap-1">
        {config.icon}
        {config.label}
      </Badge>
    )
  }

  if (loading) {
    return (
      <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
        <div className="container mx-auto px-4">
          <div className="max-w-6xl mx-auto">
            <div className="flex items-center justify-center py-16">
              <div className="text-center">
                <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4" />
                <p className="text-gray-600">Cargando reservas...</p>
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
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-2">Mis Reservas</h1>
            <p className="text-gray-600">Gestiona tus reservas como pasajero</p>
          </div>

          {/* Success Message */}
          {successMessage && (
            <div className="mb-6 flex items-start gap-2 p-4 bg-green-50 border border-green-200 rounded-lg">
              <CheckCircle className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-green-700 font-medium">{successMessage}</p>
            </div>
          )}

          {/* Error Message */}
          {error && (
            <div className="mb-6 flex items-start gap-2 p-4 bg-red-50 border border-red-200 rounded-lg">
              <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-red-700">{error}</p>
            </div>
          )}

          {/* Bookings List */}
          {bookings.length === 0 ? (
            <Card>
              <CardContent className="pt-6">
                <div className="text-center py-12">
                  <Users className="w-16 h-16 text-gray-400 mx-auto mb-4" />
                  <h2 className="text-2xl font-bold text-gray-900 mb-2">No tienes reservas</h2>
                  <p className="text-gray-600 mb-6">
                    Busca viajes disponibles y reserva tus asientos
                  </p>
                  <Link to="/search">
                    <Button>
                      <MapPin className="w-4 h-4 mr-2" />
                      Buscar Viajes
                    </Button>
                  </Link>
                </div>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-6">
              {bookings.map((booking) => (
                <Card key={booking.id} className="overflow-hidden">
                  <CardHeader className="bg-gradient-to-r from-primary-50 to-secondary-50">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                          <CardTitle className="text-2xl">
                            {booking.trip ? (
                              <>
                                {booking.trip.origin.city} → {booking.trip.destination.city}
                              </>
                            ) : (
                              `Viaje ID: ${booking.trip_id}`
                            )}
                          </CardTitle>
                          {getStatusBadge(booking.status)}
                        </div>
                        {booking.trip && (
                          <CardDescription className="text-base">
                            {booking.trip.origin.province} → {booking.trip.destination.province}
                          </CardDescription>
                        )}
                      </div>
                      <div className="text-right">
                        <div className="text-3xl font-bold text-primary">
                          {formatPrice(booking.total_price)}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          {booking.seats_requested} {booking.seats_requested === 1 ? 'asiento' : 'asientos'}
                        </div>
                      </div>
                    </div>
                  </CardHeader>

                  <CardContent className="pt-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      {/* Trip Details */}
                      {booking.trip && (
                        <>
                          <div>
                            <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                              <Calendar className="w-5 h-5 text-primary" />
                              <span>Salida</span>
                            </div>
                            <p className="ml-7 text-gray-900">
                              {formatDate(booking.trip.departure_datetime)}
                            </p>
                          </div>

                          <div>
                            <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                              <Calendar className="w-5 h-5 text-secondary" />
                              <span>Llegada estimada</span>
                            </div>
                            <p className="ml-7 text-gray-900">
                              {formatDate(booking.trip.estimated_arrival_datetime)}
                            </p>
                          </div>

                          <div>
                            <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                              <MapPin className="w-5 h-5 text-muted-foreground" />
                              <span>Origen</span>
                            </div>
                            <p className="ml-7 text-gray-900">{booking.trip.origin.address}</p>
                          </div>

                          <div>
                            <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                              <MapPin className="w-5 h-5 text-muted-foreground" />
                              <span>Destino</span>
                            </div>
                            <p className="ml-7 text-gray-900">{booking.trip.destination.address}</p>
                          </div>
                        </>
                      )}

                      {/* Booking Info */}
                      <div>
                        <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                          <DollarSign className="w-5 h-5 text-muted-foreground" />
                          <span>Precio por asiento</span>
                        </div>
                        <p className="ml-7 text-gray-900">
                          {formatPrice(booking.total_price / booking.seats_requested)}
                        </p>
                      </div>

                      <div>
                        <div className="flex items-center gap-2 text-gray-700 font-semibold mb-2">
                          <Users className="w-5 h-5 text-muted-foreground" />
                          <span>ID de Reserva</span>
                        </div>
                        <p className="ml-7 text-gray-900 text-sm font-mono">
                          {booking.booking_uuid || booking.id}
                        </p>
                      </div>
                    </div>

                    {/* Cancellation Info */}
                    {booking.status === 'cancelled' && booking.cancellation_reason && (
                      <div className="mt-6 pt-6 border-t">
                        <div className="flex items-start gap-2 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
                          <AlertCircle className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
                          <div>
                            <p className="text-sm font-semibold text-yellow-900 mb-1">
                              Motivo de cancelación:
                            </p>
                            <p className="text-sm text-yellow-700">
                              {booking.cancellation_reason}
                            </p>
                          </div>
                        </div>
                      </div>
                    )}
                  </CardContent>

                  <CardFooter className="bg-gray-50 border-t flex flex-col sm:flex-row gap-3">
                    {booking.trip && (
                      <Link to={`/trips/${booking.trip_id}`} className="flex-1">
                        <Button variant="outline" className="w-full">
                          Ver Detalles del Viaje
                        </Button>
                      </Link>
                    )}

                    {(booking.status === 'pending' || booking.status === 'confirmed') && (
                      <Button
                        variant="destructive"
                        onClick={() => handleCancelBooking(booking.id)}
                        disabled={cancellingId === booking.id}
                        className="flex-1"
                      >
                        {cancellingId === booking.id ? (
                          <>
                            <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                            Cancelando...
                          </>
                        ) : (
                          <>
                            <X className="w-4 h-4 mr-2" />
                            Cancelar Reserva
                          </>
                        )}
                      </Button>
                    )}

                    {booking.status === 'cancelled' && (
                      <div className="flex-1 flex items-center justify-center px-4 py-2 bg-gray-100 border border-gray-200 rounded-md">
                        <span className="text-sm text-gray-600 font-medium">
                          Reserva cancelada
                        </span>
                      </div>
                    )}

                    {booking.status === 'completed' && (
                      <div className="flex-1 flex items-center justify-center px-4 py-2 bg-green-50 border border-green-200 rounded-md">
                        <CheckCircle className="w-4 h-4 text-green-600 mr-2" />
                        <span className="text-sm text-green-700 font-medium">
                          Viaje completado
                        </span>
                      </div>
                    )}
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
