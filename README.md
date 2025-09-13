# Swiflet Platform Backend

Backend API untuk platform manajemen rumah walet dengan integrasi IoT, built dengan Go dan Gin framework.

## 🏗️ Arsitektur

Aplikasi ini mengimplementasikan arsitektur microservices dengan komponen:

- **REST API Server**: Gin-based HTTP server
- **PostgreSQL**: Database utama untuk data bisnis
- **TimescaleDB**: Time-series database untuk data sensor IoT
- **MQTT Broker**: Message broker untuk komunikasi IoT
- **JWT Authentication**: Token-based authentication

## 📋 Fitur

### Core Features

- ✅ User Authentication & Authorization (JWT)
- ✅ User Management CRUD
- 🔄 Article & Content Management
- 🔄 IoT Device Management
- 🔄 Sensor Data Collection (MQTT)
- 🔄 Harvest Management
- 🔄 Market & Pricing
- 🔄 Service Requests
- 🔄 Transaction & Membership

### IoT Integration

- MQTT sensor data ingestion
- Real-time temperature & humidity monitoring
- Device control commands
- TimescaleDB for time-series data

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- TimescaleDB extension
- MQTT Broker (Mosquitto)

### Environment Setup

1. Copy environment template:

```bash
cp .env.example .env
```

2. Update `.env` with your configuration:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=swiflet_db

# TimescaleDB
TIMESCALE_HOST=localhost
TIMESCALE_PORT=5432
TIMESCALE_USER=postgres
TIMESCALE_PASSWORD=your_password
TIMESCALE_DB=swiflet_timeseries

# JWT
JWT_SECRET=your-super-secret-key
JWT_EXPIRY=24h

# MQTT
MQTT_BROKER=tcp://localhost:1883
```

### Database Setup

1. Create databases:

```sql
CREATE DATABASE swiflet_db;
CREATE DATABASE swiflet_timeseries;
\c swiflet_timeseries;
CREATE EXTENSION IF NOT EXISTS timescaledb;
```

2. Run migrations:

```bash
# PostgreSQL tables
psql -h localhost -U postgres -d swiflet_db -f migrations/001_create_tables.sql

# TimescaleDB tables
psql -h localhost -U postgres -d swiflet_timeseries -f migrations/002_timescale_tables.sql
```

### Installation & Running

1. Install dependencies:

```bash
make deps
```

2. Run the application:

```bash
make run
```

Or for development with hot reload:

```bash
make dev
```

Server will start on `http://localhost:8080`

## 📖 API Documentation

### Authentication Endpoints

#### Register User

```http
POST /v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "phone": "+628123456789",
  "address": "Jakarta, Indonesia"
}
```

#### Login

```http
POST /v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

### Protected Endpoints

All protected endpoints require `Authorization: Bearer <token>` header.

#### Users

- `GET /v1/users` - List users (paginated)
- `GET /v1/users/{id}` - Get user by ID
- `PATCH /v1/users/{id}` - Update user
- `DELETE /v1/users/{id}` - Delete user

#### IoT Devices

- `GET /v1/iot-devices` - List IoT devices
- `POST /v1/iot-devices` - Create IoT device
- `GET /v1/sensors` - Get sensor data

### Health Check

```http
GET /health
```

## 🏗️ Project Structure

```
.
├── cmd/
│   └── server/          # Application entrypoint
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connections
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   └── services/        # Business logic services
├── pkg/
│   └── utils/           # Utility functions
├── migrations/          # Database migration files
├── docs/               # API documentation
├── .env.example        # Environment template
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── Makefile            # Build automation
└── README.md           # This file
```

## 🔧 Development

### Available Make Commands

```bash
make build         # Build the application
make run           # Run the application
make dev           # Run with hot reload
make test          # Run tests
make fmt           # Format code
make lint          # Lint code
make clean         # Clean build artifacts
make deps          # Download dependencies
```

### Code Style

- Follow Go standards and conventions
- Use `gofmt` for formatting
- Write tests for business logic
- Add comments for exported functions

## 🐳 Docker

Build and run with Docker:

```bash
make docker-build
make docker-run
```

## 📊 Monitoring

### Health Check

The application exposes a health check endpoint at `/health`

### MQTT Monitoring

Monitor MQTT traffic for IoT sensor data:

```bash
mosquitto_sub -t "sensors/+/data"
```

## 🔧 Configuration

Environment variables:

| Variable      | Description        | Default              |
| ------------- | ------------------ | -------------------- |
| `DB_HOST`     | PostgreSQL host    | localhost            |
| `DB_PORT`     | PostgreSQL port    | 5432                 |
| `JWT_SECRET`  | JWT signing secret | (required)           |
| `JWT_EXPIRY`  | JWT token expiry   | 24h                  |
| `MQTT_BROKER` | MQTT broker URL    | tcp://localhost:1883 |
| `SERVER_PORT` | HTTP server port   | 8080                 |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make fmt` and `make lint`
6. Submit a pull request

## 📝 License

[Add your license here]

## 📞 Support

For support and questions, please contact the development team.
