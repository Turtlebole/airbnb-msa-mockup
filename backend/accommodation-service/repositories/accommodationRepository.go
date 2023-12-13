package repositories

import (
	"accommodation-service/database"
	"accommodation-service/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

// NoSQL: AccommodationRepo struct encapsulating Mongo api client
type AccommodationRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment
func New(ctx context.Context, logger *log.Logger) (*AccommodationRepo, error) {
	dburi := os.Getenv("MONGODB_URL")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &AccommodationRepo{
		cli:    client,
		logger: logger,
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

func (ar *AccommodationRepo) GetAll() (*models.Accommodation, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accommodationsCollection := ar.getCollection()

	var accommodations models.Accommodation
	accommodationsCursor, err := accommodationsCollection.Find(ctx, bson.M{})
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}

	if err = accommodationsCursor.All(ctx, &accommodations); err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	return &accommodations, nil
}

func (ar *AccommodationRepo) Insert(accommodation *models.Accommodation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accommodationsCollection := ar.getCollection()

	result, err := accommodationsCollection.InsertOne(ctx, &accommodation)
	if err != nil {
		ar.logger.Println(err)
		return err
	}
	ar.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (ar *AccommodationRepo) Update(id string, accommodation *models.Accommodation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	//accommodationsCollection := ar.getCollection()

	var accommodationsCollection *mongo.Collection = database.OpenCollection(database.Client, "accommodations")

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	// filter := bson.M{"name": "sadradimozda"}
	update := bson.M{
		"$set": bson.M{
			"name":            accommodation.Name,
			"minguests":       accommodation.MinGuests,
			"maxguests":       accommodation.MaxGuests,
			"location":        accommodation.Location,
			"amenities":       accommodation.Amenities,
			"pricepernight": accommodation.PricePerNight,
		},
	}

	result, err := accommodationsCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no document found for update")
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
	accommodationDatabase := ar.cli.Database("accommodationServiceDB")
	accommodationsCollection := accommodationDatabase.Collection("accommodations")
	return accommodationsCollection
}
