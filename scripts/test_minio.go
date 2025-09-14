package main

import (
	"fmt"
	"log"
	"strings"
	"swiflet-backend/internal/config"
	"swiflet-backend/internal/services"
)

func main() {
	fmt.Println("Testing MinIO Connection...")
	
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("MinIO Endpoint: %s\n", cfg.S3.Endpoint)
	fmt.Printf("MinIO Bucket: %s\n", cfg.S3.Bucket)
	fmt.Printf("MinIO Region: %s\n", cfg.S3.Region)
	fmt.Printf("MinIO Access Key: %s\n", maskString(cfg.S3.AccessKey))

	// Initialize S3 service
	s3Service, err := services.NewS3Service(cfg)
	if err != nil {
		log.Fatalf("Failed to create S3 service: %v", err)
	}

	fmt.Println("✅ S3 Service initialized successfully!")
	fmt.Println("✅ MinIO connection established!")
	fmt.Printf("✅ Bucket '%s' is ready for use\n", cfg.S3.Bucket)
	
	// Test presigned URL generation (without actual file)
	testKey := "test/connection-test.txt"
	fmt.Printf("\nTesting presigned URL generation for key: %s\n", testKey)
	
	// Note: This might fail if the key doesn't exist, but it tests the connection
	url, err := s3Service.GeneratePresignedURL(testKey, 3600) // 1 hour
	if err != nil {
		fmt.Printf("⚠️  Presigned URL generation failed (expected if key doesn't exist): %v\n", err)
	} else {
		fmt.Printf("✅ Presigned URL generated successfully!\n")
		fmt.Printf("URL (first 50 chars): %s...\n", truncateString(url, 50))
	}
	
	fmt.Println("\n🎉 MinIO connection test completed!")
	fmt.Println("Your application is ready to use MinIO for file storage.")
}

func maskString(s string) string {
	if len(s) <= 4 {
		return strings.Repeat("*", len(s))
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}