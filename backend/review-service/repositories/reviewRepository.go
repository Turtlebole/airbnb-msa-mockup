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
func (rr *ReviewRepo) GetAverageRatingByAccommodation(ctx context.Context, accommodationID string) (float64, error) {
	reviewCollection := rr.getCollection(ctx)

	// Convert string accommodationID to ObjectId
	objID, err := primitive.ObjectIDFromHex(accommodationID)
	if err != nil {
		return 0, err
	}

	// Query the collection
	cursor, err := reviewCollection.Find(ctx, bson.M{"accommodation_id": objID})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	// Calculate the average rating
	var totalRating, reviewCount int
	for cursor.Next(ctx) {
		var review Review
		if err := cursor.Decode(&review); err != nil {
			return 0, err
		}
		totalRating += review.Rating
		reviewCount++
	}

	if reviewCount == 0 {
		return 0, nil // No reviews, return 0
	}

	return float64(totalRating) / float64(reviewCount), nil
}
func (rr *ReviewRepo) GetReviewsByUser(ctx context.Context, userID string) ([]*Review, error) {
	reviewCollection := rr.getCollection(ctx)

	// Convert string userID to ObjectId
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Query the collection
	cursor, err := reviewCollection.Find(ctx, bson.M{"user_id": objID})
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
func (rr *ReviewRepo) UpdateReview(ctx context.Context, reviewID string, updatedReview *Review) error {
	reviewCollection := rr.getCollection(ctx)

	// Convert string reviewID to ObjectId
	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return err
	}

	// Update the review by its ID
	_, err = reviewCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"comment": updatedReview.Comment, "rating": updatedReview.Rating, "host_comment": updatedReview.HostComment, "host_rating": updatedReview.HostRating}},
	)
	if err != nil {
		return err
	}

	return nil
}

func (rr *ReviewRepo) DeleteReview(ctx context.Context, reviewID string) error {
	reviewCollection := rr.getCollection(ctx)

	// Convert string reviewID to ObjectId
	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return err
	}

	// Delete the review by its ID
	_, err = reviewCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	return nil
}

func (rr *ReviewRepo) getCollection(ctx context.Context) *mongo.Collection {
	reviewDatabase := rr.cli.Database("reviews")
	reviewCollection := reviewDatabase.Collection("reviews")
	return reviewCollection
}
func (rr *ReviewRepo) CreateHostReview(ctx context.Context, review *Review) error {
	// Set AccommodationID to zero as it's not applicable for host reviews
	review.AccommodationID = primitive.NilObjectID

	reviewCollection := rr.getCollection(ctx)

	result, err := reviewCollection.InsertOne(ctx, review)
	if err != nil {
		return err
	}

	review.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}
func (rr *ReviewRepo) GetReviewsByUserAndHosts(ctx context.Context, userID string) ([]*Review, error) {
	reviewCollection := rr.getCollection(ctx)

	// Convert string userID to ObjectId
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Query the collection
	cursor, err := reviewCollection.Find(ctx, bson.M{"user_id": objID})
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
func (rr *ReviewRepo) GetHostReviews(ctx context.Context, hostID string) ([]*Review, error) {
	reviewCollection := rr.getCollection(ctx)

	// Convert string hostID to ObjectId
	objID, err := primitive.ObjectIDFromHex(hostID)
	if err != nil {
		return nil, err
	}

	// Query the collection
	cursor, err := reviewCollection.Find(ctx, bson.M{"host_id": objID})
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
func GetReviewsByAccommodationAndHost(client *mongo.Client, accommodationID string, hostID string) ([]Review, error) {
	collection := client.Database("reviews").Collection("reviews")

	filter := bson.M{
		"accommodationID": accommodationID,
		"hostID":          hostID,
	}

	var reviews []Review
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}
func (rr *ReviewRepo) GetAllReviews(ctx context.Context) ([]*Review, error) {
	reviewCollection := rr.getCollection(ctx)

	// Query the collection to get all reviews
	cursor, err := reviewCollection.Find(ctx, bson.M{})
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
func (rr *ReviewRepo) GetReviewByID(ctx context.Context, reviewID string) (*Review, error) {
	collection := rr.getCollection(ctx)

	// Convert string reviewID to ObjectId
	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return nil, err
	}

	var review Review
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&review)
	if err != nil {
		return nil, err
	}

	return &review, nil
}
