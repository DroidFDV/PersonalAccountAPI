package usecase

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/repository"
	"PersonalAccountAPI/internal/uploading"
	"context"
	"mime/multipart"
	"path/filepath"

	"github.com/pkg/errors"
)

type UserUsecase struct {
	repository repository.RepoProvider
}

func New(repository repository.RepoProvider) *UserUsecase {
	return &UserUsecase{
		repository: repository,
	}
}

func (u *UserUsecase) GetIDByLogin(ctx context.Context, userRequest models.UserRequest) (int, error) {
	userDTO := userRequest.ToDTO()
	id, err := u.repository.GetIDByLogin(ctx, userDTO)
	if err != nil {
		return id, errors.Wrap(err, "UserUsecase.GetUserByID pgx.QueryRow().Scan:")
	}
	return id, errors.Wrap(err, "UserUsecase.GetUserByID pgx.QueryRow().Scan:")
}

func (u *UserUsecase) GetUserByID(ctx context.Context, userRequest models.UserRequest) (string, error) {
	userDTO := userRequest.ToDTO()
	login, err := u.repository.GetUserByID(ctx, userDTO)
	if err != nil {
		return login, errors.Wrap(err, "UserUsecase.GetUserByID pgx.QueryRow().Scan:")
	}
	return login, errors.Wrap(err, "UserUsecase.GetUserByID")
}

func (u *UserUsecase) AddingUser(ctx context.Context, userRequest models.UserRequest) error {
	userDTO := userRequest.ToDTO()
	err := u.repository.AddingUser(ctx, userDTO)
	return errors.Wrap(err, "UserUsecase.AddingUser pgx.Exec:")
}

func (u *UserUsecase) UpdateUser(ctx context.Context, userRequest models.UserRequest) error {
	userDTO := userRequest.ToDTO()
	err := u.repository.UpdateUser(ctx, userDTO)
	return errors.Wrap(err, "UserUsecase.UpdateUser pgx.Exec:")
}

func (u *UserUsecase) UploadFile(ctx context.Context, file *multipart.FileHeader) func(context.Context) error {
	return func(ctx context.Context) error {
		dstFilePath := filepath.Join(models.UploadsDir, file.Filename)
		err := uploading.SaveUploadedFile(file, dstFilePath)
		if err != nil {
			return errors.Wrap(err, "UserUsecase.UploadFile utils.SaveUploadedFile:")
		}
		return nil
	}
}
