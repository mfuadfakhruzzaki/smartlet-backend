package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"swiflet-backend/internal/config"
	"swiflet-backend/internal/services"
	"time"
)

// TestFile implements multipart.File interface
type TestFile struct {
	*bytes.Reader
}

func (tf *TestFile) Close() error {
	return nil
}

func main() {
	fmt.Println("Testing MinIO Upload with Traefik endpoint...")
	
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize S3 service
	s3Service, err := services.NewS3Service(cfg)
	if err != nil {
		log.Fatalf("Failed to create S3 service: %v", err)
	}

	// Create a test file in memory
	timestamp := time.Now().Format("20060102_150405")
	testContent := fmt.Sprintf("Test file uploaded at %s via Traefik endpoint", timestamp)
	testFile := &TestFile{bytes.NewReader([]byte(testContent))}
	
	// Create multipart file header - gunakan .jpg yang diizinkan
	fileHeader := &multipart.FileHeader{
		Filename: fmt.Sprintf("traefik-test-%s.jpg", timestamp),
		Size:     int64(len(testContent)),
	}

	fmt.Printf("Uploading test file: %s (%d bytes)\n", fileHeader.Filename, fileHeader.Size)
	fmt.Printf("Content: %s\n", testContent)
	fmt.Printf("Endpoint: %s\n", cfg.S3.Endpoint)

	// Test upload
	result, err := s3Service.UploadFile(testFile, fileHeader, "traefik-test")
	if err != nil {
		log.Fatalf("❌ Upload failed: %v", err)
	}

	fmt.Printf("\n✅ Upload successful!\n")
	fmt.Printf("URL: %s\n", result.URL)
	fmt.Printf("Key: %s\n", result.Key)
	fmt.Printf("Bucket: %s\n", result.Bucket)
	fmt.Printf("Size: %d bytes\n", result.Size)
	fmt.Printf("MIME Type: %s\n", result.MimeType)

	// Test presigned URL for the uploaded file
	fmt.Printf("\nGenerating presigned URL for uploaded file...\n")
	presignedURL, err := s3Service.GeneratePresignedURL(result.Key, 3600)
	if err != nil {
		fmt.Printf("⚠️ Presigned URL generation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Presigned URL: %s\n", presignedURL)
	}

	fmt.Printf("\n🔍 Direct URL (should work with public policy): %s\n", result.URL)
	fmt.Printf("📱 MinIO Console: https://minio.fuadfakhruz.id/browser/swiftlead-storage\n")
	
	fmt.Printf("\n📋 File details:\n")
	fmt.Printf("- Path: %s\n", result.Key)
	fmt.Printf("- Size: %d bytes\n", result.Size)
	fmt.Printf("- Public URL: %s\n", result.URL)

	fmt.Println("\n🎉 Upload test with Traefik endpoint completed successfully!")
}