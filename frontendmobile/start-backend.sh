#!/bin/bash

# Script para iniciar todos los servicios del backend CarPooling
# Ejecutar en Kali Linux

set -e  # Salir si hay algún error

echo "========================================="
echo "Iniciando servicios de CarPooling Backend"
echo "========================================="

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Variables de configuración - AJUSTAR SEGÚN TU SETUP
BACKEND_DIR="${HOME}/Carpooling/backend"
USERS_API_DIR="${BACKEND_DIR}/users-api"
TRIPS_API_DIR="${BACKEND_DIR}/trips-api"
BOOKINGS_API_DIR="${BACKEND_DIR}/bookings-api"
SEARCH_API_DIR="${BACKEND_DIR}/search-api"

# Puertos
USERS_API_PORT=8001
TRIPS_API_PORT=8002
BOOKINGS_API_PORT=8003
SEARCH_API_PORT=8004

# Función para verificar si un puerto está en uso
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        return 0  # Puerto en uso
    else
        return 1  # Puerto libre
    fi
}

# Función para matar proceso en puerto
kill_port() {
    local port=$1
    local pid=$(lsof -ti:$port)
    if [ ! -z "$pid" ]; then
        echo -e "${BLUE}Matando proceso en puerto $port (PID: $pid)${NC}"
        kill -9 $pid
        sleep 1
    fi
}

echo ""
echo -e "${GREEN}[1/6] Verificando Docker...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker no está instalado${NC}"
    exit 1
fi

if ! docker info &> /dev/null; then
    echo -e "${RED}Error: Docker daemon no está corriendo${NC}"
    echo "Inicia Docker con: sudo systemctl start docker"
    exit 1
fi
echo -e "${GREEN}✓ Docker está corriendo${NC}"

echo ""
echo -e "${GREEN}[2/6] Iniciando contenedores Docker...${NC}"

# MongoDB
echo -e "${BLUE}Iniciando MongoDB...${NC}"
if docker ps -a --format '{{.Names}}' | grep -q '^carpooling-mongodb$'; then
    docker start carpooling-mongodb
else
    docker run -d \
        --name carpooling-mongodb \
        -p 27017:27017 \
        -v carpooling-mongo-data:/data/db \
        mongo:4.4
fi
echo -e "${GREEN}✓ MongoDB iniciado en puerto 27017${NC}"

# MySQL
echo -e "${BLUE}Iniciando MySQL...${NC}"
if docker ps -a --format '{{.Names}}' | grep -q '^carpooling-mysql$'; then
    docker start carpooling-mysql
else
    docker run -d \
        --name carpooling-mysql \
        -p 3306:3306 \
        -e MYSQL_ROOT_PASSWORD=root \
        -e MYSQL_DATABASE=carpooling \
        -v carpooling-mysql-data:/var/lib/mysql \
        mysql:8.0
fi
echo -e "${GREEN}✓ MySQL iniciado en puerto 3306${NC}"

# RabbitMQ
echo -e "${BLUE}Iniciando RabbitMQ...${NC}"
if docker ps -a --format '{{.Names}}' | grep -q '^carpooling-rabbitmq$'; then
    docker start carpooling-rabbitmq
else
    docker run -d \
        --name carpooling-rabbitmq \
        -p 5672:5672 \
        -p 15672:15672 \
        -e RABBITMQ_DEFAULT_USER=guest \
        -e RABBITMQ_DEFAULT_PASS=guest \
        rabbitmq:3-management
fi
echo -e "${GREEN}✓ RabbitMQ iniciado en puerto 5672 (Management: 15672)${NC}"

# Solr
echo -e "${BLUE}Iniciando Solr...${NC}"
if docker ps -a --format '{{.Names}}' | grep -q '^carpooling-solr$'; then
    docker start carpooling-solr
else
    docker run -d \
        --name carpooling-solr \
        -p 8983:8983 \
        -v carpooling-solr-data:/var/solr \
        solr:8
fi
echo -e "${GREEN}✓ Solr iniciado en puerto 8983${NC}"

# Memcached
echo -e "${BLUE}Iniciando Memcached...${NC}"
if docker ps -a --format '{{.Names}}' | grep -q '^carpooling-memcached$'; then
    docker start carpooling-memcached
