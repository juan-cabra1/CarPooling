# 🚗 CarPooling Platform

Sistema de carpooling/ridesharing desarrollado con arquitectura de microservicios, diseñado para conectar conductores y pasajeros de manera eficiente y segura.

## 📋 Descripción

CarPooling es una plataforma completa que permite a los usuarios:
- **Conductores**: Publicar viajes disponibles, gestionar reservas y recibir calificaciones
- **Pasajeros**: Buscar viajes, realizar reservas y calificar conductores
- **Administradores**: Gestionar usuarios, viajes y reservas desde un panel administrativo

El sistema está construido con una arquitectura de microservicios event-driven, garantizando escalabilidad, mantenibilidad y alta disponibilidad.

## 🏗️ Arquitectura

```
┌──────────────────────────────────────────────────────────────────────┐
│                           AWS EC2 Instance                            │
│                                                                       │
│  ┌─────────────────────────────────────────────────────────────────┐ │
│  │                       Nginx (Port 80)                            │ │
│  │  - Reverse Proxy para APIs                                       │ │
│  │  - Servidor del Frontend estático                                │ │
│  └────┬──────────┬──────────┬──────────┬───────────┬───────────────┘ │
│       │          │          │          │           │                  │
│  ┌────▼────┐ ┌──▼───┐  ┌───▼────┐ ┌───▼─────┐ ┌──▼────────┐        │
│  │Frontend │ │Users │  │ Trips  │ │Bookings │ │  Search   │        │
│  │ React   │ │ API  │  │  API   │ │   API   │ │    API    │        │
│  │ (SPA)   │ │:8001 │  │  :8002 │ │  :8003  │ │   :8004   │        │
│  └─────────┘ └──┬───┘  └───┬────┘ └───┬─────┘ └──┬────────┘        │
│                 │          │          │           │                  │
│  ┌──────────────▼──────────▼──────────▼───────────▼────────────┐    │
│  │                      RabbitMQ :5672                          │    │
│  │         Event Bus (trip.*, booking.*, user.*)                │    │
│  └──────────────────────────────────────────────────────────────┘    │
│                 │          │          │           │                  │
│  ┌──────────────▼────┐ ┌──▼──────┐ ┌─▼──────┐ ┌─▼──────────┐       │
│  │ MySQL (Users)     │ │ MongoDB │ │ MySQL  │ │   Solr     │       │
│  │ :3308             │ │ :27017  │ │ :3307  │ │   :8983    │       │
│  │ - users           │ │ - trips │ │- books │ │- ft search │       │
│  │ - ratings         │ │         │ │        │ │            │       │
│  └───────────────────┘ └─────────┘ └────────┘ └────────────┘       │
│                                                                       │
│  ┌─────────────────┐                                                 │
│  │   Memcached     │                                                 │
│  │   :11211        │                                                 │
│  │   (Cache)       │                                                 │
│  └─────────────────┘                                                 │
└──────────────────────────────────────────────────────────────────────┘
```

### Flujo de Datos

1. **Usuario → Nginx**: Todas las peticiones llegan a Nginx en el puerto 80
2. **Nginx → Microservicios**: Nginx enruta las peticiones API al servicio correspondiente
3. **Nginx → Frontend**: Sirve la SPA de React para rutas no-API
4. **APIs → RabbitMQ**: Los microservicios publican eventos (trip.created, booking.confirmed, etc.)
5. **RabbitMQ → Consumers**: Los servicios consumen eventos relevantes (event-driven)
6. **APIs → Bases de Datos**: Cada servicio accede a su propia base de datos

## 🚀 Tecnologías

### Backend
- **Lenguaje**: Go 1.21+
- **Framework HTTP**: Gin
- **Autenticación**: JWT (JSON Web Tokens)
- **ORM**: GORM (MySQL)
- **Logging**: zerolog

### Frontend
- **Framework**: React 18
- **Lenguaje**: TypeScript
- **Build Tool**: Vite
- **Routing**: React Router v6
- **Styling**: Tailwind CSS
- **Animaciones**: Framer Motion
- **Iconos**: Tabler Icons, Lucide React

### Bases de Datos
- **MySQL 8.0**: users-api (usuarios y calificaciones), bookings-api (reservas)
- **MongoDB 7.0**: trips-api (viajes), search-api (denormalización)

### Infraestructura
- **Message Broker**: RabbitMQ 3.13
- **Search Engine**: Apache Solr 9.0
- **Cache**: Memcached 1.6
- **Reverse Proxy**: Nginx
- **Orquestación**: Docker Compose
- **Deployment**: AWS EC2

## 📁 Estructura del Proyecto

