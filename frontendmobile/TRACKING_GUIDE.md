# 🗺️ Guía de Tracking en Tiempo Real

## Cómo Ver el Mapa en Vivo del Vehículo

Ya tienes todo implementado para ver el tracking en tiempo real. Aquí te explico cómo acceder a las diferentes vistas:

---

## 📱 3 Formas de Ver el Mapa

### 1️⃣ **Ver Ruta Estática** (Antes de iniciar el viaje)

**Desde:** Cualquier viaje en TripDetailScreen

**Cómo acceder:**
1. Ve a **"Buscar"** o **"Mis Viajes"**
2. Toca cualquier viaje
3. Presiona el botón **"Ver en Mapa"**

**Qué verás:**
- 📍 Origen y destino marcados
- 🛣️ Ruta completa trazada usando Google Directions API
- 📏 Distancia y duración estimada
- 🗺️ Mapa estático (sin movimiento en tiempo real)

**Pantalla:** `TripRouteMapScreen`

---

### 2️⃣ **Vista de Conductor** (Durante el viaje)

**Para:** El conductor del viaje

**Cómo acceder:**
1. Ve a **"Mis Viajes"**
2. Selecciona tu viaje publicado
3. **Inicia el viaje** (cambia status a `in_progress`)
4. Presiona **"Continuar Viaje"**

**Qué verás:**
```
┌─────────────────────────────────┐
│  [<] Trip en Progreso           │
├─────────────────────────────────┤
│                                 │
│     🗺️ MAPA PANTALLA COMPLETA  │
│                                 │
│   🚗 Tu ubicación (icono auto)  │
│   📍 Destino marcado            │
│   🛣️ Ruta trazada               │
│                                 │
├─────────────────────────────────┤
│ ⏱️ ETA: 25 min | 📏 12.5 km    │
├─────────────────────────────────┤
│  👥 Pasajeros: 2/4              │
│  • Juan Pérez                   │
│  • María González               │
├─────────────────────────────────┤
│     [🏁 Finalizar Viaje]        │
└─────────────────────────────────┘
```

**Funcionalidades:**
- ✅ Tu ubicación GPS se actualiza cada 5 segundos
- ✅ Se envía por WebSocket a todos los pasajeros
- ✅ ETA se recalcula automáticamente
- ✅ Tracking continúa en background (Android)
- ✅ Lista de pasajeros en tiempo real
- ✅ Botón para finalizar el viaje

**Pantalla:** `DriverTrackingScreen`
**Archivo:** `src/screens/DriverTrackingScreen/DriverTrackingScreen.tsx`

---

### 3️⃣ **Vista de Pasajero** (Siguiendo al conductor)

**Para:** Los pasajeros que reservaron asientos

**Cómo acceder:**
1. Ve a **"Mis Reservas"**
2. Selecciona una reserva con viaje `in_progress`
3. Presiona **"Rastrear Viaje"**

**Qué verás:**
```
┌─────────────────────────────────┐
│  [<] Siguiendo a Juan (Conductor)│
├─────────────────────────────────┤
│                                 │
│     🗺️ MAPA PANTALLA COMPLETA  │
│                                 │
│   🚗 Conductor (animado)        │
│   📍 Tu ubicación (opcional)    │
│   🛣️ Ruta trazada               │
│                                 │
├─────────────────────────────────┤
│ ⏱️ Llegará en: 25 min           │
│ 📏 Distancia: 12.5 km           │
├─────────────────────────────────┤
│ 👤 Conductor: Juan Pérez ⭐ 4.8 │
│ 🚙 Auto: Toyota Corolla 2020    │
│ 🚗 Placa: ABC-123               │
├─────────────────────────────────┤
│  [📞 Contactar] [❌ Cancelar]   │
└─────────────────────────────────┘
```

**Funcionalidades:**
- ✅ Recibes la ubicación del conductor en tiempo real vía WebSocket
- ✅ El marcador del auto se anima suavemente al moverse
- ✅ ETA se recalcula automáticamente
- ✅ Notificación cuando el conductor inicia/finaliza el viaje
- ✅ Botones para contactar o cancelar reserva

**Pantalla:** `TripTrackingScreen`
**Archivo:** `src/screens/TripTrackingScreen/TripTrackingScreen.tsx`

---

## 🎯 Flujo Completo: Paso a Paso

### Escenario: Tracking en Vivo con Modo Dev

#### **PASO 1: Preparación**

```bash
1. Activa modo dev en "Perfil" (para saltar pagos)
2. Asegúrate de tener permisos de ubicación activados
3. Ten el backend corriendo (trips-api con WebSocket)
```

#### **PASO 2: Crear Viaje (Como Conductor)**

