package usecase

import (
	"PersonalAccountAPI/internal/models"
	"context"
	"mime/multipart"
)

type UserProvider interface {
	GetIDByLogin(ctx context.Context, user models.UserRequest) (int, error)
	GetUserByID(ctx context.Context, user models.UserRequest) (string, error)
	AddingUser(ctx context.Context, user models.UserRequest) error
	UpdateUser(ctx context.Context, user models.UserRequest) error
	UploadFile(ctx context.Context, file *multipart.FileHeader) func(ctx context.Context) error
}
