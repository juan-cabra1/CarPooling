import { useEffect, useState } from 'react';
import { Users, Car, Calendar, Star } from 'lucide-react';
import StatsCard from '@/components/admin/StatsCard';
import adminService from '@/services/adminService';
import type { Trip } from '@/types/trip';

interface DashboardStats {
  totalUsers: number;
  totalTrips: number;
  activeTrips: number;
  totalBookings: number;
  avgRating: number;
}

export default function AdminDashboardPage() {
  const [stats, setStats] = useState<DashboardStats>({
    totalUsers: 0,
    totalTrips: 0,
    activeTrips: 0,
    totalBookings: 0,
    avgRating: 0
  });
  const [recentTrips, setRecentTrips] = useState<Trip[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      setLoading(true);

      // Obtener estadísticas usando el servicio de admin
      const [usersResponse, tripsResponse, bookingsResponse] = await Promise.all([
        adminService.getAllUsers(1, 1).catch(() => ({ users: [], pagination: { total: 0, page: 1, limit: 1, totalPages: 0 } })),
        adminService.getAllTrips(1, 100).catch(() => ({ trips: [], total: 0 })),
        adminService.getAllBookings(1, 1).catch(() => ({ bookings: [], pagination: { total: 0, page: 1, limit: 1, totalPages: 0 } })),
      ]);

      const activeTrips = tripsResponse.trips.filter(
        (t) => t.status === 'published' || t.status === 'in_progress'
      );

      const calculatedStats: DashboardStats = {
        totalUsers: usersResponse.pagination.total,
        totalTrips: tripsResponse.total,
        activeTrips: activeTrips.length,
        totalBookings: bookingsResponse.pagination.total,
        avgRating: 0 // TODO: Calcular promedio de ratings
      };

      setStats(calculatedStats);

      // Obtener viajes recientes (últimos 5)
      const sortedTrips = [...tripsResponse.trips].sort(
        (a, b) =>
          new Date(b.departure_datetime).getTime() -
          new Date(a.departure_datetime).getTime()
      );
      setRecentTrips(sortedTrips.slice(0, 5));
    } catch (error) {
      console.error('Error loading dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div>
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
          Dashboard
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Vista general de la plataforma
        </p>
      </div>

      {/* Stats grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatsCard
          title="Total Usuarios"
          value={stats.totalUsers || '-'}
          description="Usuarios registrados"
          icon={Users}
        />

        <StatsCard
          title="Total Viajes"
          value={stats.totalTrips}
          description="Viajes publicados"
          icon={Car}
        />

        <StatsCard
          title="Viajes Activos"
          value={stats.activeTrips}
          description="En curso o publicados"
          icon={Calendar}
        />

        <StatsCard
          title="Rating Promedio"
          value={stats.avgRating > 0 ? stats.avgRating.toFixed(1) : '-'}
          description="Calificación general"
          icon={Star}
        />
      </div>

      {/* Recent trips */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 p-6">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          Viajes Recientes
        </h3>

        {recentTrips.length === 0 ? (
          <p className="text-gray-500 dark:text-gray-400 text-center py-8">
            No hay viajes registrados
          </p>
        ) : (
          <div className="space-y-3">
            {recentTrips.map((trip) => (
              <div
                key={trip.id}
                className="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg"
              >
                <div className="flex-1">
                  <p className="font-medium text-gray-900 dark:text-white">
                    {trip.origin.city} → {trip.destination.city}
                  </p>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {new Date(trip.departure_datetime).toLocaleDateString('es-AR', {
                      day: 'numeric',
                      month: 'long',
                      year: 'numeric'
                    })}
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-sm font-medium text-gray-900 dark:text-white">
                    ${trip.price_per_seat}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {trip.available_seats}/{trip.total_seats} asientos
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
