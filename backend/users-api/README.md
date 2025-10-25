# Users API - CarPooling

API de gestión de usuarios, autenticación y calificaciones para el sistema de CarPooling.

## Características

- ✅ Registro y autenticación de usuarios con JWT
- ✅ Verificación de email vía SMTP
- ✅ Recuperación y cambio de contraseña
- ✅ Gestión de perfiles de usuario
- ✅ Sistema de calificaciones para conductores y pasajeros
- ✅ Arquitectura limpia con capas separadas (Domain, DAO, Repository, Service, Controller)
- ✅ CORS configurado
- ✅ Middleware de autenticación JWT
- ✅ Tests unitarios con testify/mock
- ✅ GORM para ORM
- ✅ Docker support

## Tecnologías

- **Go 1.21+**
- **Gin** - Framework web
- **GORM** - ORM para MySQL
- **JWT** - Autenticación
- **bcrypt** - Hashing de contraseñas (cost 10)
- **MySQL** - Base de datos
- **SMTP** - Envío de emails
- **Testify** - Testing y mocking

## Estructura del Proyecto

```
users-api/
├── cmd/
│   └── api/
│       └── main.go                 # Punto de entrada
├── internal/
│   ├── config/
│   │   └── config.go               # Configuración y variables de entorno
│   ├── dao/
│   │   ├── user.go                 # UserDAO con GORM tags
│   │   └── rating.go               # RatingDAO con GORM tags
│   ├── domain/
│   │   ├── user.go                 # DTOs de usuario
│   │   └── rating.go               # DTOs de calificaciones
│   ├── repository/
│   │   ├── user.go                 # UserRepository con GORM
│   │   └── rating.go               # RatingRepository con GORM
│   ├── service/
│   │   ├── email.go                # EmailService (SMTP)
│   │   ├── auth.go                 # AuthService (JWT, bcrypt)
│   │   ├── user.go                 # UserService
│   │   ├── rating.go               # RatingService
│   │   └── user_test.go            # Tests unitarios
│   ├── middleware/
│   │   ├── auth.go                 # Middleware JWT
│   │   ├── cors.go                 # Middleware CORS
│   │   └── error.go                # Middleware de errores
│   ├── controller/
│   │   ├── auth.go                 # AuthController
│   │   ├── user.go                 # UserController
│   │   └── rating.go               # RatingController
│   └── routes/
│       └── routes.go               # Configuración de rutas
├── scripts/
│   └── init_db.sql                 # Script de inicialización de BD
├── Dockerfile                      # Multi-stage Docker build
├── .env.example                    # Ejemplo de variables de entorno
└── go.mod                          # Dependencias

```

## Instalación

### 1. Clonar el repositorio

```bash
cd backend/users-api
```

### 2. Copiar y configurar variables de entorno

```bash
cp .env.example .env
```

Editar `.env` con tus configuraciones:
- Base de datos MySQL
- JWT secret
- Credenciales SMTP
- URL de la aplicación frontend

### 3. Instalar dependencias

```bash
go mod download
```

### 4. Crear la base de datos

```bash
mysql -u root -p < scripts/init_db.sql
```

O dejar que GORM auto-migre las tablas al iniciar.

### 5. Ejecutar la aplicación

```bash
go run cmd/api/main.go
```

La API estará disponible en `http://localhost:8001`

## Ejecutar con Docker

```bash
docker build -t users-api .
docker run -p 8001:8001 --env-file .env users-api
```

## Ejecutar Tests

```bash
go test ./internal/service/... -v
```

## API Endpoints

### Rutas Públicas (sin autenticación)

#### Registro y Login
- `POST /users` - Registro de nuevo usuario
- `POST /login` - Autenticación, retorna JWT

#### Verificación de Email
- `GET /verify-email?token=xxx` - Verificar email
- `POST /resend-verification` - Reenviar email de verificación

#### Recuperación de Contraseña
- `POST /forgot-password` - Solicitar reset de contraseña
- `POST /reset-password` - Restablecer contraseña con token

### Rutas Protegidas (requieren JWT)

Incluir header: `Authorization: Bearer <token>`

