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
	db *pgx.Conn
}

func New(conn *pgx.Conn) *UserUsecase {
	return &UserUsecase{
		db: conn,
	}
}

func (u *UserUsecase) GetIDByLoginFromDB(ctx context.Context, userRequest models.UserRequest) (int, error) {
	userDTO := userRequest.ToDTO()
	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	err := u.db.QueryRow(ctx, query, userDTO.Login, userDTO.Password).Scan(&userDTO.ID)
	// warning: бессмысленная проверка, т.к я не знаю как ее обработать
	// if errors.Is(err, pgx.ErrNoRows) {
	// 	return 0, errors.Wrap(err, "UserUsecase.GetIDByLoginFromDB pgx.QueryRow().Scan")
	// }
	if err != nil {
		return userDTO.ID, errors.Wrap(err, "UserUsecase.GetIDByLoginFromDB pgx.QueryRow().Scan")
	}
	return userDTO.ID, errors.Wrap(err, "UserUsecase.GetIDByLoginFromDB")
}

func (u *UserUsecase) GetUserByIDFromDB(ctx context.Context, userRequest models.UserRequest) (string, error) {
	userDTO := userRequest.ToDTO()
	query := `SELECT login FROM users WHERE id = $1`
	err := u.db.QueryRow(ctx, query, userDTO.ID).Scan(&userDTO.Login)
	// if errors.Is(err, pgx.ErrNoRows) {
	// 	return userDTO.Login, errors.Wrap(err, "UserUsecase.GetIDByLoginFromDB pgx.QueryRow().Scan")
	// }
	if err != nil {
		return userDTO.Login, errors.Wrap(err, "UserUsecase.GetUserByIDFromDB pgx.QueryRow().Scan:")
	}
	return userDTO.Login, errors.Wrap(err, "UserUsecase.GetUserByIDFromDB")
}

func (u *UserUsecase) AddingUserToDB(ctx context.Context, userRequest models.UserRequest) error {
	userDTO := userRequest.ToDTO()
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`
	_, err := u.db.Exec(ctx, query, userDTO.ID, userDTO.Login, userDTO.Password)
	return errors.Wrap(err, "UserUsecase.AddingUserToDB pgx.Exec:")
}

func (u *UserUsecase) UpdateUserInDB(ctx context.Context, userRequest models.UserRequest) error {
	userDTO := userRequest.ToDTO()
	query := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := u.db.Exec(ctx, query, userDTO.ID, userDTO.Login, userDTO.Password)
	return errors.Wrap(err, "UserUsecase.UpdateUserInDB pgx.Exec:")
}

func (u *UserUsecase) UploadFile(ctx context.Context, id string, file *multipart.FileHeader) func(context.Context) error {
	return func(ctx context.Context) error {
		dstPath := filepath.Join(models.UploadsDir, id)
		err := utils.SaveUploadedFile(file, dstPath)
		if err != nil {
			return errors.Wrap(err, "UserUsecase.UploadFile utils.SaveUploadedFile:")
		}
		return nil
	}
}
