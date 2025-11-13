import React from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Card } from '@/components/common';

export const ProfilePage = () => {
  const { user } = useAuth();

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('es-AR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  return (
    <div className="max-w-4xl mx-auto animate-slide-up">
      <div className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900 mb-2">Mi Perfil</h1>
        <p className="text-gray-600">Información personal y estadísticas</p>
      </div>

      <div className="grid md:grid-cols-3 gap-6">
        {/* Avatar y estadísticas principales */}
        <Card className="md:col-span-1 text-center">
          <div className="flex flex-col items-center space-y-4">
            <div className="w-32 h-32 bg-gradient-to-br from-primary-600 to-primary-700 rounded-full flex items-center justify-center text-white text-5xl font-bold shadow-xl">
              {user?.name?.charAt(0).toUpperCase() || 'U'}{user?.lastname?.charAt(0).toUpperCase() || ''}
            </div>
            <div>
              <h2 className="text-2xl font-bold text-gray-900">{user?.name} {user?.lastname}</h2>
              <p className="text-sm text-gray-500 capitalize">{user?.role || 'Usuario'}</p>
            </div>

            {/* Estadísticas */}
            <div className="w-full pt-4 border-t border-gray-200 space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600">⭐ Calificación Conductor</span>
                <span className="font-semibold text-primary-600">
                  {user?.avg_driver_rating ? user.avg_driver_rating.toFixed(1) : 'N/A'}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600">⭐ Calificación Pasajero</span>
                <span className="font-semibold text-primary-600">
                  {user?.avg_passenger_rating ? user.avg_passenger_rating.toFixed(1) : 'N/A'}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600">🚗 Viajes como conductor</span>
                <span className="font-semibold text-gray-900">{user?.total_trips_driver || 0}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600">🎒 Viajes como pasajero</span>
                <span className="font-semibold text-gray-900">{user?.total_trips_passenger || 0}</span>
              </div>
            </div>
          </div>
        </Card>

        {/* Información personal */}
        <Card className="md:col-span-2">
          <h3 className="text-xl font-bold text-gray-900 mb-6 pb-3 border-b border-gray-200">
            Información Personal
          </h3>

          <div className="grid md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-500 mb-1">Email</label>
              <p className="text-gray-900 font-medium flex items-center">
                <span className="mr-2">📧</span>
                {user?.email}
              </p>
              {user?.email_verified && (
                <span className="inline-block mt-1 text-xs bg-green-100 text-green-700 px-2 py-1 rounded">
                  ✓ Verificado
                </span>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500 mb-1">Teléfono</label>
              <p className="text-gray-900 font-medium flex items-center">
                <span className="mr-2">📱</span>
                {user?.phone || 'No proporcionado'}
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500 mb-1">Fecha de nacimiento</label>
              <p className="text-gray-900 font-medium flex items-center">
                <span className="mr-2">🎂</span>
                {user?.birthdate ? formatDate(user.birthdate) : 'No proporcionado'}
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500 mb-1">Sexo</label>
              <p className="text-gray-900 font-medium capitalize flex items-center">
                <span className="mr-2">👤</span>
                {user?.sex || 'No proporcionado'}
              </p>
            </div>

            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-500 mb-1">Dirección</label>
              <p className="text-gray-900 font-medium flex items-center">
                <span className="mr-2">🏠</span>
                {user?.street && user?.number
                  ? `${user.street} ${user.number}`
                  : 'No proporcionada'}
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500 mb-1">Miembro desde</label>
              <p className="text-gray-900 font-medium flex items-center">
                <span className="mr-2">📅</span>
                {user?.created_at ? formatDate(user.created_at) : 'N/A'}
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-500 mb-1">Última actualización</label>
              <p className="text-gray-900 font-medium flex items-center">
                <span className="mr-2">🔄</span>
                {user?.updated_at ? formatDate(user.updated_at) : 'N/A'}
              </p>
            </div>
          </div>

          <div className="mt-6 pt-6 border-t border-gray-200">
            <button className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors shadow-md hover:shadow-lg transform hover:scale-105">
              Editar Perfil
            </button>
          </div>
        </Card>
      </div>

      {/* Sección de logros (placeholder) */}
      <div className="mt-8">
        <Card>
          <h3 className="text-xl font-bold text-gray-900 mb-4">🏆 Logros</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="text-center p-4 bg-gray-50 rounded-lg">
              <div className="text-3xl mb-2">🌟</div>
              <p className="text-sm font-medium text-gray-700">Nuevo Usuario</p>
            </div>
            {user && user.total_trips_driver > 0 && (
              <div className="text-center p-4 bg-primary-50 rounded-lg">
                <div className="text-3xl mb-2">🚗</div>
                <p className="text-sm font-medium text-gray-700">Primer Viaje Conductor</p>
              </div>
            )}
            {user && user.total_trips_passenger > 0 && (
              <div className="text-center p-4 bg-blue-50 rounded-lg">
                <div className="text-3xl mb-2">🎒</div>
                <p className="text-sm font-medium text-gray-700">Primer Viaje Pasajero</p>
              </div>
            )}
            {user && user.email_verified && (
              <div className="text-center p-4 bg-green-50 rounded-lg">
                <div className="text-3xl mb-2">✅</div>
                <p className="text-sm font-medium text-gray-700">Email Verificado</p>
              </div>
            )}
          </div>
        </Card>
      </div>
    </div>
  );
};
