package serverutils

import (
	"strconv"
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
	}, []string{"path", "code"})
)

// SetMonitoringHandler sets a /monitoring route
func SetMonitoringHandler(router *gin.Engine) {
	// http localhost:3000/metrics
	p := ginprometheus.NewPrometheus("service_name")
	p.Use(router)
}

func durationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		RequestDuration.With(prometheus.Labels{"path": c.Request.URL.Path, "code": strconv.Itoa(c.Writer.Status())}).Observe(time.Since(now).Seconds())
	}
}
