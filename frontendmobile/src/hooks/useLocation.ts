/**
 * useLocation - Hook personalizado para manejo de ubicación GPS
 *
 * Proporciona una interfaz sencilla para:
 * - Obtener ubicación actual
 * - Verificar y solicitar permisos
 * - Estado de carga y errores
 */

import { useState, useCallback, useEffect } from 'react';
import { Coordinates } from '../types';
import locationService, { LocationPermissionStatus } from '../services/locationService';

interface UseLocationResult {
  // Estado
  location: Coordinates | null;
  isLoading: boolean;
  error: string | null;
  permissions: LocationPermissionStatus | null;
  isLocationEnabled: boolean;

  // Acciones
  getCurrentLocation: () => Promise<void>;
  requestPermissions: () => Promise<boolean>;
  checkLocationEnabled: () => Promise<void>;
  clearError: () => void;
}

/**
 * Hook para manejar la ubicación del usuario
 */
export const useLocation = (): UseLocationResult => {
  const [location, setLocation] = useState<Coordinates | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [permissions, setPermissions] = useState<LocationPermissionStatus | null>(null);
  const [isLocationEnabled, setIsLocationEnabled] = useState(true);

  /**
   * Obtiene la ubicación actual del usuario
   */
  const getCurrentLocation = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);

      const currentLocation = await locationService.getCurrentLocation();
      setLocation(currentLocation);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Error al obtener ubicación';
      setError(message);
      console.error('Error getting location:', err);
    } finally {
      setIsLoading(false);
    }
  }, []);

  /**
   * Solicita permisos de ubicación al usuario
   * @returns true si los permisos fueron otorgados
   */
  const requestPermissions = useCallback(async (): Promise<boolean> => {
    try {
      setIsLoading(true);
      setError(null);

      const permissionStatus = await locationService.requestPermissions();
      setPermissions(permissionStatus);

      if (!permissionStatus.granted) {
        setError('Permisos de ubicación denegados');
        return false;
      }

      return true;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Error al solicitar permisos';
      setError(message);
      console.error('Error requesting permissions:', err);
      return false;
    } finally {
      setIsLoading(false);
    }
  }, []);

  /**
   * Verifica si los servicios de ubicación están habilitados
   */
  const checkLocationEnabled = useCallback(async () => {
    try {
      const enabled = await locationService.isLocationEnabled();
      setIsLocationEnabled(enabled);

      if (!enabled) {
        setError('GPS deshabilitado. Por favor actívalo en la configuración del dispositivo.');
      }
    } catch (err) {
      console.error('Error checking if location is enabled:', err);
      setIsLocationEnabled(false);
    }
  }, []);

  /**
   * Limpia el error actual
   */
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  /**
   * Verificar permisos y estado del GPS al montar
   */
  useEffect(() => {
    const checkInitialState = async () => {
      const permissionStatus = await locationService.checkPermissions();
      setPermissions(permissionStatus);

      await checkLocationEnabled();
    };

    checkInitialState();
  }, [checkLocationEnabled]);

  return {
    location,
    isLoading,
    error,
    permissions,
    isLocationEnabled,
    getCurrentLocation,
    requestPermissions,
    checkLocationEnabled,
    clearError,
  };
};

export default useLocation;
