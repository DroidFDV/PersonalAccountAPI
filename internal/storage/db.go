package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// query := `
// 			DO $$
// 			BEGIN
// 			IF NOT EXISTS (
// 				SELECT FROM information_schema.tables
// 				WHERE table_name = 'users'
// 			) THEN
// 				CREATE TABLE users (
// 					id SERIAL PRIMARY KEY,
// 					login VARCHAR(100) NOT NULL,
// 					password VARCHAR(255) NOT NULL,
// 					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 				);
// 			END IF;
// 			END
// 			$$;
// 			`
// 	_, err = conn.Exec(context.Background(), query)
// 	if err != nil {
// 		log.Panic(errors.Wrap(err, "Ошибка выполнения запроса CREATE TABLE"))
// 		return
// 	}

func NewConn(connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, errors.Wrap(err, "pgx.Connect")
	}
	return conn, err
}
