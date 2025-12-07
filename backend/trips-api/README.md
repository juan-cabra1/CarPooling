# 🚗 Trips API - CarPooling

Microservicio de gestión de viajes para la plataforma CarPooling. Maneja la creación, actualización, eliminación y consulta de viajes, además de gestionar la disponibilidad de asientos mediante optimistic locking y comunicación event-driven con otros servicios.

## 📋 Descripción

El **Trips API** es el núcleo del sistema CarPooling, responsible de:
- Gestionar el ciclo de vida completo de los viajes
- Validar conductores contra users-api
- Gestionar la disponibilidad de asientos con optimistic locking
- Publicar eventos a RabbitMQ cuando ocurren cambios en viajes
- Consumir eventos de bookings-api para actualizar asientos reservados

### Características Principales

- ✅ CRUD completo de viajes
- ✅ Validación de conductores con users-api
- ✅ Gestión de disponibilidad de asientos (optimistic locking)
- ✅ Publicación de eventos a RabbitMQ (trip.created, trip.updated, etc.)
- ✅ Consumo de eventos de bookings-api (reservation.created, reservation.cancelled)
- ✅ Búsqueda de viajes por conductor
- ✅ Persistencia en MongoDB
- ✅ Autenticación y autorización con JWT
- ✅ Arquitectura limpia (Clean Architecture)
- ✅ Logging estructurado con zerolog

---

## 🏗️ Arquitectura

Este microservicio sigue los principios de **Clean Architecture**, separando responsabilidades en capas bien definidas:

```
trips-api/
├── cmd/
│   └── api/
│       └── main.go                # Entry point de la aplicación
├── internal/
│   ├── config/
│   │   └── config.go              # Configuración del servicio
│   ├── domain/
│   │   └── trip.go                # Entidades de negocio (Trip, Location, Car, Preferences)
│   ├── dao/
│   │   └── trip.go                # Data Access Objects (modelos de BD)
│   ├── repository/
│   │   ├── trip.go                # Capa de acceso a datos (MongoDB)
│   │   └── idempotency.go         # Gestión de idempotencia de eventos
│   ├── service/
│   │   ├── trip_service.go        # Lógica de negocio de viajes
│   │   ├── trip_event_service.go  # Procesamiento de eventos
│   │   └── idempotency_service.go # Servicio de idempotencia
│   ├── controller/
│   │   └── trip_controller.go     # HTTP handlers (Gin)
│   ├── clients/
│   │   └── users_client.go        # Cliente HTTP para users-api
│   ├── messaging/
│   │   ├── publisher.go           # Publisher de eventos a RabbitMQ
│   │   ├── consumer.go            # Consumer de eventos de RabbitMQ
│   │   └── events.go              # Definición de eventos
│   ├── middleware/
│   │   ├── auth.go                # Middleware JWT
│   │   └── cors.go                # Middleware CORS
│   └── routes/
│       └── routes.go              # Configuración de rutas
├── .env.example                   # Plantilla de variables de entorno
├── go.mod                         # Dependencias Go
└── README.md                      # Este archivo
```

### Capas de la Arquitectura

1. **Domain Layer** (`internal/domain/`)
   - Entidades de negocio puras: Trip, Location, Car, Preferences
   - DTOs de request/response
   - Sin dependencias externas

2. **DAO Layer** (`internal/dao/`)
   - Data Access Objects para MongoDB
   - Tags BSON para serialización
   - Mapeo entre domain y base de datos

3. **Repository Layer** (`internal/repository/`)
   - Abstracción del acceso a datos
   - Operaciones CRUD sobre MongoDB
   - Optimistic locking para concurrencia
   - Gestión de idempotencia de eventos

4. **Service Layer** (`internal/service/`)
   - Lógica de negocio de viajes
   - Validación de conductores
   - Gestión de eventos (publish/consume)
   - Orquestación de operaciones

5. **Controller Layer** (`internal/controller/`)
   - HTTP handlers con Gin
   - Validación de entrada
   - Transformación de datos
   - Manejo de respuestas HTTP

