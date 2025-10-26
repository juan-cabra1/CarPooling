# Bookings API (Reservations API)

Microservicio de gestión de reservas para el Sistema de Carpooling UCC.

## Características

- **Proceso concurrente**: Validaciones paralelas con Goroutines y Channels
- **Arquitectura MVC**: Patrón Model-View-Controller con inyección de dependencias
- **Eventos RabbitMQ**: Publicación de eventos de reserva
- **JWT Authentication**: Validación de tokens JWT
- **Transacciones ACID**: Gestión transaccional con MySQL

## Estructura del Proyecto

```
bookings-api/
├── cmd/api/main.go              # Entry point con inyección de dependencias
├── internal/
│   ├── config/                  # Configuración desde .env
│   ├── controller/              # Controladores HTTP
│   ├── dao/                     # Modelos de base de datos (GORM)
│   ├── domain/                  # DTOs y entidades de dominio
│   ├── middleware/              # Auth, CORS, Error handling
│   ├── repository/              # Capa de acceso a datos
│   ├── routes/                  # Definición de rutas
│   └── service/                 # Lógica de negocio
│       ├── booking.go           # Servicio principal
│       ├── validation.go        # Validaciones concurrentes
│       ├── users_client.go      # Cliente HTTP users-api
│       ├── trips_client.go      # Cliente HTTP trips-api
│       └── rabbitmq.go          # Cliente RabbitMQ
├── Dockerfile                   # Multi-stage build
├── .env                        # Variables de entorno
└── go.mod                      # Dependencias
```

## Tecnologías

- **Go 1.21**
- **Gin** - Framework web
- **GORM** - ORM para MySQL
- **JWT** - Autenticación
- **RabbitMQ** - Mensajería asíncrona
- **MySQL** - Base de datos relacional

## Endpoints

### Públicos
- `GET /health` - Health check

### Protegidos (requieren JWT)
- `POST /bookings` - Crear reserva (con validaciones concurrentes)
- `GET /bookings/:id` - Obtener reserva por ID
- `PUT /bookings/:id/cancel` - Cancelar reserva
- `POST /bookings/:id/confirm-arrival` - Confirmar llegada segura
- `GET /bookings/user/:userId` - Obtener reservas de usuario
- `GET /bookings/trip/:tripId` - Obtener reservas de viaje

### Internos (sin autenticación)
- `GET /internal/bookings/:id` - Obtener reserva (para microservicios)
- `PUT /internal/bookings/:id/complete` - Completar reserva

## Proceso Concurrente

La creación de reservas utiliza **4 goroutines en paralelo**:

1. **Validar disponibilidad del viaje** (trips-api)
2. **Validar usuario pasajero** (users-api)
3. **Verificar reservas duplicadas** (base de datos)
4. **Validar pasajero ≠ conductor**

Utiliza `sync.WaitGroup` y `channels` con timeout de 5 segundos.

## Eventos RabbitMQ

Exchange: `reservations.events` (tipo: topic)

Eventos publicados:
- `reservation.created` - Nueva reserva creada
- `reservation.cancelled` - Reserva cancelada
- `reservation.completed` - Reserva completada

## Configuración (.env)

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=carpooling_reservations

# Server
SERVER_PORT=8003

# JWT (mismo secret que users-api)
JWT_SECRET=your-secret-key

# External Services
TRIPS_API_URL=http://localhost:8002
USERS_API_URL=http://localhost:8001

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=reservations.events
```

## Ejecución

### Desarrollo
```bash
go run cmd/api/main.go
```

### Compilar
```bash
go build -o bookings-api cmd/api/main.go
./bookings-api
```

### Tests
```bash
go test ./internal/service/... -v
```

### Docker
```bash
docker build -t bookings-api .
docker run -p 8003:8003 --env-file .env bookings-api
```

## Base de Datos

Schema: `carpooling_reservations`

Tabla: `reservations`
- id (UUID, PK)
- trip_id (VARCHAR(24), indexed) - ObjectID de MongoDB
- passenger_id (BIGINT, indexed)
- driver_id (BIGINT, indexed)
- seats_reserved (INT)
- price_per_seat (DECIMAL)
- total_amount (DECIMAL)
- status (ENUM: pending, confirmed, completed, cancelled)
- payment_status (ENUM: pending, paid, refunded)
- arrived_safely (BOOLEAN)
- arrival_confirmed_at (TIMESTAMP)
- created_at, updated_at (TIMESTAMP)

Índice único: (trip_id, passenger_id)

## Validaciones

- trip_id debe existir en trips-api
- passenger_id debe existir en users-api
- available_seats >= seats_reserved
- passenger_id ≠ driver_id
- No duplicados (trip_id + passenger_id)
- Departure date > NOW()
- Cancelación: > 24h antes de la salida

## Códigos HTTP

- `200` OK
- `201` Created
- `400` Bad Request
- `401` Unauthorized
- `403` Forbidden
- `404` Not Found
- `409` Conflict
- `500` Internal Server Error

## Arquitectura

Patrón **MVC** con inyección de dependencias:

```
Controller → Service → Repository → DAO
              ↓
        Validation Service
              ↓
    TripsClient, UsersClient, RabbitMQ
```

---

**Proyecto**: Sistema de Carpooling UCC  
**Curso**: Arquitectura de Software II  
**Fecha**: Noviembre 2024
