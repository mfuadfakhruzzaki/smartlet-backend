@echo off
REM Restart MQTT Service Script for Windows

echo 🔄 Restarting MQTT service...

REM Check if docker-compose is running
docker-compose ps | findstr "mosquitto" >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️ MQTT service is not running. Starting services...
    docker-compose up -d mosquitto
) else (
    echo 🔄 Restarting MQTT service...
    docker-compose restart mosquitto
)

REM Wait a moment for service to start
echo ⏳ Waiting for service to start...
timeout /t 5 /nobreak >nul

REM Check MQTT service status
echo 📊 Checking MQTT service status...
docker-compose logs --tail=10 mosquitto

echo ✅ MQTT service restart completed!
echo 📡 MQTT is available at: localhost:1883
echo 🌐 WebSocket MQTT at: localhost:9001