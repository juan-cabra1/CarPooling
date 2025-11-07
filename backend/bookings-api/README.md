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

### Construir imagen

```bash
docker build -t bookings-api:latest .
```

### Ejecutar con Docker Compose

```bash
docker-compose up -d
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
