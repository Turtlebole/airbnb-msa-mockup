package repositories

import (
	"context"
	"errors"
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

// NoSQL: ProfileRepo struct encapsulating Mongo api client
type ProfileRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment
func New(ctx context.Context, logger *log.Logger) (*ProfileRepo, error) {
	dburi := os.Getenv("MONGODB_URL")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &ProfileRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// Disconnect from database
func (ar *ProfileRepo) Disconnect(ctx context.Context) error {
	err := ar.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (ar *ProfileRepo) Ping() {
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

func CreateProfile(client *mongo.Client, profile *Profile) error {
	profileCollection := client.Database("cluster0").Collection("user")

	result, err := profileCollection.InsertOne(context.TODO(), profile)
	if err != nil {
		return err
	}

	profile.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func GetAllProfiles(client *mongo.Client) (Profiles, error) {
	profileCollection := client.Database("cluster0").Collection("user")

	cursor, err := profileCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var profiles Profiles
	for cursor.Next(context.TODO()) {
		var profile Profile
		if err := cursor.Decode(&profile); err != nil {
			return nil, err
		}
		profiles = append(profiles, &profile)
	}

	return profiles, nil
}

func GetProfileByID(client *mongo.Client, userID string) (*Profile, error) {
	profileCollection := client.Database("cluster0").Collection("user")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var profile Profile
	err = profileCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func UpdateProfile(client *mongo.Client, profileID primitive.ObjectID, updatedProfile *Profile) error {
	profileCollection := client.Database("cluster0").Collection("user")

	//updatedProfile.Updated_On = time.Now().Format(time.RFC3339)

	result, err := profileCollection.ReplaceOne(context.TODO(), bson.M{"_id": profileID}, updatedProfile)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("profile not found")
	}

	return nil
}

func DeleteProfile(client *mongo.Client, profileID primitive.ObjectID) error {
	profileCollection := client.Database("cluster0").Collection("user")

	result, err := profileCollection.DeleteOne(context.TODO(), bson.M{"_id": profileID})
	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return errors.New("profile not found")
	}

	return nil
}

func GetProfileByEmail(client *mongo.Client, email string) (*Profile, error) {
	profileCollection := client.Database("cluster0").Collection("user")

	// Create a filter for the email
	filter := bson.D{{"email", email}}

	// Find the user in the database
	var profile Profile
	err := profileCollection.FindOne(context.TODO(), filter).Decode(&profile)
	if err != nil {
		// Handle the error (e.g., user not found)
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	return &profile, nil
}

func (ar *ProfileRepo) getCollection() *mongo.Collection {
	accommodationDatabase := ar.cli.Database("cluster0")
	accommodationsCollection := accommodationDatabase.Collection("user")
	return accommodationsCollection
}
