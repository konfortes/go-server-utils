package serverutils

import (
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
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
