package main

import (
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"os"
	"profile-service/routes"
	"time"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8088"
	}

	router := gin.New()
	router.Use(gin.Logger())

	// CORS
	router.Use(cors.Middleware(cors.Config{
		Origins:         "https://localhost:4200, *",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	routes.ProfileRoutes(router)
	// API-2
	router.GET("/api-1", func(c *gin.Context) {

		c.JSON(200, gin.H{"success": "Access granted for api-1"})

	})

	// API-1
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
