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
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Метрики
var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status"},
	)

	goRoutineCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_routines_count",
			Help: "Number of active Go routines",
		},
	)

	cacheSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_size",
			Help: "Number of users in cache",
		},
	)
)

func updateMetrics(c *cache.CacheDecorator) {
	for {
		goRoutineCount.Set(float64(runtime.NumGoroutine()))
		cacheSize.Set(float64(c.GetCacheSize()))
		time.Sleep(5 * time.Second)
	}
}

func requestReceived(c *gin.Context) {
	c.Next()
	status := c.Writer.Status()
	httpRequests.WithLabelValues(c.Request.Method, http.StatusText(status)).Inc()
}

func init() {
	// Регистрируем метрики
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(goRoutineCount)
	prometheus.MustRegister(cacheSize)
}

func newRouter(handler *handler.Handle) *gin.Engine {
	router := gin.Default()
	router.POST("/user", handler.AddUser)
	router.POST("/login", handler.Login)
	router.GET("/user/:id", handler.GetUserByID)
	router.PUT("/user", handler.UpdateUser)
	router.POST("/file/upload", handler.UploadFile)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(requestReceived)

	return router
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	conn, err := storage.GetConnect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "main db.NewConn:"))
	}
	defer conn.Close(context.Background())

	if err := database.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatal(errors.Wrap(err, "main db.NewConn:"))
	}

	if err := os.MkdirAll(models.UploadsDir, os.ModePerm); err != nil {
		log.Fatal(errors.Wrap(err, "main os.MkdirAll: failed to create uploads directory"))
	}

	workers.Run(10)

	time.Sleep(time.Duration(2) * time.Second)

	userProvider := usecase.New(conn)
	cacheProvider := cache.New(userProvider)
	workerManager := workers.Run(10)

	handle := handler.New(cacheProvider, workerManager)
	// handle := handler.New(userProvider)

	router := newRouter(handle)
	// Запускаем обновление метрик в фоне
	go updateMetrics(cacheProvider)

	router.Run(cfg.Port)

}
