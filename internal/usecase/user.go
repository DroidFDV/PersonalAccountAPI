package usecase

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/repository"
	"context"

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

func (u *UserUsecase) GetIDByLogin(ctx context.Context, userRequest models.UserRequest) (models.UserResponse, error) {
	user, err := u.repository.GetIDByLogin(ctx, userRequest.ToDTO())
	if err != nil {
		return models.UserResponse{}, errors.Wrap(err, "UserUsecase.GetUserByID pgx.QueryRow().Scan:")
	}
	return user.ToResponce(), nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, userRequest models.UserRequest) (models.UserResponse, error) {
	user, err := u.repository.GetUserByID(ctx, userRequest.ToDTO())
	if err != nil {
		return models.UserResponse{}, errors.Wrap(err, "UserUsecase.GetUserByID pgx.QueryRow().Scan:")
	}
	return user.ToResponce(), nil
}

func (u *UserUsecase) AddUser(ctx context.Context, userRequest models.UserRequest) error {
	if err := u.repository.AddUser(ctx, userRequest.ToDTO()); err != nil {
		return errors.Wrap(err, "UserUsecase.AddingUser pgx.Exec:")
	}
	return nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, userRequest models.UserRequest) error {
	if err := u.repository.UpdateUser(ctx, userRequest.ToDTO()); err != nil {
		return errors.Wrap(err, "UserUsecase.UpdateUser pgx.Exec:")
	}
	return nil
}
