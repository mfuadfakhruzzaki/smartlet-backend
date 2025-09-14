# MinIO Migration Guide

## Overview

Aplikasi ini telah dimigrasi dari AWS S3 ke MinIO untuk object storage lokal. MinIO adalah high-performance object storage yang kompatibel dengan S3 API.

## Perubahan yang Dilakukan

### 1. Service Layer (`internal/services/s3.go`)

- Menambahkan auto-detection SSL based on endpoint
- Menambahkan fungsi `ensureBucketExists()` untuk membuat bucket otomatis
- Tetap menggunakan AWS SDK v1 untuk kompatibilitas (bisa di-upgrade ke v2 nanti)

### 2. Docker Compose (`docker-compose.yml`)

- Menambahkan service MinIO
- Port 9000: MinIO API endpoint
- Port 9001: MinIO Web Console
- Auto-configuration untuk S3 environment variables

### 3. Konfigurasi Environment

- `S3_ENDPOINT`: http://minio:9000 (dalam Docker)
- `S3_ACCESS_KEY`: minioadmin (default)
- `S3_SECRET_KEY`: minioadmin123 (default)
- `S3_BUCKET`: swiftlead-storage
- `S3_REGION`: us-east-1

## Cara Menjalankan

### 1. Development (Docker Compose)

```bash
# Copy environment variables
cp .env.example .env

# Edit .env sesuai kebutuhan
# Terutama MINIO_ROOT_USER dan MINIO_ROOT_PASSWORD untuk production

# Start semua services
docker-compose up -d

# Cek MinIO Web Console
# http://localhost:9001
# Login: minioadmin / minioadmin123 (atau sesuai .env)
```

### 2. Production

```bash
# Gunakan credentials yang aman
export MINIO_ROOT_USER=your_secure_username
export MINIO_ROOT_PASSWORD=your_secure_password

# Start dengan production config
docker-compose up -d
```

## Fitur-fitur yang Didukung

### Upload Files

- ✅ User profile images: `/upload/user/profile`
- ✅ Article cover images: `/upload/article/:id/cover`
- ✅ E-book files: `/upload/ebook`
- ✅ E-book thumbnails: `/upload/ebook/:id/thumbnail`
- ✅ Harvest proof photos: `/upload/harvest/proof`

### File Management

- ✅ Auto bucket creation
- ✅ File type validation
- ✅ File size limits
- ✅ Unique filename generation
- ✅ File deletion
- ✅ Presigned URLs

## MinIO Web Console

Akses MinIO Console di: http://localhost:9001

Features:

- Browse buckets dan objects
- Upload/download files manual
- Manage bucket policies
- Monitor performance
- User management

## Troubleshooting

### 1. Bucket tidak dibuat otomatis

```bash
# Masuk ke MinIO container
docker exec -it swiflet-minio sh

# Install MC client
curl https://dl.min.io/client/mc/release/linux-amd64/mc -o /usr/local/bin/mc
chmod +x /usr/local/bin/mc

# Configure dan buat bucket
mc alias set local http://localhost:9000 minioadmin minioadmin123
mc mb local/swiftlead-storage
```

### 2. Connection refused

- Pastikan MinIO service running: `docker ps`
- Check logs: `docker logs swiflet-minio`
- Pastikan port 9000 dan 9001 tidak digunakan aplikasi lain

### 3. Access denied

- Periksa MINIO_ROOT_USER dan MINIO_ROOT_PASSWORD
- Pastikan S3_ACCESS_KEY dan S3_SECRET_KEY sesuai dengan MinIO credentials

## Migration dari AWS S3 (Opsional)

Jika Anda ingin memigrate data dari AWS S3 ke MinIO:

```bash
# Configure AWS profile
mc alias set aws https://s3.amazonaws.com AWS_ACCESS_KEY AWS_SECRET_KEY

# Configure MinIO
mc alias set minio http://localhost:9000 minioadmin minioadmin123

# Mirror bucket
mc mirror aws/your-aws-bucket minio/swiftlead-storage
```

## Keamanan Production

1. **Ganti default credentials:**

```env
MINIO_ROOT_USER=your_admin_user
MINIO_ROOT_PASSWORD=very_secure_password_here
```

2. **Use HTTPS in production:**

```env
S3_ENDPOINT=https://your-minio-domain.com
```

3. **Setup reverse proxy** dengan SSL certificate

4. **Configure bucket policies** sesuai kebutuhan access control

## Monitoring

MinIO menyediakan metrics yang bisa diintegrasikan dengan:

- Prometheus
- Grafana
- CloudWatch (jika perlu)

Endpoint metrics: http://localhost:9000/minio/v2/metrics/cluster
