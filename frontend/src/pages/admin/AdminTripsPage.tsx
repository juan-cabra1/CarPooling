import { useEffect, useState } from 'react';
import { Trash2, AlertTriangle, MapPin, Calendar, User as UserIcon } from 'lucide-react';
import tripsService from '@/services/tripsService';
import searchService from '@/services/searchService';
import type { Trip } from '@/types/trip';
import type { SearchTrip } from '@/types/search';
import { Button } from '@/components/ui/button';
import StatusBadge from '@/components/admin/StatusBadge';
import { getErrorMessage } from '@/services/api';

interface DeleteTripModal {
  isOpen: boolean;
  trip: Trip | SearchTrip | null;
}

export default function AdminTripsPage() {
  const [trips, setTrips] = useState<SearchTrip[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [deleteModal, setDeleteModal] = useState<DeleteTripModal>({ isOpen: false, trip: null });
  const [deletingTripId, setDeletingTripId] = useState<string | null>(null);

  // Load trips on mount
  useEffect(() => {
    loadTrips();
  }, []);

  const loadTrips = async () => {
    try {
      setLoading(true);

      // Use Solr search API to get all trips
      const response = await searchService.searchTrips({
        page: 1,
        limit: 100, // Get all trips for admin
      });

      console.log('✅ Search API Response:', response);
      console.log('📊 Trips loaded:', response.trips?.length);
      console.log('📈 Total trips:', response.total);

      setTrips(response.trips || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('❌ Error loading trips:', error);
      // Fallback to direct trips API if Solr fails
      try {
        const fallbackData = await tripsService.getAllTrips();
        console.log('🔄 Using fallback data:', fallbackData);
        setTrips(fallbackData as any);
      } catch (fallbackError) {
        console.error('❌ Fallback also failed:', fallbackError);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteTrip = async () => {
    if (!deleteModal.trip) return;

    try {
      // SearchTrip has both 'id' (MongoDB) and 'trip_id' (original trip ID)
      // We need to use trip_id for the trips-api
      const tripId = 'trip_id' in deleteModal.trip ? deleteModal.trip.trip_id : deleteModal.trip.id;
      setDeletingTripId(deleteModal.trip.id);
      await tripsService.deleteTrip(tripId);
      setTrips(trips.filter((t) => t.id !== deleteModal.trip!.id));
      setDeleteModal({ isOpen: false, trip: null });
      alert('✅ Viaje eliminado exitosamente');
    } catch (error: any) {
      console.error('Error deleting trip:', error);
      const errorMsg = getErrorMessage(error);
      alert(`❌ Error al eliminar el viaje: ${errorMsg}`);
    } finally {
      setDeletingTripId(null);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('es-AR', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    });
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
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
          Gestión de Viajes
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Mostrando {trips.length} viaje{trips.length !== 1 ? 's' : ''} {total > 0 && `de ${total} total`}
        </p>
      </div>


      {/* Trips table */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Ruta
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Fecha
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Precio
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Asientos
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Estado
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Acciones
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200 dark:divide-gray-600">
              {trips.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
                    No se encontraron viajes
                  </td>
                </tr>
              ) : (
                trips.map((trip) => (
                  <tr key={trip.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                    <td className="px-6 py-4">
                      <div className="text-sm font-medium text-gray-900 dark:text-white">
                        {trip.origin.city} → {trip.destination.city}
                      </div>
                      <div className="text-sm text-gray-500 dark:text-gray-400">
                        {trip.origin.province}
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900 dark:text-white">
                      {new Date(trip.departure_datetime).toLocaleDateString('es-AR', {
                        day: 'numeric',
                        month: 'short',
                        year: 'numeric'
                      })}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900 dark:text-white">
                      ${trip.price_per_seat}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900 dark:text-white">
                      {trip.available_seats}/{trip.total_seats}
                    </td>
                    <td className="px-6 py-4">
                      <StatusBadge status={trip.status} />
                    </td>
                    <td className="px-6 py-4 text-right text-sm font-medium">
                      <Button
                        onClick={() => setDeleteModal({ isOpen: true, trip })}
                        disabled={deletingTripId === trip.id}
                        className="bg-red-600 hover:bg-red-700 text-white"
                        size="sm"
                      >
                        {deletingTripId === trip.id ? (
                          <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
                        ) : (
                          <>
                            <Trash2 className="h-4 w-4 mr-1" />
                            Eliminar
                          </>
                        )}
                      </Button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Delete Confirmation Modal */}
      {deleteModal.isOpen && deleteModal.trip && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full p-6">
            <div className="flex items-center gap-3 mb-4">
              <div className="flex-shrink-0 w-12 h-12 rounded-full bg-red-100 dark:bg-red-900/20 flex items-center justify-center">
                <AlertTriangle className="h-6 w-6 text-red-600 dark:text-red-400" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                Eliminar Viaje
              </h3>
            </div>

            <p className="text-gray-600 dark:text-gray-400 mb-4">
              ¿Estás seguro de que deseas eliminar este viaje? Esta acción no se puede deshacer.
            </p>

            <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4 mb-6 space-y-2">
              <div className="flex items-start gap-2">
                <MapPin className="h-5 w-5 text-gray-400 mt-0.5 flex-shrink-0" />
                <div className="text-sm">
                  <p className="font-medium text-gray-900 dark:text-white">
                    {deleteModal.trip.origin.city} → {deleteModal.trip.destination.city}
                  </p>
                  <p className="text-gray-500 dark:text-gray-400">
                    {deleteModal.trip.origin.province}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Calendar className="h-5 w-5 text-gray-400 flex-shrink-0" />
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  {formatDate(deleteModal.trip.departure_datetime)}
                </span>
              </div>
              <div className="flex items-center gap-2">
                <UserIcon className="h-5 w-5 text-gray-400 flex-shrink-0" />
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  Conductor: {'driver' in deleteModal.trip && deleteModal.trip.driver
                    ? `${deleteModal.trip.driver.name} (ID: ${deleteModal.trip.driver.id})`
                    : `ID: ${'driver_id' in deleteModal.trip ? deleteModal.trip.driver_id : 'N/A'}`}
                </span>
              </div>
            </div>

            <div className="flex items-center justify-end gap-3">
              <Button
                variant="outline"
                onClick={() => setDeleteModal({ isOpen: false, trip: null })}
                disabled={deletingTripId !== null}
              >
                Cancelar
              </Button>
              <Button
                onClick={handleDeleteTrip}
                disabled={deletingTripId !== null}
                className="bg-red-600 hover:bg-red-700 text-white"
              >
                {deletingTripId ? 'Eliminando...' : 'Eliminar Viaje'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
