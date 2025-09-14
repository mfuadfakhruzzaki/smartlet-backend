@echo off
echo ================================
echo Smartlet Backend Setup Script
echo ================================
echo.

echo [1/4] Copying environment file...
if not exist .env (
    copy .env.example .env
    echo ✓ Environment file created
) else (
    echo ! Environment file already exists
)
echo.

echo [2/4] Checking Docker...
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ✗ Docker not found! Please install Docker first.
    echo   Download: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)
echo ✓ Docker is available
echo.

echo [3/4] Starting MinIO for S3 storage...
echo Starting MinIO container...
docker run -d ^
  --name smartlet-minio ^
  -p 9000:9000 ^
  -p 9001:9001 ^
  -e "MINIO_ROOT_USER=minioadmin" ^
  -e "MINIO_ROOT_PASSWORD=minioadmin" ^
  -v "%cd%\minio-data:/data" ^
  minio/minio server /data --console-address ":9001"

if %errorlevel% equ 0 (
    echo ✓ MinIO started successfully
) else (
    echo ! MinIO container might already exist, trying to start...
    docker start smartlet-minio
)
echo.

echo [4/4] Configuring environment for MinIO...
echo Setting up .env with MinIO configuration...

REM Update .env file with MinIO settings
powershell -Command "(Get-Content .env) -replace '^S3_ACCESS_KEY=.*', 'S3_ACCESS_KEY=minioadmin' | Set-Content .env.tmp"
powershell -Command "(Get-Content .env.tmp) -replace '^S3_SECRET_KEY=.*', 'S3_SECRET_KEY=minioadmin' | Set-Content .env.tmp2"
powershell -Command "(Get-Content .env.tmp2) -replace '^S3_ENDPOINT=.*', 'S3_ENDPOINT=http://localhost:9000' | Set-Content .env"
del .env.tmp .env.tmp2

echo ✓ Environment configured
echo.

echo ================================
echo Setup Complete!
echo ================================
echo.
echo Next steps:
echo 1. Open MinIO Console: http://localhost:9001
echo    Login: minioadmin / minioadmin
echo 2. Create bucket: 'swiflet-storage'
echo 3. Set bucket policy to 'public' for read access
echo 4. Start the server: go run cmd/server/main.go
echo 5. Test upload endpoints with Postman
echo.
echo MinIO Console: http://localhost:9001
echo MinIO API: http://localhost:9000
echo.
pause