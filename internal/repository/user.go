package repository

import (
	"PersonalAccountAPI/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db *pgx.Conn
}

func New(conn *pgx.Conn) *UserRepository {
	return &UserRepository{
		db: conn,
	}
}

func (ur *UserRepository) GetIDByLogin(ctx context.Context, user models.UserDTO) (int, error) {
	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	err := ur.db.QueryRow(ctx, query, user.Login, user.Password).Scan(&user.ID)
	if err != nil {
		return user.ID, errors.Wrap(err, "UserUsecase.GetIDByLogin pgx.QueryRow().Scan")
	}
	return user.ID, errors.Wrap(err, "UserUsecase.GetIDByLogin")
}

func (ur *UserRepository) GetUserByID(ctx context.Context, user models.UserDTO) (string, error) {
	query := `SELECT login FROM users WHERE id = $1`
	err := ur.db.QueryRow(ctx, query, user.ID).Scan(&user.Login)
	if err != nil {
		return user.Login, errors.Wrap(err, "UserUsecase.GetUserByID pgx.QueryRow().Scan:")
	}
	return user.Login, errors.Wrap(err, "UserUsecase.GetUserByID")
}

func (ur *UserRepository) AddingUser(ctx context.Context, user models.UserDTO) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`
	_, err := ur.db.Exec(ctx, query, user.ID, user.Login, user.Password)
	return errors.Wrap(err, "UserUsecase.AddingUser pgx.Exec:")
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user models.UserDTO) error {
	query := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := ur.db.Exec(ctx, query, user.ID, user.Login, user.Password)
	return errors.Wrap(err, "UserUsecase.UpdateUser pgx.Exec:")
}
