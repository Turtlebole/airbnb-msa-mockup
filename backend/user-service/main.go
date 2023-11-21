package main

import (
	"os"
<<<<<<< HEAD
=======
	"time"
>>>>>>> main

	routes "backend/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
<<<<<<< HEAD
=======
	cors "github.com/itsjamie/gin-cors"
>>>>>>> main
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

<<<<<<< HEAD
=======
	// CORS
	router.Use(cors.Middleware(cors.Config{
		Origins:         "http://localhost:4200, *",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

>>>>>>> main
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

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
