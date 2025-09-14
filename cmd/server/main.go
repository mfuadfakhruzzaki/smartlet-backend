package main

import (
	"fmt"
	"log"
	"swiflet-backend/internal/config"
	"swiflet-backend/internal/database"
	"swiflet-backend/internal/handlers"
	"swiflet-backend/internal/middleware"
	"swiflet-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database connections
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize MQTT service
	mqttService, err := services.NewMQTTService(cfg, db)
	if err != nil {
		log.Printf("Warning: Failed to initialize MQTT service: %v", err)
		log.Println("Server will continue without MQTT functionality")
		mqttService = nil
	} else {
		// Connect to MQTT broker with retry
		if err := mqttService.ConnectWithRetry(5); err != nil {
			log.Printf("Warning: Failed to connect to MQTT broker after retries: %v", err)
			log.Println("Server will continue without MQTT functionality")
			mqttService = nil
		} else {
			defer mqttService.Disconnect()
		}
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db)
	articleHandler := handlers.NewArticleHandler(db)
	iotHandler := handlers.NewIoTHandler(db)

	// Setup router
	router := setupRouter(cfg, authHandler, userHandler, articleHandler, iotHandler)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(cfg *config.Config, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, 
	articleHandler *handlers.ArticleHandler, iotHandler *handlers.IoTHandler) *gin.Engine {
	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"version": "0.0.3",
		})
	})

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Auth routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Users routes
			users := protected.Group("/users")
			{
				users.GET("", userHandler.ListUsers)
				users.POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "User created"})
				})
				users.GET("/:id", userHandler.GetUser)
				users.PATCH("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}

			// Articles routes
			articles := protected.Group("/articles")
			{
				articles.GET("", articleHandler.ListArticles)
				articles.POST("", articleHandler.CreateArticle)
				articles.GET("/:id", articleHandler.GetArticle)
				articles.PATCH("/:id", articleHandler.UpdateArticle)
				articles.DELETE("/:id", articleHandler.DeleteArticle)
			}

			// Tags routes (placeholder)
			tags := protected.Group("/tags")
			{
				tags.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				})
				tags.POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Tag created"})
				})
			}

			// Comments routes (placeholder)
			protected.GET("/articles/:id/comments", func(c *gin.Context) {
				c.JSON(200, gin.H{"data": []interface{}{}})
			})
			protected.POST("/articles/:id/comments", func(c *gin.Context) {
				c.JSON(201, gin.H{"message": "Comment created"})
			})

			// EBooks routes (placeholder)
			ebooks := protected.Group("/ebooks")
			{
				ebooks.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				})
				ebooks.POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "EBook created"})
				})
			}

			// Videos routes (placeholder)
			videos := protected.Group("/videos")
			{
				videos.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				})
				videos.POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Video created"})
				})
			}

			// Market routes (placeholder)
			protected.Group("/weekly-prices").
				GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				}).
				POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Weekly price created"})
				})

			protected.Group("/harvests").
				GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				}).
				POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Harvest created"})
				})

			protected.Group("/harvest-sales").
				GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				}).
				POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Harvest sale created"})
				})

			// IoT routes
			houses := protected.Group("/swiflet-houses")
			{
				houses.GET("", iotHandler.ListSwifletHouses)
				houses.POST("", iotHandler.CreateSwifletHouse)
			}

			devices := protected.Group("/iot-devices")
			{
				devices.GET("", iotHandler.ListIoTDevices)
				devices.POST("", iotHandler.CreateIoTDevice)
			}

			sensors := protected.Group("/sensors")
			{
				sensors.GET("", iotHandler.ListSensors)
			}

			// Request routes (placeholder)
			protected.Group("/installation-requests").
				GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				}).
				POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Installation request created"})
				})

			protected.Group("/maintenance-requests").
				GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				}).
				POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Maintenance request created"})
				})

			protected.Group("/uninstallation-requests").
				GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				}).
				POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Uninstallation request created"})
				})

			// Transaction routes (placeholder)
			protected.Group("/transactions").
				GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				}).
				POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Transaction created"})
				})

			protected.Group("/memberships").
				GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"data": []interface{}{}})
				}).
				POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Membership created"})
				})
		}
	}

	return router
}