```
CarPooling/
├── backend/
│   ├── users-api/          # Gestión de usuarios y autenticación
│   │   ├── cmd/api/        # Entry point
│   │   ├── internal/       # Lógica de negocio
│   │   └── README.md       # Documentación detallada
│   │
│   ├── trips-api/          # Gestión de viajes
│   │   ├── cmd/api/
│   │   ├── internal/
│   │   └── README.md
│   │
│   ├── bookings-api/       # Gestión de reservas
│   │   ├── cmd/api/
│   │   ├── internal/
│   │   └── README.md
│   │
│   └── search-api/         # Búsqueda avanzada de viajes
│       ├── cmd/api/
│       ├── internal/
│       └── README.md
│
├── frontend/               # Aplicación React
│   ├── src/
│   │   ├── components/    # Componentes reutilizables
│   │   ├── pages/         # Páginas de la aplicación
│   │   ├── services/      # API clients
│   │   ├── context/       # Context providers
│   │   └── types/         # TypeScript types
│   ├── nginx.conf         # Configuración de Nginx
│   ├── Dockerfile         # Build de producción
│   └── README.md
│
├── docker-compose.yml      # Orquestación de servicios
├── .env.example           # Template de variables de entorno
└── README.md              # Este archivo
```

## 🌐 Deployment en AWS EC2

El sistema está desplegado en una instancia AWS EC2 con la siguiente configuración:

### Componentes de Deployment

1. **Nginx** (Puerto 80)
   - **Función**: Reverse proxy y servidor del frontend
   - **Configuración**: [frontend/nginx.conf](frontend/nginx.conf)
   - **Rutas**:
     - `/` → Frontend React (SPA)
     - `/api/users` → users-api:8001
     - `/api/trips` → trips-api:8002
     - `/api/bookings` → bookings-api:8003
     - `/api/search` → search-api:8004

2. **Docker Compose**
   - **Función**: Orquestación de todos los servicios
   - **Archivo**: [docker-compose.yml](docker-compose.yml)
   - **Redes**: Todos los servicios en `carpooling-network`

3. **Volúmenes Persistentes**
   - `mongo_data`: Datos de MongoDB
   - `mysql_users_data`: Base de datos de usuarios
   - `mysql_bookings_data`: Base de datos de reservas
   - `rabbit_data`: Colas y mensajes de RabbitMQ
   - `solr_data`: Índices de búsqueda

### Puertos Expuestos

| Servicio | Puerto Interno | Puerto Externo | Acceso |
|----------|----------------|----------------|--------|
| Nginx (Frontend) | 80 | 80 | Público |
| users-api | 8001 | - | Via Nginx |
| trips-api | 8002 | - | Via Nginx |
| bookings-api | 8003 | - | Via Nginx |
| search-api | 8004 | - | Via Nginx |
| MongoDB | 27017 | 27017 | Interno |
| MySQL (users) | 3306 | 3308 | Interno |
| MySQL (bookings) | 3306 | 3307 | Interno |
| RabbitMQ (AMQP) | 5672 | 5672 | Interno |
| RabbitMQ (Management) | 15672 | 15672 | Interno |
| Solr | 8983 | 8983 | Interno |
| Memcached | 11211 | 11211 | Interno |

## ⚡ Inicio Rápido

### Prerrequisitos

- Docker 20.10+
- Docker Compose 2.0+
- Git

### Instalación

1. **Clonar el repositorio**
   ```bash
   git clone https://github.com/your-org/CarPooling.git
   cd CarPooling
   ```

2. **Configurar variables de entorno**
   ```bash
   cp .env.example .env
   # Editar .env con tus valores
   ```

3. **Iniciar todos los servicios**
   ```bash
   docker-compose up -d
   ```

4. **Verificar que todos los servicios estén corriendo**
   ```bash
   docker-compose ps
   ```

5. **Acceder a la aplicación**
   - Frontend: http://localhost
   - RabbitMQ Management: http://localhost:15672 (guest/guest)
   - Solr Admin: http://localhost:8983

### Comandos Útiles

```bash
# Ver logs de todos los servicios
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f users-api

# Reiniciar un servicio
docker-compose restart users-api

# Reconstruir un servicio tras cambios de código
docker-compose build users-api
docker-compose up -d users-api

# Detener todos los servicios
docker-compose down

# Detener y eliminar volúmenes (CUIDADO: elimina datos)
docker-compose down -v

# Verificar salud de los servicios
curl http://localhost/api/users/health
curl http://localhost/api/trips/health
curl http://localhost/api/bookings/health
curl http://localhost/api/search/health
```

## 📚 Microservicios

Cada microservicio tiene su propia documentación detallada:

### [Users API](backend/users-api/README.md) - Puerto 8001
- ✅ Registro y autenticación de usuarios (JWT)
- ✅ Verificación de email (SMTP)
- ✅ Recuperación de contraseña
- ✅ Gestión de perfiles
- ✅ Sistema de calificaciones para conductores y pasajeros
- **Stack**: Go, Gin, MySQL, GORM, bcrypt, JWT

