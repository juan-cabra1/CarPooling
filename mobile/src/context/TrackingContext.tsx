/**
 * TrackingContext - Contexto global para manejo de tracking de viajes
 *
 * Proporciona:
 * - Estado global del tracking activo
 * - Ubicación del conductor en tiempo real
 * - ETA y distancia restante
 * - Funciones para iniciar/detener tracking
 */

import React, { createContext, useContext, useState, useEffect, useRef, ReactNode, useCallback } from 'react';
import { Trip, Coordinates } from '../types';
import { locationService, LocationUpdate, DirectionsResponse, NavigationStep } from '../services/locationService';
import { socketService, LocationUpdateEvent } from '../services/socketService';
import { useAuth } from './AuthContext';
import AsyncStorage from '@react-native-async-storage/async-storage';

// Define a type for the result of startTracking
export interface StartTrackingResult {
  success: boolean;
  message: string;
}

interface TrackingState {
  // Estado del tracking
  activeTrip: Trip | null;
  isTracking: boolean;
  isDriver: boolean;

  // Ubicación del conductor
  driverLocation: Coordinates | null;
  driverSpeed: number | null; // km/h
  driverHeading: number | null; // degrees

  // Métricas del viaje
  eta: number | null; // minutos estimados
  distanceRemaining: number | null; // kilómetros
  distanceTraveled: number | null; // kilómetros
  routeCoordinates: Coordinates[] | null; // Coordenadas de la ruta (origen → destino)
  routeToOrigin: Coordinates[] | null; // Coordenadas de la ruta (ubicación actual → origen)

  // NUEVO: Estado de navegación
  navigationSteps: NavigationStep[];
  activeRouteType: 'toOrigin' | 'toDestination' | null;
  currentStepIndex: number;
  distanceToNextStep: number | null; // metros
  nextStepPreview: NavigationStep | null;
  isApproachingTurn: boolean; // true si < 100m al giro

  // Estado del flujo conductor
  isNearOrigin: boolean; // true cuando está a <300m del punto de salida

  // Acciones
  startTracking: (trip: Trip, asDriver: boolean) => Promise<StartTrackingResult>;
  stopTracking: () => Promise<void>;
  updateDriverLocation: (location: Coordinates) => void;
  confirmArrivalAtOrigin: () => void; // conductor confirma que llegó al punto de salida
}

const TrackingContext = createContext<TrackingState | undefined>(undefined);

const STORAGE_KEY = '@carpooling:active_trip';

interface TrackingProviderProps {
  children: ReactNode;
}

