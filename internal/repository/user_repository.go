package repository

import (
	domainModels "Ev-Charge-Hub/Server/internal/domain/models"   // สำหรับ Domain Model
	repoModels "Ev-Charge-Hub/Server/internal/repository/models" // ตรวจสอบ path
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface {
	FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domainModels.UserModel, error)
	CreateUser(ctx context.Context, user *domainModels.UserModel) error
}

// Define Class userRepository
type userRepository struct {
	collection *mongo.Collection
}

// Comfirm Protocal -> UserRepository
func NewUserRepository(db *mongo.Database) UserRepositoryInterface {
	return &userRepository{collection: db.Collection("users")} // return pointer for pass by reference
}

// FindByUsernameOrEmail implements UserRepository.
func (u *userRepository) FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domainModels.UserModel, error) {
	var userDB repoModels.UserDB

	filter := bson.M{
		"$or": []bson.M{
			{"username": usernameOrEmail},
			{"email": usernameOrEmail},
		},
	}

	err := u.collection.FindOne(ctx, filter).Decode(&userDB)
	if err != nil {
		return nil, err
	}

	return &domainModels.UserModel{
		ID:        userDB.ID.Hex(),
		Username:  userDB.Username,
		Email:     userDB.Email,
		Password:  userDB.Password,
		Role:      userDB.Role,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
	}, nil
}

// CreateUser implements UserRepository.
func (u *userRepository) CreateUser(ctx context.Context, user *domainModels.UserModel) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return errors.New("invalid object ID")
	}
	userDB := repoModels.UserDB{
		ID:        objectID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	_, err = u.collection.InsertOne(ctx, userDB)
	return err
}
