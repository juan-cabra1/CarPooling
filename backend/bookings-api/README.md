# Bookings API

Microservicio de gestión de reservas para el sistema de CarPooling.

## 📋 Descripción

El **Bookings API** es responsable de manejar todas las operaciones relacionadas con reservas de viajes en el sistema CarPooling. Permite a los usuarios crear, consultar, actualizar y cancelar reservas, además de gestionar el estado de las mismas mediante eventos asíncronos con RabbitMQ.

### Características principales

- ✅ CRUD completo de reservas (bookings)
- ✅ Autenticación y autorización con JWT
- ✅ Comunicación asíncrona con RabbitMQ
- ✅ Persistencia en MySQL con GORM
- ✅ Arquitectura limpia (Clean Architecture)
- ✅ Logging estructurado con zerolog
- ✅ API RESTful con Gin Framework

---

## 🏗️ Arquitectura

Este microservicio sigue los principios de **Clean Architecture**, separando responsabilidades en capas bien definidas:

```
bookings-api/
├── cmd/
│   └── api/
│       └── main.go              # Entry point de la aplicación
├── internal/
│   ├── config/
│   │   └── config.go            # Configuración del servicio
│   ├── domain/
│   │   └── booking.go           # Entidades de negocio (modelos de dominio)
│   ├── dao/
│   │   └── booking.go           # Data Access Objects (modelos de BD)
│   ├── repository/
│   │   └── booking.go           # Capa de acceso a datos (GORM)
│   ├── services/
│   │   └── booking.go           # Lógica de negocio
│   ├── controllers/
│   │   └── booking.go           # HTTP handlers (Gin)
│   └── middleware/
│       └── cors.go              # Middlewares (CORS, Auth, etc.)
├── .env.example                 # Plantilla de variables de entorno
├── go.mod                       # Dependencias Go
└── README.md                    # Este archivo
```

### Capas de la arquitectura

1. **Domain Layer** (`internal/domain/`)
   - Define las entidades de negocio
   - Reglas de negocio puras, sin dependencias externas
   - Modelos independientes de la base de datos

2. **DAO Layer** (`internal/dao/`)
   - Data Access Objects
   - Modelos específicos para la base de datos (GORM)
   - Mapeo entre domain y base de datos

3. **Repository Layer** (`internal/repository/`)
   - Abstracción del acceso a datos
   - Operaciones CRUD sobre la base de datos
   - Usa GORM para interactuar con MySQL

4. **Service Layer** (`internal/services/`)
   - Lógica de negocio
   - Orquestación de operaciones
   - Validaciones y reglas de negocio complejas
   - Comunicación con RabbitMQ

5. **Controller Layer** (`internal/controllers/`)
   - HTTP handlers
   - Validación de entrada
   - Transformación de datos (DTO ↔ Domain)
   - Manejo de respuestas HTTP

6. **Config Layer** (`internal/config/`)
   - Carga de configuración desde variables de entorno
   - Validación de configuración requerida

---

## 🚀 Tecnologías

| Tecnología | Versión | Propósito |
|------------|---------|-----------|
| **Go** | 1.24.1 | Lenguaje de programación |
| **Gin** | 1.11.0 | Framework HTTP/REST |
| **GORM** | 1.31.1 | ORM para MySQL |
| **MySQL** | 8.0+ | Base de datos relacional |
| **RabbitMQ** | 3.13+ | Message broker (AMQP) |
| **JWT** | v5.3.0 | Autenticación con tokens |
| **zerolog** | 1.34.0 | Logging estructurado |
| **UUID** | 1.6.0 | Generación de IDs únicos |
| **godotenv** | 1.5.1 | Carga de variables de entorno |

---

## 📦 Instalación

### Prerrequisitos

- Go 1.23 o superior
- MySQL 8.0 o superior
- RabbitMQ 3.13 o superior
- Git

### Pasos de instalación

1. **Clonar el repositorio**
   ```bash
   git clone https://github.com/your-org/CarPooling.git
   cd CarPooling/backend/bookings-api
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

4. **Configurar base de datos**
   ```bash
   # Crear base de datos en MySQL
   mysql -u root -p -e "CREATE DATABASE bookings_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
   ```

5. **Ejecutar el servicio**
   ```bash
   go run cmd/api/main.go
   ```

El servidor estará disponible en `http://localhost:8003`

---

## ⚙️ Configuración

### Variables de entorno

Todas las variables de entorno están documentadas en `.env.example`. Las principales son:

| Variable | Descripción | Requerida | Default |
|----------|-------------|-----------|---------|
| `SERVER_PORT` | Puerto HTTP del servidor | No | `8003` |
| `DATABASE_URL` | DSN de MySQL | Sí | - |
| `JWT_SECRET` | Secreto para firmar JWT | Sí | - |
| `RABBITMQ_URL` | URL de RabbitMQ | Sí | - |
| `ENVIRONMENT` | Entorno de ejecución | No | `development` |

