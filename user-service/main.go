package main

import (
	"12factorapp/routes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)

	// Load .env file to get configuration values
	e := godotenv.Load()
	if e != nil {
		logger.Fatal(e)
	}

	// Read environment variables
	port := os.Getenv("PORT")
	mongoURL := os.Getenv("MONGODB_URL")

	if len(port) == 0 {
		port = "8080"
	}

	if len(mongoURL) == 0 {
		logger.Fatal("MONGODB_URL is not defined in env file.")
	}

	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logger.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	router := gin.New()
	router.Use(gin.Logger())

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

	// CORS settings
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
	)

	// Initialize the HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)

	// Start the server in a goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	logger.Println("Received termination signal:", sig)

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Graceful shutdown error:", err)
	}

	logger.Println("Server shut down gracefully")
}
