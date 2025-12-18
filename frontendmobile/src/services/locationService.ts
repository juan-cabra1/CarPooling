/**
 * LocationService - Servicio para manejo de ubicación GPS
 *
 * Proporciona funcionalidades para:
 * - Obtener ubicación actual
 * - Tracking continuo de ubicación (foreground)
 * - Tracking en background (para conductores)
 * - Manejo de permisos
 * - Cálculo de distancias
 * - Obtención de rutas desde Google Maps
 */

import * as Location from 'expo-location';
import axios from 'axios';
import { Coordinates } from '../types';
import { GOOGLE_MAPS_API_KEY } from '../config';

export interface LocationPermissionStatus {
  granted: boolean;
  canAskAgain: boolean;
}

export interface LocationUpdate {
  coordinates: Coordinates;
  speed: number | null; // km/h
  heading: number | null; // degrees
  accuracy: number | null; // meters
  timestamp: number;
}

/**
 * Maneuver types from Google Directions API
 * Maps to visual icons and audio instructions
 */
export type ManeuverType =
  | 'turn-left'
  | 'turn-right'
  | 'turn-slight-left'
  | 'turn-slight-right'
  | 'turn-sharp-left'
  | 'turn-sharp-right'
  | 'uturn-left'
  | 'uturn-right'
  | 'keep-left'
  | 'keep-right'
  | 'merge'
  | 'ramp-left'
  | 'ramp-right'
  | 'fork-left'
  | 'fork-right'
  | 'ferry'
  | 'roundabout-left'
  | 'roundabout-right'
  | 'straight'
  | '';

/**
 * Helper to get maneuver icon
 */
export interface ManeuverIcon {
  unicode: string; // Unicode arrow character
  description: string; // Accessibility text
}

/**
 * Single navigation step with all details
 */
export interface NavigationStep {
  // Position in route
  stepIndex: number;

  // Instruction details
  instruction: string; // Plain text (HTML stripped)
  maneuver: ManeuverType;
  streetName?: string; // Extracted from instruction if available

  // Distances and duration
  distance: number; // meters for this step
  duration: number; // seconds for this step

  // Coordinates
  startLocation: Coordinates;
  endLocation: Coordinates;
  polyline: Coordinates[]; // Detailed path for this step

  // Metadata
  distanceText: string; // "200m", "1.2km"
  durationText: string; // "2 min"
}

export interface DirectionsResponse {
  routeCoordinates: Coordinates[];
  distance: number; // en kilómetros
  duration: number; // en minutos
  steps: NavigationStep[]; // ← NUEVO: Navigation steps
}

/**
 * Maps maneuver types to visual icons
 */
export const MANEUVER_ICONS: Record<ManeuverType | 'default', ManeuverIcon> = {
  'turn-left': { unicode: '↰', description: 'Gira a la izquierda' },
  'turn-right': { unicode: '↱', description: 'Gira a la derecha' },
  'turn-slight-left': { unicode: '↖', description: 'Gira levemente a la izquierda' },
  'turn-slight-right': { unicode: '↗', description: 'Gira levemente a la derecha' },
  'turn-sharp-left': { unicode: '⬅', description: 'Gira pronunciadamente a la izquierda' },
  'turn-sharp-right': { unicode: '➡', description: 'Gira pronunciadamente a la derecha' },
  'uturn-left': { unicode: '↶', description: 'Da vuelta en U' },
  'uturn-right': { unicode: '↷', description: 'Da vuelta en U' },
  'straight': { unicode: '↑', description: 'Continúa recto' },
  'keep-left': { unicode: '⬉', description: 'Mantente a la izquierda' },
  'keep-right': { unicode: '⬈', description: 'Mantente a la derecha' },
  'merge': { unicode: '⛙', description: 'Incorpórate' },
  'ramp-left': { unicode: '⬑', description: 'Toma la rampa izquierda' },
  'ramp-right': { unicode: '⬏', description: 'Toma la rampa derecha' },
  'fork-left': { unicode: '⑂', description: 'Toma el desvío izquierdo' },
  'fork-right': { unicode: '⑃', description: 'Toma el desvío derecho' },
  'roundabout-left': { unicode: '↺', description: 'Entra en la rotonda' },
  'roundabout-right': { unicode: '↻', description: 'Entra en la rotonda' },
  'ferry': { unicode: '⛴', description: 'Toma el ferry' },
  '': { unicode: '●', description: 'Continúa' },
  'default': { unicode: '↑', description: 'Continúa' },
};

