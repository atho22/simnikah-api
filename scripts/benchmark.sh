#!/bin/bash

# =====================================================
# SimNikah API - Benchmark Script
# =====================================================

echo "🔥 SimNikah API Benchmark"
echo "========================="
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "⚠️  Server not running. Starting server..."
    cd ..
    make dev &
    SERVER_PID=$!
    sleep 3
    echo "Server started (PID: $SERVER_PID)"
else
    echo "✅ Server already running"
fi

echo ""
echo "📊 Running benchmarks..."
echo ""

# Benchmark 1: Health endpoint
echo "1️⃣  Health Endpoint Benchmark"
echo "   Target: < 5ms average"
ab -n 1000 -c 50 -q http://localhost:8080/health 2>/dev/null | grep -E "Requests per second|Time per request|50%|90%|99%"
echo ""

# Benchmark 2: Sustained load test
echo "2️⃣  Sustained Load Test (5000 requests, 100 concurrent)"
echo "   Target: 0 failed requests"
ab -n 5000 -c 100 -t 30 -q http://localhost:8080/health 2>/dev/null | grep -E "Complete requests|Failed requests|Requests per second|Time per request"
echo ""

# Benchmark 3: Rate limiting test
echo "3️⃣  Rate Limiting Test (should limit at 100 req/min)"
echo "   Sending 150 requests in quick succession..."
FAILED=0
for i in {1..150}; do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
    if [ "$STATUS" == "429" ]; then
        FAILED=$((FAILED + 1))
    fi
done
echo "   Requests blocked by rate limiter: $FAILED / 150"
if [ $FAILED -gt 0 ]; then
    echo "   ✅ Rate limiting working!"
else
    echo "   ⚠️  Rate limiting might not be active"
fi
echo ""

# Benchmark 4: Memory leak test
echo "4️⃣  Memory Leak Test (1 minute sustained load)"
PID=$(pgrep -f "simnikah-api" | head -1)
if [ ! -z "$PID" ]; then
    echo "   Initial memory:"
    ps aux | grep "$PID" | grep -v grep | awk '{print "   " $6/1024 " MB"}'
    
    ab -n 10000 -c 50 -t 60 -q http://localhost:8080/health > /dev/null 2>&1 &
    BENCH_PID=$!
    
    sleep 30
    echo "   Memory after 30s:"
    ps aux | grep "$PID" | grep -v grep | awk '{print "   " $6/1024 " MB"}'
    
    wait $BENCH_PID
    echo "   Final memory:"
    ps aux | grep "$PID" | grep -v grep | awk '{print "   " $6/1024 " MB"}'
    echo "   ✅ Check if memory is stable (no significant increase)"
else
    echo "   ⚠️  Process not found"
fi
echo ""

echo "=================================="
echo "✅ Benchmark completed!"
echo ""
echo "📈 Expected Results:"
echo "  Response time: < 10ms"
echo "  Requests/sec: > 2000"
echo "  Failed requests: 0"
echo "  Memory: Stable (no leaks)"
echo ""

