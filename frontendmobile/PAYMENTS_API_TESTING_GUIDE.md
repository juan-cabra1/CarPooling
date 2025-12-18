# Guía de Pruebas - Sistema de Pagos con MercadoPago

## 📋 Índice

1. [Configuración Inicial](#configuración-inicial)
2. [Configurar MercadoPago](#configurar-mercadopago)
3. [Levantar el Backend](#levantar-el-backend)
4. [Probar desde la App Mobile](#probar-desde-la-app-mobile)
5. [Pruebas con Postman/cURL](#pruebas-con-postmancurl)
6. [Escenarios de Prueba](#escenarios-de-prueba)
7. [Solución de Problemas](#solución-de-problemas)

---

## 🔧 Configuración Inicial

### 1. Crear archivo .env

```bash
cd CarPooling
cp .env.example .env
```

### 2. Editar .env con tus credenciales

```bash
nano .env  # o vim .env
```

Configurar las variables de MercadoPago (ver siguiente sección).

---

## 💳 Configurar MercadoPago

### Paso 1: Crear Cuenta de Desarrollador

1. Ve a https://www.mercadopago.com.ar/developers
2. Crea una cuenta o inicia sesión
3. Completa la verificación de identidad

### Paso 2: Crear Aplicación

1. Ve al [Panel de Aplicaciones](https://www.mercadopago.com.ar/developers/panel/app)
2. Click en "Crear aplicación"
3. Configura:
   - **Nombre**: "CarPooling App"
   - **Modelo de negocio**: Marketplace
   - **Tipo**: Pagos online

### Paso 3: Obtener Credenciales

En el panel de tu aplicación, copia:

```bash
# Credenciales de PRUEBA (Sandbox)
MP_ACCESS_TOKEN=TEST-1234567890-123456-abcdef1234567890-123456789
MP_PUBLIC_KEY=TEST-abcd1234-5678-90ef-ghij-klmn12345678
MP_CLIENT_ID=1234567890123456
MP_CLIENT_SECRET=ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef
```

### Paso 4: Configurar Webhooks

1. En el panel, ve a "Webhooks"
2. Agrega URL: `http://tu-ip:8005/api/v1/webhooks/mercadopago`
3. Eventos a escuchar:
   - `payment` (pagos)
   - `merchant_order` (órdenes)

### Paso 5: Configurar OAuth Redirect URI

Para vincular cuentas de conductores:

1. En "OAuth", agrega:
   - **Desarrollo**: `http://localhost:3000/seller/callback`
   - **Móvil**: `carpooling://seller/callback`

### Paso 6: Crear Usuarios de Prueba

1. Ve a "Cuentas de prueba"
2. Crea 2 usuarios:
   - **Comprador** (pasajero): Para simular pagos
   - **Vendedor** (conductor): Para simular recepción de fondos

### Tarjetas de Prueba (Sandbox)

```
✅ APROBADA:
   Número: 4509 9535 6623 3704
   CVV: 123
   Venc: 11/25

✅ APROBADA (Mastercard):
   Número: 5031 7557 3453 0604
   CVV: 123
   Venc: 11/25

❌ RECHAZADA:
   Número: 4000 0000 0000 0002
   CVV: 123
   Venc: 11/25
```

---

## 🚀 Levantar el Backend

### Opción 1: Docker Compose (Recomendado)

```bash
cd CarPooling

# Levantar todos los servicios
docker-compose up -d

# Ver logs del payments-api
docker logs -f payments-api

# Verificar que esté corriendo
curl http://localhost:8005/health
```

### Opción 2: Manual (Solo payments-api)

```bash
cd CarPooling/backend/payments-api

# Compilar
go build -o payments-api ./cmd/api

# Ejecutar
./payments-api
```

### Verificar que funciona

```bash
# Health check
curl http://localhost:8005/health

# Respuesta esperada:
# {"status":"ok","database":"connected","service":"payments-api"}
```

---

## 📱 Probar desde la App Mobile

### 1. Configurar IP del Backend

Edita `frontendmobile/src/services/api.ts`:

```typescript
const API_HOST = 'TU_IP_LOCAL' // Por ejemplo: '192.168.1.100'
```

### 2. Ejecutar la App

```bash
cd frontendmobile
npm start
```

### 3. Flujo Completo de Prueba

#### A. Como Pasajero (Hacer un Pago)

1. **Login** con un usuario
2. **Buscar viaje** en la pestaña Search
3. **Seleccionar viaje** disponible
4. **Reservar asientos** (se abre BookingModal)
5. **Hacer clic en "Pagar con MercadoPago"**
6. Se abre MercadoPago checkout
7. **Pagar con tarjeta de prueba**
8. Volver a la app → Reserva confirmada ✅

#### B. Como Conductor (Recibir Pagos)

1. **Login** con un usuario conductor
2. **Ir al Perfil** → "Mi Billetera"
3. **Vincular cuenta** de MercadoPago:
   - Click en "Vincular MercadoPago"
   - Iniciar sesión con usuario vendedor de prueba
   - Autorizar la aplicación
4. **Ver saldo** y transacciones
5. **Retirar fondos** cuando haya saldo disponible

---

## 🧪 Pruebas con Postman/cURL

### 1. Crear Preferencia de Pago

```bash
curl -X POST http://localhost:8005/api/v1/payments/preference \
  -H "Authorization: Bearer TU_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "booking_id": "booking-123",
    "trip_id": "trip-456",
    "driver_id": "driver-789",
    "seats_count": 2,
    "price_per_seat": 150000,
    "origin": "Buenos Aires",
    "destination": "Rosario",
    "departure_at": "2025-12-15T10:00:00Z",
    "payer_email": "test@example.com",
    "payer_name": "Juan Perez"
  }'
```

**Respuesta esperada:**

```json
{
  "success": true,
  "data": {
    "payment_uuid": "uuid-...",
    "preference_id": "123456789-...",
    "init_point": "https://www.mercadopago.com.ar/checkout/v1/redirect?pref_id=...",
    "amount": 300000,
    "platform_fee": 45000,
    "driver_amount": 255000,
    "expires_at": "2025-12-11T..."
  }
}
```

### 2. Consultar Wallet

```bash
curl -X GET http://localhost:8005/api/v1/wallet \
  -H "Authorization: Bearer TU_JWT_TOKEN"
```

### 3. Solicitar Retiro

```bash
curl -X POST http://localhost:8005/api/v1/withdrawals \
  -H "Authorization: Bearer TU_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100000
  }'
```

### 4. Ver Estado de Vinculación MP

```bash
curl -X GET http://localhost:8005/api/v1/seller/status \
  -H "Authorization: Bearer TU_JWT_TOKEN"
```

---

## ✅ Escenarios de Prueba

### Escenario 1: Pago Exitoso

1. Crear reserva → Generar preferencia
2. Pagar con tarjeta aprobada
3. Webhook recibe `payment.approved`
4. Booking cambia a `confirmed`
5. Wallet del conductor se actualiza (held balance)

### Escenario 2: Pago Rechazado

1. Crear reserva → Generar preferencia
2. Pagar con tarjeta rechazada
3. Webhook recibe `payment.rejected`
4. Booking queda en `pending_payment`

### Escenario 3: Liberación de Fondos Automática

1. Pago exitoso (fondos retenidos)
2. Esperar 2 horas (o configuración `AUTO_RELEASE_HOURS`)
3. Job automático libera fondos
4. Saldo pasa de `held` → `available`

### Escenario 4: Confirmación Manual de Pasajero

1. Pago exitoso (fondos retenidos)
2. Pasajero confirma viaje manualmente:
   ```bash
   curl -X POST http://localhost:8005/api/v1/payments/{uuid}/confirm \
     -H "Authorization: Bearer PASSENGER_TOKEN"
   ```
3. Fondos liberados inmediatamente

### Escenario 5: Reembolso

1. Pago exitoso
2. Cancelar reserva:
   ```bash
   curl -X POST http://localhost:8005/api/v1/payments/{uuid}/refund \
     -H "Authorization: Bearer PASSENGER_TOKEN"
   ```
3. Dinero devuelto al pasajero

### Escenario 6: Retiro de Fondos

1. Tener saldo disponible
2. Vincular cuenta MP (OAuth)
3. Solicitar retiro
4. Job procesa retiro → dinero transferido

---

## 🐛 Solución de Problemas

### Error: "Database connection failed"

```bash
# Verificar que MySQL esté corriendo
docker ps | grep mysql-payments

# Ver logs
docker logs mysql-payments

# Reiniciar
docker-compose restart mysql-payments
```

### Error: "MercadoPago authentication failed"

- Verifica que `MP_ACCESS_TOKEN` sea correcto
- Asegúrate de usar credenciales de TEST (sandbox)
- Revisa los logs: `docker logs payments-api`

### Error: "Webhook signature invalid"

- Verifica `MP_WEBHOOK_SECRET` en .env
- En desarrollo, puedes desactivar validación temporalmente

### Pago no se actualiza

```bash
# Ver si el webhook llegó
docker logs payments-api | grep webhook

# Verificar que MP tenga la URL correcta
# Panel MP → Webhooks → Verificar URL
```

### OAuth no funciona

- Verifica `MP_CLIENT_ID` y `MP_CLIENT_SECRET`
- Revisa Redirect URI en panel de MP
- Logs: `docker logs payments-api | grep oauth`

### Fondos no se liberan

```bash
# Ver logs del job auto-release
docker logs payments-api | grep "auto-release"

# Forzar liberación manual (para testing)
curl -X POST http://localhost:8005/api/v1/payments/{uuid}/release \
  -H "Authorization: Bearer DRIVER_TOKEN"
```

### Retiros no se procesan

```bash
# Ver logs del job de retiros
docker logs payments-api | grep "withdrawal"

# Verificar saldo disponible
curl -X GET http://localhost:8005/api/v1/wallet \
  -H "Authorization: Bearer DRIVER_TOKEN"
```

---

## 📊 Monitoreo

### Logs importantes

```bash
# Pagos
docker logs payments-api | grep "payment"

# Webhooks
docker logs payments-api | grep "webhook"

# Jobs background
docker logs payments-api | grep "job"

# Errores
docker logs payments-api | grep "error"
```

### Base de datos

```bash
# Conectar a MySQL
docker exec -it mysql-payments mysql -u root -p payments_db

# Ver pagos
SELECT payment_uuid, status, fund_status, amount FROM payments;

# Ver wallets
SELECT wallet_uuid, available_balance, held_balance, total_balance FROM wallets;

# Ver transacciones
SELECT * FROM wallet_transactions ORDER BY created_at DESC LIMIT 10;

# Ver retiros
SELECT * FROM withdrawals ORDER BY created_at DESC;
```

### RabbitMQ

- Web UI: http://localhost:15672
- Usuario/Password: guest/guest
- Ver eventos en las queues de payments

---

## 🎯 Checklist Final

Antes de considerar el sistema listo para producción:

- [ ] Credenciales de producción configuradas
- [ ] Webhooks funcionando correctamente
- [ ] OAuth flujo completo probado
- [ ] Pagos exitosos y rechazados
- [ ] Liberación automática de fondos
- [ ] Retiros funcionando
- [ ] Reembolsos funcionando
- [ ] Jobs en background corriendo
- [ ] Logs sin errores
- [ ] Base de datos con datos correctos
- [ ] Seguridad: HTTPS en producción
- [ ] Encryption key configurado
- [ ] Backup de base de datos configurado

---

## 📚 Recursos

- [Documentación MercadoPago](https://www.mercadopago.com.ar/developers/es/docs)
- [API Reference](https://www.mercadopago.com.ar/developers/es/reference)
- [Webhooks Guide](https://www.mercadopago.com.ar/developers/es/docs/your-integrations/notifications/webhooks)
- [OAuth Flow](https://www.mercadopago.com.ar/developers/es/docs/security/oauth)

---

**¡Sistema de Pagos Listo!** 🎉
