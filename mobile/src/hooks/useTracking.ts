/**
 * useTracking - Hook personalizado para manejo de tracking de viajes
 *
 * Proporciona una interfaz simplificada al TrackingContext
 * con funcionalidades adicionales y helpers
 */

import { useCallback } from 'react';
import { Trip } from '../types';
import { useTrackingContext } from '../context/TrackingContext';

interface UseTrackingResult {
  // Estado
  activeTrip: Trip | null;
  isTracking: boolean;
  isDriver: boolean;
  driverLocation: {
    lat: number;
    lng: number;
  } | null;
  driverSpeed: number | null;
  driverHeading: number | null;
  eta: number | null;
  distanceRemaining: number | null;
  distanceTraveled: number | null;

  // Helpers
  hasActiveTrip: boolean;
  isTrackingAsDriver: boolean;
  isTrackingAsPassenger: boolean;
  formattedETA: string;
  formattedDistance: string;
  formattedSpeed: string;

  // Acciones
  startDriverTracking: (trip: Trip) => Promise<void>;
  startPassengerTracking: (trip: Trip) => Promise<void>;
  stopTracking: () => Promise<void>;
}

/**
 * Hook para manejar el tracking de viajes
 */
export const useTracking = (): UseTrackingResult => {
  const context = useTrackingContext();

  /**
   * Inicia tracking como conductor
   */
  const startDriverTracking = useCallback(
    async (trip: Trip) => {
      await context.startTracking(trip, true);
    },
    [context]
  );

  /**
   * Inicia tracking como pasajero
   */
  const startPassengerTracking = useCallback(
    async (trip: Trip) => {
      await context.startTracking(trip, false);
    },
    [context]
  );

  /**
   * Helpers computados
   */
  const hasActiveTrip = context.activeTrip !== null;
  const isTrackingAsDriver = context.isTracking && context.isDriver;
  const isTrackingAsPassenger = context.isTracking && !context.isDriver;

  /**
   * Formatea el ETA para mostrar
   * Ejemplo: "25 min", "1h 15min", "menos de 1 min"
   */
  const formattedETA = (() => {
    if (context.eta === null) return '--';
    if (context.eta < 1) return 'menos de 1 min';
    if (context.eta < 60) return `${Math.round(context.eta)} min`;

    const hours = Math.floor(context.eta / 60);
    const minutes = Math.round(context.eta % 60);

    if (minutes === 0) return `${hours}h`;
    return `${hours}h ${minutes}min`;
  })();

  /**
   * Formatea la distancia para mostrar
   * Ejemplo: "12.5 km", "850 m"
   */
  const formattedDistance = (() => {
    if (context.distanceRemaining === null) return '--';
    if (context.distanceRemaining < 1) {
      return `${Math.round(context.distanceRemaining * 1000)} m`;
    }
    return `${context.distanceRemaining.toFixed(1)} km`;
  })();

  /**
   * Formatea la velocidad para mostrar
   * Ejemplo: "65 km/h", "0 km/h"
   */
  const formattedSpeed = (() => {
    if (context.driverSpeed === null) return '--';
    return `${Math.round(context.driverSpeed)} km/h`;
  })();

  return {
    // Estado del contexto
    activeTrip: context.activeTrip,
    isTracking: context.isTracking,
    isDriver: context.isDriver,
    driverLocation: context.driverLocation,
    driverSpeed: context.driverSpeed,
    driverHeading: context.driverHeading,
    eta: context.eta,
    distanceRemaining: context.distanceRemaining,
    distanceTraveled: context.distanceTraveled,

    // Helpers
    hasActiveTrip,
    isTrackingAsDriver,
    isTrackingAsPassenger,
    formattedETA,
    formattedDistance,
    formattedSpeed,

    // Acciones
    startDriverTracking,
    startPassengerTracking,
    stopTracking: context.stopTracking,
  };
};

export default useTracking;
