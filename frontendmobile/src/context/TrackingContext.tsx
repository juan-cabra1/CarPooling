/**
 * TrackingContext - Contexto global para manejo de tracking de viajes
 *
 * Proporciona:
 * - Estado global del tracking activo
 * - Ubicación del conductor en tiempo real
 * - ETA y distancia restante
 * - Funciones para iniciar/detener tracking
 */

import React, { createContext, useContext, useState, useEffect, ReactNode, useCallback } from 'react';
import { Trip, Coordinates } from '../types';
import { locationService, LocationUpdate, DirectionsResponse, NavigationStep } from '../services/locationService';
import { socketService, LocationUpdateEvent } from '../services/socketService';
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

  // Acciones
  startTracking: (trip: Trip, asDriver: boolean) => Promise<StartTrackingResult>;
  stopTracking: () => Promise<void>;
  updateDriverLocation: (location: Coordinates) => void;
}

const TrackingContext = createContext<TrackingState | undefined>(undefined);

const STORAGE_KEY = '@carpooling:active_trip';

interface TrackingProviderProps {
  children: ReactNode;
}

export const TrackingProvider: React.FC<TrackingProviderProps> = ({ children }) => {
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

  // Cache para evitar re-fetches
  const [routeToOriginSteps, setRouteToOriginSteps] = useState<NavigationStep[]>([]);
  const [routeToDestinationSteps, setRouteToDestinationSteps] = useState<NavigationStep[]>([]);

  const calculateETA = useCallback((currentLocation: Coordinates, destination: Coordinates, speed: number | null) => {
    const distance = locationService.calculateDistance(currentLocation, destination);
    setDistanceRemaining(distance);

    if (speed && speed > 0) {
      const etaMinutes = (distance / speed) * 60;
      setEta(Math.round(etaMinutes));
    } else {
      const etaMinutes = (distance / 50) * 60;
      setEta(Math.round(etaMinutes));
    }
  }, []);

  /**
   * Determine which route is active based on driver's location
   * Switch from 'toOrigin' to 'toDestination' when driver reaches origin
   */
  const determineActiveRoute = useCallback((currentLocation: Coordinates) => {
    if (!activeTrip?.origin?.coordinates) return;

    const distanceToOrigin = locationService.calculateDistance(
      currentLocation,
      activeTrip.origin.coordinates
    );

    const ARRIVAL_THRESHOLD_KM = 0.05; // 50 metros

    if (distanceToOrigin <= ARRIVAL_THRESHOLD_KM) {
      // Llegó al origen → cambiar a ruta azul
      if (activeRouteType !== 'toDestination') {
        console.log('🎯 Conductor llegó al origen! Cambiando a ruta destino.');
        setActiveRouteType('toDestination');
        setNavigationSteps(routeToDestinationSteps);
        setCurrentStepIndex(0);
      }
    } else {
      // Todavía yendo al origen → ruta naranja
      if (activeRouteType !== 'toOrigin') {
        setActiveRouteType('toOrigin');
        setNavigationSteps(routeToOriginSteps);
        setCurrentStepIndex(0);
      }
    }
  }, [activeTrip, activeRouteType, routeToOriginSteps, routeToDestinationSteps]);

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

      if (activeTrip && socketService.isConnected()) {
        const event: LocationUpdateEvent = {
          tripId: activeTrip.id, // Use id instead of _id
          coordinates,
          speed,
          heading,
          timestamp: locationUpdate.timestamp,
        };
        socketService.emitLocationUpdate(activeTrip.id, event);
      }

      // Determinar qué ruta está activa
      determineActiveRoute(coordinates);

      // Actualizar ruta según el tipo activo
      if (activeRouteType === 'toOrigin' && activeTrip?.origin?.coordinates) {
        // Yendo al origen (ruta naranja)
        const routeToOriginDirections = await locationService.getDirections(
          coordinates,
          activeTrip.origin.coordinates
        );
        if (routeToOriginDirections) {
          setRouteToOrigin(routeToOriginDirections.routeCoordinates);
          setRouteToOriginSteps(routeToOriginDirections.steps);
          setNavigationSteps(routeToOriginDirections.steps);
          // ETA hasta el origen
          setEta(routeToOriginDirections.duration);
          setDistanceRemaining(routeToOriginDirections.distance);
        }
      } else if (activeRouteType === 'toDestination' && activeTrip?.destination.coordinates) {
        // Yendo al destino (ruta azul)
        const routeToDestDirections = await locationService.getDirections(
          coordinates,
          activeTrip.destination.coordinates
        );
        if (routeToDestDirections) {
          setRouteCoordinates(routeToDestDirections.routeCoordinates);
          setRouteToDestinationSteps(routeToDestDirections.steps);
          setNavigationSteps(routeToDestDirections.steps);
          // ETA hasta el destino
          setEta(routeToDestDirections.duration);
          setDistanceRemaining(routeToDestDirections.distance);
        }
      }

      // Calcular paso actual
      if (navigationSteps.length > 0) {
        const newStepIndex = calculateCurrentStep(coordinates, navigationSteps);

        if (newStepIndex !== currentStepIndex) {
          console.log(`📍 Avanzó al paso ${newStepIndex + 1}/${navigationSteps.length}`);
          setCurrentStepIndex(newStepIndex);
        }

        const currentStep = navigationSteps[newStepIndex];
        if (currentStep) {
          const distanceMeters = calculateDistanceToNextStep(coordinates, currentStep);
          setDistanceToNextStep(distanceMeters);
          setIsApproachingTurn(distanceMeters < 100);
          setNextStepPreview(navigationSteps[newStepIndex + 1] || null);
        }
      }
    },
    [
      activeTrip,
      routeCoordinates,
      navigationSteps,
      currentStepIndex,
      activeRouteType,
      calculateETA,
      determineActiveRoute,
      calculateCurrentStep,
      calculateDistanceToNextStep,
    ]
  );

  const handlePassengerLocationUpdate = useCallback(
    async (event: LocationUpdateEvent) => {
      setDriverLocation(event.coordinates);
      setDriverSpeed(event.speed);
      setDriverHeading(event.heading);

      if (activeTrip?.destination.coordinates && routeCoordinates) {
        const directions = await locationService.getDirections(event.coordinates, activeTrip.destination.coordinates);
        if (directions) {
          setEta(directions.duration);
          setDistanceRemaining(directions.distance);
        }
      } else if (activeTrip?.destination.coordinates) {
        calculateETA(event.coordinates, activeTrip.destination.coordinates, event.speed);
      }
    },
    [activeTrip, routeCoordinates, calculateETA]
  );

  const startTracking = useCallback(
    async (trip: Trip, asDriver: boolean): Promise<StartTrackingResult> => {
      try {
        console.log(`Starting tracking for trip ${trip.id} as ${asDriver ? 'driver' : 'passenger'}`);
        
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
            
            const bgPermissions = await locationService.requestBackgroundPermissions();
            if (bgPermissions.granted) {
              console.log('✅ Background location tracking enabled');
              await locationService.startBackgroundTracking();
            } else {
              console.warn('⚠️ Background location tracking not available (use foreground tracking instead)');
            }
          } else {
            if (socketService.isConnected()) {
              socketService.joinTripTracking(trip.id);
              socketService.onLocationUpdate(handlePassengerLocationUpdate);
            }
          }
          
          setIsTracking(true); // Set tracking to true only after all checks and setup
          return { success: true, message: 'Tracking started successfully' };
          
        } catch (error) {
          console.error('Error starting tracking:', error);
          // Cleanup state on unexpected error
          setIsTracking(false);
          setActiveTrip(null);
          setRouteCoordinates(null);
          return { success: false, message: error instanceof Error ? error.message : 'An unknown error occurred' };
        }
      },
      [handleDriverLocationUpdate, handlePassengerLocationUpdate, calculateETA]
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
      
      useEffect(() => {
        return () => {
          if (isTracking) {
            stopTracking();
          }
        };
      }, [isTracking, stopTracking]);
      
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
        startTracking,
        stopTracking,
        updateDriverLocation,
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