```bash
1. Ve a "Buscar" → Botón "+"
2. Crea un viaje:
   - Origen: Buenos Aires
   - Destino: Córdoba
   - Fecha: Hoy
   - Asientos: 4
   - Precio: $5000
3. Publica el viaje
```

#### **PASO 3: Reservar Asientos (Como Pasajero)**

Opción A - Mismo usuario:
```bash
1. Busca el viaje que acabas de crear
2. Reserva 2 asientos
3. Con modo dev: se confirma instantáneamente
```

Opción B - Otro usuario:
```bash
1. Inicia sesión con otra cuenta
2. Busca el viaje
3. Reserva asientos
```

#### **PASO 4: Iniciar el Viaje (Conductor)**

```bash
1. Ve a "Mis Viajes"
2. Toca tu viaje
3. Verás un botón "Iniciar Viaje" (si hay reservas confirmadas)
4. Presiona "Iniciar Viaje"
   → El viaje cambia a status: in_progress
5. Presiona "Continuar Viaje"
   → Se abre DriverTrackingScreen
```

**Lo que pasa internamente:**
```javascript
// TripDetailScreen.tsx - línea 366
onPress={() => navigation.navigate('DriverTracking', { tripId: trip.trip_id.toString() })}

// DriverTrackingScreen se monta y:
1. Conecta al WebSocket del viaje
2. Inicia tracking GPS (cada 5 segundos)
3. Emite ubicación por WebSocket a todos los pasajeros
4. Muestra el mapa con tu ubicación en vivo
```

#### **PASO 5: Ver Tracking (Pasajero)**

```bash
1. Ve a "Mis Reservas"
2. Toca la reserva del viaje en progreso
3. Presiona "Rastrear Viaje"
   → Se abre TripTrackingScreen
```

**Lo que pasa internamente:**
```javascript
// TripDetailScreen.tsx - línea 373
onPress={() => navigation.navigate('TripTracking', { tripId: trip.trip_id.toString() })}

// TripTrackingScreen se monta y:
1. Conecta al WebSocket del viaje
2. Escucha eventos de ubicación del conductor
3. Actualiza el marcador del auto en tiempo real
4. Recalcula ETA automáticamente
5. Anima el movimiento del auto suavemente
```

---

## 🔌 Arquitectura de WebSocket

### Cómo funciona el tracking en tiempo real:

```
┌──────────────┐          ┌──────────────┐          ┌──────────────┐
│   CONDUCTOR  │          │   BACKEND    │          │  PASAJERO 1  │
│ (Dispositivo)│          │  WebSocket   │          │ (Dispositivo)│
└──────┬───────┘          └──────┬───────┘          └──────┬───────┘
       │                         │                         │
       │ 1. Conecta al WS       │                         │
       ├────────────────────────>│                         │
       │                         │                         │
       │                         │ 2. Conecta al WS        │
       │                         │<────────────────────────┤
       │                         │                         │
       │ 3. join_trip            │                         │
       ├────────────────────────>│                         │
       │                         │                         │
       │                         │ 4. join_trip            │
       │                         │<────────────────────────┤
       │                         │                         │
       │ 5. location_update      │                         │
       │    (lat, lng, speed)    │                         │
       ├────────────────────────>│                         │
       │                         │                         │
       │                         │ 6. broadcast location   │
       │                         ├────────────────────────>│
       │                         │                         │
       │                         │ 7. Actualiza mapa       │
       │                         │                         ├─┐
       │                         │                         │ │ Anima
       │                         │                         │<┘ marcador
       │                         │                         │
```

### Eventos WebSocket:

**Emitidos por el conductor:**
```javascript
// socketService.ts
socket.emit('join_trip', { tripId: '...' })
socket.emit('location_update', {
  tripId: '...',
  coordinates: { lat: -34.603722, lng: -58.381592 },
  speed: 60, // km/h
  heading: 180, // degrees
  timestamp: Date.now()
})
socket.emit('trip_complete', { tripId: '...' })
```

**Recibidos por los pasajeros:**
```javascript
// socketService.ts
socket.on('location_update', (data) => {
  // data.coordinates, data.speed, data.heading
  // Actualizar posición del marcador en el mapa
})

socket.on('trip_started', (data) => {
  // Notificar que el viaje comenzó
})

socket.on('trip_completed', (data) => {
  // Viaje finalizado
})
```

---

## 🛠️ Componentes Clave

### 1. **LiveTrackingMap** (Componente de Mapa)

**Archivo:** `src/components/map/LiveTrackingMap.tsx`

**Funcionalidad:**
- Muestra el mapa en tiempo real
- Obtiene datos del `TrackingContext`
- Dibuja la ruta con Google Directions API
- Anima el marcador del conductor
- Centra la cámara siguiendo el movimiento

### 2. **TrackingContext** (Estado Global)

**Archivo:** `src/context/TrackingContext.tsx`

