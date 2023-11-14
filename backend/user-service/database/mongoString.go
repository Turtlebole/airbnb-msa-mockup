package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBInstance func
func db() *mongo.Client {
	// Replace the following line with your MongoDB connection string
	mongoURI := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

// Client Database instance
var Client *mongo.Client = db()

// Function that makes a connection with a collection in the database
func DbConnection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}
