#!/bin/bash

echo "================================"
echo "Smartlet Backend Setup Script"
echo "================================"
echo

echo "[1/4] Copying environment file..."
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✓ Environment file created"
else
    echo "! Environment file already exists"
fi
echo

echo "[2/4] Checking Docker..."
if ! command -v docker &> /dev/null; then
    echo "✗ Docker not found! Please install Docker first."
    echo "  Download: https://www.docker.com/products/docker-desktop"
    exit 1
fi
echo "✓ Docker is available"
echo

echo "[3/4] Starting MinIO for S3 storage..."
echo "Starting MinIO container..."
docker run -d \
  --name smartlet-minio \
  -p 9000:9000 \
  -p 9001:9001 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  -v "$(pwd)/minio-data:/data" \
  minio/minio server /data --console-address ":9001"

if [ $? -eq 0 ]; then
    echo "✓ MinIO started successfully"
else
    echo "! MinIO container might already exist, trying to start..."
    docker start smartlet-minio
fi
echo

echo "[4/4] Configuring environment for MinIO..."
echo "Setting up .env with MinIO configuration..."

# Update .env file with MinIO settings
sed -i 's/^S3_ACCESS_KEY=.*/S3_ACCESS_KEY=minioadmin/' .env
sed -i 's/^S3_SECRET_KEY=.*/S3_SECRET_KEY=minioadmin/' .env
sed -i 's/^S3_ENDPOINT=.*/S3_ENDPOINT=http:\/\/localhost:9000/' .env

echo "✓ Environment configured"
echo

echo "================================"
echo "Setup Complete!"
echo "================================"
echo
echo "Next steps:"
echo "1. Open MinIO Console: http://localhost:9001"
echo "   Login: minioadmin / minioadmin"
echo "2. Create bucket: 'swiflet-storage'"
echo "3. Set bucket policy to 'public' for read access"
echo "4. Start the server: go run cmd/server/main.go"
echo "5. Test upload endpoints with Postman"
echo
echo "MinIO Console: http://localhost:9001"
echo "MinIO API: http://localhost:9000"
echo