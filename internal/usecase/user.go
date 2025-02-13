package usecase

import (
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/utils"
	"context"
	"mime/multipart"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type UserUsecase struct {
	db         *pgx.Conn
	fileHeader *multipart.FileHeader
}

func New(conn *pgx.Conn) *UserUsecase {
	return &UserUsecase{
		db:         conn,
		fileHeader: nil,
	}
}

func (u *UserUsecase) GetIDByLoginFromDB(ctx context.Context, login, password string) (int, error) {
	var id int
	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	err := u.db.QueryRow(ctx, query, login, password).Scan(&id)
	if err != nil {
		return id, errors.Wrap(err, "UserUsecase.GetIDByLoginFromDB pgx.QueryRow().Scan:")
	}
	return id, errors.Wrap(err, "")
}

func (u *UserUsecase) GetUserByIDFromDB(ctx context.Context, id int) (string, error) {
	var login string
	query := `SELECT login FROM users WHERE id = $1`
	err := u.db.QueryRow(ctx, query, id).Scan(&login)
	if err != nil {
		return login, errors.Wrap(err, "UserUsecase.GetUserByIDFromDB pgx.QueryRow().Scan:")
	}
	return login, errors.Wrap(err, "")
}

func (u *UserUsecase) AddingUserToDB(ctx context.Context, id int, login, password string) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`
	_, err := u.db.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "UserUsecase.AddingUserToDB pgx.Exec:")
}

func (u *UserUsecase) UpdateUserInDB(ctx context.Context, id int, login, password string) error {
	query := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := u.db.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "UserUsecase.UpdateUserInDB pgx.Exec:")
}

func (u *UserUsecase) SetFile(fileHeader *multipart.FileHeader) {
	u.fileHeader = fileHeader
}

func (u *UserUsecase) UploadFile(ctx context.Context) error {
	err := utils.SaveUploadedFile(u.fileHeader, filepath.Join(models.UploadsDir, u.fileHeader.Filename))
	u.fileHeader = nil
	if err != nil {
		return errors.Wrap(err, "UserUsecase.UploadFile utils.SaveUploadedFile:")
	}
	return nil
}
