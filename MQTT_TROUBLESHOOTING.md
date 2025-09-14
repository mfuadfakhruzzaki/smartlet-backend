# MQTT Troubleshooting Guide

## 🚨 Common MQTT Issues & Solutions

### 1. "Invalid bridge configuration" Error

**Problem:** Line 33 configuration error in mosquitto.conf

```
Error: Invalid bridge configuration.
Error found at /mosquitto/config/mosquitto.conf:33.
```

**Solution:**

```bash
# Use fixed mosquitto.conf
./scripts/restart-mqtt.sh  # Linux/Mac
scripts\restart-mqtt.bat   # Windows
```

**Root Cause:**

- `protocol websockets` must be specified AFTER the websocket listener
- `message_size_limit` deprecated, use `max_packet_size`

### 2. Backend Cannot Connect to MQTT

**Problem:**

```
Failed to connect to MQTT broker: network Error : dial tcp: lookup mosquitto
```

**Solutions:**

1. **Check MQTT service status:**

   ```bash
   docker-compose logs mosquitto
   docker-compose ps
   ```

2. **Restart MQTT service:**

   ```bash
   docker-compose restart mosquitto
   ```

3. **Check network connectivity:**

   ```bash
   docker network ls
   docker network inspect backend_swiflet-network
   ```

4. **Test MQTT connection manually:**

   ```bash
   # Subscribe to test topic
   docker exec -it swiflet-mqtt mosquitto_sub -h localhost -t "test/topic"

   # Publish test message (in another terminal)
   docker exec -it swiflet-mqtt mosquitto_pub -h localhost -t "test/topic" -m "Hello"
   ```

### 3. Permission Issues

**Problem:** MQTT authentication failures in production

**Solution:**

```bash
# Create password file
mkdir -p mosquitto_config
docker run --rm eclipse-mosquitto:2.0 mosquitto_passwd -c -b /tmp/passwd username password > mosquitto_config/passwd

# Create ACL file
cat > mosquitto_config/acl << EOF
user username
topic readwrite sensors/+/data
topic readwrite control/+/command
EOF
```

### 4. WebSocket Connection Issues

**Problem:** Web clients cannot connect to MQTT WebSocket

**Solutions:**

1. **Check WebSocket listener:**

   ```bash
   # Should show port 9001 open
   docker-compose ps mosquitto
   ```

2. **Test WebSocket connection:**
   ```bash
   curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
   http://localhost:9001/
   ```

### 5. High CPU/Memory Usage

**Problem:** MQTT broker consuming too many resources

**Solutions:**

1. **Adjust connection limits in mosquitto.conf:**

   ```properties
   max_connections 100          # Reduce from 1000
   max_queued_messages 50       # Reduce from 100
   max_inflight_messages 10     # Reduce from 20
   ```

2. **Enable message expiry:**
   ```properties
   message_expiry_interval 3600  # 1 hour
   ```

## 🛠️ Debugging Commands

### Check Configuration Syntax

```bash
# Test configuration files
./scripts/test-mqtt-config.sh   # Linux/Mac
scripts\test-mqtt-config.bat    # Windows
```

### Monitor MQTT Traffic

```bash
# Real-time monitoring
docker exec -it swiflet-mqtt mosquitto_sub -h localhost -t "#" -v

# Monitor specific topics
docker exec -it swiflet-mqtt mosquitto_sub -h localhost -t "sensors/+/data" -v
```

### Check Logs

```bash
# MQTT broker logs
docker-compose logs -f mosquitto

# Backend logs (MQTT connection)
docker-compose logs -f backend | grep -i mqtt
```

### Network Diagnostics

```bash
# Check if MQTT port is accessible from backend
docker exec -it swiflet-backend nc -zv mosquitto 1883

# Check internal DNS resolution
docker exec -it swiflet-backend nslookup mosquitto
```

## 📋 Configuration Files

### Development: `mosquitto.conf`

- Anonymous connections allowed
- WebSocket enabled on port 9001
- Verbose logging

### Production: `mosquitto.prod.conf`

- Authentication required
- No WebSocket (security)
- Minimal logging
- ACL enforced

### Commands to Switch

```bash
# Use development config
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Use production config
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up
```

## 🚀 Quick Fixes

### Reset MQTT Service

```bash
# Stop all services
docker-compose down

# Remove MQTT volumes (will lose data!)
docker volume rm backend_mosquitto_data backend_mosquitto_logs

# Start fresh
docker-compose up -d
```

### Emergency Disable MQTT

If MQTT is causing issues, backend will continue running without it:

```bash
# Backend will gracefully handle MQTT unavailability
# Check logs for: "Server will continue without MQTT functionality"
docker-compose logs backend
```
