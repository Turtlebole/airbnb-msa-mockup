package repositories

import (
	"context"
	_ "errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepo struct {
	cli *mongo.Client
}

func NewReviewRepo(cli *mongo.Client) *ReviewRepo {
	return &ReviewRepo{
		cli: cli,
	}
}

func (rr *ReviewRepo) CreateReview(ctx context.Context, review *Review) error {
	reviewCollection := rr.getCollection(ctx)

	result, err := reviewCollection.InsertOne(ctx, review)
	if err != nil {
		return err
	}

	review.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (rr *ReviewRepo) GetReviewsByAccommodation(ctx context.Context, accommodationID string) ([]*Review, error) {
	reviewCollection := rr.getCollection(ctx)

	// Convert string accommodationID to ObjectId
	objID, err := primitive.ObjectIDFromHex(accommodationID)
	if err != nil {
		return nil, err
	}

	// Query the collection
	cursor, err := reviewCollection.Find(ctx, bson.M{"accommodation_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice of Review
	var reviews []*Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (rr *ReviewRepo) getCollection(ctx context.Context) *mongo.Collection {
	reviewDatabase := rr.cli.Database("reviews")
	reviewCollection := reviewDatabase.Collection("reviews")
	return reviewCollection
}