/**
 * Get icon for maneuver type
 */
export function getManeuverIcon(maneuver: ManeuverType): ManeuverIcon {
  return MANEUVER_ICONS[maneuver] || MANEUVER_ICONS.default;
}

/**
 * Format distance for display
 */
export function formatDistance(meters: number): string {
  if (meters < 1000) {
    return `${Math.round(meters)}m`;
  }
  return `${(meters / 1000).toFixed(1)}km`;
}

/**
 * Format duration for display
 */
export function formatDuration(seconds: number): string {
  const minutes = Math.round(seconds / 60);
  if (minutes < 1) return '< 1 min';
  return `${minutes} min`;
}

class LocationService {
  private watchId: Location.LocationSubscription | null = null;
  private backgroundTaskName = 'CARPOOLING_LOCATION_TRACKING';

  /**
   * Strip HTML tags from instruction text
   */
  private stripHtml(html: string): string {
    return html
      .replace(/<[^>]*>/g, '') // Remove HTML tags
      .replace(/&nbsp;/g, ' ')
      .replace(/&amp;/g, '&')
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      .replace(/&quot;/g, '"')
      .trim();
  }

  /**
   * Extract street name from instruction if available
   * Example: "Turn left onto Av. Colón" -> "Av. Colón"
   */
  private extractStreetName(instruction: string): string | undefined {
    const ontoMatch = instruction.match(/onto\s+([^,]+)/i);
    if (ontoMatch) return ontoMatch[1].trim();

    const enMatch = instruction.match(/en\s+([^,]+)/i);
    if (enMatch) return enMatch[1].trim();

    return undefined;
  }

  /**
   * Decodifica una polyline de Google Maps a un array de coordenadas
   * @param encoded La polyline codificada
   */
  private decodePolyline(encoded: string): Coordinates[] {
    if (!encoded) {
      return [];
    }

    let points: Coordinates[] = [];
    let index = 0,
      len = encoded.length;
    let lat = 0,
      lng = 0;

    while (index < len) {
      let b,
        shift = 0,
        result = 0;
      do {
        b = encoded.charCodeAt(index++) - 63;
        result |= (b & 0x1f) << shift;
        shift += 5;
      } while (b >= 0x20);
      let dlat = (result & 1) !== 0 ? ~(result >> 1) : result >> 1;
      lat += dlat;

      shift = 0;
      result = 0;
      do {
        b = encoded.charCodeAt(index++) - 63;
        result |= (b & 0x1f) << shift;
        shift += 5;
      } while (b >= 0x20);
      let dlng = (result & 1) !== 0 ? ~(result >> 1) : result >> 1;
      lng += dlng;

      points.push({ lat: lat / 1e5, lng: lng / 1e5 });
    }
    return points;
  }

