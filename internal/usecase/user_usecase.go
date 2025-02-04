package usecase

import (
	"Ev-Charge-Hub/Server/dto/request"
	"Ev-Charge-Hub/Server/internal/repository"
	"context"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, req request.RegisterUserRequest) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}