else
    docker run -d \
        --name carpooling-memcached \
        -p 11211:11211 \
        memcached:alpine
fi
echo -e "${GREEN}✓ Memcached iniciado en puerto 11211${NC}"

echo ""
echo -e "${GREEN}[3/6] Esperando a que los servicios estén listos...${NC}"
sleep 5

echo ""
echo -e "${GREEN}[4/6] Verificando puertos del backend...${NC}"
# Matar procesos en puertos si están ocupados
for port in $USERS_API_PORT $TRIPS_API_PORT $BOOKINGS_API_PORT $SEARCH_API_PORT; do
    if check_port $port; then
        echo -e "${BLUE}Puerto $port está en uso${NC}"
        kill_port $port
    fi
done

echo ""
echo -e "${GREEN}[5/6] Iniciando servicios del backend...${NC}"

# Verificar que los directorios existen
if [ ! -d "$BACKEND_DIR" ]; then
    echo -e "${RED}Error: Directorio del backend no encontrado: $BACKEND_DIR${NC}"
    echo "Por favor ajusta la variable BACKEND_DIR en el script"
    exit 1
fi

# Función para compilar y ejecutar API en Go
start_go_api() {
    local api_name=$1
    local api_dir=$2
    local port=$3

    echo -e "${BLUE}Compilando y ejecutando $api_name (puerto $port)...${NC}"

    if [ ! -d "$api_dir" ]; then
        echo -e "${RED}⨯ Directorio no encontrado: $api_dir${NC}"
        return 1
    fi

    cd "$api_dir"

    # Compilar el binario si no existe o si el código fuente es más reciente
    local binary_name=$(basename "$api_dir")
    if [ ! -f "$binary_name" ] || [ "cmd/api/main.go" -nt "$binary_name" ]; then
        echo -e "${BLUE}  Compilando $api_name...${NC}"
        go build -o "$binary_name" cmd/api/main.go
        if [ $? -ne 0 ]; then
            echo -e "${RED}⨯ Error al compilar $api_name${NC}"
            return 1
        fi
    fi

    # Ejecutar el binario en background
    ./"$binary_name" > "/tmp/${binary_name}.log" 2>&1 &
    local pid=$!
    echo $pid > "/tmp/${binary_name}.pid"
    echo -e "${GREEN}✓ $api_name iniciado (PID: $pid)${NC}"
}

# Users API
start_go_api "Users API" "$USERS_API_DIR" "$USERS_API_PORT"

# Trips API
start_go_api "Trips API" "$TRIPS_API_DIR" "$TRIPS_API_PORT"

# Bookings API
start_go_api "Bookings API" "$BOOKINGS_API_DIR" "$BOOKINGS_API_PORT"

# Search API
start_go_api "Search API" "$SEARCH_API_DIR" "$SEARCH_API_PORT"

echo ""
echo -e "${GREEN}[6/6] Esperando a que los servicios estén listos...${NC}"
sleep 5

echo ""
echo "========================================="
echo -e "${GREEN}Estado de los servicios:${NC}"
echo "========================================="
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "========================================="
echo -e "${GREEN}Backend APIs:${NC}"
echo "========================================="
echo "Users API:    http://localhost:$USERS_API_PORT"
echo "Trips API:    http://localhost:$TRIPS_API_PORT"
echo "Bookings API: http://localhost:$BOOKINGS_API_PORT"
echo "Search API:   http://localhost:$SEARCH_API_PORT"

echo ""
echo "Logs guardados en:"
echo "  - /tmp/users-api.log"
echo "  - /tmp/trips-api.log"
echo "  - /tmp/bookings-api.log"
echo "  - /tmp/search-api.log"

echo ""
echo "PIDs guardados en:"
echo "  - /tmp/users-api.pid"
echo "  - /tmp/trips-api.pid"
echo "  - /tmp/bookings-api.pid"
echo "  - /tmp/search-api.pid"

echo ""
echo -e "${GREEN}✓ Todos los servicios iniciados correctamente${NC}"
echo ""
echo "Para detener los servicios, ejecuta: ./stop-backend.sh"
