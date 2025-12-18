Plan: Implementación de Sistema de Mapas y Tracking en Tiempo Real (Estilo Uber)
Resumen Ejecutivo
Implementar un sistema completo de mapas interactivos con tracking en tiempo real para la aplicación CarPooling, similar a Uber/Cabify. El sistema permitirá visualizar rutas, rastrear la ubicación del conductor en vivo, calcular ETAs dinámicos y mostrar el progreso del viaje tanto a conductores como pasajeros.
Estado Actual del Proyecto
✅ Ya Implementado (Ventajas)
Infraestructura de Coordenadas Completa
Todos los viajes tienen origin/destination con coordenadas (lat/lng)
Backend MongoDB con índices 2dsphere para queries geoespaciales
LocationInput component con autocomplete usando Photon API (GRATIS)
Tipos TypeScript bien definidos para Location y Coordinates
Estados de Viaje
Backend tiene estados: draft, published, in_progress, completed, cancelled
Infraestructura para actualizar status existe
Arquitectura Microservicios
trips-api (Go + MongoDB) maneja ciclo de vida del viaje
RabbitMQ para eventos entre servicios
Sistema de mensajería asíncrona funcionando
Frontend Mobile
React Native + Expo
React Navigation configurado
AuthContext para manejo de sesión
❌ Faltante para Tracking en Tiempo Real
Librerías de Mapas: No hay react-native-maps ni mapbox instalado
Servicios de Ubicación: No hay expo-location
Comunicación en Tiempo Real: No hay WebSockets ni Socket.io
Permisos GPS: No configurados en app.json
UI de Tracking: No hay pantallas de mapa
Backend WebSocket: No existe servidor WebSocket para broadcast
Decisión Técnica: Google Maps vs Mapbox vs OpenStreetMap
Recomendación: Google Maps Platform
Razones:
Costo para MVP: $200 USD/mes GRATIS = ~28,000 cargas de mapa
Facilidad: react-native-maps incluido en Expo, documentación excelente
Precisión: Mejores rutas y direcciones que alternativas
Escalabilidad: Solo pagas si creces (problema de éxito)
Estimación de Costos:
100 usuarios activos/mes: $0 (dentro límite gratuito)
1,000 usuarios activos/mes: $0 (aún gratis)
10,000 usuarios activos/mes: ~$50-80 USD/mes
Arquitectura Propuesta
Frontend (React Native)
src/
├── screens/
│   ├── TripTrackingScreen/          # Nueva - Vista tracking para pasajero
│   │   ├── TripTrackingScreen.tsx
│   │   └── TripTrackingScreen.styles.ts
│   └── DriverTrackingScreen/        # Nueva - Vista tracking para conductor
│       ├── DriverTrackingScreen.tsx
│       └── DriverTrackingScreen.styles.ts
├── components/
│   └── map/                         # Nuevos componentes
│       ├── MapView.tsx              # Wrapper de react-native-maps
│       ├── TripRouteMap.tsx         # Mapa con ruta trazada
│       ├── LiveTrackingMap.tsx      # Mapa con tracking en vivo
│       └── LocationMarker.tsx       # Marcadores personalizados
├── services/
│   ├── locationService.ts           # Nueva - Manejo de GPS
│   ├── socketService.ts             # Nueva - WebSocket client
│   └── trackingService.ts           # Nueva - Lógica de tracking
├── context/
│   └── TrackingContext.tsx          # Nuevo - Estado global de tracking
└── hooks/
    ├── useLocation.ts               # Nuevo - Hook para GPS
    └── useTracking.ts               # Nuevo - Hook para tracking
Backend (Go)
Opción Recomendada: Extender trips-api
backend/trips-api/
├── internal/
│   ├── domain/
│   │   └── trip.go                  # Modificar - Agregar campos location tracking
│   ├── controller/
│   │   ├── trip_controller.go       # Modificar - Agregar endpoints tracking
│   │   └── websocket_controller.go  # Nuevo - WebSocket handler
│   ├── service/
│   │   ├── trip_service.go          # Modificar - Lógica tracking
│   │   └── tracking_service.go      # Nuevo - Servicio tracking
│   └── websocket/
│       ├── hub.go                   # Nuevo - WebSocket hub
│       └── client.go                # Nuevo - WebSocket client manager
└── cmd/api/
    └── main.go                      # Modificar - Agregar WebSocket routes
