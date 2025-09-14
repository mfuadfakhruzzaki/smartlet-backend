#!/bin/bash

# Test MQTT Configuration Script

echo "🧪 Testing MQTT Configuration..."

# Test mosquitto.conf syntax
echo "📝 Testing mosquitto.conf syntax..."
docker run --rm -v "$(pwd)/mosquitto.conf:/mosquitto/config/mosquitto.conf" eclipse-mosquitto:2.0 mosquitto -c /mosquitto/config/mosquitto.conf -t

if [ $? -eq 0 ]; then
    echo "✅ mosquitto.conf syntax is valid"
else
    echo "❌ mosquitto.conf has syntax errors"
    exit 1
fi

# Test mosquitto.dev.conf syntax
echo "📝 Testing mosquitto.dev.conf syntax..."
docker run --rm -v "$(pwd)/mosquitto.dev.conf:/mosquitto/config/mosquitto.conf" eclipse-mosquitto:2.0 mosquitto -c /mosquitto/config/mosquitto.conf -t

if [ $? -eq 0 ]; then
    echo "✅ mosquitto.dev.conf syntax is valid"
else
    echo "❌ mosquitto.dev.conf has syntax errors"
    exit 1
fi

# Test mosquitto.prod.conf syntax (if password file exists)
if [ -f "./mosquitto_config/passwd" ]; then
    echo "📝 Testing mosquitto.prod.conf syntax..."
    docker run --rm \
        -v "$(pwd)/mosquitto.prod.conf:/mosquitto/config/mosquitto.conf" \
        -v "$(pwd)/mosquitto_config/passwd:/mosquitto/config/passwd" \
        -v "$(pwd)/mosquitto_config/acl:/mosquitto/config/acl" \
        eclipse-mosquitto:2.0 mosquitto -c /mosquitto/config/mosquitto.conf -t

    if [ $? -eq 0 ]; then
        echo "✅ mosquitto.prod.conf syntax is valid"
    else
        echo "❌ mosquitto.prod.conf has syntax errors"
        exit 1
    fi
else
    echo "⚠️ Skipping mosquitto.prod.conf test (no password file found)"
fi

echo "🎉 All MQTT configurations are valid!"