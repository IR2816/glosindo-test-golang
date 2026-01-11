package main

import (
	"fmt"
	"log"

	"glosindo-backend-go/config"
	"glosindo-backend-go/database"
	"glosindo-backend-go/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 Starting Glosindo API (Go)...")

	// Load configuration
	config.LoadConfig()
	fmt.Println("✅ Configuration loaded")

	// Connect to database
	database.ConnectDatabase()

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:52302", "https://your-frontend-domain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	port := config.AppConfig.Port
	fmt.Printf("✅ Server running on http://localhost:%s\n", port)
	fmt.Printf("📖 API Docs: http://localhost:%s/api/docs\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
