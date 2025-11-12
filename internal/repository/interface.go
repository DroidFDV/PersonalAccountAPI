package repository

import (
	"PersonalAccountAPI/internal/models"
	"context"
)

type RepoProvider interface {
	GetIDByLogin(ctx context.Context, user *models.UserDTO) (*models.UserDTO, error)
	GetUserByID(ctx context.Context, user *models.UserDTO) (*models.UserDTO, error)
	AddUser(ctx context.Context, user *models.UserDTO) error
	UpdateUser(ctx context.Context, user *models.UserDTO) error
}
