package logging

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// context.WithValue is being looked up by the key type. having a key of private type avoids context misuse
type correlationIDType int

const (
	requestIDKey correlationIDType = iota
)

var logger *zap.Logger

func init() {
	// a fallback/root logger for events without context
	// TODO: handle error
	if os.Getenv("GO_ENV") == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
}

// WithRqID returns a context contains a request ID
func WithRqID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// Logger returns a zap logger with request id context if available
func Logger(ctx context.Context) *zap.Logger {
	newLogger := logger
	if ctx != nil {
		if ctxRqID, ok := ctx.Value(requestIDKey).(string); ok {
			newLogger = newLogger.With(zap.String("rqId", ctxRqID))
		}
	}
	return newLogger
}

// RequestIDMiddleware returns a gin.HandlerFunc that sets the context with request id
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := uuid.NewRandom()
		reqIDCtx := context.WithValue(c.Request.Context(), requestIDKey, requestID.String())
		c.Request = c.Request.WithContext(reqIDCtx)
		c.Next()
	}
}

// JSONLogMiddleware returns a gin HandlerFunc that logs json
func JSONLogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(formatter)
}

func formatter(params gin.LogFormatterParams) string {
	var errorMessage string
	if len(params.ErrorMessage) > 0 {
		errorMessage = fmt.Sprintf(`"error":"%s"`, params.ErrorMessage)
	}
	return fmt.Sprintf(`{"timestamp":"%s","method":"%s","path":"%s","code":%d,"took":%d%s}`,
		params.TimeStamp,
		params.Method,
		params.Path,
		params.StatusCode,
		params.Latency,
		errorMessage) + "\n"
}
