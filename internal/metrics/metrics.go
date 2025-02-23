package metrics

import (
	"PersonalAccountAPI/internal/cache"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
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

func InitMetrics() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(goRoutineCount)
	prometheus.MustRegister(cacheSize)
}

func UpdateMetrics(c *cache.CacheDecorator) {
	for {
		goRoutineCount.Set(float64(runtime.NumGoroutine()))
		cacheSize.Set(float64(c.GetCacheSize()))
		time.Sleep(1 * time.Second)
	}
}

func RequestReceived(c *gin.Context) {
	c.Next()
	status := c.Writer.Status()
	method := c.Request.Method
	httpRequests.WithLabelValues(method, http.StatusText(status)).Inc()
}

// warning: пока не понятно зачем
// func PrometheusMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Next()
// 		status := c.Writer.Status()
// 		method := c.Request.Method
// 		httpRequests.WithLabelValues(method, http.StatusText(status)).Inc()
// 	}
// }