### Ejemplo de configuración para desarrollo

```bash
SERVER_PORT=8003
DATABASE_URL=root:password@tcp(localhost:3306)/bookings_db?charset=utf8mb4&parseTime=True&loc=Local
JWT_SECRET=dev-secret-key-change-in-production
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
ENVIRONMENT=development
```

### Ejemplo de configuración para Docker

```bash
SERVER_PORT=8003
DATABASE_URL=root:password@tcp(mysql:3306)/bookings_db?charset=utf8mb4&parseTime=True&loc=Local
JWT_SECRET=your-production-secret
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
ENVIRONMENT=production
```

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
go test ./internal/services/...
```

---

## 📡 API Endpoints

### Health Check

- **GET** `/health` - Verifica el estado del servicio

### Bookings

- **POST** `/api/v1/bookings` - Crear nueva reserva (requiere auth)
- **GET** `/api/v1/bookings/:id` - Obtener reserva por ID
- **GET** `/api/v1/bookings` - Listar reservas (con filtros)
- **PUT** `/api/v1/bookings/:id` - Actualizar reserva (requiere auth)
- **DELETE** `/api/v1/bookings/:id` - Cancelar reserva (requiere auth)
- **PATCH** `/api/v1/bookings/:id/confirm` - Confirmar reserva (requiere auth)

---

## 🔧 Desarrollo

### Estructura de código recomendada

1. **Siempre empieza por el dominio**: Define tus entidades en `internal/domain/`
2. **Crea el DAO**: Modela la estructura de BD en `internal/dao/`
3. **Implementa el repository**: Operaciones de datos en `internal/repository/`
4. **Lógica de negocio en service**: Orquestación en `internal/services/`
5. **Expón vía controller**: HTTP handlers en `internal/controllers/`

### Buenas prácticas

- ✅ Usa interfaces para abstraer dependencias
- ✅ Aplica inyección de dependencias
- ✅ Escribe tests unitarios para cada capa
- ✅ Documenta funciones públicas
- ✅ Maneja errores de forma consistente
- ✅ Usa logging estructurado con zerolog
- ✅ Valida entrada de usuario siempre

---

## 🐳 Docker

### Descripción del Dockerfile

El proyecto utiliza un **Dockerfile multi-stage** optimizado para producción:

- **Stage 1 (Builder)**: Compila el binario Go estático
  - Imagen base: `golang:1.22-alpine`
  - Optimizaciones: CGO deshabilitado, binario estático
  - Layer caching para builds más rápidos

- **Stage 2 (Runtime)**: Imagen minimalista para ejecución
  - Imagen base: `alpine:latest`
  - Usuario no-root para seguridad
  - Health check integrado
  - Tamaño final reducido (~20MB)

### Prerrequisitos para Docker

- Docker 20.10 o superior
- Docker Compose 2.0 o superior

### Variables de entorno necesarias

Antes de ejecutar con Docker, asegúrate de tener un archivo `.env` en la raíz del proyecto con estas variables:

```bash
# MySQL
MYSQL_ROOT_PASSWORD=your_secure_password

# RabbitMQ
RABBITMQ_USER=guest
RABBITMQ_PASS=guest

# JWT
JWT_SECRET=your-secret-key-change-in-production

# Environment
ENVIRONMENT=development
GIN_MODE=release
LOG_LEVEL=info
```

### Comandos Docker

#### 1. Construir imagen localmente

```bash
cd backend/bookings-api
docker build -t bookings-api:latest .
```

#### 2. Ejecutar imagen standalone (no recomendado)

```bash
docker run -d \
  --name bookings-api \
  -p 8003:8003 \
  -e SERVER_PORT=8003 \
  -e DATABASE_URL="root:password@tcp(mysql:3306)/bookings_db?charset=utf8mb4&parseTime=True&loc=Local" \
  -e JWT_SECRET="your-secret" \
  -e RABBITMQ_URL="amqp://guest:guest@rabbitmq:5672/" \
  -e TRIPS_API_URL="http://trips-api:8002" \
  bookings-api:latest
```

#### 3. Ejecutar con Docker Compose (recomendado)

```bash
# Desde la raíz del proyecto
cd /path/to/CarPooling

# Iniciar todos los servicios
docker-compose up -d

# Iniciar solo bookings-api y sus dependencias
docker-compose up -d bookings-api

# Ver logs en tiempo real
docker-compose logs -f bookings-api

# Reconstruir imagen tras cambios
docker-compose up --build bookings-api

# Detener todos los servicios
docker-compose down

