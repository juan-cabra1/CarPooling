/**
 * SocketService - Comunicación WebSocket en tiempo real para tracking de viajes
 *
 * Backend: plain WebSocket (gorilla/websocket) en /ws?trip_id=<id>
 * Auth: JWT en Authorization header (React Native soporta headers en WS)
 *
 * Protocolo de mensajes:
 *   Driver → Server: { event: "location_update", trip_id, lat, lng, speed, heading, timestamp }
 *   Server → Client: { event: "location_update", payload: { lat, lng, speed, heading, timestamp } }
 */

import { Coordinates } from '../types'

const TRIPS_WS_BASE = 'wss://trips-api-production.up.railway.app'

export interface LocationUpdateEvent {
  tripId: string
  coordinates: Coordinates
  speed: number | null
  heading: number | null
  timestamp: number
}

export interface TripStatusEvent {
  tripId: string
  status: 'in_progress' | 'completed' | 'cancelled'
  timestamp: number
}

export interface PassengerJoinedEvent {
  tripId: string
  passengerId: number
  passengerName: string
  timestamp: number
}

class SocketService {
  private ws: WebSocket | null = null
  private tripId: string | null = null
  private token: string | null = null
  private locationUpdateCallback: ((data: LocationUpdateEvent) => void) | null = null
  private tripStartedCallback: ((data: TripStatusEvent) => void) | null = null
  private tripCompletedCallback: ((data: TripStatusEvent) => void) | null = null
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private manualDisconnect = false

  connect(userId: number, token: string): void {
    // No-op: connection is established per-trip via joinTripTracking
    this.token = token
  }

  disconnect(): void {
    this.manualDisconnect = true
    this._close()
  }

  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN
  }

  joinTripTracking(tripId: string): void {
    if (!this.token) {
      console.error('[WS] No token set — call connect() first')
      return
    }
    this.tripId = tripId
    this.manualDisconnect = false
    this.reconnectAttempts = 0
    this._connect()
  }

  leaveTripTracking(_tripId: string): void {
    this.manualDisconnect = true
    this._close()
  }

  emitLocationUpdate(tripId: string, location: LocationUpdateEvent): void {
    if (!this.isConnected()) return
    const msg = {
      event: 'location_update',
      trip_id: tripId,
      lat: location.coordinates.lat,
      lng: location.coordinates.lng,
      speed: location.speed ?? 0,
      heading: location.heading ?? 0,
      timestamp: location.timestamp,
    }
    this.ws!.send(JSON.stringify(msg))
  }

  emitTripStarted(tripId: string): void {
    if (!this.isConnected()) return
    this.ws!.send(JSON.stringify({ event: 'trip_start', trip_id: tripId, timestamp: Date.now() }))
  }

  emitTripCompleted(tripId: string): void {
    if (!this.isConnected()) return
    this.ws!.send(JSON.stringify({ event: 'trip_complete', trip_id: tripId, timestamp: Date.now() }))
  }

  onLocationUpdate(callback: (data: LocationUpdateEvent) => void): void {
    this.locationUpdateCallback = callback
  }

  offLocationUpdate(): void {
    this.locationUpdateCallback = null
  }

  onTripStarted(callback: (data: TripStatusEvent) => void): void {
    this.tripStartedCallback = callback
  }

  offTripStarted(): void {
    this.tripStartedCallback = null
  }

  onTripCompleted(callback: (data: TripStatusEvent) => void): void {
    this.tripCompletedCallback = callback
  }

  offTripCompleted(): void {
    this.tripCompletedCallback = null
  }

  removeAllListeners(): void {
    this.locationUpdateCallback = null
    this.tripStartedCallback = null
    this.tripCompletedCallback = null
  }

  // ─── internals ────────────────────────────────────────────────────────────

  private _connect(): void {
    if (!this.tripId || !this.token) return
    this._close()

    const url = `${TRIPS_WS_BASE}/ws?trip_id=${encodeURIComponent(this.tripId)}`

    // React Native WebSocket accepts headers as 3rd argument (not in TS types but supported natively)
    const WS = WebSocket as any
    const ws: WebSocket = new WS(url, null, {
      headers: { Authorization: `Bearer ${this.token}` },
    })
    this.ws = ws

    ws.onopen = () => {
      console.log(`[WS] Connected to trip ${this.tripId}`)
      this.reconnectAttempts = 0
    }

    ws.onmessage = (e) => {
      this._handleMessage(e.data)
    }

    ws.onerror = (e) => {
      console.warn('[WS] Error:', (e as any).message)
    }

    ws.onclose = () => {
      console.log('[WS] Closed')
      if (!this.manualDisconnect && this.reconnectAttempts < this.maxReconnectAttempts) {
        const delay = 1000 * Math.pow(2, this.reconnectAttempts) // exponential backoff
        this.reconnectAttempts++
        console.log(`[WS] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`)
        this.reconnectTimer = setTimeout(() => this._connect(), delay)
      }
    }
  }

  private _close(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    if (this.ws) {
      this.ws.onopen = null
      this.ws.onmessage = null
      this.ws.onerror = null
      this.ws.onclose = null
      if (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING) {
        this.ws.close()
      }
      this.ws = null
    }
  }

  private _handleMessage(raw: string): void {
    let msg: any
    try {
      msg = JSON.parse(raw)
    } catch {
      return
    }

    const event = msg.event
    const payload = msg.payload ?? msg

    switch (event) {
      case 'location_update': {
        if (this.locationUpdateCallback && this.tripId) {
          this.locationUpdateCallback({
            tripId: this.tripId,
            coordinates: { lat: payload.lat, lng: payload.lng },
            speed: payload.speed ?? null,
            heading: payload.heading ?? null,
            timestamp: payload.timestamp ?? Date.now(),
          })
        }
        break
      }
      case 'trip_started': {
        this.tripStartedCallback?.({
          tripId: payload.trip_id ?? this.tripId ?? '',
          status: 'in_progress',
          timestamp: payload.timestamp ?? Date.now(),
        })
        break
      }
      case 'trip_completed': {
        this.tripCompletedCallback?.({
          tripId: payload.trip_id ?? this.tripId ?? '',
          status: 'completed',
          timestamp: payload.timestamp ?? Date.now(),
        })
        break
      }
    }
  }
}

export const socketService = new SocketService()
export default socketService
