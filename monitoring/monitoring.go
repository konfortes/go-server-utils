package monitoring

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestDuration is a histogram metric to report request duration
	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "request_duration_seconds",
		Help: "Handlers request duration in seconds",
	}, []string{"path", "status", "service"})
)

// Instrument a gin app with prometheus monitoring
func Instrument(router *gin.Engine, appName string) {
	sanitizedServiceName := strings.ReplaceAll(appName, "-", "_")

	// http localhost:3000/metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(requestMiddleware(sanitizedServiceName))
}

func requestMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		before := time.Now()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsedSeconds := float64(time.Since(before)) / float64(time.Second)
		labels := prometheus.Labels{"service": serviceName, "path": c.Request.URL.Path, "status": status}
		RequestDuration.With(labels).Observe(elapsedSeconds)
	}
}
