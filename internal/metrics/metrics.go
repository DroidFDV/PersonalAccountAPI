package metrics

import (
	"PersonalAccountAPI/internal/cache"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

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

func InitMetrics(port string) {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(goRoutineCount)
	prometheus.MustRegister(cacheSize)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		logrus.Printf("Starting metrics server on port %s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			logrus.Fatalf("Failed to start metrics server: %v", err)
		}
	}()
}

func UpdateMetrics(c *cache.CacheDecorator) {
	for {
		goRoutineCount.Set(float64(runtime.NumGoroutine()))
		cacheSize.Set(float64(c.GetCacheSize()))
		time.Sleep(1 * time.Second)
	}
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		method := c.Request.Method
		httpRequests.WithLabelValues(method, http.StatusText(status)).Inc()
	}
}
