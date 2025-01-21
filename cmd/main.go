package main

import (
	"PersonalAccountAPI/internal/handler"
	db "PersonalAccountAPI/internal/storage"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewRouter(u *handler.UserHandle) *gin.Engine { return gin.Default() }

func main() {

	conn, err := db.NewConn("postgres://postgres:postgres@postgres:5432/postgres")
	if err != nil {
		log.Fatal(errors.Wrap(err, "main db.NewConn"))
	}
	defer conn.Close(context.Background())
	if err := db.CreateIfNotExistsUsers(conn); err != nil {
		log.Fatal(errors.Wrap(err, "main db.ConnectToUsers"))
	}

	userHandle := handler.NewUser(conn)

	router := NewRouter(userHandle)

	router.POST("/login", userHandle.Login)
	router.GET("/user/:id", userHandle.UserByID)
	router.POST("/user", userHandle.AddUser)
	router.PUT("/user", userHandle.UpdateUser)

	router.Run(":8080")

}
