package database

import (
	"database/sql"
	"embed"
	"log/slog"

	"github.com/go-faster/errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(connString string) error {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return errors.Wrap(err, "Migrate sql.Open")
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return errors.Wrap(err, "Migrate db.Ping")
	}

	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "Migrate goose.SetDialect")
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		return errors.Wrap(err, "Migrate goose.GetDBVersion")
	}

	if err = goose.Up(db, "migrations"); err != nil {
		if err := goose.DownTo(db, "migrations", version); err != nil {
			slog.Error(
				"Migrate goose.DownTo: cannot rollback migrations",
				slog.Any("error", err),
				slog.Any("try rollback to version", version),
			)
		}

		return errors.Wrap(err, "Migrate goose.Up")
	}

	return nil
}