# Detener y eliminar volúmenes (WARNING: borra datos)
docker-compose down -v
```

### Docker Compose - Servicios incluidos

Cuando ejecutas `docker-compose up`, se inician los siguientes servicios:

| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| **mysql** | 3306 | Base de datos MySQL 8.0 |
| **rabbitmq** | 5672, 15672 | Message broker + Management UI |
| **trips-api** | 8002 | API de viajes (dependencia) |
| **mongodb** | 27017 | Base de datos para trips-api |
| **memcached** | 11211 | Cache para search-api |
| **bookings-api** | 8003 | Este servicio |

### Health Checks

El servicio incluye health checks automáticos:

```bash
# Verificar estado del servicio
curl http://localhost:8003/health

# Ver estado de health checks en Docker
docker inspect bookings-api-container | grep -A 10 "Health"

# Ver health status en compose
docker-compose ps
```

**Configuración del health check:**
- Intervalo: 30s
- Timeout: 10s
- Retries: 3
- Start period: 40s (tiempo para inicialización)

### Acceso a servicios

```bash
# Acceder al contenedor de bookings-api
docker exec -it bookings-api-container sh

# Ver logs de bookings-api
docker-compose logs -f bookings-api

# Acceder a MySQL
docker exec -it mysql-container mysql -uroot -p

# Verificar base de datos
docker exec -it mysql-container mysql -uroot -p -e "SHOW DATABASES;"
docker exec -it mysql-container mysql -uroot -p bookings_db -e "SHOW TABLES;"

# Acceder a RabbitMQ Management
# http://localhost:15672 (usuario: guest, password: guest)

# Ver exchanges y queues
docker exec -it rabbitmq-container rabbitmqadmin list exchanges
docker exec -it rabbitmq-container rabbitmqadmin list queues
```

### Troubleshooting

#### 1. El contenedor no inicia

```bash
# Ver logs detallados
docker-compose logs bookings-api

# Verificar que MySQL esté listo
docker-compose logs mysql | grep "ready for connections"

# Verificar que RabbitMQ esté listo
docker-compose logs rabbitmq | grep "Server startup complete"
```

#### 2. Error de conexión a MySQL

```bash
# Verificar conectividad desde el contenedor
docker exec -it bookings-api-container ping mysql

# Verificar variables de entorno
docker exec -it bookings-api-container env | grep DATABASE_URL
```

#### 3. Error de conexión a RabbitMQ

```bash
# Verificar que RabbitMQ esté escuchando
docker exec -it rabbitmq-container rabbitmq-diagnostics ping

# Verificar URL de conexión
docker exec -it bookings-api-container env | grep RABBITMQ_URL
```

#### 4. Port already in use

```bash
# Verificar qué proceso usa el puerto 8003
# Windows:
netstat -ano | findstr :8003

# Linux/Mac:
lsof -i :8003

# Cambiar puerto en docker-compose.yml o detener el proceso conflictivo
```

#### 5. Reconstruir todo desde cero

```bash
# Detener y eliminar contenedores, redes, volúmenes
docker-compose down -v

# Eliminar imágenes
docker rmi bookings-api

# Reconstruir
docker-compose build --no-cache bookings-api

# Iniciar
docker-compose up -d
```

### Workflow de desarrollo con Docker

```bash
# 1. Iniciar servicios de infraestructura
docker-compose up -d mysql rabbitmq mongodb memcached

# 2. Desarrollar localmente (fuera de Docker)
cd backend/bookings-api
go run cmd/api/main.go

# 3. Cuando esté listo, probar con Docker
docker-compose up --build bookings-api

# 4. Verificar que todo funciona
curl http://localhost:8003/health
```

### Testing en Docker

```bash
# Ejecutar tests dentro del contenedor
docker-compose exec bookings-api go test ./...

# Ejecutar tests con coverage
docker-compose exec bookings-api go test -cover ./...

# Ejecutar tests de integración
docker-compose exec bookings-api go test -tags=integration ./...
```

### Optimizaciones de producción

Para producción, considera:

1. **Resource limits**: Descomentar sección `deploy.resources` en `docker-compose.yml`
2. **Multi-stage build**: Ya implementado, reduce tamaño de imagen
3. **Non-root user**: Ya implementado, mejora seguridad
4. **Health checks**: Ya implementados, mejora confiabilidad
5. **Restart policy**: `unless-stopped` configurado

### Monitoreo de recursos

```bash
# Ver uso de recursos en tiempo real
docker stats bookings-api-container

# Ver consumo de disco
docker system df

# Ver volúmenes
docker volume ls
docker volume inspect carpooling_mysql_data
```

---

## 📚 Recursos adicionales

- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

## 🤝 Contribución

1. Crea una rama desde `dev`: `git checkout -b feature/nueva-funcionalidad`
2. Realiza tus cambios siguiendo las convenciones del proyecto
3. Escribe tests para tu código
4. Crea un Pull Request a `dev`

---

## 📄 Licencia

Este proyecto es parte del sistema CarPooling desarrollado para fines educativos.

---

## 👥 Equipo

Desarrollado por el equipo de CarPooling - Software Architecture II

---

**Estado del proyecto**: 🚧 En desarrollo activo
