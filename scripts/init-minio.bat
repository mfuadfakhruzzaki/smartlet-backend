@echo off
REM MinIO Bucket Initialization Script for Windows
REM This script creates the required bucket in MinIO after it starts

echo Waiting for MinIO to be ready...
timeout /t 10 > nul

echo Setting up MinIO client...

REM Download MinIO client if not exists
if not exist "mc.exe" (
    echo Downloading MinIO client...
    curl -o mc.exe https://dl.min.io/client/mc/release/windows-amd64/mc.exe
)

REM Configure MinIO client
echo Configuring MinIO client...
set MINIO_USER=%MINIO_ROOT_USER%
set MINIO_PASS=%MINIO_ROOT_PASSWORD%
set BUCKET_NAME=%S3_BUCKET%

if "%MINIO_USER%"=="" set MINIO_USER=minioadmin
if "%MINIO_PASS%"=="" set MINIO_PASS=minioadmin123
if "%BUCKET_NAME%"=="" set BUCKET_NAME=swiftlead-storage

mc.exe alias set local http://localhost:9000 %MINIO_USER% %MINIO_PASS%

REM Create bucket if it doesn't exist
echo Creating bucket: %BUCKET_NAME%
mc.exe mb local/%BUCKET_NAME% 2>nul

if %ERRORLEVEL% EQU 0 (
    echo Bucket %BUCKET_NAME% created successfully
) else (
    echo Bucket %BUCKET_NAME% already exists or creation failed
)

REM Optional: Set public access for certain folders
REM Uncomment these lines if you want public access
REM mc.exe anonymous set download local/%BUCKET_NAME%/articles
REM mc.exe anonymous set download local/%BUCKET_NAME%/ebooks
REM mc.exe anonymous set download local/%BUCKET_NAME%/users

echo MinIO initialization completed
pause