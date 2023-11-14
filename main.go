package main

import (
	"context"
	"github.com/SimikicAleksandar/AccomodationService/data"
	"github.com/SimikicAleksandar/AccomodationService/handlers"
	protos "github.com/SimikicAleksandar/AccomodationService/protos/main"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Set your MongoDB Atlas connection string
	connectionString := "mongodb+srv://simikic:simikic@cluster0.kxo3ocx.mongodb.net/?retryWrites=true&w=majority"

	// Connect to MongoDB Atlas
	client, err := mongo.Connect(timeoutContext, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(timeoutContext); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Ensure that the MongoDB connection is established
	err = client.Ping(timeoutContext, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	lis, err := net.Listen("tcp", ":9093")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}
	serverRegister := grpc.NewServer()

	logger := log.New(os.Stdout, "[accommodation-main] ", log.LstdFlags)

	// Initialize the handler and inject the logger and MongoDB collection
	accommodationLog := log.New(os.Stdout, "[accommodation-repo-log] ", log.LstdFlags)
	accommodationRepo, err := data.New(timeoutContext, accommodationLog, client, "accommodationDB", "accommodations")
	if err != nil {
		log.Fatalf("Error creating data repository: %v", err)
	}
	service := handlers.NewServer(logger, accommodationRepo)

	protos.RegisterAccommodationServer(serverRegister, service)
	err = serverRegister.Serve(lis)
	if err != nil {
		log.Fatalf("Error serving gRPC server: %v", err)
	}
}
