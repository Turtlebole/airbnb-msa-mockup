package data

import (
	"context"
	"fmt"
	protos "github.com/SimikicAleksandar/AccomodationService/protos/main"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

// NoSQL: AccommodationRepo enkapsulira mongo api klijenta
type AccommodationRepo struct {
	cli            *mongo.Client
	logger         *log.Logger
	databaseName   string
	collectionName string
}

// NoSQL: Konstrukto za citanje konfiguracije baze
func New(ctx context.Context, logger *log.Logger, mongoClient *mongo.Client, dbName, collectionName string) (*AccommodationRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &AccommodationRepo{
		cli:            mongoClient,
		logger:         logger,
		databaseName:   dbName,
		collectionName: collectionName,
	}, nil
}

// Disconnect from database
func (ar *AccommodationRepo) Disconnect(ctx context.Context) error {
	err := ar.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (ar *AccommodationRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := ar.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		ar.logger.Println(err)
	}

	// Print available databases
	databases, err := ar.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		ar.logger.Println(err)
	}
	fmt.Println(databases)
}

func (ar *AccommodationRepo) GetAll() (Accommodations, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accommodationsCollection := ar.getCollection()

	var accommodations Accommodations
	accommodationsCursor, err := accommodationsCollection.Find(ctx, bson.M{})
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	if err = accommodationsCursor.All(ctx, &accommodations); err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	return accommodations, nil
}

func (ar *AccommodationRepo) Insert(profile *protos.AccommodationResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accommodationsCollection := ar.getCollection()

	result, err := accommodationsCollection.InsertOne(ctx, &profile)
	if err != nil {
		ar.logger.Println(err)
		return err
	}
	ar.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (ar *AccommodationRepo) Update(accommodation *protos.AccommodationResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accommodationCollection := ar.getCollection()

	filter := bson.M{"email": accommodation.GetEmail()}
	update := bson.M{"$set": bson.M{
		"name":     accommodation.GetName(),
		"price":    accommodation.GetPrice(),
		"location": accommodation.GetLocation(),
		"adress":   accommodation.GetAdress(),
	}}
	result, err := accommodationCollection.UpdateOne(ctx, filter, update)
	ar.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	ar.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		ar.logger.Println(err)
		return err
	}
	return nil
}

func (ar *AccommodationRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accommodationsCollection := ar.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	result, err := accommodationsCollection.DeleteOne(ctx, filter)
	if err != nil {
		ar.logger.Println(err)
		return err
	}
	ar.logger.Printf("Documents deleted: %v\n", result.DeletedCount)
	return nil
}

func (ar *AccommodationRepo) getCollection() *mongo.Collection {
	accommodationDatabase := ar.cli.Database("mongoDemo")
	accommodationsCollection := accommodationDatabase.Collection("accommodations")
	return accommodationsCollection
}

func (ar *AccommodationRepo) GetById(email string) ([]*protos.AccommodationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accommodationCollection := ar.getCollection()
	var accommodationsSlice []*protos.AccommodationResponse

	// Assuming you have a filter based on the email, modify the filter as needed
	filter := bson.M{"email": email}

	accommodationCursor, err := accommodationCollection.Find(ctx, filter)
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	defer func(accommodationCursor *mongo.Cursor, ctx context.Context) {
		err := accommodationCursor.Close(ctx)
		if err != nil {
			ar.logger.Println(err)
		}
	}(accommodationCursor, ctx)

	for accommodationCursor.Next(ctx) {
		var accommodation protos.AccommodationResponse
		if err := accommodationCursor.Decode(&accommodation); err != nil {
			ar.logger.Println(err)
			return nil, err
		}
		accommodationsSlice = append(accommodationsSlice, &accommodation)
	}

	if err := accommodationCursor.Err(); err != nil {
		ar.logger.Println(err)
		return nil, err
	}

	return accommodationsSlice, nil
}
