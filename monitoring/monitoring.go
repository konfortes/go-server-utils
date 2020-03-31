package monitoring

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var (
	// RequestDuration is a histogram metric to report request duration
	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "request_duration_seconds",
		Help: "Handlers request duration in seconds",
	}, []string{"path", "code", "service"})
)

// Instrument instruments a gin app with prometheus monitoring
func Instrument(router *gin.Engine, appName string) {
	sanitizedServiceName := strings.ReplaceAll(appName, "-", "_")

	// http localhost:3000/metrics
	p := ginprometheus.NewPrometheus(sanitizedServiceName)
	p.Use(router)

	router.Use(durationMiddleware(sanitizedServiceName))
}

func durationMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		before := time.Now()
		c.Next()
		RequestDuration.With(prometheus.Labels{"service": serviceName, "path": c.Request.URL.Path, "code": strconv.Itoa(c.Writer.Status())}).Observe(time.Since(before).Seconds())
	}
}
