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
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status", "path"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
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
	prometheus.MustRegister(httpRequestDuration)
}

func InitMetrics(port string) {
	MetricsRegistration()
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		slog.Info("Starting metrics server", slog.String("port", port))
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			slog.Error("InitMetrics http.ListenAndServe", slog.Any("error", err))
		}
	}()
}

func UpdateMetrics(c *cache.CacheDecorator) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cacheSize.Set(float64(c.GetCacheSize()))
		//WARNING: not work
		case <-context.Background().Done():
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

		httpRequests.WithLabelValues(method, http.StatusText(status), path).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
