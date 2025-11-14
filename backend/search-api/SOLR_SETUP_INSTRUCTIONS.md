# 🔍 Configuración de Apache Solr para Search-API

## Problema Detectado
El core `carpooling_trips` no existe en Solr, por eso las búsquedas retornan 0 resultados.

## Solución Rápida

### Opción 1: Crear el core desde la línea de comandos

```bash
# Navega al directorio de Solr (ajusta la ruta según tu instalación)
cd /path/to/solr

# Crea el core 'carpooling_trips'
bin/solr create -c carpooling_trips
```

### Opción 2: Crear el core desde el Admin UI

1. Abre http://localhost:8983/solr/#/
2. Ve a "Core Admin" en el menú lateral
3. Haz clic en "Add Core"
4. Nombre del core: `carpooling_trips`
5. Haz clic en "Add Core"

### Opción 3: Usar Docker Compose (Recomendado)

Si estás usando Docker, actualiza tu `docker-compose.yml`:

```yaml
version: '3.8'

services:
  solr:
    image: solr:9.4
    ports:
      - "8983:8983"
    command:
      - solr-precreate
      - carpooling_trips
    volumes:
      - solr_data:/var/solr
    environment:
      - SOLR_HEAP=512m

volumes:
  solr_data:
```

Luego ejecuta:

```bash
docker-compose up -d solr
```

## Verificación

Después de crear el core, verifica que existe:

```bash
curl http://localhost:8983/solr/admin/cores?action=STATUS
```

Deberías ver algo como:

```json
{
  "status": {
    "carpooling_trips": {
      "name": "carpooling_trips",
      "instanceDir": "...",
      ...
    }
  }
}
```

## Schema de Solr (Opcional - Configuración Avanzada)

Si quieres optimizar los campos, puedes configurar el schema:

```bash
# Añadir campos específicos
curl -X POST -H 'Content-type:application/json' \
  'http://localhost:8983/solr/carpooling_trips/schema' -d '{
  "add-field": [
    {"name":"trip_id", "type":"string", "stored":true, "indexed":true},
    {"name":"origin_city", "type":"text_general", "stored":true, "indexed":true},
    {"name":"destination_city", "type":"text_general", "stored":true, "indexed":true},
    {"name":"driver_name", "type":"text_general", "stored":true, "indexed":true},
    {"name":"price_per_seat", "type":"pfloat", "stored":true, "indexed":true},
    {"name":"available_seats", "type":"pint", "stored":true, "indexed":true},
    {"name":"status", "type":"string", "stored":true, "indexed":true},
    {"name":"departure_datetime", "type":"pdate", "stored":true, "indexed":true},
    {"name":"search_text", "type":"text_general", "stored":true, "indexed":true},
    {"name":"pets_allowed", "type":"boolean", "stored":true, "indexed":true},
    {"name":"smoking_allowed", "type":"boolean", "stored":true, "indexed":true},
    {"name":"music_allowed", "type":"boolean", "stored":true, "indexed":true}
  ]
}'
```

## Reiniciar el Search-API

Una vez que Solr tenga el core, reinicia el search-api:

```bash
cd backend/search-api
go run cmd/api/main.go
```

Deberías ver en los logs:

```
Connected to Apache Solr successfully
```

## Verificar la Conexión

```bash
curl http://localhost:8004/health | jq
```

Deberías ver:

```json
{
  "status": "ok",
  "services": {
    "solr": {
      "status": "healthy",
      "message": "Connected"
    },
    ...
  }
}
```

## Indexar Viajes Existentes

Si ya tienes viajes en MongoDB (trips-api), necesitas sincronizarlos a Solr:

### Opción A: Publicar eventos trip.created via RabbitMQ

El consumer del search-api escucha eventos `trip.created` y automáticamente:
1. Obtiene el viaje del trips-api
2. Obtiene el conductor del users-api
3. Crea el SearchTrip denormalizado
4. Lo guarda en MongoDB
5. Lo indexa en Solr

### Opción B: Script de migración manual

Si necesitas migrar todos los viajes existentes, crea un script que:

1. Lea todos los viajes del trips-api
2. Publique eventos trip.created en RabbitMQ
3. El consumer del search-api procesará cada evento

## Troubleshooting

### Problema: "Solr client not initialized"

**Causa:** Solr no está corriendo o el core no existe
**Solución:** Verifica que Solr esté en http://localhost:8983 y que el core exista

### Problema: Búsquedas retornan 0 resultados

**Posibles causas:**

1. **No hay viajes en MongoDB**
   ```bash
   # Verifica con mongosh
   mongosh
   > use carpooling_search
   > db.trips.countDocuments()
   ```

2. **Los viajes no tienen status='published'**
   ```bash
   > db.trips.countDocuments({status: 'published'})
   ```

3. **Los viajes no tienen available_seats > 0**
   ```bash
   > db.trips.countDocuments({status: 'published', available_seats: {$gt: 0}})
   ```

4. **Solr no tiene documentos indexados**
   ```bash
   curl "http://localhost:8983/solr/carpooling_trips/select?q=*:*&rows=0"
   ```

### Problema: Viajes en MongoDB pero no en Solr

**Causa:** Los viajes se crearon antes de que el consumer estuviera activo
**Solución:** Volver a publicar los eventos trip.created

## Logs para Depuración

Para ver los logs del search-api en tiempo real:

```bash
cd backend/search-api
go run cmd/api/main.go
```

Los logs mostrarán:

```
[INFO] Starting search-api
[INFO] MongoDB indexes created successfully
[INFO] Connected to Apache Solr successfully
[INFO] Connected to Memcached successfully
[INFO] RabbitMQ consumer started in background
[INFO] search-api server listening port=8004

# Cuando llegue una búsqueda:
[INFO] HTTP request completed method=GET url=/api/v1/search/trips status_code=200 duration_ms=45

# Si hay errores:
[ERROR] Failed to index trip in Solr trip_id=... error=...
```

## Logs del Frontend

Para ver las requests del frontend:

1. Abre DevTools (F12)
2. Ve a la pestaña "Network"
3. Filtra por "search"
4. Haz una búsqueda
5. Verás las requests a `/api/search/trips`

## Próximos Pasos

Una vez que Solr esté configurado:

1. ✅ Crear el core `carpooling_trips`
2. ✅ Reiniciar search-api
3. ✅ Verificar el health check
4. ✅ Crear algunos viajes de prueba
5. ✅ Verificar que se indexen en Solr
6. ✅ Probar búsquedas desde el frontend
