#!/bin/bash

# ==========================================
# Test Railway Endpoint
# ==========================================

# GANTI dengan URL Railway kamu!
API_URL="https://simnikah-api-production-xxxx.up.railway.app"

echo "================================================"
echo "Testing Railway Deployment"
echo "API URL: $API_URL"
echo "================================================"
echo ""

# Test 1: Health Check
echo "1️⃣  Test Health Check..."
echo "GET $API_URL/health"
curl -s "$API_URL/health" | jq '.' 2>/dev/null || curl -s "$API_URL/health"
echo ""
echo ""

# Test 2: Login (example)
echo "2️⃣  Test Login Endpoint..."
echo "POST $API_URL/simnikah/login"
curl -s -X POST "$API_URL/simnikah/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }' | jq '.' 2>/dev/null || curl -s -X POST "$API_URL/simnikah/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
echo ""
echo ""

# Test 3: Search Address (autocomplete)
echo "3️⃣  Test Address Search (Map Integration)..."
echo "GET $API_URL/simnikah/location/search?q=jakarta"
echo "Note: Butuh JWT token, akan dapat error 401 Unauthorized (normal)"
curl -s "$API_URL/simnikah/location/search?q=jakarta" | jq '.' 2>/dev/null || curl -s "$API_URL/simnikah/location/search?q=jakarta"
echo ""
echo ""

echo "================================================"
echo "✅ Test selesai!"
echo "================================================"
echo ""
echo "Kalau health check berhasil, berarti deployment SUKSES!"
echo ""
echo "CATATAN:"
echo "- Endpoint yang butuh auth akan return 401 (normal)"
echo "- Simpan URL ini untuk frontend kamu"
echo ""

