#!/bin/bash
# Solr initialization script for carpooling_trips core
# This script defines the schema for the search-api Solr core

set -e

CORE_NAME="carpooling_trips"
SOLR_URL="${SOLR_URL:-http://localhost:8983/solr}"

echo "Initializing Solr core: $CORE_NAME"

# Wait for Solr to be ready
echo "Waiting for Solr to be ready..."
until curl -sf "${SOLR_URL}/admin/cores?action=STATUS" > /dev/null; do
  echo "Solr is unavailable - sleeping"
  sleep 2
done

echo "Solr is ready!"

# Check if core exists
CORE_EXISTS=$(curl -s "${SOLR_URL}/admin/cores?action=STATUS&core=${CORE_NAME}" | grep -c "\"name\":\"${CORE_NAME}\"" || true)

if [ "$CORE_EXISTS" -eq "0" ]; then
  echo "Creating core: $CORE_NAME"
  curl -X POST "${SOLR_URL}/admin/cores?action=CREATE&name=${CORE_NAME}&configSet=_default"
  echo "Core created successfully"
else
  echo "Core $CORE_NAME already exists"
fi

# Define schema fields
echo "Defining schema fields..."

# Trip ID (unique identifier)
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "trip_id",
    "type": "string",
    "indexed": true,
    "stored": true,
    "required": true
  }
}' 2>/dev/null || true

# Driver ID
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "driver_id",
    "type": "plong",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Origin and Destination - Text fields for city names
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "origin_city",
    "type": "text_general",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "origin_province",
    "type": "text_general",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "destination_city",
    "type": "text_general",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "destination_province",
    "type": "text_general",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Geospatial fields (location type for lat/lng)
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "origin_coordinates",
    "type": "location",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "destination_coordinates",
    "type": "location",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Date field for departure
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "departure_datetime",
    "type": "pdate",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Numeric fields
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "price_per_seat",
    "type": "pdouble",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "available_seats",
    "type": "pint",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "total_seats",
    "type": "pint",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Popularity score
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "popularity_score",
    "type": "pdouble",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Driver rating
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "driver_rating",
    "type": "pdouble",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Status field
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "status",
    "type": "string",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Boolean preferences
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "pets_allowed",
    "type": "boolean",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "smoking_allowed",
    "type": "boolean",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "music_allowed",
    "type": "boolean",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Full-text search field
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "search_text",
    "type": "text_general",
    "indexed": true,
    "stored": true,
    "multiValued": false
  }
}' 2>/dev/null || true

# Description
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "description",
    "type": "text_general",
    "indexed": true,
    "stored": true
  }
}' 2>/dev/null || true

# Set unique key (if not already set)
curl -X POST "${SOLR_URL}/${CORE_NAME}/schema" -H 'Content-type:application/json' -d '{
  "add-field": {
    "name": "id",
    "type": "string",
    "indexed": true,
    "stored": true,
    "required": true
  }
}' 2>/dev/null || true

echo "Schema configuration completed!"
echo "Solr core $CORE_NAME is ready at ${SOLR_URL}/${CORE_NAME}"
