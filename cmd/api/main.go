package main

import (
	"log"
	"os"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env vars")
	}

	// Initialize Database
	database.Connect()

	// Initialize Router
	r := gin.Default()

	// Setup Routes
	routes.SetupRoutes(r)

	// Start Server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("🚀 API Server corriendo en http://localhost:%s (GO VERSION)", port)
	r.Run(":" + port)
}
