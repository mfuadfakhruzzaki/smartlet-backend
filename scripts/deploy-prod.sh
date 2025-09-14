#!/bin/bash

# Production deployment script

set -e

echo "🚀 Deploying Swiflet Backend to Production..."

# Check if .env.prod exists
if [ ! -f .env.prod ]; then
    echo "❌ .env.prod file not found. Please create it from .env.prod.example"
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Load production environment variables
set -a
source .env.prod
set +a

# Create MQTT password file for production
echo "📝 Creating MQTT authentication files..."
mkdir -p ./mosquitto_config
docker run --rm -it eclipse-mosquitto:2.0 mosquitto_passwd -c -b /tmp/passwd ${MQTT_USERNAME} ${MQTT_PASSWORD} > ./mosquitto_config/passwd

# Create MQTT ACL file
cat > ./mosquitto_config/acl << EOF
# MQTT Access Control List for Production

# Default deny all
topic read #
topic write #

# Allow swiflet backend client
user ${MQTT_USERNAME}
topic readwrite sensors/+/data
topic readwrite control/+/command

# Allow specific ESP32 devices (add as needed)
pattern read sensors/%c/data
pattern write control/%c/command
EOF

# Pull latest images
echo "📦 Pulling latest Docker images..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml pull

# Build and start production services
echo "🔨 Building and starting production services..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
sleep 30

# Check service health
echo "🔍 Checking service health..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml ps

# Test backend health endpoint
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Backend is healthy and running!"
else
    echo "❌ Backend health check failed"
    exit 1
fi

echo "🎉 Production deployment completed successfully!"
echo "🌐 Backend API: http://localhost:8080"
echo "📡 MQTT Broker: localhost:1883"
echo ""
echo "📊 Monitor logs with:"
echo "docker-compose -f docker-compose.yml -f docker-compose.prod.yml logs -f"