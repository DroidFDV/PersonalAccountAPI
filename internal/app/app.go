package app

import (
	"PersonalAccountAPI/config"
	"PersonalAccountAPI/database"
	"PersonalAccountAPI/internal/cache"
	"PersonalAccountAPI/internal/handler"
	"PersonalAccountAPI/internal/metrics"
	"PersonalAccountAPI/internal/models"
	"PersonalAccountAPI/internal/repository"
	"PersonalAccountAPI/internal/storage"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"context"
	"time"

	"github.com/go-faster/errors"
)

func Run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.Wrap(err, "main storage.GetConnect: Failed to load config")
	}
	models.UploadsDir = cfg.Storage.Path

	conn, err := storage.GetConnect(cfg.GetDSN())
	if err != nil {
		return errors.Wrap(err, "main storage.GetConnect")
	}
	defer conn.Close(context.Background())

	if err := database.Migrate(cfg.GetDSN()); err != nil {
		return errors.Wrap(err, "main database.Migrate")
	}

	workerManager := workers.Run(cfg.GetWorkersNum(), cfg.GetWorkersQueueLen())

	userRepository := repository.New(conn)
	userProvider := usecase.New(userRepository)
	cacheProvider := cache.New(userProvider, cfg.GetCacheTTL()*time.Second)
	go cacheProvider.RunCleaner(cfg.GetCacheInterval() * time.Second)

	handle := handler.New(cacheProvider, workerManager)

	router := NewRouter(handle)
	metrics.InitMetrics(cfg.Metrics.Port)
	go metrics.UpdateMetrics(cacheProvider)

	return router.Run(":" + cfg.App.Port)
}