  /**
   * Obtiene la ruta, distancia y duración desde Google Maps Directions API
   * @param origin Coordenadas de origen
   * @param destination Coordenadas de destino
   */
  async getDirections(
    origin: Coordinates,
    destination: Coordinates
  ): Promise<DirectionsResponse | null> {
    if (!GOOGLE_MAPS_API_KEY) {
      console.error('GOOGLE_MAPS_API_KEY no está configurada.');
      // Devolver una ruta en línea recta como fallback si no hay API key
      return {
        routeCoordinates: [origin, destination],
        distance: this.calculateDistance(origin, destination),
        duration: (this.calculateDistance(origin, destination) / 50) * 60, // Asume 50 km/h
        steps: [], // No steps without API
      };
    }

    const url = `https://maps.googleapis.com/maps/api/directions/json?origin=${origin.lat},${origin.lng}&destination=${destination.lat},${destination.lng}&key=${GOOGLE_MAPS_API_KEY}`;

    try {
      const response = await axios.get(url);
      const data = response.data;

      if (data.status !== 'OK' || !data.routes || data.routes.length === 0) {
        console.error('Error fetching directions:', data.error_message || data.status);
        // Fallback to a straight line if API returns an error
        return {
          routeCoordinates: [origin, destination],
          distance: this.calculateDistance(origin, destination),
          duration: (this.calculateDistance(origin, destination) / 50) * 60, // Assume 50 km/h
          steps: [],
        };
      }

      const route = data.routes[0];
      const leg = route.legs[0];

      // Parse navigation steps from API response
      const steps: NavigationStep[] = leg.steps.map((step: any, index: number) => {
        const stepPolyline = this.decodePolyline(step.polyline.points);
        const instruction = this.stripHtml(step.html_instructions);

        return {
          stepIndex: index,
          instruction,
          maneuver: (step.maneuver || '') as ManeuverType,
          streetName: this.extractStreetName(instruction),
          distance: step.distance.value, // meters
          duration: step.duration.value, // seconds
          startLocation: {
            lat: step.start_location.lat,
            lng: step.start_location.lng,
          },
          endLocation: {
            lat: step.end_location.lat,
            lng: step.end_location.lng,
          },
          polyline: stepPolyline,
          distanceText: formatDistance(step.distance.value),
          durationText: formatDuration(step.duration.value),
        };
      });

      const routeCoordinates = this.decodePolyline(route.overview_polyline.points);
      const distance = leg.distance.value / 1000; // a kilómetros
      const duration = leg.duration.value / 60; // a minutos

      console.log(`✅ Parsed ${steps.length} navigation steps from Google Directions API`);

      return {
        routeCoordinates,
        distance: parseFloat(distance.toFixed(1)),
        duration: Math.round(duration),
        steps,
      };
    } catch (error) {
      console.error('Error en la llamada a Google Directions API:', error);
      // Fallback a línea recta en caso de error de red
      return {
        routeCoordinates: [origin, destination],
        distance: this.calculateDistance(origin, destination),
        duration: (this.calculateDistance(origin, destination) / 50) * 60, // Asume 50 km/h
        steps: [],
      };
    }
  }

  /**
   * Solicita permisos de ubicación al usuario
   * @returns true si los permisos fueron otorgados
   */
  async requestPermissions(): Promise<LocationPermissionStatus> {
    try {
      const { status, canAskAgain } = await Location.requestForegroundPermissionsAsync();

      if (status === 'granted') {
        return { granted: true, canAskAgain };
      }

      console.log('Location permission denied');
      return { granted: false, canAskAgain };
    } catch (error) {
      console.error('Error requesting location permissions:', error);
      return { granted: false, canAskAgain: false };
    }
  }

  /**
   * Solicita permisos de ubicación en background (solo para conductores)
   * Debe llamarse DESPUÉS de obtener permisos de foreground
   * NOTA: En Expo Go, esto puede fallar. Solo funciona en builds nativos.
   */
  async requestBackgroundPermissions(): Promise<LocationPermissionStatus> {
    try {
      const { status, canAskAgain } = await Location.requestBackgroundPermissionsAsync();

      if (status === 'granted') {
        return { granted: true, canAskAgain };
      }

      console.log('Background location permission denied');
      return { granted: false, canAskAgain };
    } catch (error) {
      // En Expo Go, esto falla porque no tiene los NSLocation keys configurados
      // Es esperado y no es un error crítico - solo afecta tracking en background
      console.warn('Background location permission not available (Expo Go limitation):', error);
      return { granted: false, canAskAgain: false };
    }
  }

  /**
   * Verifica el estado actual de los permisos
   */
  async checkPermissions(): Promise<LocationPermissionStatus> {
    try {
      const { status, canAskAgain } = await Location.getForegroundPermissionsAsync();
      return {
        granted: status === 'granted',
        canAskAgain,
      };
    } catch (error) {
      console.error('Error checking location permissions:', error);
      return { granted: false, canAskAgain: false };
    }
  }

  /**
   * Obtiene la ubicación actual del usuario (una sola vez)
   * @returns Coordenadas actuales
   */
  async getCurrentLocation(): Promise<Coordinates> {
    try {
      const location = await Location.getCurrentPositionAsync({
        accuracy: Location.Accuracy.High,
      });

      return {
        lat: location.coords.latitude,
        lng: location.coords.longitude,
      };
    } catch (error) {
      console.error('Error getting current location:', error);
      throw new Error('No se pudo obtener la ubicación actual. Verifica los permisos y el GPS.');
    }
    
  }

