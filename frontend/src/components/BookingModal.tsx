import { useState } from 'react'
import { X, Users, DollarSign, AlertCircle, CheckCircle, CreditCard, ExternalLink } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { bookingsService, paymentsService, getErrorMessage } from '@/services'
import { useAuth } from '@/context/AuthContext'
import type { Trip, SearchTrip } from '@/types'
import type { Booking } from '@/types/booking'
import type { CreatePreferenceResponse } from '@/types/payment'

interface BookingModalProps {
  trip: Trip | SearchTrip
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

type Step = 'booking' | 'payment'

export default function BookingModal({ trip, isOpen, onClose, onSuccess }: BookingModalProps) {
  const { user } = useAuth()
  const [step, setStep] = useState<Step>('booking')
  const [seatsRequested, setSeatsRequested] = useState(1)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [booking, setBooking] = useState<Booking | null>(null)
  const [preference, setPreference] = useState<CreatePreferenceResponse | null>(null)

  // Handle both Trip and SearchTrip types
  const tripId = 'trip_id' in trip ? trip.trip_id : trip.id
  const totalPrice = seatsRequested * trip.price_per_seat

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS',
      minimumFractionDigits: 0,
    }).format(price)
  }

  // Step 1: create booking
  const handleSubmitBooking = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!user) {
      setError('Debés iniciar sesión para reservar')
      return
    }
    if (seatsRequested < 1) {
      setError('Debés reservar al menos 1 asiento')
      return
    }
    if (seatsRequested > trip.available_seats) {
      setError(`Solo hay ${trip.available_seats} asientos disponibles`)
      return
    }

    try {
      setLoading(true)
      setError('')

      const newBooking = await bookingsService.createBooking({
        trip_id: tripId,
        passenger_id: user.id,
        seats_reserved: seatsRequested,
      })
      setBooking(newBooking)

      // Step 2: create MercadoPago preference
      const pref = await paymentsService.createPreference({
        booking_id: newBooking.id,
        trip_id: tripId,
        driver_id: String(trip.driver_id),
        seats_count: seatsRequested,
        price_per_seat: Math.round(trip.price_per_seat * 100), // ARS → centavos
        origin: `${trip.origin.city}, ${trip.origin.province}`,
        destination: `${trip.destination.city}, ${trip.destination.province}`,
        departure_at: trip.departure_datetime,
        payer_email: user.email,
        payer_name: user.name,
        return_base_url: window.location.origin,
      })
      setPreference(pref)
      setStep('payment')
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  // Redirect to MercadoPago checkout
  const handleGoToPayment = () => {
    if (!preference) return
    // In production use init_point; in sandbox use sandbox_init_point
    const checkoutURL = import.meta.env.VITE_MP_SANDBOX === 'true'
      ? preference.sandbox_init_point
      : preference.init_point
    window.location.href = checkoutURL
  }

  const handleClose = () => {
    if (!loading) {
      setSeatsRequested(1)
      setError('')
      setStep('booking')
      setBooking(null)
      setPreference(null)
      onClose()
    }
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-2xl font-bold text-gray-900">
            {step === 'booking' ? 'Reservar Asientos' : 'Confirmar Pago'}
          </h2>
          <button
            onClick={handleClose}
            disabled={loading}
            className="text-gray-400 hover:text-gray-600 transition-colors disabled:opacity-50"
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          {/* ── Step 1: Booking form ── */}
          {step === 'booking' && (
            <form onSubmit={handleSubmitBooking} className="space-y-6">
              {/* Trip Info */}
              <div className="bg-gray-50 rounded-lg p-4 space-y-2">
                <div className="font-semibold text-gray-900">
                  {trip.origin.city} → {trip.destination.city}
                </div>
                <div className="text-sm text-gray-600">
                  {new Intl.DateTimeFormat('es-AR', {
                    weekday: 'short',
                    day: 'numeric',
                    month: 'short',
                    hour: '2-digit',
                    minute: '2-digit',
                  }).format(new Date(trip.departure_datetime))}
                </div>
                <div className="flex items-center gap-2 text-sm text-gray-600">
                  <Users className="w-4 h-4" />
                  <span>{trip.available_seats} asientos disponibles</span>
                </div>
              </div>

              {/* Seats Selection */}
              <div>
                <Label htmlFor="seats" className="text-base font-semibold mb-2 block">
                  Cantidad de asientos
                </Label>
                <div className="flex items-center gap-3">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setSeatsRequested(Math.max(1, seatsRequested - 1))}
                    disabled={seatsRequested <= 1 || loading}
                  >
                    -
                  </Button>
                  <Input
                    id="seats"
                    type="number"
                    min="1"
                    max={trip.available_seats}
                    value={seatsRequested}
                    onChange={(e) =>
                      setSeatsRequested(
                        Math.max(1, Math.min(trip.available_seats, parseInt(e.target.value) || 1))
                      )
                    }
                    disabled={loading}
                    className="text-center text-lg font-semibold w-24"
                  />
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() =>
                      setSeatsRequested(Math.min(trip.available_seats, seatsRequested + 1))
                    }
                    disabled={seatsRequested >= trip.available_seats || loading}
                  >
                    +
                  </Button>
                </div>
                <p className="text-sm text-gray-500 mt-2">
                  Máximo: {trip.available_seats} asientos
                </p>
              </div>

              {/* Price Summary */}
              <div className="bg-primary-50 rounded-lg p-4 space-y-2">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600">Precio por asiento</span>
                  <span className="font-medium">{formatPrice(trip.price_per_seat)}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-gray-600">Asientos</span>
                  <span className="font-medium">{seatsRequested}</span>
                </div>
                <div className="border-t border-primary-200 pt-2 mt-2">
                  <div className="flex items-center justify-between">
                    <span className="font-semibold text-gray-900">Total</span>
                    <div className="flex items-center gap-2">
                      <DollarSign className="w-5 h-5 text-primary" />
                      <span className="text-2xl font-bold text-primary">
                        {formatPrice(totalPrice)}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              {error && (
                <div className="flex items-start gap-2 p-4 bg-red-50 border border-red-200 rounded-lg">
                  <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
                  <p className="text-sm text-red-700">{error}</p>
                </div>
              )}

              <div className="flex gap-3 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleClose}
                  disabled={loading}
                  className="flex-1"
                >
                  Cancelar
                </Button>
                <Button type="submit" disabled={loading} className="flex-1">
                  {loading ? (
                    <>
                      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                      Procesando...
                    </>
                  ) : (
                    <>
                      <CheckCircle className="w-4 h-4 mr-2" />
                      Continuar al pago
                    </>
                  )}
                </Button>
              </div>
            </form>
          )}

          {/* ── Step 2: Payment confirmation ── */}
          {step === 'payment' && preference && booking && (
            <div className="space-y-6">
              <div className="text-center py-2">
                <CheckCircle className="w-12 h-12 text-green-500 mx-auto mb-3" />
                <h3 className="text-lg font-bold text-gray-900">¡Reserva creada!</h3>
                <p className="text-gray-600 text-sm mt-1">
                  Tu reserva fue registrada. Ahora completá el pago para confirmarla.
                </p>
              </div>

              {/* Payment summary */}
              <div className="bg-gray-50 rounded-lg p-4 space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">{trip.origin.city} → {trip.destination.city}</span>
                  <span className="font-medium">{seatsRequested} asiento{seatsRequested > 1 ? 's' : ''}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Total a pagar</span>
                  <span className="font-bold text-primary text-lg">
                    {formatPrice(preference.amount / 100)}
                  </span>
                </div>
                <div className="flex justify-between text-xs text-gray-500">
                  <span>Fee de plataforma</span>
                  <span>{formatPrice(preference.platform_fee / 100)}</span>
                </div>
              </div>

              <p className="text-xs text-gray-500 text-center">
                Serás redirigido a MercadoPago para completar el pago de forma segura.
              </p>

              <div className="flex gap-3">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => {
                    // Booking was already created; just close and navigate
                    onSuccess()
                    handleClose()
                  }}
                  className="flex-1"
                >
                  Pagar después
                </Button>
                <Button onClick={handleGoToPayment} className="flex-1 gap-2">
                  <CreditCard className="w-4 h-4" />
                  Pagar con MercadoPago
                  <ExternalLink className="w-3 h-3" />
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
