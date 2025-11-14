import { Link } from 'react-router-dom'
import { MapPin, Calendar, DollarSign, Users, Star, Car } from 'lucide-react'
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import type { SearchTrip } from '@/types'

interface TripCardProps {
  trip: SearchTrip
}

export default function TripCard({ trip }: TripCardProps) {
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
    <Card className="overflow-hidden hover:shadow-xl transition-shadow duration-300 border-2 hover:border-primary-200 group">
      <CardHeader className="bg-gradient-to-r from-primary-50 to-secondary-50 pb-4">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <MapPin className="w-4 h-4 text-primary flex-shrink-0" />
              <span className="font-semibold text-gray-900">
                {trip.origin.city}, {trip.origin.province}
              </span>
            </div>
            <div className="flex items-center gap-2 ml-6">
              <div className="w-px h-4 bg-gray-300" />
            </div>
            <div className="flex items-center gap-2 mt-1">
              <MapPin className="w-4 h-4 text-secondary flex-shrink-0" />
              <span className="font-semibold text-gray-900">
                {trip.destination.city}, {trip.destination.province}
              </span>
            </div>
          </div>

          <div className="flex flex-col items-end gap-1">
            <span className="text-2xl font-bold text-primary">
              {formatPrice(trip.price_per_seat)}
            </span>
            <span className="text-xs text-muted-foreground">por persona</span>
          </div>
        </div>
      </CardHeader>

      <CardContent className="pt-4 pb-3">
        {/* Driver Info */}
        <div className="flex items-center gap-3 mb-4 pb-4 border-b">
          <div className="w-12 h-12 rounded-full bg-primary-100 flex items-center justify-center text-primary font-bold text-lg">
            {trip.driver.name.charAt(0).toUpperCase()}
          </div>
          <div className="flex-1">
            <div className="font-semibold text-gray-900">{trip.driver.name}</div>
            <div className="flex items-center gap-1 text-sm">
              <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
              <span className="font-medium">{trip.driver.rating.toFixed(1)}</span>
              <span className="text-muted-foreground">
                • {trip.driver.total_trips} viajes
              </span>
            </div>
          </div>
        </div>

        {/* Trip Details */}
        <div className="space-y-3">
          <div className="flex items-center gap-2 text-sm">
            <Calendar className="w-4 h-4 text-muted-foreground" />
            <span className="text-gray-700">{formatDate(trip.departure_datetime)}</span>
          </div>

          <div className="flex items-center gap-2 text-sm">
            <Car className="w-4 h-4 text-muted-foreground" />
            <span className="text-gray-700">
              {trip.car.brand} {trip.car.model} ({trip.car.color})
            </span>
          </div>

          <div className="flex items-center gap-2 text-sm">
            <Users className="w-4 h-4 text-muted-foreground" />
            <span className="text-gray-700">
              {trip.available_seats} {trip.available_seats === 1 ? 'asiento disponible' : 'asientos disponibles'}
            </span>
          </div>
        </div>

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
        <Link to={`/trips/${trip.trip_id}`} className="w-full">
          <Button className="w-full group-hover:bg-primary-600 transition-colors">
            Ver Detalles
          </Button>
        </Link>
      </CardFooter>
    </Card>
  )
}
