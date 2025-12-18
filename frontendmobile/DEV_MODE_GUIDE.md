# 🔧 Modo Desarrollo - Guía de Uso

## Descripción

El **Modo Desarrollo** permite hacer testing de toda la aplicación **sin procesar pagos reales**. Las reservas creadas en modo dev se marcan automáticamente como `confirmed` en el backend, permitiéndote probar todo el flujo de:

- Crear viajes como conductor
- Reservar asientos como pasajero
- Tracking en tiempo real
- Completar viajes
- Ver historial de reservas

Todo esto **sin necesidad de tener implementado el microservicio de pagos**.

---

## ⚡ Inicio Rápido

### 1. Activar Modo Dev

En la app móvil:
1. Abre la pestaña **"Perfil"**
2. En la sección **"Acciones"**, activa el switch **"Modo Desarrollo"** (color naranja)
3. Verás una alerta confirmando la activación

### 2. Hacer una Reserva de Prueba

1. Busca un viaje disponible
2. Haz clic en **"Reservar Asientos"**
3. Notarás:
   - Banner amarillo: "MODO DESARROLLO: La reserva se creará sin procesar pagos"
   - Botón cambia a: **"Reservar (Sin Pago)"**
4. Completa la reserva
5. ✅ La reserva se crea instantáneamente con estado `confirmed`

### 3. Probar Tracking en Vivo

1. Con modo dev activado, crea un viaje como conductor
2. Reserva asientos en ese viaje (puedes usar otra cuenta o mismo usuario)
3. Como conductor, inicia el viaje desde "Mis Viajes"
4. Como pasajero, verás el tracking en tiempo real funcionando
5. No necesitas procesar ningún pago

---

## 🛠️ Cómo Funciona

### Frontend (React Native)

Cuando el modo dev está activado:

```typescript
// BookingModal.tsx - línea 174
const booking = await bookingsService.createBooking({
  trip_id: tripId,
  passenger_id: user.id,
  seats_reserved: seatsRequested,
  dev_mode: true, // 👈 Envía este flag al backend
})

// 🔧 DEV MODE: Skip payment flow
// No se crea preferencia de pago
// No se abre MercadoPago
// No se hace polling de estado de pago
```

### Backend (Go - Bookings API)

El backend recibe el flag `dev_mode` y actúa en consecuencia:

```go
// booking_service.go - línea 96
initialStatus := dao.BookingStatusPending
if req.DevMode {
  initialStatus = dao.BookingStatusConfirmed
  log.Warn().Msg("🔧 DEV MODE: Creating booking as CONFIRMED")
}

booking := &dao.Booking{
  // ...
  Status: initialStatus, // ✅ "confirmed" en modo dev
}

// Línea 137: Skip publicación de eventos si dev_mode
if req.DevMode {
  log.Info().Msg("🔧 DEV MODE: Skipping reservation.created event")
} else {
  // Publica evento a RabbitMQ para validación asíncrona
  s.publisher.PublishReservationCreated(...)
}
```

---

## 📊 Comparación: Modo Dev vs Producción

| Característica | Modo Desarrollo | Modo Producción |
|---------------|-----------------|-----------------|
| **Estado inicial de booking** | `confirmed` | `pending` (espera pago) |
| **Procesamiento de pago** | ❌ Saltado | ✅ MercadoPago |
| **Eventos RabbitMQ** | ❌ No se publican | ✅ Validación asíncrona |
| **Confirmación** | ⚡ Inmediata | ⏳ Después del pago |
| **Aparece en "Mis Reservas"** | ✅ Sí | ✅ Sí (después del pago) |
| **Tracking funciona** | ✅ Sí | ✅ Sí |
| **Indicador visual** | 🟡 Banner amarillo | 🟢 Banner verde MP |

---

## 🎯 Casos de Uso

### Testing de Features Completas

```bash
# Escenario: Testing del flujo completo sin pagos
1. Activar modo dev
2. Usuario A crea viaje Buenos Aires → Córdoba
3. Usuario B reserva 2 asientos (✅ confirmado instantáneamente)
4. Usuario A inicia el viaje
5. Usuario B ve tracking en vivo
6. Usuario A completa el viaje
7. Ambos ven el viaje en historial
```

### Demo para Cliente/Stakeholder

```bash
# Mostrar todas las funcionalidades sin esperar pagos reales
1. Activar modo dev antes de la demo
2. Crear viajes, reservar, hacer tracking, etc.
3. Todo funciona sin MercadoPago
4. Desactivar después de la demo
```

### Desarrollo de Nuevas Features

```bash
# Ejemplo: Implementando sistema de ratings
1. Modo dev ON
2. Crear viaje → Reservar → Completar
3. Implementar feature de calificaciones
4. Probar con viajes ya completados
5. No necesitas procesar pagos para llegar al estado "completed"
```

---

## 📁 Archivos Modificados

