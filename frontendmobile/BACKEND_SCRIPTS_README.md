# Scripts de Backend - CarPooling

Scripts para iniciar y detener todos los servicios del backend de CarPooling en Kali Linux.

## 📋 Archivos

- `start-backend.sh` - Inicia todos los servicios (Docker + APIs)
- `stop-backend.sh` - Detiene todos los servicios

## 🚀 Instalación en Kali

1. **Copia los scripts a Kali:**
   ```bash
   # Opción 1: Usando SCP desde Windows
   scp start-backend.sh stop-backend.sh usuario@kali-ip:~/

   # Opción 2: Copiar manualmente o usar Git
   ```

2. **Dale permisos de ejecución:**
   ```bash
   chmod +x start-backend.sh stop-backend.sh
   ```

## ▶️ Iniciar todos los servicios

```bash
./start-backend.sh
```

Esto hará:
1. ✅ Verifica que Docker esté corriendo
2. ✅ Inicia contenedores Docker:
   - MongoDB (puerto 27017)
   - MySQL (puerto 3306)
   - RabbitMQ (puerto 5672, management 15672)
   - Solr (puerto 8983)
   - Memcached (puerto 11211)
3. ✅ Compila y ejecuta las APIs en Go:
   - users-api (puerto 8001)
   - trips-api (puerto 8002)
   - bookings-api (puerto 8003)
   - search-api (puerto 8004)

## ⏹️ Detener todos los servicios

```bash
./stop-backend.sh
```

## 📊 Ver logs

Los logs de las APIs se guardan en:
```bash
tail -f /tmp/users-api.log
tail -f /tmp/trips-api.log
tail -f /tmp/bookings-api.log
tail -f /tmp/search-api.log
```

## 🔍 Verificar estado

### Ver contenedores Docker:
```bash
docker ps
```

### Ver procesos de las APIs:
```bash
ps aux | grep -E "users-api|trips-api|bookings-api|search-api"
```

### Ver puertos en uso:
```bash
sudo lsof -i :8001  # users-api
sudo lsof -i :8002  # trips-api
sudo lsof -i :8003  # bookings-api
sudo lsof -i :8004  # search-api
```

## ⚙️ Configuración

Si tu backend está en una ruta diferente, edita `start-backend.sh` y cambia:

```bash
BACKEND_DIR="${HOME}/Carpooling/backend"  # Línea 19
```

## 🛠️ Solución de problemas

### Docker no está corriendo:
```bash
sudo systemctl start docker
```

### Puerto ocupado:
```bash
# Matar proceso en puerto específico
sudo kill -9 $(lsof -ti:8001)
```

### Recompilar forzadamente:
```bash
cd ~/Carpooling/backend/users-api
rm users-api  # Eliminar binario
./start-backend.sh  # Volverá a compilar
```

### Ver logs de Docker:
```bash
docker logs carpooling-mongodb
docker logs carpooling-mysql
docker logs carpooling-rabbitmq
docker logs carpooling-solr
docker logs carpooling-memcached
```

## 📝 Notas

- Los binarios compilados se guardan en sus respectivos directorios de API
- Los PIDs se guardan en `/tmp/*.pid`
- Los logs se guardan en `/tmp/*.log`
- Si ya existe un contenedor Docker, el script lo inicia en lugar de crear uno nuevo
- La compilación solo ocurre si el binario no existe o si el código fuente cambió

## 🔄 Reiniciar un servicio específico

```bash
# Detener
kill $(cat /tmp/users-api.pid)

# Iniciar
cd ~/Carpooling/backend/users-api
./users-api > /tmp/users-api.log 2>&1 &
echo $! > /tmp/users-api.pid
```
