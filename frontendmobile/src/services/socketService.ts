/**
 * SocketService - Servicio para comunicación WebSocket en tiempo real
 *
 * Proporciona funcionalidades para:
 * - Conectar/desconectar del servidor WebSocket
 * - Unirse a rooms de tracking de viajes
 * - Emitir actualizaciones de ubicación (conductores)
 * - Escuchar actualizaciones de ubicación (pasajeros)
 * - Manejo de eventos de viaje (inicio, finalización, etc.)
 */

import { io, Socket } from 'socket.io-client';
import { Coordinates } from '../types';

const API_HOST = '181.85.173.171';
const TRIPS_WS_URL = `http://${API_HOST}:8002`; // Trips API WebSocket endpoint

export interface LocationUpdateEvent {
  tripId: string;
  coordinates: Coordinates;
  speed: number | null;
  heading: number | null;
  timestamp: number;
}

export interface TripStatusEvent {
  tripId: string;
  status: 'in_progress' | 'completed' | 'cancelled';
  timestamp: number;
}

export interface PassengerJoinedEvent {
  tripId: string;
  passengerId: number;
  passengerName: string;
  timestamp: number;
}

class SocketService {
  private socket: Socket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000; // 1 second
  private isManualDisconnect = false;

  /**
   * Conecta al servidor WebSocket
   * @param userId ID del usuario conectado
   * @param token Token JWT para autenticación
   */
  connect(userId: number, token: string): void {
    if (this.socket?.connected) {
      console.log('Socket already connected');
      return;
    }

    this.isManualDisconnect = false;

    this.socket = io(TRIPS_WS_URL, {
      auth: {
        token,
        userId,
      },
      transports: ['websocket', 'polling'], // Fallback a polling si websocket falla
      reconnection: true,
      reconnectionAttempts: this.maxReconnectAttempts,
      reconnectionDelay: this.reconnectDelay,
      timeout: 10000,
    });

    this.setupEventHandlers();
  }

  /**
   * Configura los manejadores de eventos del socket
   */
  private setupEventHandlers(): void {
    if (!this.socket) return;

    this.socket.on('connect', () => {
      console.log('✅ Socket connected:', this.socket?.id);
      this.reconnectAttempts = 0;
    });

    this.socket.on('disconnect', (reason) => {
      console.log('❌ Socket disconnected:', reason);

      // No intentar reconectar si fue desconexión manual
      if (this.isManualDisconnect) {
        return;
      }

      // Intentar reconectar si no fue intencional
      if (reason === 'io server disconnect') {
        // El servidor desconectó, intentar reconectar manualmente
        this.reconnect();
      }
    });

    this.socket.on('connect_error', (error) => {
      console.error('Socket connection error:', error.message);
      this.reconnectAttempts++;

      if (this.reconnectAttempts >= this.maxReconnectAttempts) {
        console.error('Max reconnection attempts reached');
        this.disconnect();
      }
    });

    this.socket.on('error', (error) => {
      console.error('Socket error:', error);
    });
  }

