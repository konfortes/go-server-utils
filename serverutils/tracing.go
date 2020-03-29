package serverutils

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

// InitJaeger inits and returns a Jaeger tracer and closer.
func InitJaeger(service string) *opentracing.Tracer {
	cfg := &jaegerConfig.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	var err error
	tracer, tracerCloser, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	ShutdownHooks = append(ShutdownHooks, func() {
		tracerCloser.Close()
	})

	opentracing.SetGlobalTracer(tracer)

	return &tracer
}

// JaegerMiddleware extracts span (if exist) from request headers and set it in the context
func jaegerMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
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
			log.Printf("error extracting span from request: %s", err)
		}

		c.Next()
	}
}
