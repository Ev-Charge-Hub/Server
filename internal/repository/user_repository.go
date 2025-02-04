package repository

import (
	domainModels "Ev-Charge-Hub/Server/internal/domain/models"   // สำหรับ Domain Model
	repoModels "Ev-Charge-Hub/Server/internal/repository/models" // ตรวจสอบ path
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domainModels.User, error)
	CreateUser(ctx context.Context, user *domainModels.User) error
}

type userRepository struct {
	collection *mongo.Collection
}

// FindByUsernameOrEmail implements UserRepository.
func (u *userRepository) FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*domainModels.User, error) {
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

	return &domainModels.User{
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
func (u *userRepository) CreateUser(ctx context.Context, user *domainModels.User) error {
	userDB := repoModels.UserDB{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	_, err := u.collection.InsertOne(ctx, userDB)
	return err
}

// Comfirm Protocal -> UserRepository
func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{collection: db.Collection("users")}
}
