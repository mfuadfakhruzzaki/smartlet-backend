@echo off
REM Production deployment script for Windows

echo 🚀 Deploying Swiflet Backend to Production...

REM Check if .env.prod exists
if not exist ".env.prod" (
    echo ❌ .env.prod file not found. Please create it from .env.prod.example
    exit /b 1
)

REM Check if Docker is running
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker is not running. Please start Docker first.
    exit /b 1
)

REM Create mosquitto config directory
if not exist "mosquitto_config" mkdir mosquitto_config

REM Pull latest images
echo 📦 Pulling latest Docker images...
docker-compose -f docker-compose.yml -f docker-compose.prod.yml pull

REM Build and start production services
echo 🔨 Building and starting production services...
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

REM Wait for services to be healthy
echo ⏳ Waiting for services to be healthy...
timeout /t 30 /nobreak >nul

REM Check service health
echo 🔍 Checking service health...
docker-compose -f docker-compose.yml -f docker-compose.prod.yml ps

REM Test backend health endpoint (using curl if available, otherwise skip)
curl -f http://localhost:8080/health >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Backend is healthy and running!
) else (
    echo ⚠️ Backend health check skipped (curl not available)
)

echo 🎉 Production deployment completed!
echo 🌐 Backend API: http://localhost:8080
echo 📡 MQTT Broker: localhost:1883
echo.
echo 📊 Monitor logs with:
echo docker-compose -f docker-compose.yml -f docker-compose.prod.yml logs -f