Plan de Implementación Detallado
FASE 1: Configuración Base (Día 1-2)
1.1 Configurar Google Maps API
Archivo: d:\frontendmobile\CarPooling\.env Acciones:
Crear proyecto en Google Cloud Console
Activar APIs necesarias:
Maps SDK for Android
Maps SDK for iOS
Directions API
Distance Matrix API (opcional)
Crear API Key con restricciones
Agregar al .env:
GOOGLE_MAPS_API_KEY=AIza...
1.2 Instalar Dependencias Frontend
Comando:
cd d:\frontendmobile
expo install expo-location
expo install react-native-maps
npm install socket.io-client
npm install zustand  # State management ligero (opcional)
1.3 Configurar Permisos
Archivo: d:\frontendmobile\app.json Modificación:
{
  "expo": {
    "plugins": [
      [
        "expo-location",
        {
          "locationAlwaysAndWhenInUsePermission": "La app necesita tu ubicación para rastrear el viaje en tiempo real.",
          "locationAlwaysPermission": "Permite que la app rastree tu ubicación incluso cuando está en segundo plano.",
          "locationWhenInUsePermission": "La app necesita tu ubicación para mostrar tu posición en el mapa.",
          "isIosBackgroundLocationEnabled": true,
          "isAndroidBackgroundLocationEnabled": true,
          "isAndroidForegroundServiceEnabled": true
        }
      ]
    ],
    "android": {
      "permissions": [
        "ACCESS_FINE_LOCATION",
        "ACCESS_COARSE_LOCATION",
        "ACCESS_BACKGROUND_LOCATION",
        "FOREGROUND_SERVICE",
        "FOREGROUND_SERVICE_LOCATION"
      ],
      "config": {
        "googleMaps": {
          "apiKey": "AIza..."
        }
      }
    },
    "ios": {
      "config": {
        "googleMapsApiKey": "AIza..."
      },
      "infoPlist": {
        "NSLocationWhenInUseUsageDescription": "Necesitamos tu ubicación para mostrarte en el mapa y calcular rutas.",
        "NSLocationAlwaysAndWhenInUseUsageDescription": "Permite que rastreemos tu ubicación durante el viaje para actualizar tu posición en tiempo real.",
        "NSLocationAlwaysUsageDescription": "La app necesita acceso a tu ubicación en segundo plano para rastrear viajes activos."
      }
    }
  }
}
FASE 2: Servicios Core (Día 3-4)
2.1 Crear Location Service
Archivo: d:\frontendmobile\src\services\locationService.ts Funcionalidades:
class LocationService {
  // Obtener ubicación actual una vez
  async getCurrentLocation(): Promise<Coordinates>

  // Iniciar tracking continuo (foreground)
  startTracking(callback: (location: Coordinates) => void): void

  // Iniciar tracking en background (solo conductor)
  startBackgroundTracking(): void

  // Detener tracking
  stopTracking(): void

  // Verificar permisos
  async requestPermissions(): Promise<boolean>

  // Calcular distancia entre dos puntos
  calculateDistance(from: Coordinates, to: Coordinates): number
}
2.2 Crear WebSocket Service
Archivo: d:\frontendmobile\src\services\socketService.ts Funcionalidades:
class SocketService {
  private socket: Socket

  // Conectar a WebSocket
  connect(userId: number, token: string): void

  // Desconectar
  disconnect(): void

  // Unirse a room de tracking de viaje
  joinTripTracking(tripId: string): void

  // Salir de room
  leaveTripTracking(tripId: string): void

  // Emitir actualización de ubicación (conductor)
  emitLocationUpdate(tripId: string, location: Coordinates): void

  // Escuchar actualizaciones de ubicación (pasajero)
  onLocationUpdate(callback: (data: LocationUpdate) => void): void

