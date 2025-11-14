# Script de diagnóstico para Search API
# Ejecuta: powershell -ExecutionPolicy Bypass -File debug-search.ps1

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "🔍 DIAGNÓSTICO DE BÚSQUEDA - Search API" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# 1. Verificar servicios
Write-Host "1. Verificando servicios..." -ForegroundColor Blue
Write-Host ""

# MongoDB
Write-Host "📦 MongoDB (puerto 27017):" -ForegroundColor Yellow
try {
    $null = Test-NetConnection -ComputerName localhost -Port 27017 -WarningAction SilentlyContinue -ErrorAction Stop
    Write-Host "✓ MongoDB está corriendo" -ForegroundColor Green
} catch {
    Write-Host "✗ MongoDB no responde" -ForegroundColor Red
}

# Solr
Write-Host "🔍 Apache Solr (puerto 8983):" -ForegroundColor Yellow
try {
    $solrResponse = Invoke-RestMethod -Uri "http://localhost:8983/solr/admin/cores?action=STATUS" -Method Get -ErrorAction Stop
    if ($solrResponse.status.carpooling_trips) {
        Write-Host "✓ Solr está corriendo con el core 'carpooling_trips'" -ForegroundColor Green
    } else {
        Write-Host "✗ Core 'carpooling_trips' no encontrado" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Solr no responde: $($_.Exception.Message)" -ForegroundColor Red
}

# Memcached
Write-Host "💾 Memcached (puerto 11211):" -ForegroundColor Yellow
try {
    $null = Test-NetConnection -ComputerName localhost -Port 11211 -WarningAction SilentlyContinue -ErrorAction Stop
    Write-Host "✓ Memcached está corriendo" -ForegroundColor Green
} catch {
    Write-Host "✗ Memcached no responde" -ForegroundColor Red
}

# RabbitMQ
Write-Host "🐰 RabbitMQ (puerto 5672):" -ForegroundColor Yellow
try {
    $null = Test-NetConnection -ComputerName localhost -Port 5672 -WarningAction SilentlyContinue -ErrorAction Stop
    Write-Host "✓ RabbitMQ está corriendo" -ForegroundColor Green
} catch {
    Write-Host "✗ RabbitMQ no responde" -ForegroundColor Red
}

# Search API
Write-Host "🚀 Search API (puerto 8004):" -ForegroundColor Yellow
try {
    $healthCheck = Invoke-RestMethod -Uri "http://localhost:8004/health" -Method Get -ErrorAction Stop
    Write-Host "✓ Search API está corriendo" -ForegroundColor Green
    Write-Host "Health check:" -ForegroundColor Gray
    $healthCheck | ConvertTo-Json -Depth 10
} catch {
    Write-Host "✗ Search API no responde: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "2. Verificando datos en MongoDB..." -ForegroundColor Blue
Write-Host ""

# Verificar si mongosh está disponible
if (Get-Command mongosh -ErrorAction SilentlyContinue) {
    Write-Host "Ejecutando consultas en MongoDB..." -ForegroundColor Yellow

    # Contar viajes
    $tripCount = mongosh --quiet --eval "db.getSiblingDB('carpooling_search').trips.countDocuments()"
    Write-Host "Total de viajes en MongoDB: $tripCount" -ForegroundColor Yellow

    if ([int]$tripCount -gt 0) {
        Write-Host "✓ Hay viajes en la base de datos" -ForegroundColor Green

        Write-Host ""
        Write-Host "Viajes por estado:" -ForegroundColor Yellow
        mongosh --quiet --eval "db.getSiblingDB('carpooling_search').trips.aggregate([{`$group: {_id: '`$status', count: {`$sum: 1}}}])"

        Write-Host ""
        Write-Host "Ejemplo de viaje:" -ForegroundColor Yellow
        mongosh --quiet --eval "db.getSiblingDB('carpooling_search').trips.findOne({}, {trip_id: 1, 'origin.city': 1, 'destination.city': 1, status: 1, available_seats: 1, departure_datetime: 1})" | ConvertFrom-Json | ConvertTo-Json -Depth 10
    } else {
        Write-Host "✗ No hay viajes en la base de datos" -ForegroundColor Red
        Write-Host "Necesitas crear viajes primero o sincronizarlos desde trips-api" -ForegroundColor Yellow
    }
} else {
    Write-Host "✗ mongosh no está instalado o no está en el PATH" -ForegroundColor Red
    Write-Host "Instala MongoDB Shell desde: https://www.mongodb.com/try/download/shell" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "3. Verificando documentos en Solr..." -ForegroundColor Blue
Write-Host ""

try {
    $solrQuery = Invoke-RestMethod -Uri "http://localhost:8983/solr/carpooling_trips/select?q=*:*&rows=0" -Method Get -ErrorAction Stop
    $solrCount = $solrQuery.response.numFound

    Write-Host "Total de documentos en Solr: $solrCount" -ForegroundColor Yellow

    if ($solrCount -gt 0) {
        Write-Host "✓ Hay documentos indexados en Solr" -ForegroundColor Green

        Write-Host ""
        Write-Host "Ejemplo de documento en Solr:" -ForegroundColor Yellow
        $solrDoc = Invoke-RestMethod -Uri "http://localhost:8983/solr/carpooling_trips/select?q=*:*&rows=1&fl=id,origin_city,destination_city,status,available_seats" -Method Get
        $solrDoc.response.docs[0] | ConvertTo-Json -Depth 10
    } else {
        Write-Host "✗ No hay documentos indexados en Solr" -ForegroundColor Red
        Write-Host "Los viajes de MongoDB no se han sincronizado a Solr" -ForegroundColor Yellow
        Write-Host "Verifica que RabbitMQ esté publicando eventos trip.created" -ForegroundColor Yellow
    }
} catch {
    Write-Host "✗ No se pudo consultar Solr: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "4. Probando búsqueda directa en la API..." -ForegroundColor Blue
Write-Host ""

# Búsqueda sin filtros
Write-Host "Búsqueda sin filtros (todos los viajes):" -ForegroundColor Yellow
try {
    $searchResult = Invoke-RestMethod -Uri "http://localhost:8004/api/v1/search/trips?page=1&limit=5" -Method Get -ErrorAction Stop

    Write-Host "✓ API respondió correctamente" -ForegroundColor Green
    Write-Host "Total de viajes encontrados: $($searchResult.data.total)" -ForegroundColor Cyan
    Write-Host "Viajes en esta página: $($searchResult.data.trips.Count)" -ForegroundColor Cyan

    Write-Host ""
    Write-Host "Respuesta completa:" -ForegroundColor Gray
    $searchResult | ConvertTo-Json -Depth 10

    if ($searchResult.data.total -eq 0) {
        Write-Host ""
        Write-Host "⚠️ PROBLEMA ENCONTRADO:" -ForegroundColor Red
        Write-Host "La API responde, pero no hay viajes en los resultados." -ForegroundColor Yellow
        Write-Host ""
        Write-Host "Posibles causas:" -ForegroundColor Yellow
        Write-Host "  1. Los viajes en MongoDB tienen status != 'published'" -ForegroundColor White
        Write-Host "  2. Los viajes tienen available_seats = 0" -ForegroundColor White
        Write-Host "  3. Hay un problema con los filtros en el backend" -ForegroundColor White
        Write-Host "  4. Solr está fallando y MongoDB también devuelve 0 resultados" -ForegroundColor White
    }
} catch {
    Write-Host "✗ La API no respondió correctamente: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "5. Consultas de depuración adicionales..." -ForegroundColor Blue
Write-Host ""

if (Get-Command mongosh -ErrorAction SilentlyContinue) {
    Write-Host "Viajes con status='published' y available_seats > 0:" -ForegroundColor Yellow
    $publishedCount = mongosh --quiet --eval "db.getSiblingDB('carpooling_search').trips.countDocuments({status: 'published', available_seats: {`$gt: 0}})"
    Write-Host "Cantidad: $publishedCount" -ForegroundColor Cyan

    if ([int]$publishedCount -eq 0) {
        Write-Host ""
        Write-Host "⚠️ PROBLEMA ENCONTRADO:" -ForegroundColor Red
        Write-Host "No hay viajes con status='published' Y available_seats > 0" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "Todos los viajes encontrados (con sus status):" -ForegroundColor Yellow
        mongosh --quiet --eval "db.getSiblingDB('carpooling_search').trips.aggregate([{`$group: {_id: {status: '`$status', has_seats: {`$cond: [{`$gt: ['`$available_seats', 0]}, 'yes', 'no']}}, count: {`$sum: 1}}}, {`$sort: {count: -1}}])"
    }
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "✅ Diagnóstico completado" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Para ver logs en tiempo real del Search API:" -ForegroundColor Yellow
Write-Host "  cd backend\search-api" -ForegroundColor White
Write-Host "  go run cmd\api\main.go" -ForegroundColor White
Write-Host ""
Write-Host "Los logs mostrarán:" -ForegroundColor Yellow
Write-Host "  - Solicitudes HTTP entrantes con parámetros" -ForegroundColor White
Write-Host "  - Consultas a MongoDB con filtros" -ForegroundColor White
Write-Host "  - Intentos de búsqueda en Solr" -ForegroundColor White
Write-Host "  - Resultados y errores" -ForegroundColor White
