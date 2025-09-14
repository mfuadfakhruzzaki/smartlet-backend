@echo off
REM Development setup script for Windows

echo 🚀 Starting Swiflet Backend Development Environment...

REM Check if Docker is running
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker is not running. Please start Docker first.
    exit /b 1
)

REM Create .env file if it doesn't exist
if not exist ".env" (
    echo 📝 Creating .env file from .env.example...
    copy ".env.example" ".env"
)

REM Build and start development services
echo 🔨 Building and starting development services...
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

echo ✅ Development environment is ready!
echo 🌐 Backend API: http://localhost:8080
echo 🗄️  PostgreSQL: localhost:5432
echo 📊 TimescaleDB: localhost:5433
echo 📡 MQTT Broker: localhost:1883
echo 🔴 Redis: localhost:6379