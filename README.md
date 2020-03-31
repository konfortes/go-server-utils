# go-server-utils

Enhancement of gin server management.  
It has some best practices for observability - logging monitoring and tracing.  

## server

- `Initialize()`: initialize a full-blown gin server instrumented with prometheus monitoring and Jaeger tracing.  
- `GracefulShutdown()`: gracefully shut down of a server

## logging

- `JSONLogMiddleware()` returns a json logging middleware.

## monitoring

- `Instrument()` instruments a gin app with prometheus monitoring and exposes the `/metrics` endpoint

## tracing

- `Instrument()`: instruments a gin app with Jaeger tracing (and adds the JaegerMiddleware).
- `JaegerMiddleware()`: A Jaeger middleware to extract span from headers and set it in context if exist
