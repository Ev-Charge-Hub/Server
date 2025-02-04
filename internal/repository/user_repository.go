package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"Ev-Charge-Hub/Server/internal/domain/models"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	collection *mongo.Collection
}