6. **Messaging Layer** (`internal/messaging/`)
   - Publisher de eventos a RabbitMQ
   - Consumer de eventos de bookings-api
   - Definición de payloads de eventos

7. **Clients Layer** (`internal/clients/`)
   - Cliente HTTP para users-api
   - Validación de conductores

---

## 🚀 Tecnologías

| Tecnología | Versión | Propósito |
|------------|---------|-----------|
| **Go** | 1.21+ | Lenguaje de programación |
| **Gin** | 1.11.0 | Framework HTTP/REST |
| **MongoDB** | 7.0+ | Base de datos NoSQL |
| **MongoDB Driver** | 1.17+ | Driver oficial de MongoDB para Go |
| **RabbitMQ** | 3.13+ | Message broker (AMQP) |
| **JWT** | v5.3.0 | Autenticación con tokens |
| **zerolog** | 1.34.0 | Logging estructurado |
| **godotenv** | 1.5.1 | Carga de variables de entorno |

---

## 📦 Instalación

### Prerrequisitos

- Go 1.21 o superior
- MongoDB 7.0 o superior
- RabbitMQ 3.13 o superior
- Git

### Pasos de Instalación

1. **Clonar el repositorio**
   ```bash
   git clone https://github.com/your-org/CarPooling.git
   cd CarPooling/backend/trips-api
   ```

2. **Instalar dependencias**
   ```bash
   go mod download
   ```

3. **Configurar variables de entorno**
   ```bash
   cp .env.example .env
   # Editar .env con tus valores reales
   ```

4. **Configurar MongoDB**
   ```bash
   # MongoDB debe estar corriendo en localhost:27017
   # O usar Docker:
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   ```

5. **Ejecutar el servicio**
   ```bash
   go run cmd/api/main.go
   ```

El servidor estará disponible en `http://localhost:8002`

---

## ⚙️ Configuración

### Variables de Entorno

Todas las variables de entorno están documentadas en `.env.example`. Las principales son:

| Variable | Descripción | Requerida | Default |
|----------|-------------|-----------|---------|
| `SERVER_PORT` | Puerto HTTP del servidor | No | `8002` |
| `MONGO_URI` | URI de conexión a MongoDB | Sí | - |
| `MONGO_DB` | Nombre de la base de datos | No | `carpooling` |
| `JWT_SECRET` | Secreto para firmar JWT | Sí | - |
| `RABBITMQ_URL` | URL de RabbitMQ | Sí | - |
| `USERS_API_URL` | URL del users-api | Sí | - |
| `ENVIRONMENT` | Entorno de ejecución | No | `development` |

### Ejemplo de Configuración para Desarrollo

```bash
SERVER_PORT=8002
MONGO_URI=mongodb://localhost:27017
MONGO_DB=carpooling
JWT_SECRET=dev-secret-key-change-in-production
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
USERS_API_URL=http://localhost:8001
ENVIRONMENT=development
```

### Ejemplo de Configuración para Docker

```bash
SERVER_PORT=8002
MONGO_URI=mongodb://mongo:27017
MONGO_DB=carpooling
JWT_SECRET=your-production-secret
RABBITMQ_URL=amqp://guest:guest@rabbit:5672/
USERS_API_URL=http://users-api:8001
ENVIRONMENT=production
```

---

## 📡 API Endpoints

### Health Check

- **GET** `/health` - Verifica el estado del servicio

### Trips

Todos los endpoints de trips requieren autenticación JWT (excepto GET públicos).