export const TrackingProvider: React.FC<TrackingProviderProps> = ({ children }) => {
  const { token, user } = useAuth();
  const [activeTrip, setActiveTrip] = useState<Trip | null>(null);
  const [isTracking, setIsTracking] = useState(false);
  const [isDriver, setIsDriver] = useState(false);
  const [driverLocation, setDriverLocation] = useState<Coordinates | null>(null);
  const [driverSpeed, setDriverSpeed] = useState<number | null>(null);
  const [driverHeading, setDriverHeading] = useState<number | null>(null);
  const [eta, setEta] = useState<number | null>(null);
  const [distanceRemaining, setDistanceRemaining] = useState<number | null>(null);
  const [distanceTraveled, setDistanceTraveled] = useState<number | null>(null);
  const [routeCoordinates, setRouteCoordinates] = useState<Coordinates[] | null>(null);
  const [routeToOrigin, setRouteToOrigin] = useState<Coordinates[] | null>(null);

  // NUEVO: Estado de navegación
  const [navigationSteps, setNavigationSteps] = useState<NavigationStep[]>([]);
  const [activeRouteType, setActiveRouteType] = useState<'toOrigin' | 'toDestination' | null>(null);
  const [currentStepIndex, setCurrentStepIndex] = useState(0);
  const [distanceToNextStep, setDistanceToNextStep] = useState<number | null>(null);
  const [nextStepPreview, setNextStepPreview] = useState<NavigationStep | null>(null);
  const [isApproachingTurn, setIsApproachingTurn] = useState(false);

  // Estado del flujo conductor
  const [isNearOrigin, setIsNearOrigin] = useState(false);

  // Cache de pasos por fase
  const [routeToOriginSteps, setRouteToOriginSteps] = useState<NavigationStep[]>([]);
  const [routeToDestinationSteps, setRouteToDestinationSteps] = useState<NavigationStep[]>([]);

  // Cooldown para recálculo (evita flood de API calls)
  const lastRecalcRef = useRef<number>(0);
  const RECALC_COOLDOWN_MS = 10000; // mínimo 10s entre recálculos

  // Refs para acceder a rutas actuales dentro del callback sin dependencias stale
  const routeToOriginRef = useRef<Coordinates[] | null>(null);
  const routeCoordinatesRef = useRef<Coordinates[] | null>(null);
  useEffect(() => { routeToOriginRef.current = routeToOrigin; }, [routeToOrigin]);
  useEffect(() => { routeCoordinatesRef.current = routeCoordinates; }, [routeCoordinates]);

  /**
   * El conductor confirma que llegó al punto de salida.
   * Cambia la fase de navegación de "toOrigin" a "toDestination".
   */
  const confirmArrivalAtOrigin = useCallback(() => {
    console.log('🎯 Conductor confirmó llegada al punto de salida. Cambiando a ruta destino.');
    setActiveRouteType('toDestination');
    setIsNearOrigin(false);
    setCurrentStepIndex(0);
    lastRecalcRef.current = 0; // forzar recálculo inmediato al cambiar de fase
    if (routeToDestinationSteps.length > 0) {
      setNavigationSteps(routeToDestinationSteps);
    }
  }, [routeToDestinationSteps]);

  /**
   * Calculate which step the driver is currently on
   * Uses Haversine distance to find closest step
   */
  const calculateCurrentStep = useCallback((
    currentLocation: Coordinates,
    steps: NavigationStep[]
  ): number => {
    if (steps.length === 0) return 0;

    let closestStepIndex = 0;
    let minDistance = Infinity;

    steps.forEach((step, index) => {
      const distanceToEnd = locationService.calculateDistance(currentLocation, step.endLocation);
      if (distanceToEnd < minDistance) {
        minDistance = distanceToEnd;
        closestStepIndex = index;
      }
    });

    // Si está muy cerca del final del paso (<20m), avanzar al siguiente
    if (minDistance <= 0.02 && closestStepIndex < steps.length - 1) {
      return closestStepIndex + 1;
    }

    return closestStepIndex;
  }, []);

  /**
   * Calculate distance to next turn
   */
  const calculateDistanceToNextStep = useCallback((
    currentLocation: Coordinates,
    currentStep: NavigationStep
  ): number => {
    const distanceKm = locationService.calculateDistance(currentLocation, currentStep.endLocation);
    return Math.round(distanceKm * 1000); // convertir a metros
  }, []);

  const handleDriverLocationUpdate = useCallback(
    async (locationUpdate: LocationUpdate) => {
      const { coordinates, speed, heading } = locationUpdate;
      setDriverLocation(coordinates);
      setDriverSpeed(speed);
      setDriverHeading(heading);

      // Emitir ubicación a pasajeros via WebSocket
      if (activeTrip) {
        if (!socketService.isConnected() && token) {
          socketService.connect(0, token);
          socketService.joinTripTracking(activeTrip.id);
        }
        socketService.emitLocationUpdate(activeTrip.id, {
          tripId: activeTrip.id,
          coordinates,
          speed,
          heading,
          timestamp: locationUpdate.timestamp,
        });
      }

      // Mostrar botón "Llegué al punto de salida" cuando está cerca (<300m)
      if (activeRouteType === 'toOrigin' && activeTrip?.origin?.coordinates) {
        const distToOrigin = locationService.calculateDistance(coordinates, activeTrip.origin.coordinates);
        setIsNearOrigin(distToOrigin <= 0.3);
      }

      // Actualizar paso de navegación actual
      if (navigationSteps.length > 0) {
        const newStepIndex = calculateCurrentStep(coordinates, navigationSteps);
        if (newStepIndex !== currentStepIndex) {
          console.log(`📍 Paso ${newStepIndex + 1}/${navigationSteps.length}`);
          setCurrentStepIndex(newStepIndex);
        }
        const currentStep = navigationSteps[newStepIndex];
        if (currentStep) {
          const distMeters = calculateDistanceToNextStep(coordinates, currentStep);
          setDistanceToNextStep(distMeters);
          setIsApproachingTurn(distMeters < 100);
          setNextStepPreview(navigationSteps[newStepIndex + 1] || null);
        }
      }

      // Recálculo por desviación: si el conductor está a >50m de la ruta → recalcular
      const currentRoute =
        activeRouteType === 'toOrigin' ? routeToOriginRef.current : routeCoordinatesRef.current;

      const isOffRoute = (route: Coordinates[]): boolean => {
        if (route.length === 0) return false;
        const minDist = route.reduce((min, point) => {
          const d = locationService.calculateDistance(coordinates, point);
          return d < min ? d : min;
        }, Infinity);
        return minDist > 0.05; // más de 50 metros de cualquier punto de la ruta
      };

      const now = Date.now();
      const needsRecalc =
        !currentRoute ||
        currentRoute.length === 0 ||
        (isOffRoute(currentRoute) && now - lastRecalcRef.current >= RECALC_COOLDOWN_MS);

      if (!needsRecalc) return;
      lastRecalcRef.current = now;

      if (activeRouteType === 'toOrigin' && activeTrip?.origin?.coordinates) {
        console.log('🔄 Recalculando ruta → punto de salida');
        const dirs = await locationService.getDirections(coordinates, activeTrip.origin.coordinates);
        if (dirs) {
          setRouteToOrigin(dirs.routeCoordinates);
          setRouteToOriginSteps(dirs.steps);
          setNavigationSteps(dirs.steps);
          setEta(dirs.duration);
          setDistanceRemaining(dirs.distance);
        }
      } else if (activeRouteType === 'toDestination' && activeTrip?.destination?.coordinates) {
        console.log('🔄 Recalculando ruta → destino');
        const dirs = await locationService.getDirections(coordinates, activeTrip.destination.coordinates);
        if (dirs) {
          setRouteCoordinates(dirs.routeCoordinates);
          setRouteToDestinationSteps(dirs.steps);
          setNavigationSteps(dirs.steps);
          setEta(dirs.duration);
          setDistanceRemaining(dirs.distance);
        }
      }
    },
    [
      activeTrip,
      token,
      activeRouteType,
      navigationSteps,
      currentStepIndex,
      calculateCurrentStep,
      calculateDistanceToNextStep,
    ]
  );

  const handlePassengerLocationUpdate = useCallback(
    async (event: LocationUpdateEvent) => {
      setDriverLocation(event.coordinates);
      setDriverSpeed(event.speed);
      setDriverHeading(event.heading);

      if (activeTrip?.destination?.coordinates) {
        const directions = await locationService.getDirections(event.coordinates, activeTrip.destination.coordinates);
        if (directions) {
          setEta(directions.duration);
          setDistanceRemaining(directions.distance);
          setRouteCoordinates(directions.routeCoordinates);
        }
      }
    },
    [activeTrip]
  );

  const startTracking = useCallback(
    async (trip: Trip, asDriver: boolean): Promise<StartTrackingResult> => {
      try {
        console.log(`Starting tracking for trip ${trip.id} as ${asDriver ? 'driver' : 'passenger'}`);
        lastRecalcRef.current = 0; // forzar recálculo inmediato al iniciar

        // Set state early so the UI can react
        setActiveTrip(trip);
        setIsDriver(asDriver);

        await AsyncStorage.setItem(
          STORAGE_KEY,
          JSON.stringify({
            tripId: trip.id,
            isDriver: asDriver,
          })
        );

        if (trip.origin?.coordinates && trip.destination?.coordinates) {
          console.log('Fetching initial directions (origin → destination)...');
          const initialDirections = await locationService.getDirections(
            trip.origin.coordinates,
            trip.destination.coordinates
            );

            if (initialDirections) {
              setRouteCoordinates(initialDirections.routeCoordinates);
              setRouteToDestinationSteps(initialDirections.steps); // ← Guardar steps
              console.log(`✅ ${initialDirections.steps.length} pasos para ruta destino`);
            } else {
              console.warn('Could not fetch initial directions.');
            }
          }

          // Si soy conductor, obtener también la ruta desde mi ubicación actual hasta el origen
          if (asDriver) {
            try {
              const currentLocation = await locationService.getCurrentLocation();
              setDriverLocation(currentLocation);

              if (trip.origin?.coordinates) {
                console.log('Fetching route to origin (current location → origin)...');
                const routeToOriginDirections = await locationService.getDirections(
                  currentLocation,
                  trip.origin.coordinates
                );

                if (routeToOriginDirections) {
                  setRouteToOrigin(routeToOriginDirections.routeCoordinates);
                  setRouteToOriginSteps(routeToOriginDirections.steps); // ← Guardar steps

                  // Inicializar navegación en ruta naranja con su ETA correcto
                  setActiveRouteType('toOrigin');
                  setNavigationSteps(routeToOriginDirections.steps);
                  setCurrentStepIndex(0);
                  setEta(routeToOriginDirections.duration); // ← ETA hasta origen
                  setDistanceRemaining(routeToOriginDirections.distance); // ← Distancia hasta origen

                  console.log(`✅ ${routeToOriginDirections.steps.length} pasos para ruta origen`);
                  console.log(`📍 ETA inicial: ${routeToOriginDirections.duration} min, ${routeToOriginDirections.distance} km`);
                } else {
                  console.warn('Could not fetch route to origin.');
                }
              }
            } catch (error) {
              console.error('Error fetching current location or route to origin:', error);
            }
          }

          if (asDriver) {
            const permissions = await locationService.requestPermissions();
            if (!permissions.granted) {
              const message = 'Permisos de ubicación denegados. La app necesita acceso a tu ubicación para iniciar el rastreo.';
              setIsTracking(false); // Reset tracking state
              setActiveTrip(null);
              return { success: false, message };
            }

            const isEnabled = await locationService.isLocationEnabled();
            if (!isEnabled) {
              const message = 'GPS deshabilitado. Por favor, activa los servicios de ubicación en tu dispositivo.';
              setIsTracking(false);
              setActiveTrip(null);
              return { success: false, message };
            }

            await locationService.startTracking(handleDriverLocationUpdate, {
              distanceInterval: 10,
              timeInterval: 5000,
            });

            const initialLocation = await locationService.getCurrentLocation();
            setDriverLocation(initialLocation);

            try {
              const bgPermissions = await locationService.requestBackgroundPermissions();
              if (bgPermissions.granted) {
                await locationService.startBackgroundTracking();
                console.log('✅ Background location tracking enabled');
              } else {
                console.warn('⚠️ Background location not available — using foreground only');
              }
            } catch {
              console.warn('⚠️ Background tracking unavailable — using foreground only');
            }
          } else {
            // Passenger: connect WebSocket to receive driver location updates
            if (token) {
              socketService.connect(user?.id ?? 0, token);
            }
            socketService.joinTripTracking(trip.id);
            socketService.onLocationUpdate(handlePassengerLocationUpdate);
          }

          setIsTracking(true); // Set tracking to true only after all checks and setup
          return { success: true, message: 'Tracking started successfully' };

        } catch (error) {
          console.error('Error starting tracking:', error);
          setIsTracking(false);
          // Don't clear activeTrip — keeps map visible even if routes/location fail
          return { success: false, message: error instanceof Error ? error.message : 'An unknown error occurred' };
        }
      },
      [handleDriverLocationUpdate, handlePassengerLocationUpdate]
      );

      const stopTracking = useCallback(async () => {
        try {
          console.log('Stopping tracking');

          if (isDriver) {
            await locationService.stopTracking();
            await locationService.stopBackgroundTracking();
          } else {
            if (activeTrip && socketService.isConnected()) {
              socketService.leaveTripTracking(activeTrip.id);
              socketService.offLocationUpdate();
            }
          }

          setActiveTrip(null);
          setIsTracking(false);
          setIsDriver(false);
          setDriverLocation(null);
          setDriverSpeed(null);
          setDriverHeading(null);
          setEta(null);
          setDistanceRemaining(null);
          setDistanceTraveled(null);
          setRouteCoordinates(null);
          setRouteToOrigin(null);

          // Limpiar estado de navegación
          setNavigationSteps([]);
          setActiveRouteType(null);
          setCurrentStepIndex(0);
          setDistanceToNextStep(null);
          setNextStepPreview(null);
          setIsApproachingTurn(false);
          setRouteToOriginSteps([]);
          setRouteToDestinationSteps([]);
          setIsNearOrigin(false);

          await AsyncStorage.removeItem(STORAGE_KEY);
        } catch (error) {
          console.error('Error stopping tracking:', error);
        }
      }, [activeTrip, isDriver]);

      const updateDriverLocation = useCallback(async (location: Coordinates) => {
        setDriverLocation(location);

        // Actualizar la ruta desde ubicación actual hasta el origen
        if (activeTrip?.origin?.coordinates) {
          const routeToOriginDirections = await locationService.getDirections(location, activeTrip.origin.coordinates);
          if (routeToOriginDirections) {
            setRouteToOrigin(routeToOriginDirections.routeCoordinates);
          }
        }

        if (activeTrip?.destination.coordinates) {
          const directions = await locationService.getDirections(location, activeTrip.destination.coordinates);
          if (directions) {
            setEta(directions.duration);
            setDistanceRemaining(directions.distance);
            setRouteCoordinates(directions.routeCoordinates);
          } else {
            const distance = locationService.calculateDistance(location, activeTrip.destination.coordinates);
            setDistanceRemaining(distance);
            setRouteCoordinates([location, activeTrip.destination.coordinates]);
          }
        }
      }, [activeTrip]);

      useEffect(() => {
        const restoreTracking = async () => {
          try {
            const stored = await AsyncStorage.getItem(STORAGE_KEY);
            if (stored) {
              await AsyncStorage.removeItem(STORAGE_KEY);
            }
          } catch (error) {
            console.error('Error restoring tracking:', error);
          }
        };

        restoreTracking();
      }, []);

      // No cleanup effect here — TrackingProvider never unmounts, and having
      // stopTracking in deps caused it to fire whenever activeTrip changed mid-tracking.

      const value: TrackingState = {
        activeTrip,
        isTracking,
        isDriver,
        driverLocation,
        driverSpeed,
        driverHeading,
        eta,
        distanceRemaining,
        distanceTraveled,
        routeCoordinates,
        routeToOrigin,
        // Estado de navegación
        navigationSteps,
        activeRouteType,
        currentStepIndex,
        distanceToNextStep,
        nextStepPreview,
        isApproachingTurn,
        isNearOrigin,
        startTracking,
        stopTracking,
        updateDriverLocation,
        confirmArrivalAtOrigin,
      };

      return <TrackingContext.Provider value={value}>{children}</TrackingContext.Provider>;
    };

    export const useTrackingContext = (): TrackingState => {
      const context = useContext(TrackingContext);

      if (context === undefined) {
        throw new Error('useTrackingContext must be used within a TrackingProvider');
      }

      return context;
    };

    export default TrackingContext;
