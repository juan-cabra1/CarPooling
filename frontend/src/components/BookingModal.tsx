import { useState } from 'react'
import { X, Users, DollarSign, AlertCircle, CheckCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { bookingsService, getErrorMessage } from '@/services'
import { useAuth } from '@/context/AuthContext'
import type { Trip, SearchTrip } from '@/types'

interface BookingModalProps {
  trip: Trip | SearchTrip
  isOpen: boolean
  onClose: () => void
  onSuccess: () => void
}

export default function BookingModal({ trip, isOpen, onClose, onSuccess }: BookingModalProps) {
  const { user } = useAuth()
  const [seatsRequested, setSeatsRequested] = useState(1)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!user) {
      setError('Debes iniciar sesi�n para reservar')
      return
    }

    if (seatsRequested < 1) {
      setError('Debes reservar al menos 1 asiento')
      return
    }

    if (seatsRequested > trip.available_seats) {
      setError(`Solo hay ${trip.available_seats} asientos disponibles`)
      return
    }

    try {
      setLoading(true)
      setError('')

      await bookingsService.createBooking({
        trip_id: tripId,
        passenger_id: user.id,
        seats_reserved: seatsRequested,
      })

      setSuccess(true)
      setTimeout(() => {
        onSuccess()
        handleClose()
      }, 2000)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleClose = () => {
    if (!loading) {
      setSeatsRequested(1)
      setError('')
      setSuccess(false)
      onClose()
    }
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-2xl font-bold text-gray-900">Reservar Asientos</h2>
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
          {success ? (
            <div className="text-center py-8">
              <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
              <h3 className="text-xl font-bold text-gray-900 mb-2">
                �Reserva exitosa!
              </h3>
              <p className="text-gray-600">
                Tu reserva ha sido creada. Redirigiendo...
              </p>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Trip Info */}
              <div className="bg-gray-50 rounded-lg p-4 space-y-2">
                <div className="font-semibold text-gray-900">
                  {trip.origin.city} � {trip.destination.city}
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
                    onChange={(e) => setSeatsRequested(Math.max(1, Math.min(trip.available_seats, parseInt(e.target.value) || 1)))}
                    disabled={loading}
                    className="text-center text-lg font-semibold w-24"
                  />
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setSeatsRequested(Math.min(trip.available_seats, seatsRequested + 1))}
                    disabled={seatsRequested >= trip.available_seats || loading}
                  >
                    +
                  </Button>
                </div>
                <p className="text-sm text-gray-500 mt-2">
                  M�ximo: {trip.available_seats} asientos
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

              {/* Error Message */}
              {error && (
                <div className="flex items-start gap-2 p-4 bg-red-50 border border-red-200 rounded-lg">
                  <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
                  <p className="text-sm text-red-700">{error}</p>
                </div>
              )}

              {/* Action Buttons */}
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
                <Button
                  type="submit"
                  disabled={loading}
                  className="flex-1"
                >
                  {loading ? (
                    <>
                      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                      Reservando...
                    </>
                  ) : (
                    <>
                      <CheckCircle className="w-4 h-4 mr-2" />
                      Confirmar Reserva
                    </>
                  )}
                </Button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  )
}
