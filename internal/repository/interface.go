package repository

import (
	"PersonalAccountAPI/internal/models"
	"context"
)

type RepoProvider interface {
	GetIDByLogin(ctx context.Context, user models.UserDTO) (int, error)
	GetUserByID(ctx context.Context, user models.UserDTO) (string, error)
	AddingUser(ctx context.Context, user models.UserDTO) error
	UpdateUser(ctx context.Context, user models.UserDTO) error
}
