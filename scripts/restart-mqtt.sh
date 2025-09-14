#!/bin/bash

# Restart MQTT Service Script

echo "🔄 Restarting MQTT service..."

# Check if docker-compose is running
if ! docker-compose ps | grep -q "mosquitto"; then
    echo "⚠️ MQTT service is not running. Starting services..."
    docker-compose up -d mosquitto
else
    echo "🔄 Restarting MQTT service..."
    docker-compose restart mosquitto
fi

# Wait a moment for service to start
sleep 5

# Check MQTT service status
echo "📊 Checking MQTT service status..."
docker-compose logs --tail=10 mosquitto

echo "✅ MQTT service restart completed!"
echo "📡 MQTT is available at: localhost:1883"
echo "🌐 WebSocket MQTT at: localhost:9001"