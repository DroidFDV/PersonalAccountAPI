package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func NewConn(connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, errors.Wrap(err, "pgx.Connect")
	}
	return conn, err
}

func CreateIfNotExistsUsers(conn *pgx.Conn) error {
	query := `
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				login VARCHAR(100) NOT NULL,
				password VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
			`
	_, err := conn.Exec(context.Background(), query)
	return errors.Wrap(err, "Ошибка выполнения запроса CREATE TABLE")
}