#### Gestión de Usuario
- `GET /users/me` - Obtener perfil del usuario autenticado
- `GET /users/:id` - Obtener información de un usuario por ID
- `PUT /users/:id` - Actualizar perfil (solo el propio usuario)
- `DELETE /users/:id` - Eliminar cuenta (solo el propio usuario)
- `POST /change-password` - Cambiar contraseña

#### Calificaciones
- `GET /users/:id/ratings?page=1&limit=10` - Obtener calificaciones de un usuario (paginado)

### Rutas Internas (comunicación entre servicios)

- `POST /internal/ratings` - Crear calificación (llamado desde trips-api)

### Health Check

- `GET /health` - Verificar estado del servicio

## Formato de Respuestas

Todas las respuestas siguen el formato:

```json
{
  "success": true,
  "data": { ... }
}
```

O en caso de error:

```json
{
  "success": false,
  "error": "mensaje de error"
}
```

## Códigos HTTP

- `200 OK` - Operación exitosa
- `201 Created` - Recurso creado exitosamente
- `400 Bad Request` - Datos inválidos
- `401 Unauthorized` - No autenticado
- `403 Forbidden` - No autorizado
- `404 Not Found` - Recurso no encontrado
- `409 Conflict` - Conflicto (ej: email ya existe)
- `500 Internal Server Error` - Error del servidor

## Ejemplos de Uso

### Registro de Usuario

```bash
curl -X POST http://localhost:8001/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@example.com",
    "password": "password123",
    "name": "Juan",
    "lastname": "Pérez",
    "phone": "123456789",
    "street": "Calle Falsa",
    "number": 123,
    "sex": "hombre",
    "birthdate": "1990-01-15"
  }'
```

### Login

```bash
curl -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@example.com",
    "password": "password123"
  }'
```

### Obtener Perfil (con JWT)

```bash
curl -X GET http://localhost:8001/users/me \
  -H "Authorization: Bearer <tu-jwt-token>"
```

### Crear Calificación (Internal API)

```bash
curl -X POST http://localhost:8001/internal/ratings \
  -H "Content-Type: application/json" \
  -d '{
    "rater_id": 1,
    "rated_user_id": 2,
    "trip_id": "trip123",
    "role_rated": "conductor",
    "score": 5,
    "comment": "Excelente conductor"
  }'
```

## Validaciones Críticas

- ✅ **Passwords**: Hasheadas con bcrypt cost 10
- ✅ **JWT**: Expira en 24 horas, contiene: user_id, email, role
- ✅ **Email**: Validación de formato
- ✅ **Score de rating**: Entre 1-5
- ✅ **Ratings duplicados**: No se permiten
- ✅ **Actualización de perfil**: Solo el propio usuario
- ✅ **TripID**: String (VARCHAR 24) para compatibilidad con MongoDB

## Seguridad

- Contraseñas hasheadas con bcrypt (cost 10)
- JWT con expiración de 24 horas
- CORS configurado
- No se revela información sensible en errores
- Prevención de enumeration attacks en reset de contraseña
- Validación de permisos en actualización/eliminación de perfil

## Desarrollo

### Agregar nueva migración
GORM auto-migra automáticamente. Para cambios manuales:

```bash
mysql -u root -p carpooling_users < scripts/migration.sql
```

### Ejecutar tests con coverage

```bash
go test ./... -cover
```

## Troubleshooting

### Error de conexión a MySQL
- Verificar que MySQL esté corriendo
- Verificar credenciales en `.env`
- Verificar que la base de datos existe

### Error de envío de email
- Verificar configuración SMTP en `.env`
- Para Gmail: habilitar "App Passwords"
- Verificar firewall/red

### JWT inválido
- Verificar que JWT_SECRET coincida
- Verificar que el token no haya expirado (24h)

## Contribuir

1. Fork el proyecto
2. Crear una rama (`git checkout -b feature/nueva-caracteristica`)
3. Commit cambios (`git commit -am 'Agregar nueva característica'`)
4. Push a la rama (`git push origin feature/nueva-caracteristica`)
5. Crear un Pull Request

## Licencia

Este proyecto es parte del sistema CarPooling.
