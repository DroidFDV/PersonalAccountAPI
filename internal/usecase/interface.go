package usecase

import (
	"PersonalAccountAPI/internal/models"
	"context"
)

type UserProvider interface {
	GetIDByLogin(ctx context.Context, user models.UserRequest) (models.UserResponse, error)
	GetUserByID(ctx context.Context, user models.UserRequest) (models.UserResponse, error)
	AddUser(ctx context.Context, user models.UserRequest) error
	UpdateUser(ctx context.Context, user models.UserRequest) error
}
