# Trips API

Microservicio de gestión de viajes para la plataforma CarPooling.

## Fase 1: Project Setup ✅

### Completado

- ✅ Go module inicializado (`trips-api`)
- ✅ Configuración con `.env` y `.env.example`
- ✅ Servidor Gin básico en puerto 8002
- ✅ Endpoint `/health` implementado
- ✅ Conexión a MongoDB configurada
- ✅ Logging estructurado con zerolog
- ✅ Dependencias instaladas

### Tecnologías

- **Go**: 1.21+
- **Framework**: Gin
- **Base de datos**: MongoDB
- **Logging**: zerolog
- **Configuración**: godotenv

## Cómo ejecutar

### Prerrequisitos

1. **MongoDB** debe estar corriendo en `localhost:27017`
   ```bash
   # Con Docker:
   docker run -d -p 27017:27017 --name mongodb mongo:latest

   # O con MongoDB local instalado
   mongod
   ```

2. **Verificar archivo .env**
   ```bash
   cp .env.example .env
   # Editar valores si es necesario
   ```

### Ejecutar el servidor

```bash
cd backend/trips-api
go run cmd/api/main.go
```

### Probar el endpoint /health

```bash
curl http://localhost:8002/health
```

**Respuesta esperada:**
```json
{
  "status": "ok",
  "service": "trips-api",
  "timestamp": "2024-11-06T14:30:00Z"
}
```

## Estructura del proyecto

```
backend/trips-api/
├── cmd/
│   └── api/
│       └── main.go           # Punto de entrada
├── internal/
│   └── config/
│       └── config.go         # Configuración y carga de .env
├── .env                      # Variables de entorno (no commitear)
├── .env.example              # Plantilla de variables
├── go.mod                    # Dependencias
├── go.sum                    # Checksums de dependencias
└── README.md                 # Este archivo
```

## Próximos pasos (Fase 2+)

MongoDB Connection y Domain Models
