package usecase

import (
	"Ev-Charge-Hub/Server/internal/dto/request"
	"Ev-Charge-Hub/Server/internal/dto/response"
	"Ev-Charge-Hub/Server/internal/domain/models"
	"Ev-Charge-Hub/Server/internal/repository"
	"Ev-Charge-Hub/Server/utils"
	"context"
	"errors"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
//go:generate mockgen -source=user_usecase.go -destination=../mocks/mock_user_usecase.go -package=mocks
type UserUsecaseInterface interface {
	RegisterUser(ctx context.Context, req request.RegisterUserRequest) (error)
	LoginUser(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error)
}

type userUsecase struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserUsecase(userRepo repository.UserRepositoryInterface) UserUsecaseInterface {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) RegisterUser(ctx context.Context, req request.RegisterUserRequest) error {
	existingUser, _ := u.userRepo.FindByUsernameOrEmail(ctx, req.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return err
	}

	newUser := &models.UserModel{
		ID:        primitive.NewObjectID().Hex(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return u.userRepo.CreateUser(ctx, newUser)
}

func (u *userUsecase) LoginUser(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	user, err := u.userRepo.FindByUsernameOrEmail(ctx, req.UsernameOrEmail)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.ComparePassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password it wrong")
	}

	// Create JWT Token
	token, err := utils.CreateToken(user.ID, user.Username,user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &response.LoginResponse{Token: token}, nil
}