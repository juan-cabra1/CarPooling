# 🔍 Diagnóstico: Por qué no obtienes viajes en las búsquedas

## ❌ Problemas Encontrados

### 1. **Solr NO tiene el core configurado**
```json
// Health check del search-api:
{
  "status": "degraded",
  "services": {
    "solr": {
      "status": "unhealthy",
      "message": "Solr client not initialized"
    }
  }
}
```

**Causa:** El core `carpooling_trips` no existe en Solr

### 2. **Búsquedas retornan 0 resultados**
```json
{
  "data": {
    "total": 0,
    "trips": []
  }
}
```

**Causa:** Sin Solr funcionando, el backend intenta usar MongoDB como fallback, pero probablemente no hay viajes indexados allí tampoco.

---

## ✅ SOLUCIÓN PASO A PASO

### Paso 1: Verificar cómo está corriendo Solr

¿Estás usando Docker o Solr local?

#### Si usas **Docker Compose** (RECOMENDADO):

1. Detén Solr si está corriendo:
   ```bash
   docker-compose down solr
   ```

2. Actualiza tu `docker-compose.yml` para crear el core automáticamente:

   ```yaml
   version: '3.8'

   services:
     solr:
       image: solr:9.4
       container_name: carpooling-solr
       ports:
         - "8983:8983"
       volumes:
         - solr_data:/var/solr
       environment:
         - SOLR_HEAP=512m
       command:
         - solr-precreate
         - carpooling_trips  # Esto crea el core automáticamente

   volumes:
     solr_data:
   ```

3. Inicia Solr:
   ```bash
   docker-compose up -d solr
   ```

4. Verifica que el core exista:
   ```bash
   curl http://localhost:8983/solr/admin/cores?action=STATUS | jq
   ```

   Deberías ver:
   ```json
   {
     "status": {
       "carpooling_trips": {
         "name": "carpooling_trips",
         ...
       }
     }
   }
   ```

#### Si usas **Solr local** (sin Docker):

1. Navega al directorio de instalación de Solr:
   ```bash
   cd /path/to/solr-9.x.x
   ```

2. Crea el core:
   ```bash
   bin/solr create -c carpooling_trips
   ```

3. Verifica:
   ```bash
   curl http://localhost:8983/solr/admin/cores?action=STATUS | jq
   ```

---

### Paso 2: Reiniciar el Search-API

Una vez que Solr tenga el core configurado:

```bash
cd backend/search-api
go run cmd/api/main.go
```

**Logs esperados:**

```
INFO  Starting search-api
INFO  Configuration loaded successfully
INFO  MongoDB indexes created successfully
INFO  Connected to Apache Solr successfully  ✅ <- Esto es lo importante
INFO  Connected to Memcached successfully
INFO  RabbitMQ consumer started in background
INFO  search-api server listening port=8004
```

---

### Paso 3: Verificar el Health Check

```bash
curl http://localhost:8004/health | jq
```

**Respuesta esperada:**

```json
{
  "status": "ok",  ✅ <- Debe ser "ok", no "degraded"
  "service": "search-api",
  "port": "8004",
  "services": {
    "mongodb": {
      "status": "healthy",
      "message": "Connected"
    },
    "solr": {
      "status": "healthy",  ✅ <- Esto debe ser "healthy"
      "message": "Connected"
    },
    "memcached": {
      "status": "healthy",
      "message": "Connected"
    }
  }
}
```

---

### Paso 4: Crear Viajes de Prueba

Ahora necesitas viajes en la base de datos. Tienes 2 opciones:

#### Opción A: Crear viajes desde el frontend/API

1. **Inicia el trips-api** (puerto 8002)
   ```bash
   cd backend/trips-api
   go run cmd/api/main.go
   ```

2. **Inicia el users-api** (puerto 8001)
   ```bash
   cd backend/users-api
   go run cmd/api/main.go
   ```

