import { useState, useEffect, useRef } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { MapPin, Calendar, Users, DollarSign, Car, ArrowLeft, Edit, Trash2, AlertCircle, Star, User, Send } from 'lucide-react'
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
import { tripsService, searchService, chatService, getErrorMessage } from '@/services'
import type { Message } from '@/services'
import { useAuth } from '@/context/AuthContext'
import BookingModal from '@/components/BookingModal'
import type { SearchTrip } from '@/types'

export default function TripDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()
  const [trip, setTrip] = useState<SearchTrip | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [deleting, setDeleting] = useState(false)
  const [isBookingModalOpen, setIsBookingModalOpen] = useState(false)

  // Chat states
  const [messages, setMessages] = useState<Message[]>([])
  const [newMessage, setNewMessage] = useState('')
  const [isSending, setIsSending] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

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
      // Use search-api to get denormalized trip data with driver information
      const data = await searchService.getTripDetails(id)
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

  const handleBookingSuccess = () => {
    // Refresh trip data to update available seats
    fetchTrip()
    // Navigate to bookings page
    navigate('/my-bookings', { state: { message: 'Reserva creada exitosamente' } })
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

  // Chat functions
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  const loadMessages = async () => {
    if (!id || !user) return

    try {
      const { messages: msgs } = await chatService.getMessages(id)
      setMessages(msgs)
    } catch (error) {
      console.error('Failed to load messages:', error)
    }
  }

  const handleSendMessage = async () => {
    if (!newMessage.trim() || isSending || !id) return

    setIsSending(true)
    try {
      await chatService.sendMessage(id, newMessage)
      setNewMessage('')
      await loadMessages() // Refresh messages immediately
    } catch (error) {
      console.error('Failed to send message:', error)
      alert('Error al enviar mensaje')
    } finally {
      setIsSending(false)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  // Load messages on mount
  useEffect(() => {
    if (user && id) {
      loadMessages()
    }
  }, [id, user])

  // Auto-refresh messages every 5 seconds
  useEffect(() => {
    if (!user || !id) return

    const interval = setInterval(() => {
      loadMessages()
    }, 5000)

    return () => clearInterval(interval)
  }, [id, user])

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    scrollToBottom()
  }, [messages])

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
                    <CardTitle className="text-3xl text-gray-900">
                      {trip.origin.city} → {trip.destination.city}
                    </CardTitle>
                    {getStatusBadge(trip.status)}
                  </div>
                  <CardDescription className="text-lg text-gray-700">
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
                    {trip.total_seats - trip.available_seats} reservados
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

              {/* Driver Information - Only show if driver data is available */}
              {trip.driver && (
                <div className="pt-6 border-t">
                  <h3 className="font-semibold text-gray-900 mb-4 text-lg flex items-center gap-2">
                    <User className="w-5 h-5 text-primary" />
                    Conductor
                  </h3>
                  <div className="flex items-center gap-4 p-4 bg-gradient-to-r from-primary-50 to-secondary-50 rounded-lg">
                    <div className="w-16 h-16 rounded-full bg-primary-100 flex items-center justify-center text-primary font-bold text-2xl flex-shrink-0">
                      {trip.driver.name.charAt(0).toUpperCase()}
                    </div>
                    <div className="flex-1">
                      <h4 className="font-semibold text-gray-900 text-lg">{trip.driver.name}</h4>
                      <div className="flex items-center gap-3 mt-1">
                        <div className="flex items-center gap-1">
                          <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                          <span className="font-medium text-gray-900">{trip.driver.rating.toFixed(1)}</span>
                        </div>
                        <span className="text-gray-500">•</span>
                        <span className="text-gray-600">{trip.driver.total_trips} viajes realizados</span>
                      </div>
                      <p className="text-sm text-gray-600 mt-1">{trip.driver.email}</p>
                    </div>
                  </div>
                </div>
              )}

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
                  {trip.status === 'published' && (trip.total_seats - trip.available_seats) === 0 && (
                    <>
                      <Link to={`/trips/${trip.trip_id}/edit`} className="flex-1">
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
                  {(trip.total_seats - trip.available_seats) > 0 && (
                    <div className="w-full flex items-center justify-center px-4 py-3 bg-yellow-50 border border-yellow-200 rounded-md">
                      <AlertCircle className="w-5 h-5 text-yellow-600 mr-2" />
                      <span className="text-sm text-yellow-700 font-medium">
                        Este viaje tiene reservas activas y no puede ser editado o eliminado
                      </span>
                    </div>
                  )}
                </>
              ) : (
                <>
                  {trip.status === 'published' && trip.available_seats > 0 ? (
                    <Button
                      className="w-full"
                      size="lg"
                      onClick={() => setIsBookingModalOpen(true)}
                    >
                      <Users className="w-5 h-5 mr-2" />
                      Reservar Asientos
                    </Button>
                  ) : (
                    <div className="w-full flex items-center justify-center px-4 py-3 bg-gray-100 border border-gray-200 rounded-md">
                      <AlertCircle className="w-5 h-5 text-gray-600 mr-2" />
                      <span className="text-sm text-gray-700 font-medium">
                        {trip.available_seats === 0 ? 'No hay asientos disponibles' : 'Viaje no disponible para reservas'}
                      </span>
                    </div>
                  )}
                </>
              )}
            </CardFooter>
          </Card>

          {/* Booking Modal */}
          {trip && (
            <BookingModal
              trip={trip}
              isOpen={isBookingModalOpen}
              onClose={() => setIsBookingModalOpen(false)}
              onSuccess={handleBookingSuccess}
            />
          )}

          {/* Chat Section - Only show for logged-in users */}
          {user && trip && (
            <Card className="mt-8">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <User className="w-5 h-5" />
                  Chat del Viaje
                </CardTitle>
                <CardDescription>
                  Comunícate con el conductor y otros pasajeros
                </CardDescription>
              </CardHeader>
              <CardContent>
                {/* Messages container */}
                <div className="border rounded-lg p-4 h-96 overflow-y-auto bg-gray-50 mb-4 space-y-3">
                  {messages.length === 0 ? (
                    <p className="text-gray-500 text-center py-8">
                      No hay mensajes aún. ¡Sé el primero en escribir!
                    </p>
                  ) : (
                    messages.map((msg) => (
                      <div
                        key={msg.id}
                        className={`p-3 rounded-lg ${
                          msg.user_id === user?.id
                            ? 'bg-primary-100 ml-auto max-w-[80%]'
                            : 'bg-white max-w-[80%] shadow-sm'
                        }`}
                      >
                        <div className="flex items-center gap-2 mb-1">
                          <span className="font-semibold text-sm text-gray-900">
                            {msg.user_name}
                          </span>
                          <span className="text-xs text-gray-500">
                            {new Date(msg.created_at).toLocaleTimeString('es-AR', {
                              hour: '2-digit',
                              minute: '2-digit',
                            })}
                          </span>
                        </div>
                        <p className="text-gray-800">{msg.message}</p>
                      </div>
                    ))
                  )}
                  <div ref={messagesEndRef} />
                </div>

                {/* Input area */}
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    onKeyPress={handleKeyPress}
                    placeholder="Escribe un mensaje..."
                    className="flex-1 border rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500"
                    disabled={isSending}
                  />
                  <Button
                    onClick={handleSendMessage}
                    disabled={!newMessage.trim() || isSending}
                    size="lg"
                  >
                    <Send className="w-4 h-4" />
                  </Button>
                </div>

                <p className="text-xs text-gray-500 mt-2">
                  💡 Los mensajes se actualizan automáticamente cada 5 segundos
                </p>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}
