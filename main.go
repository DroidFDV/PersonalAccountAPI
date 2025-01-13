package main

import (
	"context"
	"log"
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

// какая гадость этот ваш глобальный коннект
var conn *pgx.Conn

func getIdByLoginFromDb(cntxt *gin.Context, login, password string) {
	var id int

	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	rows, err := conn.Query(context.Background(), query, login, password)
	if err != nil {
		log.Panic(errors.Wrap(err, "POST /login Query"))
		rows.Close()
		return
	}

	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Panic(errors.Wrap(err, "POST /login Scan"))
			rows.Close()
			return
		}
		cntxt.JSON(http.StatusOK, gin.H{"Id": id})
		rows.Close()
		return
	}
	rows.Close()

	cntxt.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}

func getUserByIdFromDb(cntxt *gin.Context, id int) {
	var login string

	query := `SELECT login FROM users WHERE id = $1`
	rows, err := conn.Query(context.Background(), query, id)
	if err != nil {
		log.Panic(errors.Wrap(err, "GET /user/:id Query\n"))
		rows.Close()
		return
	}

	if rows.Next() {
		err := rows.Scan(&login)
		if err != nil {
			log.Panic(errors.Wrap(err, "POST /login Scan"))
			rows.Close()
			return
		}
		cntxt.JSON(http.StatusOK, gin.H{"user": login})
		rows.Close()
		return
	}
	rows.Close()

	cntxt.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}

func addingUserToDb(cntxt *gin.Context, id int, login, password string) {
	query := `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`
	// был вариант передавать структуру user
	// _, err := conn.Exec(context.Background(), query, user.Id, user.Login, user.Password)
	_, err := conn.Exec(context.Background(), query, id, login, password)
	if err != nil {
		log.Panic(errors.Wrap(err, "POST /user Exec\n"))
		return
	}

	cntxt.JSON(http.StatusOK, gin.H{"user": "added"})
}

// доделываю
func updateUserInDb(cntxt *gin.Context, user user) {
	query := `UPDATE users SET login = $2, password = $3 WHERE id = $1`
	_, err := conn.Exec(context.Background(), query, user.Id, user.Login, user.Password)
	if err != nil {
		log.Panic(errors.Wrap(err, "POST /user Exec\n"))
		return
	}

	cntxt.JSON(http.StatusOK, gin.H{"user " + strconv.FormatInt(int64(user.Id), 10): "updated"})
}

func loginHandler(cntxt *gin.Context) {
	var user user

	if err := cntxt.ShouldBind(&user); err != nil {
		log.Panic(errors.Wrap(err, "POST /login ShouldBind\n"))
		return
	}

	getIdByLoginFromDb(cntxt, user.Login, user.Password)
}

func userByIdHandler(cntxt *gin.Context) {
	idParam := cntxt.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Panic(errors.Wrap(err, "GET /user Atoi\n"))
		return
	}

	getUserByIdFromDb(cntxt, id)
}

func addingHandler(cntxt *gin.Context) {
	var user user

	if err := cntxt.ShouldBind(&user); err != nil {
		log.Panic(errors.Wrap(err, "POST /login ShouldBind\n"))
		return
	}

	addingUserToDb(cntxt, user.Id, user.Login, user.Password)
}

// доделываю
func updateHandler(cntxt *gin.Context) {
	var user user

	if err := cntxt.ShouldBind(&user); err != nil {
		log.Panic(errors.Wrap(err, "POST /login ShouldBind\n"))
		return
	}

	updateUserInDb(cntxt, user)
}

func main() {
	var err error
	// конект к базе
	ctx := context.Background()
	conn, err = pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Panic(errors.Wrap(err, "main pgx.Connect\n"))
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
		log.Panic(errors.Wrap(err, "Ошибка выполнения запроса CREATE TABLE\n"))
		return
	}

	router := gin.Default()

	router.POST("/login", loginHandler)

	router.GET("/user/:id", userByIdHandler)

	router.POST("/user", addingHandler)

	router.PUT("/user", updateHandler)

	router.Run(":8080")
}
