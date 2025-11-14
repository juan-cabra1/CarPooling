import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeft, Calendar, MapPin, DollarSign, Car, Users, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { tripsService, getErrorMessage } from '@/services'
import type { CreateTripRequest, Location, Car as CarType, Preferences } from '@/types'

export default function CreateTripPage() {
  const navigate = useNavigate()
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  const [formData, setFormData] = useState<CreateTripRequest>({
    origin: {
      city: '',
      province: '',
      address: '',
      coordinates: { lat: 0, lng: 0 },
    },
    destination: {
      city: '',
      province: '',
      address: '',
      coordinates: { lat: 0, lng: 0 },
    },
    departure_datetime: '',
    estimated_arrival_datetime: '',
    price_per_seat: 0,
    total_seats: 1,
    car: {
      brand: '',
      model: '',
      year: new Date().getFullYear(),
      color: '',
      plate: '',
    },
    preferences: {
      pets_allowed: false,
      smoking_allowed: false,
      music_allowed: false,
    },
    description: '',
  })

  // Coordenadas aproximadas de ciudades argentinas (en el futuro usar API de geocodificación)
  const getCityCoordinates = (city: string): { lat: number; lng: number } => {
    const cityCoords: Record<string, { lat: number; lng: number }> = {
      'cordoba': { lat: -31.4201, lng: -64.1888 },
      'buenos aires': { lat: -34.6037, lng: -58.3816 },
      'rosario': { lat: -32.9442, lng: -60.6505 },
      'mendoza': { lat: -32.8895, lng: -68.8458 },
      'tucuman': { lat: -26.8083, lng: -65.2176 },
      'salta': { lat: -24.7859, lng: -65.4117 },
      'santa fe': { lat: -31.6333, lng: -60.7000 },
      'mar del plata': { lat: -38.0055, lng: -57.5426 },
      'catamarca': { lat: -28.4696, lng: -65.7795 },
      'villa carlos paz': { lat: -31.4204, lng: -64.4975 },
    }

    const normalizedCity = city.toLowerCase().trim()
    return cityCoords[normalizedCity] || { lat: -34.6037, lng: -58.3816 } // Default: Buenos Aires
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      setSaving(true)
      setError('')

      // Validación de fechas
      const departureDate = new Date(formData.departure_datetime)
      const arrivalDate = new Date(formData.estimated_arrival_datetime)
      const now = new Date()

      if (departureDate < now) {
        setError('La fecha de salida debe ser en el futuro')
        return
      }

      if (arrivalDate <= departureDate) {
        setError('La fecha de llegada debe ser posterior a la fecha de salida')
        return
      }

      // Obtener coordenadas basadas en las ciudades
      const originCoords = getCityCoordinates(formData.origin.city)
      const destCoords = getCityCoordinates(formData.destination.city)

      // Convertir fechas a ISO 8601 y agregar coordenadas
      const createData: CreateTripRequest = {
        ...formData,
        origin: {
          ...formData.origin,
          coordinates: originCoords,
        },
        destination: {
          ...formData.destination,
          coordinates: destCoords,
        },
        departure_datetime: departureDate.toISOString(),
        estimated_arrival_datetime: arrivalDate.toISOString(),
      }

      const newTrip = await tripsService.createTrip(createData)
      navigate(`/trips/${newTrip.id}`, {
        state: { message: 'Viaje publicado exitosamente' }
      })
    } catch (err) {
      setError(getErrorMessage(err))
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } finally {
      setSaving(false)
    }
  }

  const handleChange = (field: keyof CreateTripRequest, value: any) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }))
  }

  const handleNestedChange = (parent: keyof CreateTripRequest, field: string, value: any) => {
    setFormData(prev => ({
      ...prev,
      [parent]: {
        ...(prev[parent] as any),
        [field]: value,
      },
    }))
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

          <form onSubmit={handleSubmit}>
            <Card>
              <CardHeader>
                <CardTitle className="text-3xl">Publicar Nuevo Viaje</CardTitle>
                <CardDescription>
                  Completa los detalles de tu viaje para que otros puedan unirse
                </CardDescription>
              </CardHeader>

              <CardContent className="space-y-8">
                {/* Error Message */}
                {error && (
                  <div className="p-4 rounded-lg bg-destructive/10 border border-destructive/20 flex items-start gap-3">
                    <AlertCircle className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
                    <div>
                      <p className="font-semibold text-destructive">Error al crear viaje</p>
                      <p className="text-sm text-destructive/80">{error}</p>
                    </div>
                  </div>
                )}

                {/* Origen y Destino */}
                <div className="space-y-4">
                  <h3 className="text-xl font-semibold flex items-center gap-2">
                    <MapPin className="w-5 h-5 text-primary" />
                    Ubicaciones
                  </h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Origen */}
                    <div className="space-y-4 p-4 border rounded-lg bg-white">
                      <h4 className="font-semibold text-primary">Origen</h4>

                      <div className="space-y-2">
                        <Label htmlFor="origin-city">
                          Ciudad <span className="text-destructive">*</span>
                        </Label>
                        <Input
                          id="origin-city"
                          value={formData.origin.city}
                          onChange={(e) => handleNestedChange('origin', 'city', e.target.value)}
                          placeholder="Ej: Córdoba"
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="origin-province">
                          Provincia <span className="text-destructive">*</span>
                        </Label>
                        <Input
                          id="origin-province"
                          value={formData.origin.province}
                          onChange={(e) => handleNestedChange('origin', 'province', e.target.value)}
                          placeholder="Ej: Córdoba"
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="origin-address">
                          Dirección <span className="text-destructive">*</span>
                        </Label>
                        <Input
                          id="origin-address"
                          value={formData.origin.address}
                          onChange={(e) => handleNestedChange('origin', 'address', e.target.value)}
                          placeholder="Ej: Av. Colón 1234"
                          required
                        />
                      </div>
                    </div>

                    {/* Destino */}
                    <div className="space-y-4 p-4 border rounded-lg bg-white">
                      <h4 className="font-semibold text-secondary">Destino</h4>

                      <div className="space-y-2">
                        <Label htmlFor="destination-city">
                          Ciudad <span className="text-destructive">*</span>
                        </Label>
                        <Input
                          id="destination-city"
                          value={formData.destination.city}
                          onChange={(e) => handleNestedChange('destination', 'city', e.target.value)}
                          placeholder="Ej: Buenos Aires"
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="destination-province">
                          Provincia <span className="text-destructive">*</span>
                        </Label>
                        <Input
                          id="destination-province"
                          value={formData.destination.province}
                          onChange={(e) => handleNestedChange('destination', 'province', e.target.value)}
                          placeholder="Ej: Buenos Aires"
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="destination-address">
                          Dirección <span className="text-destructive">*</span>
                        </Label>
                        <Input
                          id="destination-address"
                          value={formData.destination.address}
                          onChange={(e) => handleNestedChange('destination', 'address', e.target.value)}
                          placeholder="Ej: Av. 9 de Julio 1000"
                          required
                        />
                      </div>
                    </div>
                  </div>
                </div>

                {/* Fechas */}
                <div className="space-y-4">
                  <h3 className="text-xl font-semibold flex items-center gap-2">
                    <Calendar className="w-5 h-5 text-primary" />
                    Fechas y Horarios
                  </h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="departure">
                        Fecha y hora de salida <span className="text-destructive">*</span>
                      </Label>
                      <Input
                        id="departure"
                        type="datetime-local"
                        value={formData.departure_datetime}
                        onChange={(e) => handleChange('departure_datetime', e.target.value)}
                        min={new Date().toISOString().slice(0, 16)}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="arrival">
                        Fecha y hora estimada de llegada <span className="text-destructive">*</span>
                      </Label>
                      <Input
                        id="arrival"
                        type="datetime-local"
                        value={formData.estimated_arrival_datetime}
                        onChange={(e) => handleChange('estimated_arrival_datetime', e.target.value)}
                        min={formData.departure_datetime || new Date().toISOString().slice(0, 16)}
                        required
                      />
                    </div>
                  </div>
                </div>

                {/* Precio y Asientos */}
                <div className="space-y-4">
                  <h3 className="text-xl font-semibold flex items-center gap-2">
                    <DollarSign className="w-5 h-5 text-primary" />
                    Precio y Capacidad
                  </h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="price">
                        Precio por asiento (ARS) <span className="text-destructive">*</span>
                      </Label>
                      <Input
                        id="price"
                        type="number"
                        min="0"
                        step="100"
                        value={formData.price_per_seat}
                        onChange={(e) => handleChange('price_per_seat', parseFloat(e.target.value))}
                        placeholder="Ej: 5000"
                        required
                      />
                      <p className="text-xs text-muted-foreground">
                        Precio que cada pasajero pagará por su asiento
                      </p>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="seats">
                        Total de asientos disponibles <span className="text-destructive">*</span>
                      </Label>
                      <Input
                        id="seats"
                        type="number"
                        min="1"
                        max="8"
                        value={formData.total_seats}
                        onChange={(e) => handleChange('total_seats', parseInt(e.target.value))}
                        required
                      />
                      <p className="text-xs text-muted-foreground">
                        Cantidad de asientos que ofreces (máximo 8)
                      </p>
                    </div>
                  </div>
                </div>

                {/* Vehículo */}
                <div className="space-y-4">
                  <h3 className="text-xl font-semibold flex items-center gap-2">
                    <Car className="w-5 h-5 text-primary" />
                    Información del Vehículo
                  </h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="car-brand">
                        Marca <span className="text-destructive">*</span>
                      </Label>
                      <Input
                        id="car-brand"
                        value={formData.car.brand}
                        onChange={(e) => handleNestedChange('car', 'brand', e.target.value)}
                        placeholder="Ej: Toyota"
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="car-model">
                        Modelo <span className="text-destructive">*</span>
                      </Label>
                      <Input
                        id="car-model"
                        value={formData.car.model}
                        onChange={(e) => handleNestedChange('car', 'model', e.target.value)}
                        placeholder="Ej: Corolla"
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="car-year">
                        Año <span className="text-destructive">*</span>
                      </Label>
                      <Input
                        id="car-year"
                        type="number"
                        min="1990"
                        max={new Date().getFullYear() + 1}
                        value={formData.car.year}
                        onChange={(e) => handleNestedChange('car', 'year', parseInt(e.target.value))}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="car-color">
                        Color <span className="text-destructive">*</span>
                      </Label>
                      <Input
                        id="car-color"
                        value={formData.car.color}
                        onChange={(e) => handleNestedChange('car', 'color', e.target.value)}
                        placeholder="Ej: Blanco"
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="car-plate">
                        Patente <span className="text-destructive">*</span>
                      </Label>
                      <Input
                        id="car-plate"
                        value={formData.car.plate}
                        onChange={(e) => handleNestedChange('car', 'plate', e.target.value.toUpperCase())}
                        placeholder="Ej: ABC123"
                        maxLength={7}
                        required
                      />
                    </div>
                  </div>
                </div>

                {/* Preferencias */}
                <div className="space-y-4">
                  <h3 className="text-xl font-semibold flex items-center gap-2">
                    <Users className="w-5 h-5 text-primary" />
                    Preferencias del Viaje
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Selecciona qué está permitido durante el viaje
                  </p>

                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="flex items-center space-x-2 p-3 border rounded-lg hover:bg-gray-50 transition-colors">
                      <input
                        type="checkbox"
                        id="pets"
                        checked={formData.preferences.pets_allowed}
                        onChange={(e) => handleNestedChange('preferences', 'pets_allowed', e.target.checked)}
                        className="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary"
                      />
                      <Label htmlFor="pets" className="cursor-pointer flex-1">
                        🐕 Permitir mascotas
                      </Label>
                    </div>

                    <div className="flex items-center space-x-2 p-3 border rounded-lg hover:bg-gray-50 transition-colors">
                      <input
                        type="checkbox"
                        id="smoking"
                        checked={formData.preferences.smoking_allowed}
                        onChange={(e) => handleNestedChange('preferences', 'smoking_allowed', e.target.checked)}
                        className="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary"
                      />
                      <Label htmlFor="smoking" className="cursor-pointer flex-1">
                        🚬 Permitir fumar
                      </Label>
                    </div>

                    <div className="flex items-center space-x-2 p-3 border rounded-lg hover:bg-gray-50 transition-colors">
                      <input
                        type="checkbox"
                        id="music"
                        checked={formData.preferences.music_allowed}
                        onChange={(e) => handleNestedChange('preferences', 'music_allowed', e.target.checked)}
                        className="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary"
                      />
                      <Label htmlFor="music" className="cursor-pointer flex-1">
                        🎵 Permitir música
                      </Label>
                    </div>
                  </div>
                </div>

                {/* Descripción */}
                <div className="space-y-2">
                  <Label htmlFor="description">Descripción del viaje (opcional)</Label>
                  <Textarea
                    id="description"
                    placeholder="Cuéntales a los pasajeros más detalles sobre el viaje, puntos de encuentro, paradas, etc."
                    value={formData.description}
                    onChange={(e) => handleChange('description', e.target.value)}
                    rows={4}
                  />
                  <p className="text-xs text-muted-foreground">
                    Incluye información útil como puntos de encuentro, paradas intermedias, forma de pago, etc.
                  </p>
                </div>
              </CardContent>

              <CardFooter className="flex gap-3 bg-gray-50">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate(-1)}
                  className="flex-1"
                  disabled={saving}
                >
                  Cancelar
                </Button>
                <Button
                  type="submit"
                  disabled={saving}
                  className="flex-1"
                  size="lg"
                >
                  {saving ? (
                    <>
                      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                      Publicando...
                    </>
                  ) : (
                    <>
                      <Car className="w-5 h-5 mr-2" />
                      Publicar Viaje
                    </>
                  )}
                </Button>
              </CardFooter>
            </Card>
          </form>
        </div>
      </div>
    </div>
  )
}
