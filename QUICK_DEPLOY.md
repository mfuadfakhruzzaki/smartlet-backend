# Quick Deployment Guide

## 🚀 Single Command Deployment

### For Dokploy/Production Server:

1. **Clone repo:**

   ```bash
   git clone https://github.com/mfuadfakhruzzaki/smartlet-backend.git
   cd smartlet-backend
   ```

2. **Deploy immediately:**
   ```bash
   docker-compose up -d --build
   ```

That's it! The app will be available on port 8080.

## 🔧 Environment Variables (Optional)

Create `.env` file if you need custom settings:

```bash
# Basic required settings
DB_PASSWORD=your-secure-db-password
JWT_SECRET=your-32-char-secret-key-here
```

## 📡 Default Services

- **Backend API:** http://localhost:8080
- **MQTT Broker:** localhost:1883
- **MQTT WebSocket:** localhost:9001

## 🔍 Health Check

```bash
curl http://localhost:8080/health
```

## 📊 Monitor Logs

```bash
docker-compose logs -f backend
docker-compose logs -f mosquitto
```

## 🛠️ Troubleshooting

If MQTT has issues:

```bash
docker-compose restart mosquitto
```

If database needs reset:

```bash
docker-compose down -v
docker-compose up -d --build
```

## 🔄 Updates

```bash
git pull
docker-compose up -d --build
```
