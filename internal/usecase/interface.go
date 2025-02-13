package usecase

import (
	"context"
	"mime/multipart"
)

type Provider interface {
	GetIDByLoginFromDB(ctx context.Context, login, password string) (int, error)
	GetUserByIDFromDB(ctx context.Context, id int) (string, error)
	AddingUserToDB(ctx context.Context, id int, login, password string) error
	UpdateUserInDB(ctx context.Context, id int, login, password string) error
	SetFile(file *multipart.FileHeader)
	UploadFile(ctx context.Context) error
}
