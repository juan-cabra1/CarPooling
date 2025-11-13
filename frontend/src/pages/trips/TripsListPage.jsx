import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { tripsService } from '@/services/api/trips.service';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Card, Loading } from '@/components/common';

export const TripsListPage = () => {
  const { isAuthenticated } = useAuth();
  const [trips, setTrips] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    origin_city: '',
    destination_city: '',
    status: 'published',
    page: 1,
    limit: 12,
  });
  const [pagination, setPagination] = useState({
    total: 0,
    page: 1,
    limit: 12,
  });

  useEffect(() => {
    loadTrips();
  }, [filters.page]);

  const loadTrips = async () => {
    try {
      setLoading(true);
      setError(null);
      const activeFilters = Object.fromEntries(
        Object.entries(filters).filter(([_, value]) => value !== '')
      );
      const data = await tripsService.getTrips(activeFilters);
      setTrips(data.trips || []);
      setPagination({
        total: data.total,
        page: data.page,
        limit: data.limit,
      });
    } catch (err) {
      console.error('Error loading trips:', err);
      setError('Error al cargar los viajes. Por favor intenta de nuevo.');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    setFilters({ ...filters, page: 1 });
    loadTrips();
  };

  const handleFilterChange = (key, value) => {
    setFilters({ ...filters, [key]: value });
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-AR', {
      day: '2-digit',
      month: 'short',
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
      <span className={`px-2 py-1 text-xs font-semibold rounded-full ${config.className}`}>
        {config.label}
      </span>
    );
  };

  if (loading && trips.length === 0) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <Loading />
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Buscar Viajes</h1>
          <p className="text-gray-600 mt-1">
            Encuentra el viaje perfecto para tu destino
          </p>
        </div>
        {isAuthenticated && (
          <Link to="/trips/new">
            <Button variant="primary" size="md">
              + Crear Viaje
            </Button>
          </Link>
        )}
      </div>

      {/* Search and Filters */}
      <Card className="mb-6">
        <form onSubmit={handleSearch} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Origen
              </label>
              <input
                type="text"
                placeholder="Ciudad de origen"
                value={filters.origin_city}
                onChange={(e) => handleFilterChange('origin_city', e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Destino
              </label>
              <input
                type="text"
                placeholder="Ciudad de destino"
                value={filters.destination_city}
                onChange={(e) => handleFilterChange('destination_city', e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Estado
              </label>
              <select
                value={filters.status}
                onChange={(e) => handleFilterChange('status', e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">Todos</option>
                <option value="published">Disponibles</option>
                <option value="full">Completos</option>
                <option value="in_progress">En curso</option>
                <option value="completed">Completados</option>
              </select>
            </div>
          </div>
          <div className="flex justify-end">
            <Button type="submit" variant="primary" size="md">
              Buscar
            </Button>
          </div>
        </form>
      </Card>

      {/* Error Message */}
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
          {error}
        </div>
      )}

      {/* Results Count */}
      {!loading && (
        <div className="mb-4 text-sm text-gray-600">
          Mostrando {trips.length} de {pagination.total} viajes
        </div>
      )}

      {/* Trips Grid */}
      {trips.length === 0 && !loading ? (
        <Card className="text-center py-12">
          <div className="text-gray-400 text-5xl mb-4">🚗</div>
          <h3 className="text-xl font-semibold text-gray-700 mb-2">
            No se encontraron viajes
          </h3>
          <p className="text-gray-500 mb-4">
            Intenta ajustar tus filtros de búsqueda
          </p>
          {isAuthenticated && (
            <Link to="/trips/new">
              <Button variant="primary" size="md">
                Crear un viaje
              </Button>
            </Link>
          )}
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {trips.map((trip) => (
            <Card
              key={trip.id}
              className="hover:shadow-xl transition-shadow duration-300 cursor-pointer"
            >
              <Link to={`/trips/${trip.id}`} className="block">
                {/* Status Badge */}
                <div className="flex justify-between items-start mb-4">
                  {getStatusBadge(trip.status)}
                  <div className="text-right">
                    <div className="text-2xl font-bold text-blue-600">
                      ${trip.price_per_seat}
                    </div>
                    <div className="text-xs text-gray-500">por asiento</div>
                  </div>
                </div>

                {/* Route */}
                <div className="mb-4">
                  <div className="flex items-start mb-2">
                    <div className="text-green-500 mr-2 mt-1">●</div>
                    <div className="flex-1">
                      <div className="font-semibold text-gray-900">
                        {trip.origin.city}
                      </div>
                      <div className="text-xs text-gray-500">
                        {trip.origin.province}
                      </div>
                    </div>
                  </div>
                  <div className="border-l-2 border-gray-300 ml-2 h-6"></div>
                  <div className="flex items-start">
                    <div className="text-red-500 mr-2 mt-1">●</div>
                    <div className="flex-1">
                      <div className="font-semibold text-gray-900">
                        {trip.destination.city}
                      </div>
                      <div className="text-xs text-gray-500">
                        {trip.destination.province}
                      </div>
                    </div>
                  </div>
                </div>

                {/* Date and Time */}
                <div className="mb-4 bg-gray-50 rounded-lg p-3">
                  <div className="flex items-center text-sm text-gray-600 mb-1">
                    <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    {formatDate(trip.departure_datetime)}
                  </div>
                  <div className="flex items-center text-sm text-gray-600">
                    <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {formatTime(trip.departure_datetime)}
                  </div>
                </div>

                {/* Car Info */}
                <div className="mb-4 text-sm text-gray-600">
                  <div className="flex items-center">
                    <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {trip.car.brand} {trip.car.model} ({trip.car.year})
                  </div>
                </div>

                {/* Available Seats */}
                <div className="flex items-center justify-between pt-4 border-t border-gray-200">
                  <div className="flex items-center text-sm">
                    <svg className="w-5 h-5 mr-1 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                    <span className="font-medium text-gray-700">
                      {trip.available_seats} asientos disponibles
                    </span>
                  </div>
                </div>

                {/* Preferences Icons */}
                {(trip.preferences.pets_allowed || trip.preferences.smoking_allowed || trip.preferences.music_allowed) && (
                  <div className="flex gap-2 mt-3 pt-3 border-t border-gray-100">
                    {trip.preferences.pets_allowed && (
                      <span className="text-xs bg-blue-50 text-blue-700 px-2 py-1 rounded" title="Mascotas permitidas">
                        🐕
                      </span>
                    )}
                    {trip.preferences.smoking_allowed && (
                      <span className="text-xs bg-orange-50 text-orange-700 px-2 py-1 rounded" title="Se permite fumar">
                        🚬
                      </span>
                    )}
                    {trip.preferences.music_allowed && (
                      <span className="text-xs bg-purple-50 text-purple-700 px-2 py-1 rounded" title="Música permitida">
                        🎵
                      </span>
                    )}
                  </div>
                )}
              </Link>
            </Card>
          ))}
        </div>
      )}

      {/* Pagination */}
      {pagination.total > pagination.limit && (
        <div className="mt-8 flex justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setFilters({ ...filters, page: filters.page - 1 })}
            disabled={filters.page === 1}
          >
            Anterior
          </Button>
          <div className="flex items-center px-4 py-2 text-sm text-gray-700">
            Página {pagination.page} de {Math.ceil(pagination.total / pagination.limit)}
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setFilters({ ...filters, page: filters.page + 1 })}
            disabled={filters.page >= Math.ceil(pagination.total / pagination.limit)}
          >
            Siguiente
          </Button>
        </div>
      )}
    </div>
  );
};
