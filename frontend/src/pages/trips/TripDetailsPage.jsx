import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { tripsService } from '@/services/api/trips.service';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Card, Loading } from '@/components/common';

export const TripDetailsPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuth();
  const [trip, setTrip] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [cancelReason, setCancelReason] = useState('');
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    loadTrip();
  }, [id]);

  const loadTrip = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await tripsService.getTripById(id);
      setTrip(data);
    } catch (err) {
      console.error('Error loading trip:', err);
      setError('Error al cargar el viaje. Por favor intenta de nuevo.');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('¿Estás seguro de que deseas eliminar este viaje?')) {
      return;
    }

    try {
      setActionLoading(true);
      await tripsService.deleteTrip(id);
      navigate('/trips');
    } catch (err) {
      console.error('Error deleting trip:', err);
      alert('Error al eliminar el viaje. Por favor intenta de nuevo.');
    } finally {
      setActionLoading(false);
    }
  };

  const handleCancel = async () => {
    if (!cancelReason.trim()) {
      alert('Por favor ingresa una razón para la cancelación');
      return;
    }

    try {
      setActionLoading(true);
      await tripsService.cancelTrip(id, cancelReason);
      setShowCancelModal(false);
      loadTrip(); // Reload trip to show updated status
    } catch (err) {
      console.error('Error cancelling trip:', err);
      alert('Error al cancelar el viaje. Por favor intenta de nuevo.');
    } finally {
      setActionLoading(false);
    }
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-AR', {
      weekday: 'long',
      day: '2-digit',
      month: 'long',
      year: 'numeric',
    });
  };

  const formatTime = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString('es-AR', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getStatusBadge = (status) => {
    const statusConfig = {
      published: { label: 'Disponible', className: 'bg-green-100 text-green-800' },
      draft: { label: 'Borrador', className: 'bg-gray-100 text-gray-800' },
      full: { label: 'Completo', className: 'bg-yellow-100 text-yellow-800' },
      in_progress: { label: 'En curso', className: 'bg-blue-100 text-blue-800' },
      completed: { label: 'Completado', className: 'bg-gray-100 text-gray-600' },
      cancelled: { label: 'Cancelado', className: 'bg-red-100 text-red-800' },
    };

    const config = statusConfig[status] || statusConfig.published;
    return (
      <span className={`px-3 py-1 text-sm font-semibold rounded-full ${config.className}`}>
        {config.label}
      </span>
    );
  };

  const isOwner = isAuthenticated && user && trip && user.id === trip.driver_id;
  const canBook = isAuthenticated && !isOwner && trip?.status === 'published' && trip?.available_seats > 0;

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <Loading />
      </div>
    );
  }

  if (error || !trip) {
    return (
      <Card className="text-center py-12">
        <div className="text-red-400 text-5xl mb-4">⚠️</div>
        <h3 className="text-xl font-semibold text-gray-700 mb-2">
          {error || 'Viaje no encontrado'}
        </h3>
        <Link to="/trips">
          <Button variant="primary" size="md">
            Volver a viajes
          </Button>
        </Link>
      </Card>
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
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">Detalles del Viaje</h1>
            {getStatusBadge(trip.status)}
          </div>
          {isOwner && (
            <div className="flex gap-2">
              <Link to={`/trips/${id}/edit`}>
                <Button variant="outline" size="sm" disabled={trip.status === 'cancelled' || trip.status === 'completed'}>
                  Editar
                </Button>
              </Link>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowCancelModal(true)}
                disabled={trip.status === 'cancelled' || trip.status === 'completed' || actionLoading}
                className="text-red-600 hover:text-red-700"
              >
                Cancelar Viaje
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={handleDelete}
                disabled={actionLoading}
                className="text-red-600 hover:text-red-700"
              >
                Eliminar
              </Button>
            </div>
          )}
        </div>
      </div>

      {/* Main Content */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column - Trip Details */}
        <div className="lg:col-span-2 space-y-6">
          {/* Route Card */}
          <Card>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Ruta</h2>
            <div className="space-y-4">
              <div className="flex items-start">
                <div className="flex-shrink-0 w-10 h-10 rounded-full bg-green-100 flex items-center justify-center text-green-600 font-semibold">
                  A
                </div>
                <div className="ml-4 flex-1">
                  <div className="text-sm text-gray-500">Origen</div>
                  <div className="font-semibold text-gray-900 text-lg">{trip.origin.city}, {trip.origin.province}</div>
                  <div className="text-sm text-gray-600">{trip.origin.address}</div>
                </div>
              </div>

              <div className="ml-5 border-l-2 border-dashed border-gray-300 h-12"></div>

              <div className="flex items-start">
                <div className="flex-shrink-0 w-10 h-10 rounded-full bg-red-100 flex items-center justify-center text-red-600 font-semibold">
                  B
                </div>
                <div className="ml-4 flex-1">
                  <div className="text-sm text-gray-500">Destino</div>
                  <div className="font-semibold text-gray-900 text-lg">{trip.destination.city}, {trip.destination.province}</div>
                  <div className="text-sm text-gray-600">{trip.destination.address}</div>
                </div>
              </div>
            </div>
          </Card>

          {/* Schedule Card */}
          <Card>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Horarios</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="bg-blue-50 rounded-lg p-4">
                <div className="text-sm text-blue-600 mb-1">Salida</div>
                <div className="font-semibold text-gray-900">{formatDate(trip.departure_datetime)}</div>
                <div className="text-lg font-bold text-blue-600 mt-1">{formatTime(trip.departure_datetime)}</div>
              </div>
              <div className="bg-purple-50 rounded-lg p-4">
                <div className="text-sm text-purple-600 mb-1">Llegada estimada</div>
                <div className="font-semibold text-gray-900">{formatDate(trip.estimated_arrival_datetime)}</div>
                <div className="text-lg font-bold text-purple-600 mt-1">{formatTime(trip.estimated_arrival_datetime)}</div>
              </div>
            </div>
          </Card>

          {/* Car Details */}
          <Card>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Vehículo</h2>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div>
                <div className="text-sm text-gray-500">Marca</div>
                <div className="font-semibold text-gray-900">{trip.car.brand}</div>
              </div>
              <div>
                <div className="text-sm text-gray-500">Modelo</div>
                <div className="font-semibold text-gray-900">{trip.car.model}</div>
              </div>
              <div>
                <div className="text-sm text-gray-500">Año</div>
                <div className="font-semibold text-gray-900">{trip.car.year}</div>
              </div>
              <div>
                <div className="text-sm text-gray-500">Color</div>
                <div className="font-semibold text-gray-900">{trip.car.color}</div>
              </div>
              <div>
                <div className="text-sm text-gray-500">Patente</div>
                <div className="font-semibold text-gray-900">{trip.car.plate}</div>
              </div>
            </div>
          </Card>

          {/* Preferences */}
          <Card>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Preferencias</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="flex items-center">
                <span className="text-2xl mr-3">🐕</span>
                <div>
                  <div className="font-medium text-gray-900">Mascotas</div>
                  <div className={`text-sm ${trip.preferences.pets_allowed ? 'text-green-600' : 'text-red-600'}`}>
                    {trip.preferences.pets_allowed ? 'Permitidas' : 'No permitidas'}
                  </div>
                </div>
              </div>
              <div className="flex items-center">
                <span className="text-2xl mr-3">🚬</span>
                <div>
                  <div className="font-medium text-gray-900">Fumar</div>
                  <div className={`text-sm ${trip.preferences.smoking_allowed ? 'text-green-600' : 'text-red-600'}`}>
                    {trip.preferences.smoking_allowed ? 'Permitido' : 'No permitido'}
                  </div>
                </div>
              </div>
              <div className="flex items-center">
                <span className="text-2xl mr-3">🎵</span>
                <div>
                  <div className="font-medium text-gray-900">Música</div>
                  <div className={`text-sm ${trip.preferences.music_allowed ? 'text-green-600' : 'text-red-600'}`}>
                    {trip.preferences.music_allowed ? 'Permitida' : 'No permitida'}
                  </div>
                </div>
              </div>
            </div>
          </Card>

          {/* Description */}
          {trip.description && (
            <Card>
              <h2 className="text-xl font-semibold text-gray-900 mb-4">Descripción</h2>
              <p className="text-gray-700 whitespace-pre-wrap">{trip.description}</p>
            </Card>
          )}

          {/* Cancellation Info */}
          {trip.status === 'cancelled' && trip.cancellation_reason && (
            <Card className="bg-red-50 border-red-200">
              <h2 className="text-xl font-semibold text-red-900 mb-2">Viaje Cancelado</h2>
              <p className="text-red-700">{trip.cancellation_reason}</p>
              {trip.cancelled_at && (
                <p className="text-sm text-red-600 mt-2">
                  Cancelado el {formatDate(trip.cancelled_at)} a las {formatTime(trip.cancelled_at)}
                </p>
              )}
            </Card>
          )}
        </div>

        {/* Right Column - Booking Info */}
        <div className="space-y-6">
          {/* Price Card */}
          <Card>
            <div className="text-center mb-4">
              <div className="text-4xl font-bold text-blue-600">${trip.price_per_seat}</div>
              <div className="text-sm text-gray-500">por asiento</div>
            </div>

            <div className="space-y-3 mb-4">
              <div className="flex justify-between py-2 border-b border-gray-200">
                <span className="text-gray-600">Asientos totales</span>
                <span className="font-semibold text-gray-900">{trip.total_seats}</span>
              </div>
              <div className="flex justify-between py-2 border-b border-gray-200">
                <span className="text-gray-600">Asientos reservados</span>
                <span className="font-semibold text-gray-900">{trip.reserved_seats}</span>
              </div>
              <div className="flex justify-between py-2">
                <span className="text-gray-600">Asientos disponibles</span>
                <span className="font-semibold text-green-600">{trip.available_seats}</span>
              </div>
            </div>

            {canBook && (
              <Button variant="primary" size="lg" fullWidth>
                Reservar Asiento
              </Button>
            )}

            {!isAuthenticated && trip.status === 'published' && (
              <Link to="/login">
                <Button variant="primary" size="lg" fullWidth>
                  Iniciar sesión para reservar
                </Button>
              </Link>
            )}

            {isOwner && (
              <div className="bg-blue-50 text-blue-700 text-sm p-3 rounded-lg text-center">
                Este es tu viaje
              </div>
            )}
          </Card>

          {/* Trip Info */}
          <Card>
            <h3 className="font-semibold text-gray-900 mb-3">Información del viaje</h3>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">ID del viaje</span>
                <span className="font-mono text-xs text-gray-900">{trip.id.slice(0, 8)}...</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Creado</span>
                <span className="text-gray-900">{formatDate(trip.created_at)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Actualizado</span>
                <span className="text-gray-900">{formatDate(trip.updated_at)}</span>
              </div>
            </div>
          </Card>
        </div>
      </div>

      {/* Cancel Modal */}
      {showCancelModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h3 className="text-xl font-semibold text-gray-900 mb-4">Cancelar Viaje</h3>
            <p className="text-gray-600 mb-4">
              Por favor ingresa una razón para la cancelación del viaje.
            </p>
            <textarea
              value={cancelReason}
              onChange={(e) => setCancelReason(e.target.value)}
              placeholder="Razón de cancelación..."
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent mb-4"
              rows={4}
            />
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="md"
                fullWidth
                onClick={() => setShowCancelModal(false)}
                disabled={actionLoading}
              >
                Cancelar
              </Button>
              <Button
                variant="primary"
                size="md"
                fullWidth
                onClick={handleCancel}
                disabled={actionLoading || !cancelReason.trim()}
              >
                {actionLoading ? 'Cancelando...' : 'Confirmar Cancelación'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
