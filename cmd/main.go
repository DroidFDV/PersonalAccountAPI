package main

import (
	"PersonalAccountAPI/database"
	"PersonalAccountAPI/internal/cache"
	"PersonalAccountAPI/internal/handler"
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/storage"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func newRouter(handler *handler.Handle) *gin.Engine {
	router := gin.Default()
	router.POST("/user", handler.AddUser)
	router.POST("/login", handler.Login)
	router.GET("/user/:id", handler.GetUserByID)
	router.PUT("/user", handler.UpdateUser)
	router.POST("/user/upload/:id", handler.UploadFile)
	return router
}

func main() {
	conn, err := storage.GetConnect("postgres://postgres:postgres@postgres:5432/postgres")
	if err != nil {
		log.Fatal(errors.Wrap(err, "main db.NewConn"))
	}
	defer conn.Close(context.Background())

	if err := database.Migrate("postgres://postgres:postgres@postgres:5432/postgres"); err != nil {
		log.Fatal(errors.Wrap(err, "main db.NewConn"))
	}

	if err := os.MkdirAll(models.UploadsDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	workManager := workers.New(10, 100)
	workManager.StartPool()
	for i := range 5 {
		workManager.AddJob(workers.Job{ID: i})
		time.Sleep(time.Duration(10) * time.Millisecond)
	}

	time.Sleep(time.Duration(2) * time.Second)

	userProvider := usecase.New(conn)
	cacheProvider := cache.New(userProvider)

	handle := handler.New(cacheProvider)
	// handle := handler.New(userProvider)

	router := newRouter(handle)

	router.Run(":8080")

}
