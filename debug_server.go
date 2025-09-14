package main

import (
	"fmt"
	"log"
	"net/http"
	"swiflet-backend/internal/config"
	"swiflet-backend/internal/database"
	"swiflet-backend/internal/handlers"
	"swiflet-backend/internal/middleware"
	"swiflet-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("=== Debugging Server with Detailed Logging ===")
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Print S3 config (masked)
	fmt.Printf("S3 Config:\n")
	fmt.Printf("  AccessKey: %s\n", maskString(cfg.S3.AccessKey))
	fmt.Printf("  SecretKey: %s\n", maskString(cfg.S3.SecretKey))
	fmt.Printf("  Bucket: %s\n", cfg.S3.Bucket)
	fmt.Printf("  Region: %s\n", cfg.S3.Region)
	fmt.Printf("  Endpoint: %s\n", cfg.S3.Endpoint)
	fmt.Println()

	// Initialize database connections
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize S3 service with debug
	fmt.Println("Initializing S3 service...")
	s3Service, err := services.NewS3Service(cfg)
	if err != nil {
		log.Fatalf("Failed to create S3 service: %v", err)
	}
	fmt.Println("✅ S3 service initialized successfully")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)

	// Setup router with debug mode
	gin.SetMode("debug")
	router := gin.New()

	// Custom logging middleware for upload endpoints
	router.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/v1/upload/profile" {
			fmt.Printf("\n=== UPLOAD DEBUG ===\n")
			fmt.Printf("Method: %s\n", c.Request.Method)
			fmt.Printf("Path: %s\n", c.Request.URL.Path)
			fmt.Printf("Content-Type: %s\n", c.Request.Header.Get("Content-Type"))
			fmt.Printf("Authorization: %s\n", maskAuthHeader(c.Request.Header.Get("Authorization")))
			fmt.Printf("Content-Length: %s\n", c.Request.Header.Get("Content-Length"))
		}
		c.Next()
		if c.Request.URL.Path == "/v1/upload/profile" {
			fmt.Printf("Response Status: %d\n", c.Writer.Status())
			fmt.Printf("=== END DEBUG ===\n\n")
		}
	})

	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Auth routes
	auth := router.Group("/v1/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/v1")
	protected.Use(middleware.AuthMiddleware(cfg))

	// Upload routes with enhanced error handling
	uploads := protected.Group("/upload")
	uploads.POST("/profile", func(c *gin.Context) {
		fmt.Println("\n🔍 Starting profile upload handler...")
		
		// Check if user is authenticated
		userID, exists := c.Get("user_id")
		if !exists {
			fmt.Println("❌ User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		fmt.Printf("✅ User authenticated: ID=%v\n", userID)

		// Check if file exists in request
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			fmt.Printf("❌ No file in request: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded: " + err.Error()})
			return
		}
		defer file.Close()
		fmt.Printf("✅ File received: %s (size: %d bytes)\n", header.Filename, header.Size)

		// Try to upload to S3
		fmt.Println("🔄 Attempting S3 upload...")
		result, err := s3Service.UploadUserProfileImage(file, header, userID.(int))
		if err != nil {
			fmt.Printf("❌ S3 Upload failed: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
			return
		}
		fmt.Printf("✅ S3 Upload successful: %s\n", result.URL)

		// Update database
		fmt.Println("🔄 Updating user profile in database...")
		_, err = db.PostgreSQL.Exec("UPDATE users SET img_profile = $1 WHERE id = $2", result.URL, userID)
		if err != nil {
			fmt.Printf("❌ Database update failed: %v\n", err)
			// Try to delete uploaded file
			s3Service.DeleteFile(result.Key)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
			return
		}
		fmt.Println("✅ Database updated successfully")

		c.JSON(http.StatusOK, gin.H{
			"message": "Profile image uploaded successfully",
			"url":     result.URL,
			"size":    result.Size,
		})
	})

	// Start server
	fmt.Println("\n🚀 Starting debug server on :8080...")
	fmt.Println("💡 Upload test: POST /v1/upload/profile with 'image' file")
	fmt.Println("💡 Auth test: POST /v1/auth/login first to get token")
	
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

func maskAuthHeader(auth string) string {
	if auth == "" {
		return "None"
	}
	if len(auth) > 20 {
		return auth[:20] + "..."
	}
	return auth
}