package main

import (
	"log"
	"os"

	"admin_statistics_api/config"
	"admin_statistics_api/handlers"
	"admin_statistics_api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	// Connect to databases
	config.ConnectDatabase()
	defer config.DisconnectDatabase()

	config.ConnectRedis()
	defer config.DisconnectRedis()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize handlers
	statsHandler := handlers.NewStatisticsHandler()

	// Public routes (no auth required)
	router.GET("/health", statsHandler.HealthCheck)

	// Protected routes (require authentication)
	api := router.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/gross_gaming_rev", statsHandler.GetGrossGamingRevenue)
		api.GET("/daily_wager_volume", statsHandler.GetDailyWagerVolume)
		api.GET("/user/:user_id/wager_percentile", statsHandler.GetUserWagerPercentile)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	log.Printf("API endpoints require Authorization header")
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}