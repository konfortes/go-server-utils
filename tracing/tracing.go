package tracing

import (
	"context"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
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

// JaegerMiddleware extracts span (if exist) from request headers and set it in the context with tags
func JaegerMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		upstreamCtx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)

		var span opentracing.Span
		if err == nil {
			span = tracer.StartSpan(c.Request.URL.Path, ext.RPCServerOption(upstreamCtx))

		} else {
			span = opentracing.StartSpan(c.Request.URL.Path)
		}

		defer span.Finish()

		span.SetTag("component", "gin")
		span.SetTag("http.method", c.Request.Method)
		span.SetTag("http.url", c.Request.URL)

		ctx := context.Background()
		ctx = opentracing.ContextWithSpan(ctx, span)
		c.Request = c.Request.Clone(ctx)

		c.Next()
		span.SetTag("http.status_code", c.Writer.Status())
	}
}

// Error adds semantic convention tags to a span
func Error(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogFields(
		traceLog.String("event", "error"),
		traceLog.String("message", err.Error()),
	)
}
