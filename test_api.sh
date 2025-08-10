#!/bin/bash

echo "========================================"
echo "Testing Admin Statistics API"
echo "========================================"

echo ""
echo "1. Testing Health Check..."
curl http://localhost:8090/health
echo ""

echo "2. Testing Gross Gaming Revenue..."
curl -H "Authorization: admin-secret-token-2024" "http://localhost:8090/gross_gaming_rev?from=2024-01-01&to=2024-12-31"
echo ""

echo "3. Testing Daily Wager Volume..."
curl -H "Authorization: admin-secret-token-2024" "http://localhost:8090/daily_wager_volume?from=2024-01-01&to=2024-12-31"
echo ""

echo "4. Testing User Wager Percentile..."
curl -H "Authorization: admin-secret-token-2024" "http://localhost:8090/user/68982b9fb3890672dd066862/wager_percentile?from=2024-01-01&to=2024-12-31"
echo ""

echo "5. Testing Error Scenarios..."
echo "Testing without auth token (should fail):"
curl "http://localhost:8090/gross_gaming_rev?from=2024-01-01&to=2024-12-31"
echo ""

echo "Testing with invalid date:"
curl -H "Authorization: admin-secret-token-2024" "http://localhost:8090/gross_gaming_rev?from=invalid-date&to=2024-12-31"
echo ""

echo "6. Running Unit Tests..."
docker-compose exec app go test ./...

echo ""
echo "========================================"
echo "Testing Complete!"
echo "========================================"