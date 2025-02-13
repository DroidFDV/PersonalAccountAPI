package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func GetConnect(connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, errors.Wrap(err, "GetConnect pgx.Connect")
	}
	return conn, err
}
