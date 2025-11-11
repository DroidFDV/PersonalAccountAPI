package metrics

import (
	"PersonalAccountAPI/internal/cache"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	//NOTE: что лучше для этого выбрать
	httpRequests = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_requests_total",
			Help:    "Total number of HTTP requests with duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status", "path"},
	)

	cacheSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_size",
			Help: "Number of users in cache",
		},
	)
)

func MetricsRegistration() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(cacheSize)
}

func InitMetrics(port string) {
	MetricsRegistration()
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		slog.Info("Starting metrics server", slog.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("InitMetrics server.ListenAndServe", slog.Any("error", err))
		}
	}()
}

func UpdateMetrics(ctx context.Context, c *cache.CacheDecorator) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cacheSize.Set(float64(c.GetCacheSize()))
		//NOTE: think again
		case <-ctx.Done():
			return
		}
	}
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		duration := time.Since(start).Seconds()

		httpRequests.WithLabelValues(method, http.StatusText(status), path).Observe(duration)
	}
}
