#!/bin/bash

# Test script for MQTT sensor data simulation

BROKER="mosquitto"
PORT="1883"

echo "Starting MQTT sensor data simulation..."

# Simulate sensor data from different devices
DEVICES=("DEVICE001" "DEVICE002" "DEVICE003")

while true; do
    for device in "${DEVICES[@]}"; do
        # Generate random sensor data
        TEMP=$(awk 'BEGIN{printf "%.2f", 25 + rand() * 10}')  # 25-35°C
        HUMID=$(awk 'BEGIN{printf "%.2f", 60 + rand() * 20}') # 60-80%
        TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
        
        # Create JSON payload
        PAYLOAD=$(cat <<EOF
{
    "install_code": "$device",
    "suhu": $TEMP,
    "kelembaban": $HUMID,
    "timestamp": "$TIMESTAMP"
}
EOF
)
        
        # Publish to MQTT topic
        echo "Publishing data for $device: Temp=$TEMP°C, Humidity=$HUMID%"
        mosquitto_pub -h $BROKER -p $PORT -t "sensors/$device/data" -m "$PAYLOAD"
        
        # Wait 2 seconds between devices
        sleep 2
    done
    
    # Wait 10 seconds before next round
    echo "Waiting 10 seconds before next round..."
    sleep 10
done