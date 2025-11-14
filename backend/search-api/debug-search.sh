#!/bin/bash

echo "=========================================="
echo "🔍 DIAGNÓSTICO DE BÚSQUEDA - Search API"
echo "=========================================="
echo ""

# Colores para mejor visualización
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}1. Verificando servicios...${NC}"
echo ""

# Verificar MongoDB
echo -e "${YELLOW}📦 MongoDB (puerto 27017):${NC}"
if curl -s http://localhost:27017 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ MongoDB está corriendo${NC}"
else
    echo -e "${RED}✗ MongoDB no responde${NC}"
fi

# Verificar Solr
echo -e "${YELLOW}🔍 Apache Solr (puerto 8983):${NC}"
SOLR_RESPONSE=$(curl -s http://localhost:8983/solr/admin/cores?action=STATUS 2>&1)
if echo "$SOLR_RESPONSE" | grep -q "carpooling_trips"; then
    echo -e "${GREEN}✓ Solr está corriendo con el core 'carpooling_trips'${NC}"
else
    echo -e "${RED}✗ Solr no responde o el core no existe${NC}"
    echo "Respuesta: $SOLR_RESPONSE"
fi

# Verificar Memcached
echo -e "${YELLOW}💾 Memcached (puerto 11211):${NC}"
if nc -z localhost 11211 2>/dev/null; then
    echo -e "${GREEN}✓ Memcached está corriendo${NC}"
else
    echo -e "${RED}✗ Memcached no responde${NC}"
fi

# Verificar RabbitMQ
echo -e "${YELLOW}🐰 RabbitMQ (puerto 5672):${NC}"
if nc -z localhost 5672 2>/dev/null; then
    echo -e "${GREEN}✓ RabbitMQ está corriendo${NC}"
else
    echo -e "${RED}✗ RabbitMQ no responde${NC}"
fi

# Verificar Search API
echo -e "${YELLOW}🚀 Search API (puerto 8004):${NC}"
HEALTH_CHECK=$(curl -s http://localhost:8004/health 2>&1)
if echo "$HEALTH_CHECK" | grep -q "ok\|degraded"; then
    echo -e "${GREEN}✓ Search API está corriendo${NC}"
    echo "Health check: $HEALTH_CHECK" | jq '.' 2>/dev/null || echo "$HEALTH_CHECK"
else
    echo -e "${RED}✗ Search API no responde${NC}"
fi

echo ""
echo -e "${BLUE}2. Verificando datos en MongoDB...${NC}"
echo ""

# Contar viajes en MongoDB
TRIP_COUNT=$(mongosh --quiet --eval "db.getSiblingDB('carpooling_search').trips.countDocuments()" 2>/dev/null)
if [ ! -z "$TRIP_COUNT" ]; then
    echo -e "${YELLOW}Total de viajes en MongoDB:${NC} $TRIP_COUNT"

    if [ "$TRIP_COUNT" -gt 0 ]; then
        echo -e "${GREEN}✓ Hay viajes en la base de datos${NC}"

        # Mostrar un viaje de ejemplo
        echo ""
        echo -e "${YELLOW}Ejemplo de viaje:${NC}"
        mongosh --quiet --eval "db.getSiblingDB('carpooling_search').trips.findOne({}, {trip_id: 1, 'origin.city': 1, 'destination.city': 1, status: 1, available_seats: 1})" 2>/dev/null | head -20
    else
        echo -e "${RED}✗ No hay viajes en la base de datos${NC}"
    fi
else
    echo -e "${RED}✗ No se pudo conectar a MongoDB${NC}"
fi

echo ""
echo -e "${BLUE}3. Verificando documentos en Solr...${NC}"
echo ""

# Contar documentos en Solr
SOLR_COUNT=$(curl -s "http://localhost:8983/solr/carpooling_trips/select?q=*:*&rows=0" | jq -r '.response.numFound' 2>/dev/null)
if [ ! -z "$SOLR_COUNT" ] && [ "$SOLR_COUNT" != "null" ]; then
    echo -e "${YELLOW}Total de documentos en Solr:${NC} $SOLR_COUNT"

    if [ "$SOLR_COUNT" -gt 0 ]; then
        echo -e "${GREEN}✓ Hay documentos indexados en Solr${NC}"

        # Mostrar un documento de ejemplo
        echo ""
        echo -e "${YELLOW}Ejemplo de documento en Solr:${NC}"
        curl -s "http://localhost:8983/solr/carpooling_trips/select?q=*:*&rows=1&fl=id,origin_city,destination_city,status,available_seats" | jq '.response.docs[0]' 2>/dev/null
    else
        echo -e "${RED}✗ No hay documentos indexados en Solr${NC}"
        echo -e "${YELLOW}Esto significa que los viajes no se han sincronizado de MongoDB a Solr${NC}"
    fi
else
    echo -e "${RED}✗ No se pudo consultar Solr${NC}"
fi

echo ""
echo -e "${BLUE}4. Probando búsqueda directa en la API...${NC}"
echo ""

# Probar búsqueda simple
echo -e "${YELLOW}Búsqueda sin filtros (todos los viajes):${NC}"
SEARCH_RESULT=$(curl -s "http://localhost:8004/api/v1/search/trips?page=1&limit=5" 2>&1)
if echo "$SEARCH_RESULT" | grep -q "trips"; then
    TRIP_COUNT_API=$(echo "$SEARCH_RESULT" | jq -r '.data.total' 2>/dev/null)
    echo -e "${GREEN}✓ API respondió correctamente${NC}"
    echo "Total de viajes encontrados: $TRIP_COUNT_API"
    echo ""
    echo "Respuesta completa:"
    echo "$SEARCH_RESULT" | jq '.' 2>/dev/null || echo "$SEARCH_RESULT"
else
    echo -e "${RED}✗ La API no respondió correctamente${NC}"
    echo "Respuesta: $SEARCH_RESULT"
fi

echo ""
echo -e "${BLUE}5. Verificando logs recientes del Search API...${NC}"
echo ""
echo -e "${YELLOW}Últimas 50 líneas de logs (si el servicio está corriendo en terminal):${NC}"
echo "Para ver logs en tiempo real, ejecuta:"
echo -e "${GREEN}cd backend/search-api && go run cmd/api/main.go${NC}"
echo ""
echo "Los logs incluirán:"
echo "  - Solicitudes HTTP entrantes"
echo "  - Consultas a MongoDB"
echo "  - Consultas a Solr"
echo "  - Resultados de búsqueda"
echo "  - Errores (si los hay)"

echo ""
echo "=========================================="
echo "✅ Diagnóstico completado"
echo "=========================================="