  /**
   * Intenta reconectar al servidor
   */
  private reconnect(): void {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      setTimeout(() => {
        console.log(`Attempting to reconnect... (${this.reconnectAttempts + 1}/${this.maxReconnectAttempts})`);
        this.socket?.connect();
        this.reconnectAttempts++;
      }, this.reconnectDelay * this.reconnectAttempts);
    }
  }

  /**
   * Desconecta del servidor WebSocket
   */
  disconnect(): void {
    if (this.socket) {
      this.isManualDisconnect = true;
      this.socket.disconnect();
      this.socket = null;
      console.log('Socket manually disconnected');
    }
  }

  /**
   * Verifica si está conectado
   */
  isConnected(): boolean {
    return this.socket?.connected ?? false;
  }

  /**
   * Une al usuario a una room de tracking de viaje
   * @param tripId ID del viaje a rastrear
   */
  joinTripTracking(tripId: string): void {
    if (!this.socket?.connected) {
      console.error('Cannot join trip tracking: Socket not connected');
      return;
    }

    this.socket.emit('join_trip', { tripId });
    console.log(`Joined trip tracking room: ${tripId}`);
  }

  /**
   * Sale de una room de tracking de viaje
   * @param tripId ID del viaje
   */
  leaveTripTracking(tripId: string): void {
    if (!this.socket?.connected) {
      return;
    }

    this.socket.emit('leave_trip', { tripId });
    console.log(`Left trip tracking room: ${tripId}`);
  }

  /**
   * Emite una actualización de ubicación (solo conductores)
   * @param tripId ID del viaje
   * @param location Datos de ubicación
   */
  emitLocationUpdate(tripId: string, location: LocationUpdateEvent): void {
    if (!this.socket?.connected) {
      console.error('Cannot emit location: Socket not connected');
      return;
    }

    this.socket.emit('location_update', {
      tripId,
      ...location,
    });
  }

  /**
   * Escucha actualizaciones de ubicación del conductor
   * @param callback Función a ejecutar cuando se recibe una actualización
   */
  onLocationUpdate(callback: (data: LocationUpdateEvent) => void): void {
    if (!this.socket) {
      console.error('Cannot listen to location updates: Socket not initialized');
      return;
    }

    this.socket.on('location_update', callback);
  }

  /**
   * Deja de escuchar actualizaciones de ubicación
   */
  offLocationUpdate(): void {
    if (this.socket) {
      this.socket.off('location_update');
    }
  }

  /**
   * Escucha cuando un viaje inicia
   * @param callback Función a ejecutar cuando el viaje inicia
   */
  onTripStarted(callback: (data: TripStatusEvent) => void): void {
    if (!this.socket) {
      console.error('Cannot listen to trip started: Socket not initialized');
      return;
    }

    this.socket.on('trip_started', callback);
  }

  /**
   * Deja de escuchar inicio de viajes
   */
  offTripStarted(): void {
    if (this.socket) {
      this.socket.off('trip_started');
    }
  }

  /**
   * Escucha cuando un viaje finaliza
   * @param callback Función a ejecutar cuando el viaje finaliza
   */
  onTripCompleted(callback: (data: TripStatusEvent) => void): void {
    if (!this.socket) {
      console.error('Cannot listen to trip completed: Socket not initialized');
      return;
    }

    this.socket.on('trip_completed', callback);
  }

  /**
   * Deja de escuchar finalización de viajes
   */
  offTripCompleted(): void {
    if (this.socket) {
      this.socket.off('trip_completed');
    }
  }

  /**
   * Escucha cuando un pasajero se une al viaje
   * @param callback Función a ejecutar cuando un pasajero se une
   */
  onPassengerJoined(callback: (data: PassengerJoinedEvent) => void): void {
    if (!this.socket) {
      console.error('Cannot listen to passenger joined: Socket not initialized');
      return;
    }

    this.socket.on('passenger_joined', callback);
  }

  /**
   * Deja de escuchar cuando pasajeros se unen
   */
  offPassengerJoined(): void {
    if (this.socket) {
      this.socket.off('passenger_joined');
    }
  }

  /**
   * Emite evento de inicio de viaje (solo conductores)
   * @param tripId ID del viaje
   */
  emitTripStarted(tripId: string): void {
    if (!this.socket?.connected) {
      console.error('Cannot emit trip start: Socket not connected');
      return;
    }

    this.socket.emit('trip_start', { tripId, timestamp: Date.now() });
  }

  /**
   * Emite evento de finalización de viaje (solo conductores)
   * @param tripId ID del viaje
   */
  emitTripCompleted(tripId: string): void {
    if (!this.socket?.connected) {
      console.error('Cannot emit trip completion: Socket not connected');
      return;
    }

    this.socket.emit('trip_complete', { tripId, timestamp: Date.now() });
  }

  /**
   * Limpia todos los listeners
   */
  removeAllListeners(): void {
    if (this.socket) {
      this.socket.removeAllListeners();
      this.setupEventHandlers(); // Re-setup core handlers
    }
  }
}

// Exportar instancia singleton
export const socketService = new SocketService();
export default socketService;
