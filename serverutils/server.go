package serverutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var (
	// ShutdownHooks is a hooks slice that can be appended to execute code on shutdown
	ShutdownHooks []func()
)

// SetMiddlewares ... sets jaeger and request time middlewares
func SetMiddlewares(router *gin.Engine, tracer *opentracing.Tracer) {
	router.Use(durationMiddleware())

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

	// http localhost:3000/metrics
	p := ginprometheus.NewPrometheus(serviceName)
	p.Use(router)
}
