package main

import (
	"PersonalAccountAPI/database"
	"PersonalAccountAPI/internal/cache"
	"PersonalAccountAPI/internal/handler"
	"PersonalAccountAPI/internal/metrics"
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/storage"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newRouter(handler *handler.Handle) *gin.Engine {
	router := gin.Default()
	router.POST("/user", handler.AddUser)
	router.POST("/login", handler.Login)
	router.GET("/user/:id", handler.GetUserByID)
	router.PUT("/user", handler.UpdateUser)
	router.POST("/file/upload", handler.UploadFile)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(metrics.RequestReceived)

	return router
}

// func newMetricsRouter() *gin.Engine {
// 	metricsRouter := gin.Default()
// 	metricsRouter.Use(metrics.PrometheusMiddleware())
// 	return metricsRouter
// }

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "main storage.GetConnect: Failed to load config"))
	}
	models.UploadsDir = cfg.FileStoragePatg

	conn, err := storage.GetConnect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "main storage.GetConnect"))
	}
	defer conn.Close(context.Background())

	if err := database.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatal(errors.Wrap(err, "main database.Migrate"))
	}

	workerManager := workers.Run(10)

	userProvider := usecase.New(conn)
	cacheProvider := cache.New(userProvider)

	handle := handler.New(cacheProvider, workerManager)
	// handle := handler.New(userProvider)

	router := newRouter(handle)
	metrics.InitMetrics()
	go metrics.UpdateMetrics(cacheProvider)

	router.Run(cfg.Port)

}