### [Trips API](backend/trips-api/README.md) - Puerto 8002
- ✅ Crear, editar y eliminar viajes
- ✅ Gestión de disponibilidad de asientos
- ✅ Publicación de eventos a RabbitMQ
- ✅ Validación de conductores contra users-api
- **Stack**: Go, Gin, MongoDB, RabbitMQ, JWT

### [Bookings API](backend/bookings-api/README.md) - Puerto 8003
- ✅ Crear y gestionar reservas
- ✅ Confirmación/cancelación de reservas
- ✅ Comunicación event-driven con trips-api
- ✅ Optimistic locking para concurrencia
- **Stack**: Go, Gin, MySQL, GORM, RabbitMQ, JWT

### [Search API](backend/search-api/README.md) - Puerto 8004
- ✅ Búsqueda full-text con Apache Solr
- ✅ Búsqueda geoespacial con MongoDB
- ✅ Denormalización de datos (trips + drivers)
- ✅ Cache con Memcached
- ✅ Consumer de eventos RabbitMQ
- **Stack**: Go, Gin, MongoDB, Solr, Memcached, RabbitMQ

### [Frontend](frontend/README.md)
- ✅ SPA con React y TypeScript
- ✅ Autenticación con JWT
- ✅ Rutas protegidas y públicas
- ✅ Panel de administración
- ✅ Sistema de búsqueda y reservas
- **Stack**: React 18, TypeScript, Vite, Tailwind CSS, React Router

## ⚙️ Variables de Entorno

Las variables de entorno están documentadas en [.env.example](.env.example).

### Variables Críticas

```bash
# JWT Secret (DEBE ser igual en todos los microservicios)
JWT_SECRET=your-secret-key-here

# URLs de las aplicaciones
APP_URL=http://localhost:3000

# SMTP para emails
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_FROM=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Bases de datos
MYSQL_USERS_ROOT_PASSWORD=strong-password
MYSQL_BOOKINGS_ROOT_PASSWORD=strong-password
MONGO_ROOT_PASSWORD=strong-password

# RabbitMQ
RABBITMQ_USER=carpool_user
RABBITMQ_PASS=strong-password
```

## 🔧 Desarrollo

### Desarrollo Local (Sin Docker)

Para desarrollar un microservicio localmente:

```bash
# 1. Iniciar infraestructura con Docker
docker-compose up -d mongo mysql-users mysql-bookings rabbit memcached solr

# 2. Configurar .env local del servicio
cd backend/users-api
cp .env.example .env
# Editar .env con valores locales (localhost en lugar de nombres de servicio)

# 3. Ejecutar el servicio
go run cmd/api/main.go

# 4. Para el frontend
cd frontend
npm install
npm run dev
```

### Testing

```bash
# Backend - ejecutar tests de un microservicio
cd backend/users-api
go test ./... -v

# Frontend - ejecutar tests
cd frontend
npm test

# Tests con coverage
go test -cover ./...
```

### Build de Producción

```bash
# Build de todos los servicios
docker-compose build

# Build de un servicio específico
docker-compose build users-api

# Frontend - build de producción
cd frontend
npm run build
# Los archivos se generan en frontend/dist
```

## 🔒 Seguridad

- **Autenticación**: JWT con expiración de 24 horas
- **Contraseñas**: Hasheadas con bcrypt (cost 10)
- **CORS**: Configurado en todos los microservicios
- **Rate Limiting**: Implementado en Nginx
- **Validación**: Validación de entrada en todos los endpoints
- **Secrets**: Gestionados mediante variables de entorno
- **TLS**: Recomendado para producción (configurar en Nginx)

## 📊 Monitoreo

### Health Checks

Todos los servicios exponen un endpoint `/health`:

```bash
curl http://localhost/api/users/health
curl http://localhost/api/trips/health
curl http://localhost/api/bookings/health
curl http://localhost/api/search/health
```

### Logs

```bash
# Ver logs de un servicio
docker-compose logs -f users-api

# Ver logs de todos los servicios
docker-compose logs -f

# Ver logs de Nginx
docker-compose logs -f nginx
```

### RabbitMQ Management

Accede a http://localhost:15672 para:
- Monitorear colas y mensajes
- Ver exchanges y bindings
- Estadísticas de throughput

## 🤝 Contribución

1. Fork el proyecto
2. Crear una rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit los cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear un Pull Request

### Convenciones

- **Commits**: Usar conventional commits (feat, fix, docs, etc.)
- **Código Go**: Seguir [Effective Go](https://golang.org/doc/effective_go)
- **Código TypeScript**: Seguir guía de estilo de Airbnb
- **Tests**: Escribir tests para nuevas funcionalidades

## 📝 Licencia

Este proyecto fue desarrollado con fines educativos para el curso de Arquitectura de Software II.

## 👥 Equipo

Desarrollado por el equipo de CarPooling - Arquitectura de Software II

---

**Versión**: 1.0.0
**Última actualización**: 2025-12-07
**Estado**: ✅ En producción en AWS EC2