#### Crear Viaje
- **POST** `/trips`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Body**:
  ```json
  {
    "origin": {
      "city": "Bogotá",
      "province": "Cundinamarca",
      "address": "Calle 100 #10-20",
      "coordinates": {
        "type": "Point",
        "coordinates": [-74.0721, 4.7110]
      }
    },
    "destination": {
      "city": "Medellín",
      "province": "Antioquia",
      "address": "Carrera 43A #1-50",
      "coordinates": {
        "type": "Point",
        "coordinates": [-75.5636, 6.2476]
      }
    },
    "departure_datetime": "2025-12-15T08:00:00Z",
    "estimated_arrival_datetime": "2025-12-15T14:00:00Z",
    "price_per_seat": 50000,
    "total_seats": 3,
    "car": {
      "model": "Toyota Corolla",
      "color": "Gris",
      "license_plate": "ABC123"
    },
    "preferences": {
      "pets_allowed": false,
      "smoking_allowed": false,
      "music_allowed": true
    },
    "description": "Viaje cómodo a Medellín, salida temprano"
  }
  ```
- **Response**: `201 Created`

#### Obtener Viaje por ID
- **GET** `/trips/:id`
- **Response**: `200 OK`

#### Listar Viajes
- **GET** `/trips?driver_id=123&page=1&limit=20`
- **Query Parameters**:
  - `driver_id` (opcional): Filtrar por conductor
  - `page` (opcional): Número de página
  - `limit` (opcional): Resultados por página
- **Response**: `200 OK`

#### Actualizar Viaje
- **PUT** `/trips/:id`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Body**: Campos a actualizar (parcial)
- **Response**: `200 OK`
- **Nota**: Solo el dueño del viaje o admin puede actualizar

#### Eliminar Viaje
- **DELETE** `/trips/:id`
- **Headers**: `Authorization: Bearer <jwt_token>`
- **Response**: `200 OK`
- **Nota**: Solo el dueño del viaje o admin puede eliminar

---

## 🔄 Event-Driven Architecture

### Eventos Publicados

El trips-api publica los siguientes eventos a RabbitMQ:

#### trip.created
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.created",
  "timestamp": "2025-12-07T10:00:00Z",
  "trip_id": "mongodb-object-id",
  "driver_id": 123,
  "origin": { ... },
  "destination": { ... },
  "departure_datetime": "2025-12-15T08:00:00Z",
  "total_seats": 3,
  "available_seats": 3,
  "price_per_seat": 50000
}
```

#### trip.updated
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.updated",
  "timestamp": "2025-12-07T11:00:00Z",
  "trip_id": "mongodb-object-id",
  "driver_id": 123,
  "available_seats": 2,
  "updated_fields": ["available_seats"]
}
```

#### trip.deleted
```json
{
  "event_id": "uuid-v4",
  "event_type": "trip.deleted",
  "timestamp": "2025-12-07T12:00:00Z",
  "trip_id": "mongodb-object-id",
  "driver_id": 123
}
```

### Eventos Consumidos

El trips-api consume eventos del bookings-api:

#### reservation.created
- **Acción**: Decrementa `available_seats` y aumenta `reserved_seats`
- **Validación**: Verifica que haya asientos disponibles
- **Optimistic Locking**: Usa `availability_version` para evitar race conditions
- **Compensación**: Publica evento de fallo si no hay asientos

#### reservation.cancelled
- **Acción**: Incrementa `available_seats` y decrementa `reserved_seats`
- **Validación**: Verifica que el viaje exista
- **Optimistic Locking**: Usa `availability_version`

---

## 🔐 Domain Models

### Trip
```go
type Trip struct {
    ID                       primitive.ObjectID
    DriverID                 int64
    Origin                   Location
    Destination              Location
    DepartureDatetime        time.Time
    EstimatedArrivalDatetime time.Time
    PricePerSeat             float64
    TotalSeats               int
    ReservedSeats            int
    AvailableSeats           int
    AvailabilityVersion      int  // Para optimistic locking
    Car                      Car
    Preferences              Preferences
    Status                   string  // published, full, cancelled, etc.
    Description              string
    CreatedAt                time.Time
    UpdatedAt                time.Time
}
```

### Location
```go
type Location struct {
    City        string
    Province    string
    Address     string
    Coordinates GeoJSONPoint  // MongoDB 2dsphere
}

type GeoJSONPoint struct {
    Type        string     // "Point"
    Coordinates []float64  // [lng, lat]
}
```

