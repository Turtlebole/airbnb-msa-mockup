package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	// Create a new router
	router := mux.NewRouter()

	// router.HandleFunc("/users")
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", usersHandler.GetUsers)

	getAllRouter := router.Methods(http.MethodGet).Subrouter()
	getAllRouter.HandleFunc("/all", usersHandler.GetAllUsers)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", usersHandler.PostUsers)
	postRouter.Use(usersHandler.MiddlewareUsersValidation)

	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", usersHandler.PutUsers)
	putRouter.Use(usersHandler.MiddlewareUsersValidation)

	deleteHandler := router.Methods(http.MethodDelete).Subrouter()
	deleteHandler.HandleFunc("/{id:[0-9]+}", usersHandler.DeleteUsers)

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
