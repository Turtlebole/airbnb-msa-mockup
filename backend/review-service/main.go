package main

import (
	"context"
	cors "github.com/itsjamie/gin-cors"
	"log"
	"os"
	"review-service/routes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var reviewCollection *mongo.Collection

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Get the value of the PORT variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8010"
	}

	// Initialize MongoDB client
	mongoClient, err := initMongoClient()
	if err != nil {
		log.Fatal("Error initializing MongoDB client: ", err)
	}

	// Initialize reviewCollection
	reviewCollection = mongoClient.Database("reviews").Collection("reviews")

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

	// Routes initialization
	routes.ReviewRoutes(router, mongoClient)

	// Run the server
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Error running the server: ", err)
	}
}

func initMongoClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
