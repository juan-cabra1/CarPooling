#!/bin/bash

# Script para detener todos los servicios del backend CarPooling

set -e

echo "========================================="
echo "Deteniendo servicios de CarPooling Backend"
echo "========================================="

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo ""
echo -e "${GREEN}[1/2] Deteniendo servicios del backend...${NC}"

# Función para detener un servicio
stop_service() {
    local service_name=$1
    local pid_file=$2

    if [ -f "$pid_file" ]; then
        PID=$(cat "$pid_file")
        if ps -p $PID > /dev/null 2>&1; then
            echo -e "${BLUE}Deteniendo $service_name (PID: $PID)...${NC}"
            kill $PID
            rm "$pid_file"
            echo -e "${GREEN}✓ $service_name detenido${NC}"
        else
            echo -e "${BLUE}$service_name ya no está corriendo (PID obsoleto)${NC}"
            rm "$pid_file"
        fi
    fi
}

# Detener todas las APIs
stop_service "Users API" "/tmp/users-api.pid"
stop_service "Trips API" "/tmp/trips-api.pid"
stop_service "Bookings API" "/tmp/bookings-api.pid"
stop_service "Search API" "/tmp/search-api.pid"

# Matar cualquier proceso en los puertos por si acaso
echo -e "${BLUE}Verificando puertos...${NC}"
for port in 8001 8002 8003 8004; do
    PID=$(lsof -ti:$port 2>/dev/null || true)
    if [ ! -z "$PID" ]; then
        echo -e "${BLUE}Matando proceso en puerto $port (PID: $PID)...${NC}"
        kill -9 $PID
    fi
done

echo ""
echo -e "${GREEN}[2/2] Deteniendo contenedores Docker...${NC}"

# Detener contenedores
for container in carpooling-mongodb carpooling-mysql carpooling-rabbitmq carpooling-solr carpooling-memcached; do
    if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
        echo -e "${BLUE}Deteniendo $container...${NC}"
        docker stop $container
        echo -e "${GREEN}✓ $container detenido${NC}"
    fi
done

echo ""
echo "========================================="
echo -e "${GREEN}✓ Todos los servicios detenidos${NC}"
echo "========================================="
echo ""
echo "Para iniciar nuevamente, ejecuta: ./start-backend.sh"
