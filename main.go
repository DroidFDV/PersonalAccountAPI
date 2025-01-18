package main

import (
	"context"
	"log"
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

// ./internal/storage
// делаем функцию func NewConn(connString string) (*pgx.Conn, error)
var conn *pgx.Conn

// ./internal/handler/user.go
//
//	type UserHandle struct {
//		db *pgx.Conn
//	}
//
//	func New(db *pgx.Conn) *UserHandle {
//		...
//	}
//
// func (u *UserHandle) getIdByLoginFromDb(...
//
//	rows, err := u.db.Query(...
//
// func (u *UserHandle) Login(...

func getIDByLoginFromDB(ctx context.Context, login, password string) (int, error) {
	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	rows, err := conn.Query(ctx, query, login, password)
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

func getUserByIDFromDB(ctx context.Context, id int) (string, error) {

	query := `SELECT login FROM users WHERE id = $1`
	rows, err := conn.Query(ctx, query, id)
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

func addingUserToDB(ctx context.Context, id int, login, password string) error {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`
	_, err := conn.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "POST /user Exec")
}

func updateUserInDB(ctx context.Context, id int, login, password string) error {
	query := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := conn.Exec(ctx, query, id, login, password)
	return errors.Wrap(err, "POST /user Exec")
}

func loginHandler(c *gin.Context) {
	var user user
	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "POST /login ShouldBind").Error())
		return
	}

	id, err := getIDByLoginFromDB(c, user.Login, user.Password)
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

func userByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error(errors.Wrap(err, "GET /user/:id Atoi").Error())
		return
	}

	login, err := getUserByIDFromDB(c, id)
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

func addingHandler(c *gin.Context) {
	var user user

	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "POST /login ShouldBind").Error())
		return
	}

	if err := addingUserToDB(c, user.Id, user.Login, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect user data"})
		slog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": "added"})
}

func updateHandler(c *gin.Context) {
	var user user

	if err := c.ShouldBind(&user); err != nil {
		slog.Error(errors.Wrap(err, "PUT /login ShouldBind").Error())
		return
	}

	if err := updateUserInDB(c, user.Id, user.Login, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		slog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"user by id: " + strconv.FormatInt(int64(user.Id), 10): "updated"})
}

func main() {
	var err error
	// конект к базе
	ctx := context.Background()
	// заменить на conn, err := storage.NewConn(ctx, connString string) (*pgx.Conn, error)
	conn, err = pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatal(errors.Wrap(err, "main pgx.Connect"))
	}
	defer conn.Close(ctx)

	query := `
			DO $$
			BEGIN
			IF NOT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_name = 'users'
			) THEN
				CREATE TABLE users (
					id SERIAL PRIMARY KEY,
					login VARCHAR(100) NOT NULL,
					password VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			END IF;
			END
			$$;
			`
	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		log.Panic(errors.Wrap(err, "Ошибка выполнения запроса CREATE TABLE"))
		return
	}

	// userHandle := handler.New(conn)
	//
	// ./cmd/main.go
	// func NewRouter(userHandle) *gin.Engine
	router := gin.Default()
	// router.POST("/login", handler.Login)
	router.POST("/login", loginHandler)
	router.GET("/user/:id", userByIDHandler)
	router.POST("/user", addingHandler)
	router.PUT("/user", updateHandler)
	// оставляй тут
	router.Run(":8080")

}
