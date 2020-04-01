# go-server-utils

Enhancement of gin server management.  
It has some best practices for observability - logging monitoring and tracing.  
It is as easy as:
```go
package main

import (
	"log"
	"net/http"

	"github.com/konfortes/go-server-utils/server"
	"github.com/konfortes/go-server-utils/utils"
)

const (
	appName = "my-app-name"
)

func main() {
	serverConfig := server.Config{
		AppName:     "my-app-name",
		Port:        utils.GetEnvOr("PORT", "3000"),
		Env:         utils.GetEnvOr("ENV", "development"),
        Handlers:    handlers(),
        ShutdownHooks: []func(){func() { log.Println("bye bye") }},
		WithTracing: utils.GetEnvOr("TRACING_ENABLED", "false") == "true",
	}

	srv := server.Initialize(serverConfig)

	go func() {
		log.Println("listening on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	server.GracefulShutdown(srv)
}
```

## server

- `Initialize()` initialize a full-blown gin server instrumented with prometheus monitoring and Jaeger tracing.  
- `GracefulShutdown()` gracefully shut down of a server

## logging

- `JSONLogMiddleware()` returns a json logging middleware.

## monitoring

- `Instrument()` instruments a gin app with prometheus monitoring and exposes the `/metrics` endpoint

## tracing

- `Instrument()` instruments a gin app with Jaeger tracing (and adds the JaegerMiddleware).
- `JaegerMiddleware()` A Jaeger middleware to extract span from headers and set it in context if exist

## utils

- `GetEnvOr()`: get env variable or fallback value