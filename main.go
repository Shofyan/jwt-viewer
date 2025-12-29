package main

import (
	"jwt-viewer/handlers"
	"jwt-viewer/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get configuration from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	if port[0] != ':' {
		port = ":" + port
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize services
	jwtService := services.NewJWTService()

	// Initialize handlers
	jwtHandler := handlers.NewJWTHandler(jwtService)

	// Setup Gin router
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// API routes
	api := router.Group("/api")
	{
		api.POST("/decode", jwtHandler.DecodeHandler)
		api.POST("/encode", jwtHandler.EncodeHandler)
		api.POST("/verify", jwtHandler.VerifyHandler)
	}

	// Start server
	log.Printf("🚀 JWT Debugger server starting on http://localhost%s", port)
	log.Printf("📝 Open your browser and navigate to http://localhost%s", port)
	log.Printf("🔧 Mode: %s", gin.Mode())

	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
