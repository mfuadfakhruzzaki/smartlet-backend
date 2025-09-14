# API Testing & Troubleshooting Summary

## 🚨 Error yang Anda Alami

```
Error: "Failed to upload file: EmptyStaticCreds: static credentials are empty"
Status: 500 Internal Server Error
Endpoint: POST /v1/upload/profile
```

## ✅ Solusi Cepat

### Option 1: Automatic Setup (RECOMMENDED)
```bash
# Windows
.\setup-dev.bat

# Linux/Mac
chmod +x setup-dev.sh
./setup-dev.sh
```

### Option 2: Manual Setup
1. **Copy environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Configure MinIO in `.env`:**
   ```env
   S3_ACCESS_KEY=minioadmin
   S3_SECRET_KEY=minioadmin
   S3_BUCKET=swiflet-storage
   S3_REGION=us-east-1
   S3_ENDPOINT=http://localhost:9000
   ```

3. **Start MinIO:**
   ```bash
   docker run -d --name smartlet-minio -p 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"
   ```

4. **Create bucket:**
   - Open: http://localhost:9001
   - Login: `minioadmin` / `minioadmin`
   - Create bucket: `swiflet-storage`
   - Set policy: Public read access

5. **Restart server:**
   ```bash
   go run cmd/server/main.go
   ```

## 🧪 Testing Upload Endpoints

Setelah konfigurasi S3, test endpoint berikut di Postman:

### 1. Upload Profile Image
```http
POST {{base_url}}/v1/upload/profile
Authorization: Bearer {{auth_token}}
Content-Type: multipart/form-data

Body: form-data
- key: image
- value: [select image file]
```

### 2. Upload Article Cover
```http
POST {{base_url}}/v1/upload/article
Authorization: Bearer {{auth_token}}
Content-Type: multipart/form-data

Body: form-data
- key: cover
- value: [select image file]
```

### 3. Upload EBook File
```http
POST {{base_url}}/v1/upload/ebook
Authorization: Bearer {{auth_token}}
Content-Type: multipart/form-data

Body: form-data
- key: ebook
- value: [select PDF file]
```

## 📝 File Restrictions

- **Images**: .jpg, .jpeg, .png, .gif, .webp (max 10MB)
- **Documents**: .pdf, .epub, .mobi, .doc, .docx (max 50MB)

## ✅ Expected Response (Success)

```json
{
  "message": "Profile image uploaded successfully",
  "url": "http://localhost:9000/swiflet-storage/profiles/user-1/filename.jpg",
  "size": 1024000
}
```

## 🔍 Validation Checklist

- [ ] `.env` file exists with S3 credentials
- [ ] MinIO container is running (check: `docker ps`)
- [ ] MinIO bucket `swiflet-storage` exists
- [ ] Bucket has public read policy
- [ ] Server restarted after `.env` changes
- [ ] Valid JWT token in Authorization header
- [ ] File type is allowed (.jpg, .png, .pdf, etc.)
- [ ] File size under limits (10MB images, 50MB docs)

## 🔧 Debugging Commands

```bash
# Check if MinIO is running
docker ps | grep minio

# Check server logs for S3 errors
go run cmd/server/main.go 2>&1 | grep -i s3

# Test MinIO connectivity
curl http://localhost:9000/minio/health/live

# Check environment variables
cat .env | grep S3
```

## 📁 Files untuk Reference

- `S3_UPLOAD_CONFIGURATION.md` - Detailed S3 setup guide
- `POSTMAN_COLLECTION_README.md` - API testing guide
- `Smartlet_Backend_API_Collection.postman_collection.json` - Postman collection
- `.env.example` - Environment template
- `setup-dev.bat` / `setup-dev.sh` - Auto setup scripts

## 🆘 Jika Masih Error

1. **Check server logs** saat hit upload endpoint
2. **Verify MinIO Console** dapat diakses di http://localhost:9001
3. **Test other endpoints** untuk pastikan token valid
4. **Check file format** dan ukuran sesuai restrictions
5. **Restart everything** (MinIO + Go server)

## 🎯 Working Flow untuk Testing

1. ✅ **Setup S3** (gunakan setup script)
2. ✅ **Start server** (`go run cmd/server/main.go`)
3. ✅ **Import Postman collection**
4. ✅ **Register/Login** untuk dapat token
5. ✅ **Test non-upload endpoints** dulu
6. ✅ **Test upload endpoints** dengan file valid
7. ✅ **Verify uploaded files** di MinIO console