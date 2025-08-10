#!/bin/bash
echo "========================================"
echo "Stopping Admin Statistics API"
echo "========================================"

echo ""
echo "Stopping Docker Compose services..."
docker-compose down

echo ""
echo "========================================"
echo "Services Stopped!"
echo "========================================"
read -p "Press any key to continue..."