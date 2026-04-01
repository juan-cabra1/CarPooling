import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Calendar, MapPin, Users, DollarSign, Edit, Trash2, Car, AlertCircle, UserCheck } from 'lucide-react'
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
import { tripsService, bookingsService, getErrorMessage } from '@/services'
import type { Trip, Booking } from '@/types'

type TabType = 'conductor' | 'pasajero'

export default function MyTripsPage() {
  const [activeTab, setActiveTab] = useState<TabType>('conductor')

  // Conductor state
  const [driverTrips, setDriverTrips] = useState<Trip[]>([])
  const [driverLoading, setDriverLoading] = useState(true)
  const [driverError, setDriverError] = useState('')
  const [deletingId, setDeletingId] = useState<string | null>(null)
  const [tripToDelete, setTripToDelete] = useState<string | null>(null)
  const [deleteError, setDeleteError] = useState('')

  // Pasajero state
  const [bookings, setBookings] = useState<Booking[]>([])
  const [passengerLoading, setPassengerLoading] = useState(false)
  const [passengerError, setPassengerError] = useState('')
  const [passengerLoaded, setPassengerLoaded] = useState(false)

  useEffect(() => {
    fetchDriverTrips()
  }, [])

  useEffect(() => {
    if (activeTab === 'pasajero' && !passengerLoaded) {
      fetchPassengerBookings()
    }
  }, [activeTab])

  const fetchDriverTrips = async () => {
    try {
      setDriverLoading(true)
      setDriverError('')
      const response = await tripsService.getMyDriverTrips()
      setDriverTrips(response.trips)
    } catch (err) {
      setDriverError(getErrorMessage(err))
    } finally {
      setDriverLoading(false)
    }
  }

  const fetchPassengerBookings = async () => {
    try {
      setPassengerLoading(true)
      setPassengerError('')
      const response = await bookingsService.getMyBookings()
      setBookings(response.bookings)
      setPassengerLoaded(true)
    } catch (err) {
      setPassengerError(getErrorMessage(err))
    } finally {
      setPassengerLoading(false)
    }
  }

  const handleDeleteTrip = async () => {
    if (!tripToDelete) return
    setDeleteError('')
    try {
      setDeletingId(tripToDelete)
      await tripsService.deleteTrip(tripToDelete)
      setTripToDelete(null)
      await fetchDriverTrips()
    } catch (err) {
      setDeleteError(getErrorMessage(err))
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

  const getTripStatusBadge = (status: string) => {
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

  const getBookingStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }> = {
      pending: { label: 'Pendiente', variant: 'outline' },
      confirmed: { label: 'Confirmada', variant: 'default' },
      cancelled: { label: 'Cancelada', variant: 'destructive' },
      completed: { label: 'Completada', variant: 'secondary' },
      failed: { label: 'Fallida', variant: 'destructive' },
    }
    const config = statusConfig[status] || { label: status, variant: 'outline' as const }
    return <Badge variant={config.variant}>{config.label}</Badge>
  }

  return (
    <>
    <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
      <div className="container mx-auto px-4">
        <div className="max-w-6xl mx-auto">
          {/* Header */}
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-8 gap-4">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">Mis Viajes</h1>
              <p className="text-gray-600">Gestioná tus viajes como conductor y tus reservas como pasajero</p>
            </div>
            <Link to="/create-trip">
              <Button size="lg" className="w-full sm:w-auto">
                <Plus className="w-5 h-5 mr-2" />
                Publicar Viaje
              </Button>
            </Link>
          </div>

          {/* Tabs */}
          <div className="flex border-b border-gray-200 mb-6">
            <button
              onClick={() => setActiveTab('conductor')}
              className={`px-6 py-3 text-sm font-medium border-b-2 transition-colors ${
                activeTab === 'conductor'
                  ? 'border-primary text-primary'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <Car className="w-4 h-4 inline mr-2" />
              Como Conductor
            </button>
            <button
              onClick={() => setActiveTab('pasajero')}
              className={`px-6 py-3 text-sm font-medium border-b-2 transition-colors ${
                activeTab === 'pasajero'
                  ? 'border-primary text-primary'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <UserCheck className="w-4 h-4 inline mr-2" />
              Como Pasajero
            </button>
          </div>

          {/* ===== TAB: CONDUCTOR ===== */}
          {activeTab === 'conductor' && (
            <>
              {driverLoading ? (
                <div className="flex items-center justify-center py-16">
                  <div className="text-center">
                    <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4" />
                    <p className="text-gray-600">Cargando tus viajes...</p>
                  </div>
                </div>
              ) : (
                <>
                  {driverError && (
                    <div className="mb-6 p-4 rounded-lg bg-destructive/10 border border-destructive/20 flex items-start gap-3">
                      <AlertCircle className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
                      <div>
                        <p className="font-semibold text-destructive">Error al cargar viajes</p>
                        <p className="text-sm text-destructive/80">{driverError}</p>
                      </div>
                    </div>
                  )}

                  {driverTrips.length === 0 ? (
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
                      {driverTrips.map((trip) => (
                        <Card key={trip.id} className="overflow-hidden hover:shadow-lg transition-shadow">
                          <CardHeader className="bg-gradient-to-r from-primary-50 to-secondary-50 pb-4">
                            <div className="flex items-start justify-between">
                              <div className="flex-1">
                                <div className="flex items-center gap-3 mb-1">
                                  <CardTitle className="text-xl text-gray-900">
                                    {trip.origin.city} → {trip.destination.city}
                                  </CardTitle>
                                  {getTripStatusBadge(trip.status)}
                                </div>
                                <CardDescription className="text-base text-gray-700">
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
                              <div className="space-y-3">
                                <div className="flex items-start gap-2 text-sm">
                                  <Calendar className="w-4 h-4 text-muted-foreground mt-0.5 flex-shrink-0" />
                                  <div>
                                    <div className="font-medium text-gray-900">Llegada estimada</div>
                                    <div className="text-gray-700">{formatDate(trip.estimated_arrival_datetime)}</div>
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

                            <div className="grid grid-cols-2 md:grid-cols-3 gap-4 pt-4 border-t">
                              <div className="flex items-center gap-2 text-sm">
                                <Car className="w-4 h-4 text-muted-foreground" />
                                <div>
                                  <div className="font-medium text-gray-900">Vehículo</div>
                                  <div className="text-gray-700">{trip.car.brand} {trip.car.model}</div>
                                </div>
                              </div>
                              <div className="flex items-center gap-2 text-sm">
                                <Users className="w-4 h-4 text-muted-foreground" />
                                <div>
                                  <div className="font-medium text-gray-900">Asientos</div>
                                  <div className="text-gray-700">{trip.available_seats}/{trip.total_seats} disponibles</div>
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

                            {trip.description && (
                              <div className="mt-4 pt-4 border-t">
                                <p className="text-sm text-gray-700">{trip.description}</p>
                              </div>
                            )}

                            {(trip.preferences.pets_allowed || trip.preferences.smoking_allowed || trip.preferences.music_allowed) && (
                              <div className="flex flex-wrap gap-2 mt-4">
                                {trip.preferences.pets_allowed && <Badge variant="secondary" className="text-xs">🐕 Mascotas</Badge>}
                                {trip.preferences.smoking_allowed && <Badge variant="secondary" className="text-xs">🚬 Fumar</Badge>}
                                {trip.preferences.music_allowed && <Badge variant="secondary" className="text-xs">🎵 Música</Badge>}
                              </div>
                            )}
                          </CardContent>

                          <CardFooter className="bg-gray-50 border-t">
                            <div className="w-full flex flex-col sm:flex-row gap-3">
                              <Link to={`/trips/${trip.id}`} className="flex-1">
                                <Button variant="outline" className="w-full">Ver Detalles</Button>
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
                                    onClick={() => { setTripToDelete(trip.id); setDeleteError('') }}
                                    disabled={deletingId === trip.id}
                                    className="flex-1"
                                  >
                                    {deletingId === trip.id ? (
                                      <><div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />Eliminando...</>
                                    ) : (
                                      <><Trash2 className="w-4 h-4 mr-2" />Eliminar</>
                                    )}
                                  </Button>
                                </>
                              )}

                              {trip.reserved_seats > 0 && (
                                <div className="flex-1 flex items-center justify-center px-4 py-2 bg-yellow-50 border border-yellow-200 rounded-md">
                                  <AlertCircle className="w-4 h-4 text-yellow-600 mr-2" />
                                  <span className="text-sm text-yellow-700 font-medium">No se puede editar (tiene reservas)</span>
                                </div>
                              )}
                            </div>
                          </CardFooter>
                        </Card>
                      ))}
                    </div>
                  )}
                </>
              )}
            </>
          )}

          {/* ===== TAB: PASAJERO ===== */}
          {activeTab === 'pasajero' && (
            <>
              {passengerLoading ? (
                <div className="flex items-center justify-center py-16">
                  <div className="text-center">
                    <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4" />
                    <p className="text-gray-600">Cargando tus reservas...</p>
                  </div>
                </div>
              ) : (
                <>
                  {passengerError && (
                    <div className="mb-6 p-4 rounded-lg bg-destructive/10 border border-destructive/20 flex items-start gap-3">
                      <AlertCircle className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
                      <div>
                        <p className="font-semibold text-destructive">Error al cargar reservas</p>
                        <p className="text-sm text-destructive/80">{passengerError}</p>
                      </div>
                    </div>
                  )}

                  {bookings.length === 0 ? (
                    <Card className="text-center py-16">
                      <CardContent>
                        <div className="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                          <UserCheck className="w-10 h-10 text-gray-400" />
                        </div>
                        <CardTitle className="text-2xl mb-2">No tienes reservas activas</CardTitle>
                        <CardDescription className="text-lg mb-6">
                          Buscá un viaje disponible y reservá tu lugar
                        </CardDescription>
                        <Link to="/search">
                          <Button size="lg">Buscar Viajes</Button>
                        </Link>
                      </CardContent>
                    </Card>
                  ) : (
                    <div className="grid grid-cols-1 gap-6">
                      {bookings.map((booking) => (
                        <Card key={booking.id} className="overflow-hidden hover:shadow-lg transition-shadow">
                          <CardHeader className="bg-gradient-to-r from-secondary-50 to-primary-50 pb-4">
                            <div className="flex items-start justify-between">
                              <div className="flex-1">
                                <div className="flex items-center gap-3 mb-1">
                                  <CardTitle className="text-xl text-gray-900">
                                    {booking.trip
                                      ? `${booking.trip.origin.city} → ${booking.trip.destination.city}`
                                      : `Reserva #${booking.id.slice(0, 8)}`}
                                  </CardTitle>
                                  {getBookingStatusBadge(booking.status)}
                                </div>
                                {booking.trip && (
                                  <CardDescription className="text-base text-gray-700">
                                    {booking.trip.origin.province} → {booking.trip.destination.province}
                                  </CardDescription>
                                )}
                              </div>
                              <div className="text-right">
                                <div className="text-2xl font-bold text-secondary">
                                  {formatPrice(booking.total_price)}
                                </div>
                                <div className="text-xs text-muted-foreground">total</div>
                              </div>
                            </div>
                          </CardHeader>

                          <CardContent className="pt-6">
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                              {booking.trip && (
                                <>
                                  <div className="flex items-start gap-2 text-sm">
                                    <Calendar className="w-4 h-4 text-muted-foreground mt-0.5 flex-shrink-0" />
                                    <div>
                                      <div className="font-medium text-gray-900">Salida</div>
                                      <div className="text-gray-700">{formatDate(booking.trip.departure_datetime)}</div>
                                    </div>
                                  </div>
                                  <div className="flex items-start gap-2 text-sm">
                                    <MapPin className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
                                    <div>
                                      <div className="font-medium text-gray-900">Punto de encuentro</div>
                                      <div className="text-gray-700">{booking.trip.origin.address}</div>
                                    </div>
                                  </div>
                                </>
                              )}
                              <div className="flex items-center gap-2 text-sm">
                                <Users className="w-4 h-4 text-muted-foreground" />
                                <div>
                                  <div className="font-medium text-gray-900">Asientos reservados</div>
                                  <div className="text-gray-700">{booking.seats_requested}</div>
                                </div>
                              </div>
                              <div className="flex items-center gap-2 text-sm">
                                <Calendar className="w-4 h-4 text-muted-foreground" />
                                <div>
                                  <div className="font-medium text-gray-900">Reservado el</div>
                                  <div className="text-gray-700">{formatDate(booking.created_at)}</div>
                                </div>
                              </div>
                            </div>

                            {booking.cancellation_reason && (
                              <div className="mt-4 pt-4 border-t">
                                <p className="text-sm text-gray-500">
                                  <span className="font-medium">Motivo de cancelación:</span> {booking.cancellation_reason}
                                </p>
                              </div>
                            )}
                          </CardContent>

                          {booking.trip && booking.status !== 'cancelled' && booking.status !== 'completed' && (
                            <CardFooter className="bg-gray-50 border-t">
                              <div className="w-full flex gap-3">
                                <Link to={`/trips/${booking.trip_id}`} className="flex-1">
                                  <Button variant="outline" className="w-full">Ver Viaje</Button>
                                </Link>
                              </div>
                            </CardFooter>
                          )}
                        </Card>
                      ))}
                    </div>
                  )}
                </>
              )}
            </>
          )}
        </div>
      </div>
    </div>

    {/* Delete Confirmation Modal */}
    {tripToDelete && (
      <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div className="bg-white rounded-xl shadow-xl p-6 max-w-sm w-full">
          <h3 className="text-lg font-semibold mb-2">Eliminar viaje</h3>
          <p className="text-sm text-gray-600 mb-4">
            ¿Estás seguro de que querés eliminar este viaje? Esta acción no se puede deshacer.
          </p>
          {deleteError && (
            <p className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg px-3 py-2 mb-4">{deleteError}</p>
          )}
          <div className="flex gap-3">
            <Button
              variant="outline"
              onClick={() => { setTripToDelete(null); setDeleteError('') }}
              disabled={!!deletingId}
              className="flex-1"
            >
              Cancelar
            </Button>
            <Button
              variant="destructive"
              onClick={handleDeleteTrip}
              disabled={!!deletingId}
              className="flex-1"
            >
              {deletingId ? (
                <><div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />Eliminando...</>
              ) : 'Eliminar'}
            </Button>
          </div>
        </div>
      </div>
    )}
    </>
  )
}