3. **Crea un viaje via API:**
   ```bash
   curl -X POST http://localhost:8002/trips \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d '{
       "origin": {
         "city": "Bogotá",
         "province": "Cundinamarca",
         "address": "Calle 100",
         "coordinates": [-74.0721, 4.7110]
       },
       "destination": {
         "city": "Medellín",
         "province": "Antioquia",
         "address": "Calle 50",
         "coordinates": [-75.5636, 6.2442]
       },
       "departure_datetime": "2025-11-20T08:00:00Z",
       "estimated_arrival_datetime": "2025-11-20T16:00:00Z",
       "price_per_seat": 45000,
       "total_seats": 4,
       "car": {
         "brand": "Toyota",
         "model": "Corolla",
         "year": 2020,
         "color": "Blanco",
         "license_plate": "ABC123"
       },
       "preferences": {
         "pets_allowed": false,
         "smoking_allowed": false,
         "music_allowed": true
       },
       "description": "Viaje cómodo y seguro"
     }'
   ```

4. **El flujo automático:**
   - trips-api publica evento `trip.created` en RabbitMQ
   - search-api consumer recibe el evento
   - Obtiene los datos del viaje (trips-api)
   - Obtiene los datos del conductor (users-api)
   - Crea SearchTrip denormalizado en MongoDB
   - **Indexa en Solr automáticamente** ✅

#### Opción B: Insertar datos directamente en MongoDB (para testing)

Si solo quieres probar rápido:

```bash
mongosh
```

```javascript
use carpooling_search

// Insertar un viaje de prueba
db.trips.insertOne({
  trip_id: "67890abcdef",
  driver_id: 1,
  driver: {
    id: 1,
    name: "Juan Pérez",
    email: "juan@example.com",
    rating: 4.8,
    total_trips: 25
  },
  origin: {
    city: "Bogotá",
    province: "Cundinamarca",
    address: "Calle 100",
    coordinates: {
      type: "Point",
      coordinates: [-74.0721, 4.7110]
    }
  },
  destination: {
    city: "Medellín",
    province: "Antioquia",
    address: "Calle 50",
    coordinates: {
      type: "Point",
      coordinates: [-75.5636, 6.2442]
    }
  },
  departure_datetime: new Date("2025-11-20T08:00:00Z"),
  estimated_arrival_datetime: new Date("2025-11-20T16:00:00Z"),
  price_per_seat: 45000,
  total_seats: 4,
  available_seats: 4,  // IMPORTANTE: Debe ser > 0
  car: {
    brand: "Toyota",
    model: "Corolla",
    year: 2020,
    color: "Blanco",
    license_plate: "ABC123"
  },
  preferences: {
    pets_allowed: false,
    smoking_allowed: false,
    music_allowed: true
  },
  status: "published",  // IMPORTANTE: Debe ser "published"
  description: "Viaje cómodo y seguro",
  search_text: "Bogotá Cundinamarca Medellín Antioquia Juan Pérez Toyota Corolla",
  popularity_score: 0,
  created_at: new Date(),
  updated_at: new Date()
})

// Verificar que se creó
db.trips.countDocuments({status: "published", available_seats: {$gt: 0}})
// Debe retornar 1
```

Luego, indexa manualmente en Solr:

```bash
# Obtener el trip_id del viaje
TRIP_ID="67890abcdef"

# Indexar en Solr (el search-api hará esto automáticamente via RabbitMQ)
curl -X POST "http://localhost:8983/solr/carpooling_trips/update?commit=true" \
  -H "Content-Type: application/json" \
  -d '[{
    "id": "'$TRIP_ID'",
    "origin_city": "Bogotá",
    "destination_city": "Medellín",
    "driver_name": "Juan Pérez",
    "price_per_seat": 45000,
    "available_seats": 4,
    "status": "published",
    "search_text": "Bogotá Cundinamarca Medellín Antioquia Juan Pérez"
  }]'
```

---

### Paso 5: Probar la Búsqueda

#### Desde la API directamente:

```bash
# Búsqueda simple (todos los viajes)
curl "http://localhost:8004/api/v1/search/trips?page=1&limit=10" | jq

# Búsqueda por ciudad
curl "http://localhost:8004/api/v1/search/trips?origin_city=Bogotá&destination_city=Medellín" | jq

# Búsqueda con texto libre
curl "http://localhost:8004/api/v1/search/trips?q=Toyota" | jq
```

