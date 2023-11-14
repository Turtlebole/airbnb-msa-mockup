package services

import (
	"backend/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transformed the only implementation we had into an interface so we can standardize the way our code works
// and allow for easy swapping between multiple implementations.
type UsersRepo interface {
	GetUser(ctx context.Context, email *string) (models.User, error)
	GetAllUsers(ctx context.Context, username string) ([]*models.User, error)
	GetUserById(ctx context.Context, id string) (models.User, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, username string) (string, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	DeleteAllUsers(ctx context.Context) error
	UpdatePassword(ctx context.Context, username string, password string) (string, error)
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
	CountByEmail(ctx context.Context, email string) (int, error)
}
