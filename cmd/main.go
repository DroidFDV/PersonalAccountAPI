package main

import (
	"PersonalAccountAPI/internal/handler"
	db "PersonalAccountAPI/internal/storage"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ./internal/storage
// делаем функцию func NewConn(connString string) (*pgx.Conn, error)

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

func NewRouter(u *handler.UserHandle) *gin.Engine { return gin.Default() }

func main() {

	conn, err := db.NewConn("postgres://postgres:postgres@postgres:5432/postgres")
	if err != nil {
		log.Fatal(errors.Wrap(err, "main pgx.Connect"))
	}
	defer conn.Close(context.Background())

	userHandle := handler.NewUser(conn)

	router := NewRouter(userHandle)

	router.POST("/login", userHandle.Login)
	router.GET("/user/:id", userHandle.UserByID)
	router.POST("/user", userHandle.AddUser)
	router.PUT("/user", userHandle.UpdateUser)

	router.Run(":8080")

}
