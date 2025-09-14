package main

import (
	"fmt"
	"log"
	"strings"
	"swiflet-backend/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	fmt.Println("=== AWS S3 Connectivity Test ===")
	fmt.Println()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create AWS session
	awsConfig := &aws.Config{
		Region:      aws.String(cfg.S3.Region),
		Credentials: credentials.NewStaticCredentials(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	fmt.Printf("Testing connection to AWS S3 in region: %s\n", cfg.S3.Region)
	fmt.Printf("Bucket: %s\n", cfg.S3.Bucket)
	fmt.Println()

	// Test 1: List all buckets (check if credentials work)
	fmt.Println("Test 1: List S3 buckets (credential test)...")
	listResult, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Printf("❌ AWS Error: %s - %s\n", awsErr.Code(), awsErr.Message())
			
			switch awsErr.Code() {
			case "InvalidAccessKeyId":
				fmt.Println("   🔧 Fix: S3_ACCESS_KEY is invalid")
			case "SignatureDoesNotMatch":
				fmt.Println("   🔧 Fix: S3_SECRET_KEY is invalid")
			case "AccessDenied":
				fmt.Println("   🔧 Fix: IAM user doesn't have ListBuckets permission")
			default:
				fmt.Printf("   🔧 Fix: Check AWS credentials and permissions\n")
			}
		} else {
			fmt.Printf("❌ Error: %v\n", err)
		}
		return
	}

	fmt.Printf("✅ Successfully connected to AWS S3\n")
	fmt.Printf("   Found %d buckets accessible to this IAM user\n", len(listResult.Buckets))
	
	// Check if our target bucket exists
	bucketExists := false
	for _, bucket := range listResult.Buckets {
		if *bucket.Name == cfg.S3.Bucket {
			bucketExists = true
			fmt.Printf("   ✅ Target bucket '%s' found\n", cfg.S3.Bucket)
			break
		}
	}

	if !bucketExists {
		fmt.Printf("   ❌ Target bucket '%s' NOT found\n", cfg.S3.Bucket)
		fmt.Println("   Available buckets:")
		for _, bucket := range listResult.Buckets {
			fmt.Printf("     - %s\n", *bucket.Name)
		}
		fmt.Println()
		fmt.Println("   🔧 Fix options:")
		fmt.Println("     1. Create bucket 'swiflet-storage' in AWS S3 console")
		fmt.Println("     2. Change S3_BUCKET in .env to an existing bucket")
		return
	}

	fmt.Println()

	// Test 2: Get bucket location
	fmt.Println("Test 2: Check bucket region...")
	locationResult, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(cfg.S3.Bucket),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Printf("❌ AWS Error: %s - %s\n", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Printf("❌ Error: %v\n", err)
		}
		return
	}

	bucketRegion := ""
	if locationResult.LocationConstraint == nil {
		bucketRegion = "us-east-1" // Default region when nil
	} else {
		bucketRegion = *locationResult.LocationConstraint
	}

	fmt.Printf("✅ Bucket region: %s\n", bucketRegion)
	
	if bucketRegion != cfg.S3.Region {
		fmt.Printf("❌ Region mismatch!\n")
		fmt.Printf("   Bucket is in: %s\n", bucketRegion)
		fmt.Printf("   Config expects: %s\n", cfg.S3.Region)
		fmt.Printf("   🔧 Fix: Update S3_REGION in .env to '%s'\n", bucketRegion)
		return
	}

	fmt.Println()

	// Test 3: Test upload permission with a small test object
	fmt.Println("Test 3: Test upload permission...")
	testKey := "test-upload-permission.txt"
	testContent := "This is a test upload to verify permissions"

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(cfg.S3.Bucket),
		Key:    aws.String(testKey),
		Body:   strings.NewReader(testContent),
		// Remove ACL since bucket doesn't support it
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Printf("❌ Upload failed: %s - %s\n", awsErr.Code(), awsErr.Message())
			
			switch awsErr.Code() {
			case "AccessDenied":
				fmt.Println("   🔧 Fix: IAM user needs s3:PutObject permission")
			case "InvalidBucketName":
				fmt.Println("   🔧 Fix: Bucket name is invalid")
			case "NoSuchBucket":
				fmt.Println("   🔧 Fix: Bucket doesn't exist")
			default:
				fmt.Println("   🔧 Fix: Check IAM permissions for S3 upload")
			}
		} else {
			fmt.Printf("❌ Error: %v\n", err)
		}
		return
	}

	fmt.Printf("✅ Upload test successful!\n")
	fmt.Printf("   Test file uploaded: s3://%s/%s\n", cfg.S3.Bucket, testKey)

	// Clean up test file
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(cfg.S3.Bucket),
		Key:    aws.String(testKey),
	})
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to clean up test file: %v\n", err)
	} else {
		fmt.Printf("✅ Test file cleaned up\n")
	}

	fmt.Println()
	fmt.Println("🎉 All AWS S3 tests passed!")
	fmt.Println("   Your S3 configuration is working correctly.")
	fmt.Println("   Upload endpoints should work now.")
	fmt.Println()
	fmt.Println("   If you still get 'EmptyStaticCreds' error:")
	fmt.Println("   1. Restart your Go server after changing .env")
	fmt.Println("   2. Make sure .env file is in the correct directory")
	fmt.Println("   3. Check for spaces or invisible characters in .env")
}