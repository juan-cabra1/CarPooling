import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
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
import { useAuth } from '@/context/AuthContext'
import type { Trip, UpdateTripRequest } from '@/types'

export default function EditTripPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()
  const [trip, setTrip] = useState<Trip | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  const [formData, setFormData] = useState<UpdateTripRequest>({})

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

      // Verificar que el usuario sea el propietario
      if (user && data.driver_id !== user.id) {
        setError('No tienes permiso para editar este viaje')
        return
      }

      // Verificar que no tenga reservas
      if (data.reserved_seats > 0) {
        setError('No puedes editar un viaje que tiene reservas activas')
        return
      }

      setTrip(data)

      // Inicializar el formulario con los datos del viaje
      setFormData({
        origin: data.origin,
        destination: data.destination,
        departure_datetime: formatDateForInput(data.departure_datetime),
        estimated_arrival_datetime: formatDateForInput(data.estimated_arrival_datetime),
        price_per_seat: data.price_per_seat,
        total_seats: data.total_seats,
        car: data.car,
        preferences: data.preferences,
        description: data.description,
      })
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const formatDateForInput = (dateString: string) => {
    const date = new Date(dateString)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    return `${year}-${month}-${day}T${hours}:${minutes}`
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!id || !trip) return

    try {
      setSaving(true)
      setError('')

      // Convertir fechas a ISO 8601
      const updateData: UpdateTripRequest = {
        ...formData,
        departure_datetime: formData.departure_datetime
          ? new Date(formData.departure_datetime).toISOString()
          : undefined,
        estimated_arrival_datetime: formData.estimated_arrival_datetime
          ? new Date(formData.estimated_arrival_datetime).toISOString()
          : undefined,
      }

      await tripsService.updateTrip(id, updateData)
      navigate(`/trips/${id}`, {
        state: { message: 'Viaje actualizado exitosamente' }
      })
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setSaving(false)
    }
  }

  const handleChange = (field: string, value: any) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }))
  }

  const handleNestedChange = (parent: string, field: string, value: any) => {
    setFormData(prev => ({
      ...prev,
      [parent]: {
        ...(prev[parent as keyof UpdateTripRequest] as any),
        [field]: value,
      },
    }))
  }

  if (loading) {
    return (
      <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <div className="flex items-center justify-center py-16">
              <div className="text-center">
                <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4" />
                <p className="text-gray-600">Cargando viaje...</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error && !trip) {
    return (
      <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <Card>
              <CardContent className="pt-6">
                <div className="text-center py-8">
                  <AlertCircle className="w-16 h-16 text-destructive mx-auto mb-4" />
                  <h2 className="text-2xl font-bold text-gray-900 mb-2">No se puede editar el viaje</h2>
                  <p className="text-gray-600 mb-6">{error}</p>
                  <Button onClick={() => navigate(-1)}>
                    <ArrowLeft className="w-4 h-4 mr-2" />
                    Volver
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    )
  }

  if (!trip || !formData.origin) return null

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
                <CardTitle className="text-3xl">Editar Viaje</CardTitle>
                <CardDescription>
                  Actualiza los detalles de tu viaje
                </CardDescription>
              </CardHeader>

              <CardContent className="space-y-8">
                {/* Error Message */}
                {error && (
                  <div className="p-4 rounded-lg bg-destructive/10 border border-destructive/20 flex items-start gap-3">
                    <AlertCircle className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
                    <p className="text-sm text-destructive">{error}</p>
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
                        <Label htmlFor="origin-city">Ciudad</Label>
                        <Input
                          id="origin-city"
                          value={formData.origin?.city || ''}
                          onChange={(e) => handleNestedChange('origin', 'city', e.target.value)}
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="origin-province">Provincia</Label>
                        <Input
                          id="origin-province"
                          value={formData.origin?.province || ''}
                          onChange={(e) => handleNestedChange('origin', 'province', e.target.value)}
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="origin-address">Dirección</Label>
                        <Input
                          id="origin-address"
                          value={formData.origin?.address || ''}
                          onChange={(e) => handleNestedChange('origin', 'address', e.target.value)}
                          required
                        />
                      </div>
                    </div>

                    {/* Destino */}
                    <div className="space-y-4 p-4 border rounded-lg bg-white">
                      <h4 className="font-semibold text-secondary">Destino</h4>

                      <div className="space-y-2">
                        <Label htmlFor="destination-city">Ciudad</Label>
                        <Input
                          id="destination-city"
                          value={formData.destination?.city || ''}
                          onChange={(e) => handleNestedChange('destination', 'city', e.target.value)}
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="destination-province">Provincia</Label>
                        <Input
                          id="destination-province"
                          value={formData.destination?.province || ''}
                          onChange={(e) => handleNestedChange('destination', 'province', e.target.value)}
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="destination-address">Dirección</Label>
                        <Input
                          id="destination-address"
                          value={formData.destination?.address || ''}
                          onChange={(e) => handleNestedChange('destination', 'address', e.target.value)}
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
                      <Label htmlFor="departure">Fecha y hora de salida</Label>
                      <Input
                        id="departure"
                        type="datetime-local"
                        value={formData.departure_datetime || ''}
                        onChange={(e) => handleChange('departure_datetime', e.target.value)}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="arrival">Fecha y hora estimada de llegada</Label>
                      <Input
                        id="arrival"
                        type="datetime-local"
                        value={formData.estimated_arrival_datetime || ''}
                        onChange={(e) => handleChange('estimated_arrival_datetime', e.target.value)}
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
                      <Label htmlFor="price">Precio por asiento (ARS)</Label>
                      <Input
                        id="price"
                        type="number"
                        min="0"
                        step="100"
                        value={formData.price_per_seat || ''}
                        onChange={(e) => handleChange('price_per_seat', parseFloat(e.target.value))}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="seats">Total de asientos disponibles</Label>
                      <Input
                        id="seats"
                        type="number"
                        min="1"
                        max="8"
                        value={formData.total_seats || ''}
                        onChange={(e) => handleChange('total_seats', parseInt(e.target.value))}
                        required
                      />
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
                      <Label htmlFor="car-brand">Marca</Label>
                      <Input
                        id="car-brand"
                        value={formData.car?.brand || ''}
                        onChange={(e) => handleNestedChange('car', 'brand', e.target.value)}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="car-model">Modelo</Label>
                      <Input
                        id="car-model"
                        value={formData.car?.model || ''}
                        onChange={(e) => handleNestedChange('car', 'model', e.target.value)}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="car-year">Año</Label>
                      <Input
                        id="car-year"
                        type="number"
                        min="1990"
                        max={new Date().getFullYear() + 1}
                        value={formData.car?.year || ''}
                        onChange={(e) => handleNestedChange('car', 'year', parseInt(e.target.value))}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="car-color">Color</Label>
                      <Input
                        id="car-color"
                        value={formData.car?.color || ''}
                        onChange={(e) => handleNestedChange('car', 'color', e.target.value)}
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="car-plate">Patente</Label>
                      <Input
                        id="car-plate"
                        value={formData.car?.plate || ''}
                        onChange={(e) => handleNestedChange('car', 'plate', e.target.value.toUpperCase())}
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

                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        id="pets"
                        checked={formData.preferences?.pets_allowed || false}
                        onChange={(e) => handleNestedChange('preferences', 'pets_allowed', e.target.checked)}
                        className="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary"
                      />
                      <Label htmlFor="pets" className="cursor-pointer">
                        🐕 Permitir mascotas
                      </Label>
                    </div>

                    <div className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        id="smoking"
                        checked={formData.preferences?.smoking_allowed || false}
                        onChange={(e) => handleNestedChange('preferences', 'smoking_allowed', e.target.checked)}
                        className="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary"
                      />
                      <Label htmlFor="smoking" className="cursor-pointer">
                        🚬 Permitir fumar
                      </Label>
                    </div>

                    <div className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        id="music"
                        checked={formData.preferences?.music_allowed || false}
                        onChange={(e) => handleNestedChange('preferences', 'music_allowed', e.target.checked)}
                        className="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary"
                      />
                      <Label htmlFor="music" className="cursor-pointer">
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
                    placeholder="Cuéntales a los pasajeros más detalles sobre el viaje..."
                    value={formData.description || ''}
                    onChange={(e) => handleChange('description', e.target.value)}
                    rows={4}
                  />
                </div>
              </CardContent>

              <CardFooter className="flex gap-3">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate(-1)}
                  className="flex-1"
                >
                  Cancelar
                </Button>
                <Button
                  type="submit"
                  disabled={saving}
                  className="flex-1"
                >
                  {saving ? (
                    <>
                      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                      Guardando...
                    </>
                  ) : (
                    'Guardar Cambios'
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
