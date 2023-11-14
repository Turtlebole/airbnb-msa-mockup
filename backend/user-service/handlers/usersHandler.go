package services

import (
	"backend/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUsersRepo struct {
	Collection *mongo.Collection
}

func NewMongoUsersRepo(client *mongo.Client, dbName, collectionName string) *MongoUsersRepo {
	return &MongoUsersRepo{
		Collection: client.Database(dbName).Collection(collectionName),
	}
}

func (m *MongoUsersRepo) GetUser(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := m.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (m *MongoUsersRepo) GetAllUsers(ctx context.Context, username string) ([]*models.User, error) {
	cursor, err := m.Collection.Find(ctx, bson.M{"name": username})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (m *MongoUsersRepo) GetUserById(ctx context.Context, id string) (models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("invalid ObjectID")
	}

	var user models.User
	err = m.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (m *MongoUsersRepo) UpdateUser(ctx context.Context, id primitive.ObjectID, username string) (string, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"name": username}}

	result, err := m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	if result.ModifiedCount == 0 {
		return "", errors.New("user not found")
	}

	return "success", nil
}

func (m *MongoUsersRepo) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := m.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (m *MongoUsersRepo) DeleteAllUsers(ctx context.Context) error {
	_, err := m.Collection.DeleteMany(ctx, bson.M{})
	return err
}

func (m *MongoUsersRepo) UpdatePassword(ctx context.Context, username string, password string) (string, error) {
	filter := bson.M{"name": username}
	update := bson.M{"$set": bson.M{"password": password}}

	result, err := m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	if result.ModifiedCount == 0 {
		return "", errors.New("user not found")
	}

	return "success", nil
}

// CountByEmail counts the number of users with a given email
func (m *MongoUsersRepo) CountByEmail(ctx context.Context, email string) (int64, error) {
	count, err := m.Collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return 0, err
	}
	return count, nil
}
