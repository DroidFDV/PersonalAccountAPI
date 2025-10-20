package app

import (
	"PersonalAccountAPI/config"
	"PersonalAccountAPI/database"
	"PersonalAccountAPI/internal/cache"
	"PersonalAccountAPI/internal/handler"
	"PersonalAccountAPI/internal/metrics"
	"PersonalAccountAPI/internal/repository"
	"PersonalAccountAPI/internal/storage"
	"PersonalAccountAPI/internal/uploading"
	"PersonalAccountAPI/internal/usecase"
	"PersonalAccountAPI/internal/workers"
	"context"
	"time"

	"github.com/go-faster/errors"
)

func Run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.Wrap(err, "Run config.LoadConfig: Failed to load config")
	}

	conn, err := storage.GetConnectDB(cfg.GetDSN())
	if err != nil {
		return errors.Wrap(err, "Run storage.GetConnectDB")
	}
	defer conn.Close(context.Background())

	if err := database.Migrate(cfg.GetDSN()); err != nil {
		return errors.Wrap(err, "Run database.Migrate")
	}

	s3Client, err := storage.InitS3Client(cfg.GetS3Config())
	if err != nil {
		return errors.Wrap(err, "Run storage.InitS3Client")
	}

	userRepository := repository.New(conn)
	userProvider := usecase.New(userRepository)
	cacheProvider := cache.New(userProvider, cfg.GetCacheTTL()*time.Second)
	go cacheProvider.RunCleaner(context.Background(), cfg.GetCacheInterval()*time.Second)
	workerManager := workers.Run(cfg.GetWorkersNum(), cfg.GetWorkersQueueLen())
	uploadProvider := uploading.New(s3Client, cfg.GetS3BucketName())

	handle := handler.New(cacheProvider, workerManager, uploadProvider)

	router := NewRouter(handle)
	metrics.InitMetrics(cfg.Metrics.Port)
	go metrics.UpdateMetrics(context.Background(), cacheProvider)

	return router.Run(":" + cfg.App.Port)
}