  /**
   * Inicia el tracking continuo de ubicación (foreground)
   * @param callback Función a ejecutar cada vez que cambia la ubicación
   * @param options Opciones de tracking
   */
  async startTracking(
    callback: (location: LocationUpdate) => void,
    options?: {
      accuracy?: Location.Accuracy;
      distanceInterval?: number; // meters
      timeInterval?: number; // milliseconds
    }
  ): Promise<void> {
    try {
      // Detener tracking anterior si existe
      if (this.watchId) {
        await this.stopTracking();
      }

      const { accuracy = Location.Accuracy.High, distanceInterval = 10, timeInterval = 5000 } = options || {};

      this.watchId = await Location.watchPositionAsync(
        {
          accuracy,
          distanceInterval,
          timeInterval,
        },
        (location) => {
          const update: LocationUpdate = {
            coordinates: {
              lat: location.coords.latitude,
              lng: location.coords.longitude,
            },
            speed: location.coords.speed ? location.coords.speed * 3.6 : null, // m/s to km/h
            heading: location.coords.heading,
            accuracy: location.coords.accuracy,
            timestamp: location.timestamp,
          };

          callback(update);
        }
      );

      console.log('Location tracking started');
    } catch (error) {
      console.error('Error starting location tracking:', error);
      throw new Error('No se pudo iniciar el tracking de ubicación');
    }
  }

  /**
   * Detiene el tracking de ubicación
   */
  async stopTracking(): Promise<void> {
    try {
      if (this.watchId) {
        this.watchId.remove();
        this.watchId = null;
        console.log('Location tracking stopped');
      }
    } catch (error) {
      console.error('Error stopping location tracking:', error);
    }
  }

  /**
   * Inicia tracking en background (solo para conductores durante viajes activos)
   * IMPORTANTE: Requiere permisos de background previamente otorgados
   */
  async startBackgroundTracking(): Promise<void> {
    try {
      // Verificar si ya está registrada la tarea
      const isRegistered = await Location.hasStartedLocationUpdatesAsync(this.backgroundTaskName);

      if (isRegistered) {
        console.log('Background tracking already running');
        return;
      }

      await Location.startLocationUpdatesAsync(this.backgroundTaskName, {
        accuracy: Location.Accuracy.High,
        distanceInterval: 20, // Update every 20 meters
        timeInterval: 10000, // Update every 10 seconds minimum
        foregroundService: {
          notificationTitle: 'CarPooling - Viaje en progreso',
          notificationBody: 'Tu ubicación está siendo compartida con los pasajeros',
          notificationColor: '#4A90E2',
        },
      });

      console.log('Background location tracking started');
    } catch (error) {
      console.error('Error starting background tracking:', error);
      throw new Error('No se pudo iniciar el tracking en background');
    }
  }

  /**
   * Detiene el tracking en background
   */
  async stopBackgroundTracking(): Promise<void> {
    try {
      const isRegistered = await Location.hasStartedLocationUpdatesAsync(this.backgroundTaskName);

      if (isRegistered) {
        await Location.stopLocationUpdatesAsync(this.backgroundTaskName);
        console.log('Background location tracking stopped');
      }
    } catch (error) {
      console.error('Error stopping background tracking:', error);
    }
  }

  /**
   * Calcula la distancia entre dos puntos (en kilómetros)
   * Usa la fórmula de Haversine
   */
  calculateDistance(from: Coordinates, to: Coordinates): number {
    const R = 6371; // Radio de la Tierra en km
    const dLat = this.toRad(to.lat - from.lat);
    const dLon = this.toRad(to.lng - from.lng);

    const a =
      Math.sin(dLat / 2) * Math.sin(dLat / 2) +
      Math.cos(this.toRad(from.lat)) *
      Math.cos(this.toRad(to.lat)) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);

    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    const distance = R * c;

    return Math.round(distance * 100) / 100; // Redondear a 2 decimales
  }

  /**
   * Convierte grados a radianes
   */
  private toRad(degrees: number): number {
    return degrees * (Math.PI / 180);
  }

  /**
   * Verifica si los servicios de ubicación están habilitados en el dispositivo
   */
  async isLocationEnabled(): Promise<boolean> {
    try {
      return await Location.hasServicesEnabledAsync();
    } catch (error) {
      console.error('Error checking if location is enabled:', error);
      return false;
    }
  }
}

// Exportar instancia singleton
export const locationService = new LocationService();
export default locationService;