#### Desde el Frontend:

1. Inicia el frontend:
   ```bash
   cd frontend
   npm run dev
   ```

2. Abre http://localhost:3000

3. Busca viajes:
   - Origen: Bogotá
   - Destino: Medellín

---

## 📊 Logs para Monitorear

### En el Search-API:

```bash
cd backend/search-api
go run cmd/api/main.go
```

**Busca estos logs:**

```
# Al iniciar:
INFO  Connected to Apache Solr successfully

# Al recibir una búsqueda:
INFO  HTTP request completed method=GET url=/api/v1/search/trips status_code=200 duration_ms=45

# Si hay un evento de RabbitMQ:
INFO  Processing trip.created event event_id=...
INFO  Trip created in MongoDB successfully trip_id=...
INFO  Trip indexed in Solr successfully trip_id=...
```

### En el Frontend (DevTools):

1. Abre F12 → Network
2. Filtra por "search"
3. Haz una búsqueda
4. Verás la request a `/api/search/trips`
5. Revisa la respuesta para ver si hay trips

---

## 🔍 Verificación Final

```bash
# 1. Verificar Solr tiene el core
curl http://localhost:8983/solr/admin/cores?action=STATUS | jq '.status | keys'
# Debe incluir "carpooling_trips"

# 2. Verificar Solr tiene documentos
curl "http://localhost:8983/solr/carpooling_trips/select?q=*:*&rows=0" | jq '.response.numFound'
# Debe ser > 0

# 3. Verificar MongoDB tiene viajes
mongosh --eval "db.getSiblingDB('carpooling_search').trips.countDocuments({status: 'published', available_seats: {\$gt: 0}})"
# Debe ser > 0

# 4. Verificar health check
curl http://localhost:8004/health | jq '.status'
# Debe ser "ok"

# 5. Probar búsqueda
curl "http://localhost:8004/api/v1/search/trips" | jq '.data.total'
# Debe ser > 0
```

---

## ⚠️ Problemas Comunes

### "Solr client not initialized"
**Solución:** El core `carpooling_trips` no existe. Sigue el Paso 1.

### "total": 0, pero hay viajes en MongoDB
**Posibles causas:**
1. Los viajes tienen `status != "published"` → Actualiza el status
2. Los viajes tienen `available_seats = 0` → Actualiza available_seats
3. Solr no tiene documentos → Reindexar

### Viajes en MongoDB pero no en Solr
**Causa:** Los viajes se crearon antes del consumer
**Solución:** Volver a publicar los eventos `trip.created` en RabbitMQ

### Frontend retorna []
**Causa:** El proxy de Vite no está funcionando
**Verificar:**
```bash
curl "http://localhost:3000/api/search/trips"
# Debe redirigir a localhost:8004
```

---

## 📝 Resumen de Comandos

```bash
# 1. Crear core en Solr (Docker)
docker-compose up -d solr

# 2. Verificar core
curl http://localhost:8983/solr/admin/cores?action=STATUS | jq '.status | keys'

# 3. Iniciar search-api
cd backend/search-api && go run cmd/api/main.go

# 4. Verificar health
curl http://localhost:8004/health | jq

# 5. Crear viaje de prueba (ver Paso 4)

# 6. Probar búsqueda
curl "http://localhost:8004/api/v1/search/trips" | jq

# 7. Iniciar frontend
cd frontend && npm run dev
```

---

## 🎯 Checklist

- [ ] Solr corriendo en localhost:8983
- [ ] Core `carpooling_trips` creado
- [ ] Search-API conectado a Solr (health check = "ok")
- [ ] MongoDB tiene viajes con status="published" y available_seats > 0
- [ ] Solr tiene documentos indexados
- [ ] Búsqueda API retorna trips
- [ ] Frontend muestra resultados

---

Si sigues estos pasos, las búsquedas deberían funcionar! 🚀
