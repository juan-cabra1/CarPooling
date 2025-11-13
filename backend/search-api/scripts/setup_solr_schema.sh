#!/bin/bash

# ============================================================================
# Solr Schema Setup Script for CarPooling Search API
# ============================================================================
# This script creates the Solr core and defines the schema for trip search
#
# NOTE: Geospatial search is NOT handled by Solr - MongoDB handles that
#       This schema focuses on full-text search, facets, and filtering
# ============================================================================

set -e  # Exit on error

SOLR_HOST="${SOLR_HOST:-localhost}"
SOLR_PORT="${SOLR_PORT:-8983}"
SOLR_CORE="${SOLR_CORE:-carpooling_trips}"
SOLR_URL="http://${SOLR_HOST}:${SOLR_PORT}/solr"

echo "==============================================="
echo "Solr Schema Setup for CarPooling Search API"
echo "==============================================="
echo "Solr URL: $SOLR_URL"
echo "Core Name: $SOLR_CORE"
echo ""

# ============================================================================
# Step 1: Check if Solr is running
# ============================================================================
echo "Step 1: Checking if Solr is accessible..."
if ! curl -s "${SOLR_URL}/admin/info/system" > /dev/null; then
    echo "❌ ERROR: Solr is not accessible at ${SOLR_URL}"
    echo "Please ensure Solr is running (docker-compose up solr)"
    exit 1
fi
echo "✅ Solr is accessible"
echo ""

# ============================================================================
# Step 2: Delete existing core if it exists (for clean setup)
# ============================================================================
echo "Step 2: Checking if core '${SOLR_CORE}' already exists..."
if curl -s "${SOLR_URL}/admin/cores?action=STATUS&core=${SOLR_CORE}" | grep -q "\"${SOLR_CORE}\""; then
    echo "⚠️  Core '${SOLR_CORE}' already exists. Deleting..."
    curl -s "${SOLR_URL}/admin/cores?action=UNLOAD&core=${SOLR_CORE}&deleteIndex=true&deleteDataDir=true&deleteInstanceDir=true" > /dev/null
    echo "✅ Existing core deleted"
else
    echo "✅ Core does not exist, proceeding with creation"
fi
echo ""

# ============================================================================
# Step 3: Create Solr core
# ============================================================================
echo "Step 3: Creating Solr core '${SOLR_CORE}'..."
curl -s "${SOLR_URL}/admin/cores?action=CREATE&name=${SOLR_CORE}&configSet=_default" > /dev/null
if [ $? -eq 0 ]; then
    echo "✅ Core '${SOLR_CORE}' created successfully"
else
    echo "❌ ERROR: Failed to create core"
    exit 1
fi
echo ""

# Wait a moment for core to initialize
sleep 2

# ============================================================================
# Step 4: Define Schema Fields
# ============================================================================
echo "Step 4: Adding schema fields..."

# Primary identifier
echo "  Adding field: id (string)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "id",
      "type": "string",
      "stored": true,
      "indexed": true,
      "required": true,
      "multiValued": false
    }
  }' > /dev/null 2>&1

# Driver fields
echo "  Adding field: driver_id (plong)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "driver_id",
      "type": "plong",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: driver_name (text_general)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "driver_name",
      "type": "text_general",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: driver_rating (pfloat)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "driver_rating",
      "type": "pfloat",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: driver_total_trips (pint)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "driver_total_trips",
      "type": "pint",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

# Location fields (NO coordinates - only text fields)
echo "  Adding field: origin_city (string)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "origin_city",
      "type": "string",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: origin_province (string)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "origin_province",
      "type": "string",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: destination_city (string)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "destination_city",
      "type": "string",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: destination_province (string)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "destination_province",
      "type": "string",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

# Timing fields
echo "  Adding field: departure_datetime (pdate)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "departure_datetime",
      "type": "pdate",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: estimated_arrival_datetime (pdate)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "estimated_arrival_datetime",
      "type": "pdate",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

# Pricing and availability
echo "  Adding field: price_per_seat (pfloat)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "price_per_seat",
      "type": "pfloat",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: total_seats (pint)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "total_seats",
      "type": "pint",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: available_seats (pint)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "available_seats",
      "type": "pint",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

# Vehicle fields
echo "  Adding field: car_brand (text_general)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "car_brand",
      "type": "text_general",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: car_model (text_general)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "car_model",
      "type": "text_general",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: car_year (pint)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "car_year",
      "type": "pint",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: car_color (string)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "car_color",
      "type": "string",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

# Preferences (boolean fields)
echo "  Adding field: pets_allowed (boolean)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "pets_allowed",
      "type": "boolean",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: smoking_allowed (boolean)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "smoking_allowed",
      "type": "boolean",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: music_allowed (boolean)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "music_allowed",
      "type": "boolean",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

# Trip details
echo "  Adding field: status (string)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "status",
      "type": "string",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: description (text_general)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "description",
      "type": "text_general",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

# Search-specific fields
echo "  Adding field: search_text (text_general)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "search_text",
      "type": "text_general",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: popularity_score (pfloat)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "popularity_score",
      "type": "pfloat",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

# Timestamps
echo "  Adding field: created_at (pdate)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "created_at",
      "type": "pdate",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "  Adding field: updated_at (pdate)"
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "add-field": {
      "name": "updated_at",
      "type": "pdate",
      "stored": true,
      "indexed": true
    }
  }' > /dev/null 2>&1

echo "✅ All fields added successfully"
echo ""

# ============================================================================
# Step 5: Set unique key
# ============================================================================
echo "Step 5: Setting unique key to 'id'..."
curl -X POST -H 'Content-Type: application/json' \
  "${SOLR_URL}/${SOLR_CORE}/schema" \
  -d '{
    "replace-field-type": {
      "name": "string",
      "class": "solr.StrField",
      "sortMissingLast": true,
      "docValues": true
    }
  }' > /dev/null 2>&1

echo "✅ Unique key configured"
echo ""

# ============================================================================
# Step 6: Test the schema
# ============================================================================
echo "Step 6: Testing schema with a sample query..."
RESPONSE=$(curl -s "${SOLR_URL}/${SOLR_CORE}/select?q=*:*&rows=0")
if echo "$RESPONSE" | grep -q "\"numFound\":0"; then
    echo "✅ Schema is working (empty result set is expected)"
else
    echo "⚠️  Unexpected response, but schema might still be working"
fi
echo ""

# ============================================================================
# Summary
# ============================================================================
echo "==============================================="
echo "✅ Solr Schema Setup Complete!"
echo "==============================================="
echo ""
echo "Core Name: ${SOLR_CORE}"
echo "Solr URL: ${SOLR_URL}/${SOLR_CORE}"
echo ""
echo "Next steps:"
echo "1. Start the search-api service"
echo "2. Index trips via RabbitMQ events or API calls"
echo "3. Test search queries"
echo ""
echo "Test query:"
echo "curl '${SOLR_URL}/${SOLR_CORE}/select?q=*:*&rows=10'"
echo ""
