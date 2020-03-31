package tracing

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

var (
	tracerCloser io.Closer
)

// CloseFunc is a tracer close function that should be executed when the app shutsdown
type CloseFunc func()

// Instrument instruments a gin app with Jaeger tracing
func Instrument(router *gin.Engine, appName string) CloseFunc {
	tracer, closer := initJaeger(appName)
	opentracing.SetGlobalTracer(tracer)
	router.Use(JaegerMiddleware(tracer))

	return func() {
		closer.Close()
	}
}

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &jaegerConfig.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, tracerCloser, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, tracerCloser
}

// jaegerMiddleware extracts span (if exist) from request headers and set it in the context
func JaegerMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		upstreamCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)

		if err == nil {
			span := tracer.StartSpan(c.Request.URL.Path, ext.RPCServerOption(upstreamCtx))
			defer span.Finish()

			ctx := context.Background()
			ctx = opentracing.ContextWithSpan(ctx, span)
			c.Request = c.Request.Clone(ctx)
		} else {
			log.Printf("could not extract span from request: %s", err)
		}

		c.Next()
	}
}
