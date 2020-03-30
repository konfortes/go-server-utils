package serverutils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var (
	// ShutdownHooks is a hooks slice that can be appended to execute code on shutdown
	ShutdownHooks []func()
)

// SetMiddlewares ... sets jaeger and request time middlewares
func SetMiddlewares(router *gin.Engine, tracer *opentracing.Tracer, serviceName string) {
	router.Use(durationMiddleware(serviceName))

	if tracer != nil {
		router.Use(jaegerMiddleware(*tracer))
	}
}

// SetRoutes sets /metrics and /health routes
func SetRoutes(router *gin.Engine, serviceName string) {
	// http localhost:8080/health
	router.GET("/health", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte("OK"))
	})
	sanitizedServiceName := strings.ReplaceAll(serviceName, "-", "_")
	// http localhost:3000/metrics
	p := ginprometheus.NewPrometheus(sanitizedServiceName)
	p.Use(router)
}
