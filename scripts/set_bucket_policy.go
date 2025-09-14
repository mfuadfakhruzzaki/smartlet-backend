package main

import (
	"fmt"
	"log"
	"swiflet-backend/internal/config"
	"swiflet-backend/internal/services"
)

func main() {
	fmt.Println("Setting MinIO bucket policy for public read access...")
	
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

	fmt.Printf("Bucket: %s\n", cfg.S3.Bucket)
	fmt.Printf("Endpoint: %s\n", cfg.S3.Endpoint)

	// Check current policy
	fmt.Println("\n📋 Checking current bucket policy...")
	currentPolicy, err := s3Service.GetBucketPolicy()
	if err != nil {
		fmt.Printf("No existing policy or error: %v\n", err)
	} else {
		fmt.Printf("Current policy: %s\n", currentPolicy)
	}

	// Set public read policy
	fmt.Println("\n🔧 Setting bucket policy for public read access...")
	err = s3Service.SetBucketPolicyPublicRead()
	if err != nil {
		log.Fatalf("Failed to set bucket policy: %v", err)
	}

	fmt.Println("✅ Bucket policy set successfully!")

	// Verify the policy was set
	fmt.Println("\n✅ Verifying new policy...")
	newPolicy, err := s3Service.GetBucketPolicy()
	if err != nil {
		fmt.Printf("Error getting new policy: %v\n", err)
	} else {
		fmt.Printf("New policy set successfully!\n")
		fmt.Printf("Policy details: %s\n", newPolicy)
	}

	fmt.Println("\n🎉 Bucket is now configured for public read access!")
	fmt.Println("Files uploaded to this bucket should now be publicly accessible via direct URL.")
}