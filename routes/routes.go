package routes

import (
	"glosindo-backend-go/controllers"
	"glosindo-backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Public routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.Register)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth
			protected.GET("/auth/me", controllers.GetMe)
			protected.PUT("/auth/me", controllers.UpdateMe)

			// Presensi
			presensi := protected.Group("/presensi")
			{
				presensi.POST("/check-in", controllers.CheckIn)
				presensi.POST("/check-out", controllers.CheckOut)
				presensi.GET("/today", controllers.GetTodayPresensi)
				presensi.GET("/history", controllers.GetPresensiHistory)
				presensi.GET("/stats", controllers.GetPresensiStats)
			}
		}
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Glosindo Connect API (Go)",
			"version": "1.0.0",
			"status":  "running",
			"docs":    "/api/docs",
		})
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})
}