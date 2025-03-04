package main

import (
	"PersonalAccountAPI/database"
	"PersonalAccountAPI/internal/app"
	"PersonalAccountAPI/internal/cache"
	"PersonalAccountAPI/internal/handler"
	"PersonalAccountAPI/internal/metrics"
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/repository"
	"PersonalAccountAPI/internal/storage"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
)

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

	UserRepository := repository.New(conn)
	userProvider := usecase.New(UserRepository)
	ttl, checkInterval := 5*time.Second, 10*time.Second
	cacheProvider := cache.New(userProvider, ttl, checkInterval)

	handle := handler.New(cacheProvider, workerManager)

	router := app.NewRouter(handle)
	metrics.InitMetrics(cfg.MetricsPort)
	go metrics.UpdateMetrics(cacheProvider)

	router.Run(":" + cfg.Port)

}
