package main

import (
	"distributed-chat-system/internal/apis/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Check if running inside Docker
	if os.Getenv("IS_DOCKER") == "true" {
		log.Println("Running in Docker, Env values satisfied by default.")
	} else {
		// Load environment variables from .env file if not in Docker
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	// Create a Gin router instance
	router := gin.Default()

	// Setup routes
	routes.Setup(router)

	// Get the port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback to default port if not set
		log.Printf("No PORT environment variable detected, using default port: %s", port)
	}

	// Start the Gin server
	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
