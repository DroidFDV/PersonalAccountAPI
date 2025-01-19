package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type user struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserHandle struct {
	db *pgx.Conn
}

func (u *UserHandle) getIDByLoginFromDB(ctx context.Context, login, password string) (int, error) {
	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	rows, err := u.db.Query(ctx, query, login, password)
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

func (u *UserHandle) getUserByIDFromDB(ctx context.Context, id int) (string, error) {

	query := `SELECT login FROM users WHERE id = $1`
	rows, err := u.db.Query(ctx, query, id)
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

func (u *UserHandle) addingUserToDB(ctx context.Context, id int, login, password string) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`
	_, err := u.db.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "POST /user Exec")
}

func (u *UserHandle) updateUserInDB(ctx context.Context, id int, login, password string) error {
	query := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := u.db.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "POST /user Exec")
}

func (u *UserHandle) Login(c *gin.Context) {
	var user user
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "POST /login ShouldBind").Error())
		return
	}

	id, err := (*UserHandle).getIDByLoginFromDB(u, c, user.Login, user.Password)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}
	if id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (u *UserHandle) UserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error(errors.Wrap(err, "GET /user/:id Atoi").Error())
		return
	}

	login, err := (*UserHandle).getUserByIDFromDB(u, c, id)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect query"})
		return
	}
	if login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": login})
}

func (u *UserHandle) AddUser(c *gin.Context) {
	var user user

	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "POST /login ShouldBind").Error())
		return
	}

	if err := (*UserHandle).addingUserToDB(u, c, user.Id, user.Login, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect user data"})
		slog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": "added"})
}

func (u *UserHandle) UpdateUser(c *gin.Context) {
	var user user

	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "PUT /login ShouldBind").Error())
		return
	}

	if err := (*UserHandle).updateUserInDB(u, c, user.Id, user.Login, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"user by id: " + strconv.FormatInt(int64(user.Id), 10): "updated"})
}

func NewUser(conn *pgx.Conn) *UserHandle {
	return &UserHandle{
		db: conn,
	}
}
