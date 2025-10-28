#!/bin/bash

# =====================================================
# SimNikah API - Performance Testing Script
# =====================================================

echo "🚀 SimNikah API Performance Test"
echo "=================================="
echo ""

# Check if server is running
SERVER_URL="http://localhost:8080"
if ! curl -s "$SERVER_URL/health" > /dev/null 2>&1; then
    echo "❌ Server tidak running di $SERVER_URL"
    echo "Start server dulu dengan: make dev"
    exit 1
fi

echo "✅ Server detected at $SERVER_URL"
echo ""

# Test 1: Health Check Response Time
echo "📊 Test 1: Health Check Response Time"
echo "------------------------------------"
for i in {1..5}; do
    time curl -s "$SERVER_URL/health" > /dev/null
done
echo ""

# Test 2: Concurrent Requests Test
echo "📊 Test 2: Concurrent Requests (100 requests, 10 concurrent)"
echo "------------------------------------------------------------"
if command -v ab > /dev/null 2>&1; then
    ab -n 100 -c 10 -q "$SERVER_URL/health" | grep -E "Requests per second|Time per request|Failed requests"
else
    echo "⚠️  Apache Bench (ab) not installed"
    echo "Install: sudo apt install apache2-utils"
    echo ""
    echo "Testing with curl loop instead..."
    start_time=$(date +%s.%N)
    for i in {1..10}; do
        curl -s "$SERVER_URL/health" > /dev/null &
    done
    wait
    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)
    echo "10 concurrent requests completed in: ${duration}s"
fi
echo ""

# Test 3: Database Query Performance
echo "📊 Test 3: Database-Heavy Endpoint (if available)"
echo "------------------------------------------------"
# This requires authentication, so we'll skip for now
echo "⚠️  Requires authentication - skipping"
echo "To test: curl -H 'Authorization: Bearer <token>' $SERVER_URL/simnikah/pendaftaran"
echo ""

# Test 4: Memory Usage
echo "📊 Test 4: Memory Usage"
echo "----------------------"
if command -v ps > /dev/null 2>&1; then
    PID=$(pgrep -f "simnikah-api" | head -1)
    if [ ! -z "$PID" ]; then
        echo "Process ID: $PID"
        ps aux | grep "$PID" | grep -v grep | awk '{print "Memory: " $6/1024 " MB"}'
        ps aux | grep "$PID" | grep -v grep | awk '{print "CPU: " $3 "%"}'
    else
        echo "⚠️  Process not found"
    fi
else
    echo "⚠️  ps command not available"
fi
echo ""

# Test 5: Geocoding Cache Test
echo "📊 Test 5: Geocoding Cache Performance"
echo "-------------------------------------"
echo "First request (cache MISS - hits API):"
time curl -s -X POST "$SERVER_URL/simnikah/location/geocode" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer test" \
    -d '{"alamat":"Jl. Merdeka No. 123, Jakarta"}' 2>&1 | head -1

echo ""
echo "Second request (cache HIT - from memory):"
time curl -s -X POST "$SERVER_URL/simnikah/location/geocode" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer test" \
    -d '{"alamat":"Jl. Merdeka No. 123, Jakarta"}' 2>&1 | head -1
echo ""

echo "=================================="
echo "✅ Performance tests completed!"
echo ""
echo "📝 Summary:"
echo "  - Check response times above"
echo "  - Lower is better!"
echo "  - Cache should show significant improvement on 2nd request"
echo ""