  // Escuchar inicio de viaje
  onTripStarted(callback: (tripId: string) => void): void

  // Escuchar fin de viaje
  onTripCompleted(callback: (tripId: string) => void): void
}
2.3 Crear Tracking Context
Archivo: d:\frontendmobile\src\context\TrackingContext.tsx Estado Global:
interface TrackingState {
  activeTrip: Trip | null
  driverLocation: Coordinates | null
  isTracking: boolean
  eta: number | null  // minutos
  distanceRemaining: number | null  // km

  // Acciones
  startTracking: (tripId: string) => Promise<void>
  stopTracking: () => void
  updateDriverLocation: (location: Coordinates) => void
}
FASE 3: Backend WebSocket (Día 5-6)
3.1 Agregar Dependencia WebSocket
Archivo: d:\frontendmobile\CarPooling\backend\trips-api\go.mod
cd d:\frontendmobile\CarPooling\backend\trips-api
go get github.com/gorilla/websocket
3.2 Modificar Modelo Trip
Archivo: d:\frontendmobile\CarPooling\backend\trips-api\internal\domain\trip.go Agregar campos:
type Trip struct {
    // ... campos existentes ...

    // Nuevos campos para tracking
    CurrentLocation  *Coordinates `bson:"current_location,omitempty" json:"current_location,omitempty"`
    LocationHistory  []LocationPoint `bson:"location_history,omitempty" json:"location_history,omitempty"`
    TripProgress     *TripProgress `bson:"trip_progress,omitempty" json:"trip_progress,omitempty"`
    StartedAt        *time.Time `bson:"started_at,omitempty" json:"started_at,omitempty"`
    CompletedAt      *time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

type LocationPoint struct {
    Lat       float64   `bson:"lat" json:"lat"`
    Lng       float64   `bson:"lng" json:"lng"`
    Timestamp time.Time `bson:"timestamp" json:"timestamp"`
    Speed     float64   `bson:"speed,omitempty" json:"speed,omitempty"` // km/h
}

type TripProgress struct {
    DistanceTraveled       float64   `bson:"distance_traveled" json:"distance_traveled"` // km
    EstimatedTimeRemaining int       `bson:"estimated_time_remaining" json:"estimated_time_remaining"` // minutes
    LastUpdated            time.Time `bson:"last_updated" json:"last_updated"`
}
3.3 Crear WebSocket Hub
Archivo: d:\frontendmobile\CarPooling\backend\trips-api\internal\websocket\hub.go Responsabilidades:
Gestionar conexiones activas por trip_id
Broadcast de ubicación a todos los pasajeros de un viaje
Manejo de conexiones/desconexiones
3.4 Nuevos Endpoints
Archivo: d:\frontendmobile\CarPooling\backend\trips-api\internal\controller\trip_controller.go Agregar:
// HTTP REST
POST   /api/v1/trips/:id/start           // Iniciar viaje (driver only)
POST   /api/v1/trips/:id/location        // Actualizar ubicación (driver only)
POST   /api/v1/trips/:id/complete        // Finalizar viaje (driver only)
GET    /api/v1/trips/:id/tracking        // Obtener estado tracking

// WebSocket
GET    /api/v1/trips/:id/ws              // WebSocket connection
FASE 4: UI - Componentes de Mapa (Día 7-9)
4.1 Componente Base de Mapa
Archivo: d:\frontendmobile\src\components\map\MapView.tsx Props:
interface MapViewProps {
  region?: Region
  markers?: Marker[]
  polylines?: Polyline[]
  onRegionChange?: (region: Region) => void
  showsUserLocation?: boolean
  followsUserLocation?: boolean
}
4.2 Mapa con Ruta Estática
Archivo: d:\frontendmobile\src\components\map\TripRouteMap.tsx Funcionalidad:
Mostrar origen y destino con markers
Trazar ruta entre puntos usando Directions API
Calcular distancia y duración
Auto-zoom para mostrar ruta completa
4.3 Mapa con Tracking en Vivo
Archivo: d:\frontendmobile\src\components\map\LiveTrackingMap.tsx Funcionalidad:
Mostrar ubicación del conductor en tiempo real
Animar movimiento del marcador
Mostrar ruta completa con polyline
Mostrar ETA actualizado
Centrar cámara siguiendo al conductor
FASE 5: Pantallas de Usuario (Día 10-12)
5.1 Pantalla Tracking Conductor
Archivo: d:\frontendmobile\src\screens\DriverTrackingScreen\DriverTrackingScreen.tsx Elementos UI:
┌─────────────────────────────────┐
│  [<] Trip en Progreso           │
├─────────────────────────────────┤
│                                 │
│         MAPA FULL SCREEN        │
│   - Mi ubicación (car icon)     │
│   - Ruta trazada                │
│   - Destino marcado             │
│                                 │
├─────────────────────────────────┤
│ ETA: 25 min | 12.5 km restantes │
├─────────────────────────────────┤
│  Pasajeros: 2/4                 │
│  • Juan Pérez                   │
│  • María González               │
├─────────────────────────────────┤
│     [Finalizar Viaje]           │
└─────────────────────────────────┘
Funcionalidades:
Botón "Iniciar Viaje" (cambia status a in_progress)
Tracking GPS continuo en foreground
Envío de ubicación cada 5 segundos vía WebSocket
Notificación cuando pasajero se une
Botón "Finalizar Viaje" (cambia status a completed)
5.2 Pantalla Tracking Pasajero
Archivo: d:\frontendmobile\src\screens\TripTrackingScreen\TripTrackingScreen.tsx Elementos UI:
┌─────────────────────────────────┐
│  [<] Siguiendo a Juan (Conductor)│
├─────────────────────────────────┤
│                                 │
│         MAPA FULL SCREEN        │
│   - Conductor (car icon animado)│
│   - Ruta trazada                │
│   - Tu ubicación (opcional)     │
│                                 │
├─────────────────────────────────┤
│ Llegará en: 25 min              │
│ Distancia: 12.5 km              │
├─────────────────────────────────┤
│ Conductor: Juan Pérez ⭐ 4.8    │
│ Auto: Toyota Corolla 2020       │
│ Placa: ABC-123                  │
├─────────────────────────────────┤
│  [Contactar] [Cancelar Reserva] │
└─────────────────────────────────┘
Funcionalidades:
Conectar a WebSocket al abrir pantalla
Recibir actualizaciones de ubicación del conductor
Actualizar marcador suavemente (animación)
Recalcular ETA cuando cambia ubicación
Notificación cuando conductor inicia/finaliza viaje
FASE 6: Integración y Navegación (Día 13)
6.1 Actualizar Navigation Types
Archivo: d:\frontendmobile\src\navigation\types.ts Agregar:
export type RootStackParamList = {
  // ... existentes ...
  DriverTracking: { tripId: string }
  TripTracking: { tripId: string, bookingId: string }
}
6.2 Modificar TripDetailScreen
Archivo: d:\frontendmobile\src\screens\TripDetailScreen\TripDetailScreen.tsx Agregar:
Botón "Ver en Mapa" (navega a mapa estático con ruta)
Si trip.status === 'in_progress':
Para conductor: Botón "Continuar Viaje" → DriverTrackingScreen
Para pasajero: Auto-navegar a TripTrackingScreen
6.3 Modificar MyTripsScreen
Archivo: d:\frontendmobile\src\screens\MyTripsScreen\MyTripsScreen.tsx Agregar:
Badge "EN VIVO" para viajes in_progress
Botón rápido "Rastrear" que abre DriverTrackingScreen
FASE 7: Testing y Optimización (Día 14-15)
7.1 Testing Funcional
 Permisos de ubicación funcionan en Android e iOS
 Mapa carga correctamente con Google Maps
 WebSocket se conecta/desconecta correctamente
 Ubicación se actualiza cada 5 segundos
 Animación del marcador es suave
 ETA se recalcula correctamente
 Background tracking funciona (Android)
7.2 Testing de Escenarios
 Conductor inicia viaje → todos los pasajeros reciben notificación
 Pasajero se une a viaje en progreso → ve ubicación actual
 Conductor pierde conexión → reconexión automática
 App va a background → tracking continúa (conductor)
 Batería baja → reducir frecuencia de updates
7.3 Optimizaciones
 Lazy loading de MapView (solo cargar cuando necesario)
 Throttle de actualizaciones de ubicación
 Cache de polylines calculadas
 Manejo de errores de GPS (sin señal, precisión baja)
 Fallback si Google Maps falla
Archivos Críticos a Modificar/Crear
Frontend (React Native)
Nuevos Archivos:
src/screens/DriverTrackingScreen/DriverTrackingScreen.tsx
src/screens/TripTrackingScreen/TripTrackingScreen.tsx
src/components/map/MapView.tsx
src/components/map/TripRouteMap.tsx
src/components/map/LiveTrackingMap.tsx
src/services/locationService.ts
src/services/socketService.ts
src/services/trackingService.ts
src/context/TrackingContext.tsx
src/hooks/useLocation.ts
src/hooks/useTracking.ts
Modificar:
app.json - Agregar permisos y config Google Maps
package.json - Nuevas dependencias
src/navigation/types.ts - Nuevas rutas
src/navigation/AppNavigator.tsx - Registrar nuevas screens
src/screens/TripDetailScreen/TripDetailScreen.tsx - Botones tracking
src/screens/MyTripsScreen/MyTripsScreen.tsx - Badge "EN VIVO"
CarPooling/.env - GOOGLE_MAPS_API_KEY
Backend (Go)
Nuevos Archivos:
backend/trips-api/internal/websocket/hub.go
backend/trips-api/internal/websocket/client.go
backend/trips-api/internal/service/tracking_service.go
backend/trips-api/internal/controller/websocket_controller.go
Modificar:
backend/trips-api/internal/domain/trip.go - Campos tracking
backend/trips-api/internal/controller/trip_controller.go - Endpoints tracking
backend/trips-api/internal/service/trip_service.go - Lógica tracking
backend/trips-api/cmd/api/main.go - WebSocket routes
backend/trips-api/go.mod - Dependencia gorilla/websocket
Estimación de Esfuerzo
Fase	Tarea	Tiempo	Complejidad
1	Configuración base	2 días	⭐⭐
2	Servicios core	2 días	⭐⭐⭐
3	Backend WebSocket	2 días	⭐⭐⭐⭐
4	Componentes mapa	3 días	⭐⭐⭐
5	Pantallas usuario	3 días	⭐⭐⭐⭐
6	Integración	1 día	⭐⭐
7	Testing	2 días	⭐⭐⭐
TOTAL: 15 días (3 semanas laborables)
Costos Estimados
Desarrollo (MVP)
Google Maps API: $0 USD (dentro de límite gratuito $200/mes)
Infraestructura: $0 USD (usa backend existente)
Librerías: $0 USD (todas open source)
Producción (estimaciones)
1,000 usuarios/mes: ~$0 USD
10,000 usuarios/mes: ~$50-80 USD/mes (Google Maps)
100,000 usuarios/mes: ~$500-800 USD/mes
Riesgos y Mitigaciones
Riesgo	Probabilidad	Impacto	Mitigación
Costos Google Maps inesperados	Baja	Alto	Monitorear uso, implementar caching, límites de requests
Batería se agota rápido	Media	Alto	Optimizar frecuencia updates, usar geofencing
Pérdida de señal GPS	Alta	Medio	Manejo de errores, fallback a última ubicación conocida
WebSocket desconexiones	Media	Alto	Reconexión automática, buffer de updates
Performance en devices antiguos	Media	Medio	Optimizar renders, lazy loading, reducir polylines complexity
Próximos Pasos Inmediatos
Crear proyecto Google Cloud y obtener API Key
Instalar dependencias (expo-location, react-native-maps, socket.io-client)
Configurar permisos en app.json
Crear LocationService básico
Probar mapa estático con ruta de viaje existente