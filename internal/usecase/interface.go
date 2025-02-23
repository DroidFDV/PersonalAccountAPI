package usecase

import (
	"PersonalAccountAPI/internal/models"
	"context"
	"mime/multipart"
)

type Provider interface {
	GetIDByLoginFromDB(ctx context.Context, user models.UserRequest) (int, error)
	GetUserByIDFromDB(ctx context.Context, user models.UserRequest) (string, error)
	AddingUserToDB(ctx context.Context, user models.UserRequest) error
	UpdateUserInDB(ctx context.Context, user models.UserRequest) error
	UploadFile(ctx context.Context, id string, file *multipart.FileHeader) func(ctx context.Context) error
}
