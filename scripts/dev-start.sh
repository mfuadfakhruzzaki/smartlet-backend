#!/bin/bash

# Development setup script

echo "🚀 Starting Swiflet Backend Development Environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
fi

# Build and start development services
echo "🔨 Building and starting development services..."
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

echo "✅ Development environment is ready!"
echo "🌐 Backend API: http://localhost:8080"
echo "🗄️  PostgreSQL: localhost:5432"
echo "📊 TimescaleDB: localhost:5433"
echo "📡 MQTT Broker: localhost:1883"
echo "🔴 Redis: localhost:6379"