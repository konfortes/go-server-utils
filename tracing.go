package main

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func initJaeger(service string, shutdownHooks []func()) (opentracing.Tracer, io.Closer) {
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

	shutdownHooks = append(shutdownHooks, func() {
		tracerCloser.Close()
	})

	opentracing.SetGlobalTracer(tracer)

	return tracer, tracerCloser
}