**Provee:**
```typescript
{
  activeTrip: Trip | null,
  driverLocation: Coordinates | null,
  isTracking: boolean,
  isDriver: boolean,
  eta: number | null,
  distanceRemaining: number | null,

  startTracking: (trip: Trip, asDriver: boolean) => void,
  stopTracking: () => void,
  updateDriverLocation: (location: Coordinates) => void,
}
```

### 3. **LocationService** (GPS)

**Archivo:** `src/services/locationService.ts`

**Métodos:**
```typescript
// Obtener ubicación actual
await locationService.getCurrentLocation()

// Iniciar tracking continuo
await locationService.startTracking((locationUpdate) => {
  // Se ejecuta cada 5 segundos con nueva ubicación
})

// Detener tracking
await locationService.stopTracking()
```

### 4. **SocketService** (WebSocket)

**Archivo:** `src/services/socketService.ts`

**Métodos:**
```typescript
// Conectar
socketService.connect(userId, token)

// Unirse al viaje
socketService.joinTripTracking(tripId)

// Emitir ubicación (conductor)
socketService.emitLocationUpdate(tripId, locationData)

// Escuchar ubicación (pasajero)
socketService.onLocationUpdate((data) => {
  // Actualizar mapa
})
```

---

## 📊 Estados del Viaje

El tracking solo funciona cuando el viaje está en estado `in_progress`:

```
draft → published → in_progress → completed
                         ↑
                    (Tracking activo)
```

### Cambiar estado a `in_progress`:

**Backend - trips-api:**
```bash
POST /api/v1/trips/:id/start
```

Este endpoint:
- Cambia el status del viaje a `in_progress`
- Establece `started_at` timestamp
- Emite evento WebSocket `trip_started`

---

## 🐛 Troubleshooting

### "No veo el botón de tracking"

**Causa:** El viaje no está en estado `in_progress`

**Solución:**
1. Asegúrate de que el viaje tenga al menos una reserva confirmada
2. Como conductor, presiona "Iniciar Viaje" en TripDetailScreen
3. Verifica que el backend haya cambiado el status

**Verificar en backend:**
```bash
# Ver status del viaje
GET http://localhost:8002/api/v1/trips/:id
```

---

### "El mapa no muestra mi ubicación"

**Causa:** Permisos de ubicación no otorgados

**Solución:**
1. Ve a Configuración del dispositivo → Apps → CarPooling
2. Permisos → Ubicación → Permitir siempre (para background)
3. Reinicia la app
4. Cuando abras DriverTrackingScreen, debe pedir permisos

**Verificar en código:**
```typescript
// locationService.ts:36
const { status } = await Location.requestForegroundPermissionsAsync()
```

---

### "La ubicación no se actualiza en el pasajero"

**Causas posibles:**

1. **WebSocket no conectado:**
```bash
# Verificar logs del backend
level=info msg="Client 123 joined trip 507f1f77..."
```

2. **Backend no corriendo:**
```bash
# Iniciar trips-api
cd CarPooling/backend/trips-api
go run cmd/api/main.go
```

3. **IP incorrecta en socketService:**
```typescript
// src/services/socketService.ts:15
const API_HOST = '181.85.173.171' // ← Verifica que sea la IP correcta
```

---

### "El marcador del auto no se mueve suavemente"

**Causa:** Animación no configurada

**Solución:**
El componente `LiveTrackingMap` debe usar `animateMarkerToCoordinate`:

```typescript
// LiveTrackingMap.tsx
driverMarkerRef.current?.animateMarkerToCoordinate(
  newCoordinate,
  500 // duración en ms
)
```

---

## 📱 Testing en Dispositivo Real

### Recomendaciones:

1. **Usa un dispositivo físico** (no emulador) para GPS real
2. **Sal a la calle** para ver el tracking con movimiento real
3. **Dos dispositivos:** Uno como conductor, otro como pasajero
4. **Modo dev activado:** Para saltar pagos y probar rápido

### Flujo de testing:

```bash
Dispositivo 1 (Conductor):
1. Crear viaje
2. Iniciar viaje
3. Caminar/conducir
4. Ver tu ubicación actualizándose

Dispositivo 2 (Pasajero):
1. Reservar asiento en el viaje
2. Abrir tracking
3. Ver al conductor moverse en tiempo real
```

---

## 🎉 ¡Listo!

Ahora sabes cómo ver el mapa en vivo del vehículo. El sistema completo está implementado:

✅ GPS tracking con LocationService
✅ WebSocket en tiempo real con SocketService
✅ Pantallas de conductor y pasajero
✅ Componente LiveTrackingMap con animaciones
✅ Contexto global TrackingContext
✅ Integración con navegación

**Próximo paso:** Sal a probar en un dispositivo real y verás cómo funciona el tracking en vivo! 🚗📍
