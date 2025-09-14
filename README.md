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

# S3 Storage (Required for file upload)
S3_ACCESS_KEY=your-s3-access-key
S3_SECRET_KEY=your-s3-secret-key
S3_BUCKET=your-bucket-name
S3_REGION=us-east-1
```

### Quick Development Setup

For easy development setup with MinIO (S3-compatible storage):

**Windows:**

```cmd
.\setup-dev.bat
```

**Linux/Mac:**

```bash
chmod +x setup-dev.sh
./setup-dev.sh
```

This script will:

- Copy `.env.example` to `.env`
- Start MinIO container for S3 storage
- Configure environment for local development
- Setup ready for file upload testing

After running the script:

1. Open MinIO Console: http://localhost:9001
2. Login with `minioadmin` / `minioadmin`
3. Create bucket `swiflet-storage`
4. Set bucket policy to public read access

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

## 🔧 Troubleshooting

### S3 Upload Error: "EmptyStaticCreds"

If you get error: `"Failed to upload file: EmptyStaticCreds: static credentials are empty"` when testing upload endpoints:

1. **Quick Fix - Use setup script:**

   ```bash
   # Windows
   .\setup-dev.bat

   # Linux/Mac
   ./setup-dev.sh
   ```

2. **Manual Fix:**

   ```bash
   # Configure S3 credentials in .env
   S3_ACCESS_KEY=minioadmin
   S3_SECRET_KEY=minioadmin
   S3_ENDPOINT=http://localhost:9000

   # Restart server
   go run cmd/server/main.go
   ```

3. **Create MinIO bucket:**
   - Open http://localhost:9001
   - Login: minioadmin/minioadmin
   - Create bucket: `swiflet-storage`
   - Set policy: Public read access

📖 **Detailed guide**: See `S3_UPLOAD_CONFIGURATION.md`

### Common Issues

- **Database connection failed**: Check PostgreSQL is running and credentials in `.env`
- **MQTT broker unavailable**: Check Mosquitto is running on port 1883
- **JWT token invalid**: Check `JWT_SECRET` in `.env` matches between requests
- **Permission denied**: Check file upload permissions and S3 bucket policy

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
