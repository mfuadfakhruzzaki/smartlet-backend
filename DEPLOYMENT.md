# Swiflet Backend Deployment Guide

## 🚀 Production Deployment

### Prerequisites

- Docker & Docker Compose installed
- `.env.prod` file created from `.env.prod.example`

### Quick Deployment

#### Linux/Mac:

```bash
chmod +x scripts/deploy-prod.sh
./scripts/deploy-prod.sh
```

#### Windows:

```cmd
scripts\deploy-prod.bat
```

### Manual Deployment

1. **Create production environment file:**

   ```bash
   cp .env.prod.example .env.prod
   # Edit .env.prod with your production values
   ```

2. **Create MQTT authentication (if using Docker MQTT):**

   ```bash
   mkdir -p mosquitto_config
   docker run --rm eclipse-mosquitto:2.0 mosquitto_passwd -c -b /tmp/passwd your_mqtt_user your_mqtt_password > mosquitto_config/passwd
   ```

3. **Deploy services:**
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
   ```

### 🛠️ Development Setup

#### Quick Start:

```bash
# Linux/Mac
chmod +x scripts/dev-start.sh
./scripts/dev-start.sh

# Windows
scripts\dev-start.bat
```

#### Manual Development:

```bash
cp .env.example .env
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

## 📊 Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f mosquitto
```

### Service Status

```bash
docker-compose ps
```

## 🔧 Troubleshooting

### MQTT Connection Issues

1. **Check MQTT broker status:**

   ```bash
   docker-compose logs mosquitto
   ```

2. **Test MQTT connection:**

   ```bash
   # Subscribe to test topic
   docker exec -it swiflet-mqtt mosquitto_sub -h localhost -t "test/topic"

   # Publish test message
   docker exec -it swiflet-mqtt mosquitto_pub -h localhost -t "test/topic" -m "Hello MQTT"
   ```

3. **Common fixes:**
   - Ensure `mosquitto.conf` has correct syntax
   - Check network connectivity between containers
   - Verify MQTT credentials in `.env` file

### Database Connection Issues

1. **Check database status:**

   ```bash
   docker-compose logs postgres
   docker-compose logs timescaledb
   ```

2. **Connect to database:**

   ```bash
   # PostgreSQL
   docker exec -it swiflet-postgres psql -U postgres -d swiflet_db

   # TimescaleDB
   docker exec -it swiflet-timescaledb psql -U postgres -d swiflet_timeseries
   ```

### Backend Issues

1. **Check backend logs:**

   ```bash
   docker-compose logs backend
   ```

2. **Restart backend only:**
   ```bash
   docker-compose restart backend
   ```

## 🔐 Production Security Checklist

- [ ] Change all default passwords in `.env.prod`
- [ ] Use strong JWT secret (>32 characters)
- [ ] Enable SSL/TLS for database connections (`DB_SSLMODE=require`)
- [ ] Configure MQTT authentication (`mosquitto_config/passwd`)
- [ ] Set up proper firewall rules
- [ ] Use reverse proxy (nginx) for HTTPS
- [ ] Regular backup of database volumes
- [ ] Monitor resource usage and logs

## 🌐 Environment URLs

### Development:

- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432
- TimescaleDB: localhost:5433
- MQTT Broker: localhost:1883
- Redis: localhost:6379

### Production:

- Backend API: http://localhost:8080
- MQTT Broker: localhost:1883
- (Database ports not exposed for security)

## 📝 Environment Variables

See `.env.example` for development and `.env.prod.example` for production configuration options.
