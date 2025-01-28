package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserUsecase struct {
	db *pgx.Conn
}

func NewUser(conn *pgx.Conn) *UserUsecase {
	return &UserUsecase{
		db: conn,
	}
}

func (user *UserUsecase) GetIDByLoginFromDB(ctx context.Context, login, password string) (int, error) {
	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	rows, err := user.db.Query(ctx, query, login, password)
	if err != nil {
		return 0, errors.Wrap(err, "POST /login Query")
	}
	defer rows.Close()

	var id int
	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, errors.Wrap(err, "POST /login Scan")
		}
		return id, errors.Wrap(err, "POST /login")
	}
	return id, errors.Wrap(err, "")
}

func (user *UserUsecase) GetUserByIDFromDB(ctx context.Context, id int) (string, error) {

	query := `SELECT login FROM users WHERE id = $1`
	rows, err := user.db.Query(ctx, query, id)
	if err != nil {
		return "", errors.Wrap(err, "GET /user/:id Query")
	}
	defer rows.Close()

	var login string
	if rows.Next() {
		if err := rows.Scan(&login); err != nil {
			return "", errors.Wrap(err, "GET /user/:id Scan")
		}
		return login, errors.Wrap(err, "GET /user/:id")
	}
	return login, errors.Wrap(err, "")
}

func (user *UserUsecase) AddingUserToDB(ctx context.Context, id int, login, password string) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`
	_, err := user.db.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "POST /user Exec")
}

func (user *UserUsecase) UpdateUserInDB(ctx context.Context, id int, login, password string) error {
	query := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := user.db.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "POST /user Exec")
}
