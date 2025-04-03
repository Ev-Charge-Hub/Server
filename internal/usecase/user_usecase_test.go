package usecase_test

import (
	"context"
	"testing"
	"time"

	"Ev-Charge-Hub/Server/internal/domain/models"
	"Ev-Charge-Hub/Server/internal/dto/request"
	"Ev-Charge-Hub/Server/internal/mocks"
	"Ev-Charge-Hub/Server/internal/usecase"
	"Ev-Charge-Hub/Server/utils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	uc := usecase.NewUserUsecase(mockRepo)

	ctx := context.TODO()
	req := request.RegisterUserRequest{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "USER",
	}

	mockRepo.EXPECT().
		FindByUsernameOrEmail(ctx, req.Email).
		Return(nil, nil)

	mockRepo.EXPECT().
		CreateUser(ctx, gomock.Any()).
		Return(nil)

	// Act
	err := uc.RegisterUser(ctx, req)

	// Assert
	assert.NoError(t, err)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	uc := usecase.NewUserUsecase(mockRepo)

	ctx := context.TODO()
	req := request.RegisterUserRequest{
		Username: "test",
		Email:    "existing@example.com",
		Password: "password123",
		Role:     "USER",
	}

	mockRepo.EXPECT().
		FindByUsernameOrEmail(ctx, req.Email).
		Return(&models.UserModel{Email: req.Email}, nil)

	// Act
	err := uc.RegisterUser(ctx, req)

	// Assert
	assert.EqualError(t, err, "email already exists")
}

func TestLoginUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	uc := usecase.NewUserUsecase(mockRepo)

	ctx := context.TODO()
	plainPassword := "password123"
	hashedPassword, _ := utils.EncryptPassword(plainPassword)

	user := &models.UserModel{
		ID:        "1",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		Role:      "USER",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.EXPECT().
		FindByUsernameOrEmail(ctx, user.Email).
		Return(user, nil)

	req := request.LoginRequest{
		UsernameOrEmail: user.Email,
		Password:        plainPassword,
	}

	// Act
	resp, err := uc.LoginUser(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token, "token should not be empty")
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	uc := usecase.NewUserUsecase(mockRepo)

	ctx := context.TODO()

	user := &models.UserModel{
		ID:       "1",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "$2a$10$invalidhashvaluehere1234567890", // invalid bcrypt
		Role:     "USER",
	}

	mockRepo.EXPECT().
		FindByUsernameOrEmail(ctx, user.Email).
		Return(user, nil)

	req := request.LoginRequest{
		UsernameOrEmail: user.Email,
		Password:        "wrongpassword",
	}

	// Act
	resp, err := uc.LoginUser(ctx, req)

	// Assert
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid email or password it wrong")
}
