#!/bin/bash

# Admin Statistics API - Quick Start Script
# This script helps you get the API running quickly

set -e

echo "🚀 Admin Statistics API - Quick Start"
echo "====================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Check if MongoDB is running
if ! command -v mongosh &> /dev/null && ! command -v mongo &> /dev/null; then
    echo "⚠️  MongoDB CLI not found. Please ensure MongoDB is installed and running."
    echo "   MongoDB installation: https://docs.mongodb.com/manual/installation/"
else
    echo "✅ MongoDB CLI found"
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created. You can edit it to customize your configuration."
else
    echo "✅ .env file already exists"
fi

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod tidy
echo "✅ Dependencies installed"

# Ask user if they want to generate sample data
echo ""
read -p "🎲 Do you want to generate sample data? (2M+ transactions, may take 10-30 minutes) [y/N]: " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🔄 Generating sample data..."
    echo "   This will create 2,000,000+ transactions with 500+ unique users"
    echo "   Please be patient, this may take a while..."
    go run scripts/generate_data.go
    echo "✅ Sample data generated successfully!"
else
    echo "⏭️  Skipping data generation. You can run 'go run scripts/generate_data.go' later."
fi

echo ""
echo "🎉 Setup complete! Here's what you can do next:"
echo ""
echo "1. Start the API server:"
echo "   go run main.go"
echo ""
echo "2. Test the health endpoint:"
echo "   curl http://localhost:8080/health"
echo ""
echo "3. Test an authenticated endpoint:"
echo "   curl -H \"Authorization: admin-secret-token-2024\" \\"
echo "        \"http://localhost:8080/gross_gaming_rev?from=2024-01-01&to=2024-12-31\""
echo ""
echo "📚 For detailed documentation, see README.md"
echo ""
echo "🐳 Recommended: Use Docker for easier setup:"
echo "   docker-compose up -d"
echo ""