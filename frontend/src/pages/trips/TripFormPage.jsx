import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { tripsService } from '@/services/api/trips.service';
import { Button, Card, Loading } from '@/components/common';

export const TripFormPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const isEditing = Boolean(id);

  const [loading, setLoading] = useState(isEditing);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState(null);

  const [formData, setFormData] = useState({
    // Origin
    origin_city: '',
    origin_province: '',
    origin_address: '',
    origin_lat: '',
    origin_lng: '',
    // Destination
    destination_city: '',
    destination_province: '',
    destination_address: '',
    destination_lat: '',
    destination_lng: '',
    // Schedule
    departure_datetime: '',
    estimated_arrival_datetime: '',
    // Pricing and capacity
    price_per_seat: '',
    total_seats: '4',
    // Car
    car_brand: '',
    car_model: '',
    car_year: new Date().getFullYear(),
    car_color: '',
    car_plate: '',
    // Preferences
    pets_allowed: false,
    smoking_allowed: false,
    music_allowed: true,
    // Description
    description: '',
  });

  const [errors, setErrors] = useState({});

  useEffect(() => {
    if (isEditing) {
      loadTrip();
    }
  }, [id]);

  const loadTrip = async () => {
    try {
      setLoading(true);
      const trip = await tripsService.getTripById(id);

      // Convert trip data to form format
      setFormData({
        origin_city: trip.origin.city,
        origin_province: trip.origin.province,
        origin_address: trip.origin.address,
        origin_lat: trip.origin.coordinates.lat,
        origin_lng: trip.origin.coordinates.lng,
        destination_city: trip.destination.city,
        destination_province: trip.destination.province,
        destination_address: trip.destination.address,
        destination_lat: trip.destination.coordinates.lat,
        destination_lng: trip.destination.coordinates.lng,
        departure_datetime: trip.departure_datetime.slice(0, 16),
        estimated_arrival_datetime: trip.estimated_arrival_datetime.slice(0, 16),
        price_per_seat: trip.price_per_seat,
        total_seats: trip.total_seats,
        car_brand: trip.car.brand,
        car_model: trip.car.model,
        car_year: trip.car.year,
        car_color: trip.car.color,
        car_plate: trip.car.plate,
        pets_allowed: trip.preferences.pets_allowed,
        smoking_allowed: trip.preferences.smoking_allowed,
        music_allowed: trip.preferences.music_allowed,
        description: trip.description || '',
      });
    } catch (err) {
      console.error('Error loading trip:', err);
      setError('Error al cargar el viaje');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData({
      ...formData,
      [name]: type === 'checkbox' ? checked : value,
    });
    // Clear error for this field
    if (errors[name]) {
      setErrors({ ...errors, [name]: null });
    }
  };

  const validateForm = () => {
    const newErrors = {};

    // Origin validation
    if (!formData.origin_city.trim()) newErrors.origin_city = 'Ciudad de origen requerida';
    if (!formData.origin_province.trim()) newErrors.origin_province = 'Provincia de origen requerida';
    if (!formData.origin_address.trim()) newErrors.origin_address = 'Dirección de origen requerida';
    if (!formData.origin_lat || isNaN(formData.origin_lat)) newErrors.origin_lat = 'Latitud inválida';
    if (!formData.origin_lng || isNaN(formData.origin_lng)) newErrors.origin_lng = 'Longitud inválida';

    // Destination validation
    if (!formData.destination_city.trim()) newErrors.destination_city = 'Ciudad de destino requerida';
    if (!formData.destination_province.trim()) newErrors.destination_province = 'Provincia de destino requerida';
    if (!formData.destination_address.trim()) newErrors.destination_address = 'Dirección de destino requerida';
    if (!formData.destination_lat || isNaN(formData.destination_lat)) newErrors.destination_lat = 'Latitud inválida';
    if (!formData.destination_lng || isNaN(formData.destination_lng)) newErrors.destination_lng = 'Longitud inválida';

    // Schedule validation
    if (!formData.departure_datetime) newErrors.departure_datetime = 'Fecha de salida requerida';
    if (!formData.estimated_arrival_datetime) newErrors.estimated_arrival_datetime = 'Fecha de llegada requerida';

    if (formData.departure_datetime && formData.estimated_arrival_datetime) {
      const departure = new Date(formData.departure_datetime);
      const arrival = new Date(formData.estimated_arrival_datetime);
      if (arrival <= departure) {
        newErrors.estimated_arrival_datetime = 'La llegada debe ser posterior a la salida';
      }
    }

    // Pricing validation
    if (!formData.price_per_seat || formData.price_per_seat <= 0) {
      newErrors.price_per_seat = 'Precio por asiento debe ser mayor a 0';
    }
    if (!formData.total_seats || formData.total_seats < 1 || formData.total_seats > 8) {
      newErrors.total_seats = 'Total de asientos debe estar entre 1 y 8';
    }

    // Car validation
    if (!formData.car_brand.trim()) newErrors.car_brand = 'Marca del auto requerida';
    if (!formData.car_model.trim()) newErrors.car_model = 'Modelo del auto requerido';
    if (!formData.car_year || formData.car_year < 1900 || formData.car_year > new Date().getFullYear() + 1) {
      newErrors.car_year = 'Año del auto inválido';
    }
    if (!formData.car_color.trim()) newErrors.car_color = 'Color del auto requerido';
    if (!formData.car_plate.trim()) newErrors.car_plate = 'Patente del auto requerida';

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!validateForm()) {
      setError('Por favor corrige los errores en el formulario');
      return;
    }

    try {
      setSubmitting(true);
      setError(null);

      // Convert form data to API format
      const tripData = {
        origin: {
          city: formData.origin_city,
          province: formData.origin_province,
          address: formData.origin_address,
          coordinates: {
            lat: parseFloat(formData.origin_lat),
            lng: parseFloat(formData.origin_lng),
          },
        },
        destination: {
          city: formData.destination_city,
          province: formData.destination_province,
          address: formData.destination_address,
          coordinates: {
            lat: parseFloat(formData.destination_lat),
            lng: parseFloat(formData.destination_lng),
          },
        },
        departure_datetime: new Date(formData.departure_datetime).toISOString(),
        estimated_arrival_datetime: new Date(formData.estimated_arrival_datetime).toISOString(),
        price_per_seat: parseFloat(formData.price_per_seat),
        total_seats: parseInt(formData.total_seats),
        car: {
          brand: formData.car_brand,
          model: formData.car_model,
          year: parseInt(formData.car_year),
          color: formData.car_color,
          plate: formData.car_plate,
        },
        preferences: {
          pets_allowed: formData.pets_allowed,
          smoking_allowed: formData.smoking_allowed,
          music_allowed: formData.music_allowed,
        },
        description: formData.description,
      };

      let trip;
      if (isEditing) {
        trip = await tripsService.updateTrip(id, tripData);
      } else {
        trip = await tripsService.createTrip(tripData);
      }

      navigate(`/trips/${trip.id}`);
    } catch (err) {
      console.error('Error saving trip:', err);
      setError(err.response?.data?.error || 'Error al guardar el viaje. Por favor intenta de nuevo.');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <Loading />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      {/* Header */}
      <div className="mb-6">
        <Link to="/trips" className="text-blue-600 hover:text-blue-700 flex items-center mb-4">
          <svg className="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Volver a viajes
        </Link>
        <h1 className="text-3xl font-bold text-gray-900">
          {isEditing ? 'Editar Viaje' : 'Crear Nuevo Viaje'}
        </h1>
        <p className="text-gray-600 mt-1">
          {isEditing ? 'Actualiza la información de tu viaje' : 'Completa los detalles de tu viaje'}
        </p>
      </div>

      {/* Error Message */}
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
          {error}
        </div>
      )}

      {/* Form */}
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Origin Section */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Origen</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Ciudad *
              </label>
              <input
                type="text"
                name="origin_city"
                value={formData.origin_city}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.origin_city ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: Córdoba"
              />
              {errors.origin_city && <p className="text-red-500 text-xs mt-1">{errors.origin_city}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Provincia *
              </label>
              <input
                type="text"
                name="origin_province"
                value={formData.origin_province}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.origin_province ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: Córdoba"
              />
              {errors.origin_province && <p className="text-red-500 text-xs mt-1">{errors.origin_province}</p>}
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Dirección *
              </label>
              <input
                type="text"
                name="origin_address"
                value={formData.origin_address}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.origin_address ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: Av. Haya de la Torre 2000"
              />
              {errors.origin_address && <p className="text-red-500 text-xs mt-1">{errors.origin_address}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Latitud *
              </label>
              <input
                type="number"
                step="any"
                name="origin_lat"
                value={formData.origin_lat}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.origin_lat ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="-31.4201"
              />
              {errors.origin_lat && <p className="text-red-500 text-xs mt-1">{errors.origin_lat}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Longitud *
              </label>
              <input
                type="number"
                step="any"
                name="origin_lng"
                value={formData.origin_lng}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.origin_lng ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="-64.1888"
              />
              {errors.origin_lng && <p className="text-red-500 text-xs mt-1">{errors.origin_lng}</p>}
            </div>
          </div>
        </Card>

        {/* Destination Section */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Destino</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Ciudad *
              </label>
              <input
                type="text"
                name="destination_city"
                value={formData.destination_city}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.destination_city ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: Buenos Aires"
              />
              {errors.destination_city && <p className="text-red-500 text-xs mt-1">{errors.destination_city}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Provincia *
              </label>
              <input
                type="text"
                name="destination_province"
                value={formData.destination_province}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.destination_province ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: Buenos Aires"
              />
              {errors.destination_province && <p className="text-red-500 text-xs mt-1">{errors.destination_province}</p>}
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Dirección *
              </label>
              <input
                type="text"
                name="destination_address"
                value={formData.destination_address}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.destination_address ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: Av. 9 de Julio 1000"
              />
              {errors.destination_address && <p className="text-red-500 text-xs mt-1">{errors.destination_address}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Latitud *
              </label>
              <input
                type="number"
                step="any"
                name="destination_lat"
                value={formData.destination_lat}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.destination_lat ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="-34.6037"
              />
              {errors.destination_lat && <p className="text-red-500 text-xs mt-1">{errors.destination_lat}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Longitud *
              </label>
              <input
                type="number"
                step="any"
                name="destination_lng"
                value={formData.destination_lng}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.destination_lng ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="-58.3816"
              />
              {errors.destination_lng && <p className="text-red-500 text-xs mt-1">{errors.destination_lng}</p>}
            </div>
          </div>
        </Card>

        {/* Schedule Section */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Horarios</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Fecha y hora de salida *
              </label>
              <input
                type="datetime-local"
                name="departure_datetime"
                value={formData.departure_datetime}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.departure_datetime ? 'border-red-500' : 'border-gray-300'
                }`}
              />
              {errors.departure_datetime && <p className="text-red-500 text-xs mt-1">{errors.departure_datetime}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Fecha y hora de llegada estimada *
              </label>
              <input
                type="datetime-local"
                name="estimated_arrival_datetime"
                value={formData.estimated_arrival_datetime}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.estimated_arrival_datetime ? 'border-red-500' : 'border-gray-300'
                }`}
              />
              {errors.estimated_arrival_datetime && <p className="text-red-500 text-xs mt-1">{errors.estimated_arrival_datetime}</p>}
            </div>
          </div>
        </Card>

        {/* Pricing and Capacity */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Precio y Capacidad</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Precio por asiento ($) *
              </label>
              <input
                type="number"
                step="0.01"
                min="0"
                name="price_per_seat"
                value={formData.price_per_seat}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.price_per_seat ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="1000.00"
              />
              {errors.price_per_seat && <p className="text-red-500 text-xs mt-1">{errors.price_per_seat}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Total de asientos disponibles (1-8) *
              </label>
              <input
                type="number"
                min="1"
                max="8"
                name="total_seats"
                value={formData.total_seats}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.total_seats ? 'border-red-500' : 'border-gray-300'
                }`}
              />
              {errors.total_seats && <p className="text-red-500 text-xs mt-1">{errors.total_seats}</p>}
            </div>
          </div>
        </Card>

        {/* Car Information */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Información del Vehículo</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Marca *
              </label>
              <input
                type="text"
                name="car_brand"
                value={formData.car_brand}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.car_brand ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: Toyota"
              />
              {errors.car_brand && <p className="text-red-500 text-xs mt-1">{errors.car_brand}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Modelo *
              </label>
              <input
                type="text"
                name="car_model"
                value={formData.car_model}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.car_model ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: Corolla"
              />
              {errors.car_model && <p className="text-red-500 text-xs mt-1">{errors.car_model}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Año *
              </label>
              <input
                type="number"
                min="1900"
                max={new Date().getFullYear() + 1}
                name="car_year"
                value={formData.car_year}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.car_year ? 'border-red-500' : 'border-gray-300'
                }`}
              />
              {errors.car_year && <p className="text-red-500 text-xs mt-1">{errors.car_year}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Color *
              </label>
              <input
                type="text"
                name="car_color"
                value={formData.car_color}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.car_color ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: Blanco"
              />
              {errors.car_color && <p className="text-red-500 text-xs mt-1">{errors.car_color}</p>}
            </div>
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Patente *
              </label>
              <input
                type="text"
                name="car_plate"
                value={formData.car_plate}
                onChange={handleChange}
                className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${
                  errors.car_plate ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Ej: ABC123"
              />
              {errors.car_plate && <p className="text-red-500 text-xs mt-1">{errors.car_plate}</p>}
            </div>
          </div>
        </Card>

        {/* Preferences */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Preferencias</h2>
          <div className="space-y-3">
            <label className="flex items-center space-x-3 cursor-pointer">
              <input
                type="checkbox"
                name="pets_allowed"
                checked={formData.pets_allowed}
                onChange={handleChange}
                className="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span className="text-gray-700">Permitir mascotas 🐕</span>
            </label>
            <label className="flex items-center space-x-3 cursor-pointer">
              <input
                type="checkbox"
                name="smoking_allowed"
                checked={formData.smoking_allowed}
                onChange={handleChange}
                className="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span className="text-gray-700">Permitir fumar 🚬</span>
            </label>
            <label className="flex items-center space-x-3 cursor-pointer">
              <input
                type="checkbox"
                name="music_allowed"
                checked={formData.music_allowed}
                onChange={handleChange}
                className="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span className="text-gray-700">Permitir música 🎵</span>
            </label>
          </div>
        </Card>

        {/* Description */}
        <Card>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Descripción (Opcional)</h2>
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            rows={4}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Agrega información adicional sobre tu viaje..."
          />
        </Card>

        {/* Submit Buttons */}
        <div className="flex gap-4">
          <Button
            type="button"
            variant="outline"
            size="lg"
            fullWidth
            onClick={() => navigate('/trips')}
            disabled={submitting}
          >
            Cancelar
          </Button>
          <Button
            type="submit"
            variant="primary"
            size="lg"
            fullWidth
            disabled={submitting}
          >
            {submitting ? 'Guardando...' : isEditing ? 'Actualizar Viaje' : 'Crear Viaje'}
          </Button>
        </div>
      </form>
    </div>
  );
};