### Car
```go
type Car struct {
    Model        string
    Color        string
    LicensePlate string
}
```

### Preferences
```go
type Preferences struct {
    PetsAllowed    bool
    SmokingAllowed bool
    MusicAllowed   bool
}
```

---

## 🔒 Optimistic Locking

El trips-api implementa **optimistic locking** para evitar race conditions cuando múltiples reservas intentan modificar los asientos simultáneamente.

### Cómo Funciona

1. Cada viaje tiene un campo `availability_version` que se incrementa en cada actualización
2. Al procesar un evento `reservation.created`:
   - Se obtiene el viaje actual con su versión
   - Se valida que haya asientos disponibles
   - Se actualiza usando la versión como condición:
     ```go
     filter := bson.M{
         "_id": tripID,
         "availability_version": currentVersion,
         "available_seats": bson.M{"$gte": seatsRequested},
     }
     update := bson.M{
         "$inc": {
             "reserved_seats": seatsRequested,
             "available_seats": -seatsRequested,
             "availability_version": 1,
         },
     }
     ```
   - Si la versión cambió (otra reserva se procesó), la actualización falla
   - Se publica un evento de compensación si falla

### Ventajas

- ✅ Evita double-booking de asientos
- ✅ No requiere locks pesados
- ✅ Soporta alta concurrencia
- ✅ Fácil de implementar con MongoDB

---

## 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./...

# Tests con coverage
go test -cover ./...

# Tests con reporte detallado
go test -v ./...

# Tests de una capa específica
go test ./internal/service/...
```

---

## 🐳 Docker

### Build de la Imagen

```bash
cd backend/trips-api
docker build -t trips-api:latest .
```

### Ejecutar con Docker Compose (Recomendado)

```bash
# Desde la raíz del proyecto
docker-compose up -d trips-api

# Ver logs
docker-compose logs -f trips-api

# Reconstruir tras cambios
docker-compose build trips-api
docker-compose up -d trips-api
```

---

## 🔧 Desarrollo

### Estructura de Código Recomendada

1. **Siempre empieza por el dominio**: Define entidades en `internal/domain/`
2. **Crea el DAO**: Modela la estructura de BD en `internal/dao/`
3. **Implementa el repository**: Operaciones de datos en `internal/repository/`
4. **Lógica de negocio en service**: Orquestación en `internal/service/`
5. **Expón vía controller**: HTTP handlers en `internal/controller/`

### Buenas Prácticas

- ✅ Usa interfaces para abstraer dependencias
- ✅ Aplica inyección de dependencias
- ✅ Escribe tests unitarios para cada capa
- ✅ Documenta funciones públicas
- ✅ Maneja errores de forma consistente
- ✅ Usa logging estructurado con zerolog
- ✅ Valida entrada de usuario siempre
- ✅ Usa optimistic locking para concurrencia

---

## 🤝 Relación con Otros Servicios

### users-api
- **Validación de conductores**: Al crear un viaje, se valida que el conductor exista y sea válido
- **Endpoint usado**: `GET /internal/users/:id`

### bookings-api
- **Consume eventos**: `reservation.created`, `reservation.cancelled`
- **Actualiza asientos**: Modifica `reserved_seats` y `available_seats` basado en eventos

### search-api
- **Consume eventos**: `trip.created`, `trip.updated`, `trip.deleted`
- **Denormaliza datos**: search-api mantiene una copia del viaje + info del conductor

---

## 📚 Recursos Adicionales

- [Gin Documentation](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

## 📄 Licencia

Este proyecto es parte del sistema CarPooling desarrollado para fines educativos.

---

## 👥 Equipo

Desarrollado por el equipo de CarPooling - Arquitectura de Software II

---

**Estado del proyecto**: ✅ En producción en AWS EC2
**Puerto**: 8002
**Base de datos**: MongoDB
