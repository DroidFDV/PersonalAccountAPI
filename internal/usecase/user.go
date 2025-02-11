package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type UserUsecase struct {
	db *pgx.Conn
}

func New(conn *pgx.Conn) *UserUsecase {
	return &UserUsecase{
		db: conn,
	}
}

func (user *UserUsecase) GetIDByLoginFromDB(ctx context.Context, login, password string) (int, error) {
	var id int
	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	err := user.db.QueryRow(ctx, query, login, password).Scan(&id)
	if err != nil {
		return id, errors.Wrap(err, "GetIDByLoginFromDB QueryRow")
	}
	return id, errors.Wrap(err, "")
}

func (user *UserUsecase) GetUserByIDFromDB(ctx context.Context, id int) (string, error) {
	var login string
	query := `SELECT login FROM users WHERE id = $1`
	err := user.db.QueryRow(ctx, query, id).Scan(&login)
	if err != nil {
		return login, errors.Wrap(err, "GetUserByIDFromDB QueryRow")
	}
	return login, errors.Wrap(err, "")
}

func (user *UserUsecase) AddingUserToDB(ctx context.Context, id int, login, password string) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`
	_, err := user.db.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "AddingUserToDB Exec")
}

func (user *UserUsecase) UpdateUserInDB(ctx context.Context, id int, login, password string) error {
	query := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := user.db.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "UpdateUserInDB Exec")
}
