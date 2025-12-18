#!/bin/bash

# Script para resetear la verificación de email de todos los usuarios
# Ejecutar en Kali Linux donde corre el contenedor MySQL

set -e

echo "========================================="
echo "Resetear verificación de email"
echo "========================================="

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo ""
echo -e "${YELLOW}⚠️  ADVERTENCIA: Esto marcará TODOS los usuarios como NO verificados${NC}"
echo -e "${YELLOW}⚠️  ¿Estás seguro de continuar? (y/n)${NC}"
read -r response

if [[ ! "$response" =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}Operación cancelada${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}Conectando a la base de datos...${NC}"

# Verificar que el contenedor MySQL esté corriendo
if ! docker ps --format '{{.Names}}' | grep -q '^carpooling-mysql$'; then
    echo -e "${RED}Error: El contenedor carpooling-mysql no está corriendo${NC}"
    echo "Inicia el contenedor con: docker start carpooling-mysql"
    exit 1
fi

echo -e "${BLUE}Mostrando estado actual de usuarios...${NC}"
docker exec -it carpooling-mysql mysql -uroot -proot carpooling -e "
SELECT
    COUNT(*) as total_usuarios,
    SUM(CASE WHEN email_verified = 1 THEN 1 ELSE 0 END) as verificados,
    SUM(CASE WHEN email_verified = 0 THEN 1 ELSE 0 END) as no_verificados
FROM users;"

echo ""
echo -e "${BLUE}Actualizando usuarios...${NC}"
docker exec -it carpooling-mysql mysql -uroot -proot carpooling -e "
UPDATE users SET email_verified = false;"

echo ""
echo -e "${GREEN}✓ Actualización completada${NC}"

echo ""
echo -e "${BLUE}Mostrando estado actualizado...${NC}"
docker exec -it carpooling-mysql mysql -uroot -proot carpooling -e "
SELECT
    COUNT(*) as total_usuarios,
    SUM(CASE WHEN email_verified = 1 THEN 1 ELSE 0 END) as verificados,
    SUM(CASE WHEN email_verified = 0 THEN 1 ELSE 0 END) as no_verificados
FROM users;"

echo ""
echo -e "${BLUE}Primeros 10 usuarios:${NC}"
docker exec -it carpooling-mysql mysql -uroot -proot carpooling -e "
SELECT id, email, email_verified, created_at
FROM users
ORDER BY created_at DESC
LIMIT 10;"

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}✓ Todos los usuarios marcados como NO verificados${NC}"
echo -e "${GREEN}=========================================${NC}"