### Frontend

```
src/
├── config/
│   └── dev.ts                          # ✨ Nuevo - Gestión modo dev
├── components/modals/
│   ├── BookingModal.tsx                # ✏️ Modificado - Soporte dev_mode
│   └── BookingModal.styles.ts          # ✏️ Modificado - Estilos banner
├── screens/ProfileScreen/
│   └── ProfileScreen.tsx               # ✏️ Modificado - Toggle UI
└── types/
    └── booking.ts                      # ✏️ Modificado - dev_mode field
```

### Backend

```
backend/bookings-api/internal/
├── domain/
│   └── booking.go                      # ✏️ Modificado - DevMode field
└── service/
    └── booking_service.go              # ✏️ Modificado - Lógica dev mode
```

---

## 🔍 Logs y Debugging

### Frontend (Console/DevTools)

```javascript
// Cuando modo dev está activado
🔧 DEV MODE: Skipping payment flow
✅ Dev mode enabled - payments will be skipped

// Cuando creas una reserva
Booking created: { id: "...", status: "confirmed", ... }
```

### Backend (Bookings API Logs)

```bash
# Al crear booking en modo dev
level=warn msg="🔧 DEV MODE: Creating booking as CONFIRMED (skipping payment flow)" trip_id="..." passenger_id=123

level=info msg="✅ Booking created successfully" booking_id="..." status="confirmed" dev_mode=true

level=info msg="🔧 DEV MODE: Skipping reservation.created event (no async validation needed)" booking_id="..."
```

---

## ⚠️ Limitaciones y Advertencias

### NO usar en Producción

```diff
- ❌ NUNCA dejar modo dev activado en producción
- ❌ Las reservas dev no generan ingresos reales
- ❌ No se validan asientos disponibles en modo dev
- ❌ No se emiten eventos a otros microservicios
```

### Consideraciones Técnicas

1. **Total Price = 0**: En modo dev, el `total_price` se queda en 0 porque no hay procesamiento de pago
2. **Sin validación de RabbitMQ**: No se publican eventos, por lo que trips-api no actualiza asientos
3. **AsyncStorage**: El estado del modo dev se guarda localmente y se pierde al desinstalar la app
4. **Backend no valida disponibilidad**: En modo dev, podrías reservar más asientos de los disponibles

---

## 🚀 Próximos Pasos

### Cuando implementes el Payments API

1. **Mantén el modo dev** para testing continuo
2. El modo producción funcionará con MercadoPago
3. Podrás alternar entre ambos para diferentes tipos de testing

### Features Futuras del Modo Dev

```typescript
// Posibles mejoras
- Simular fallos de pago para testing
- Fast-forward de tiempos de espera
- Panel de admin para ver/limpiar reservas dev
- Logs más detallados de requests
- Modo "slow network" para testing de UX
```

---

## 🐛 Troubleshooting

### "El modo dev no se activa"

```bash
Solución:
1. Verifica logs en console
2. Asegúrate de tener permisos de AsyncStorage
3. Intenta: Settings > Clear App Data
4. Reinstala la app si persiste
```

### "Las reservas siguen pidiendo pago"

```bash
Solución:
1. Verifica que el switch esté en ON (naranja)
2. Cierra y reabre el BookingModal
3. Verifica console: debe decir "🔧 DEV MODE: Skipping payment flow"
4. Si usa expo-go, recarga la app
```

### "Las reservas no aparecen en Mis Reservas"

```bash
Posibles causas:
1. El backend no está corriendo
2. Filtros en MyBookingsScreen ocultan estado "confirmed"
3. Error de red al crear la reserva

Debug:
1. Verifica logs del backend
2. Revisa respuesta del POST /api/v1/bookings
3. Verifica que bookingsService.getMyBookings() incluya confirmed
```

### "El backend no reconoce dev_mode"

```bash
Solución:
1. Asegúrate de tener la última versión del backend
2. Verifica que domain/booking.go tenga el campo DevMode
3. Reinicia el servidor Go
4. Verifica logs: debe decir "🔧 DEV MODE: Creating booking as CONFIRMED"
```

---

## 📞 Soporte

Si tienes problemas con el modo dev:

1. Revisa los logs del frontend (console) y backend
2. Verifica que todos los archivos estén actualizados
3. Asegúrate de que el backend esté corriendo
4. Prueba desactivar/activar el modo dev

---

**¡Happy Testing!** 🎉

---

## 📝 Changelog

### v1.0.0 - 2024-12-14
- ✨ Implementación inicial del modo desarrollo
- ✅ Frontend: Toggle en ProfileScreen
- ✅ Frontend: Banner y lógica en BookingModal
- ✅ Backend: Soporte para dev_mode flag
- ✅ Backend: Creación de bookings como "confirmed" en modo dev
- ✅ Backend: Skip de eventos RabbitMQ en modo dev
