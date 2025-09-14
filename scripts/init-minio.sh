#!/bin/bash

# MinIO Bucket Initialization Script
# This script creates the required bucket in MinIO after it starts

echo "Waiting for MinIO to be ready..."
sleep 10

# Install MinIO client if not exists
if ! command -v mc &> /dev/null; then
    echo "Installing MinIO client..."
    curl https://dl.min.io/client/mc/release/linux-amd64/mc \
      --create-dirs \
      -o /usr/local/bin/mc
    chmod +x /usr/local/bin/mc
fi

# Configure MinIO client
echo "Configuring MinIO client..."
mc alias set local http://minio:9000 ${MINIO_ROOT_USER:-minioadmin} ${MINIO_ROOT_PASSWORD:-minioadmin123}

# Create bucket if it doesn't exist
BUCKET_NAME=${S3_BUCKET:-swiftlead-storage}
echo "Creating bucket: $BUCKET_NAME"

if ! mc ls local/$BUCKET_NAME > /dev/null 2>&1; then
    mc mb local/$BUCKET_NAME
    echo "Bucket $BUCKET_NAME created successfully"
    
    # Set bucket policy to allow public read for certain folders (optional)
    # Uncomment the lines below if you want public access to uploaded files
    # mc anonymous set download local/$BUCKET_NAME/articles
    # mc anonymous set download local/$BUCKET_NAME/ebooks
    # mc anonymous set download local/$BUCKET_NAME/users
else
    echo "Bucket $BUCKET_NAME already exists"
fi

echo "MinIO initialization completed"