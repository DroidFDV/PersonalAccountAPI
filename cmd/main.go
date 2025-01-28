package main

import (
	"PersonalAccountAPI/database"
	"PersonalAccountAPI/internal/handler"
	"PersonalAccountAPI/internal/storage"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewRouter(handler *handler.Handle) *gin.Engine {
	router := gin.Default()
	router.POST("/login", handler.Login)
	router.GET("/user/:id", handler.GetUserByID)
	router.POST("/user", handler.AddUser)
	router.PUT("/user", handler.UpdateUser)
	return router
}

func main() {

	conn, err := storage.NewConn("postgres://postgres:postgres@postgres:5432/postgres")
	if err != nil {
		log.Fatal(errors.Wrap(err, "main db.NewConn"))
	}
	defer conn.Close(context.Background())

	if err := database.Migrate("postgres://postgres:postgres@postgres:5432/postgres"); err != nil {
		log.Fatal(errors.Wrap(err, "main db.NewConn"))
	}

	userHandle := handler.NewHandle(conn)

	router := NewRouter(userHandle)

	router.Run(":8080")

}
