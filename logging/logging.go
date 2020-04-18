package logging

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// context.WithValue is being looked up by the key type. having a key of private type avoids context misuse
type correlationIDType int

const (
	requestIDKey     correlationIDType = iota
	requestIDKeyName                   = "reqId"
)

// a fallback/root logger for events without context
// it is a sugared zap logger but can be de-sugared by the app when performance is critical
var logger *zap.SugaredLogger

func init() {
	var l *zap.Logger
	var err error

	if os.Getenv("ENV") == "production" {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Panic(err)
	}
	logger = l.Sugar()
}

// WithRqID returns a context contains a request ID
func WithRqID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// Logger returns a zap logger with request id context if available
func Logger(ctx context.Context) *zap.SugaredLogger {
	newLogger := logger
	if ctx != nil {
		if ctxRqID, ok := ctx.Value(requestIDKey).(string); ok {
			newLogger = newLogger.With(zap.String(requestIDKeyName, ctxRqID))
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

		if c.Keys == nil {
			c.Keys = make(map[string]interface{})
		}
		c.Keys[requestIDKeyName] = requestID.String()

		c.Next()
	}
}

// RequestMiddleware returns a middleware to log requests in a json format
// func RequestMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()
// 		c.Next()
// 		timeStamp := time.Now()
// 		latency := timeStamp.Sub(start)

// 		// errorMessage := c.Errors.ByType(ErrorTypePrivate).String()

// 		logger := Logger(c.Request.Context())
// 		defer logger.Sync()

// 		logger.With(
// 			"ts", timeStamp,
// 			"method", c.Request.Method,
// 			"path", c.Request.URL.Path,
// 			"duration", latency,
// 			"statusCode", c.Writer.Status(),
// 		).Info("got request")
// 	}
// }

// JSONLogMiddleware returns a gin HandlerFunc that logs json using gin's infrastructure
func JSONLogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(formatter)
}

func formatter(params gin.LogFormatterParams) string {
	var errorMessageEntry string
	if len(params.ErrorMessage) > 0 {
		errorMessageEntry = fmt.Sprintf(`,"error":"%s"`, params.ErrorMessage)
	}

	var reqIDEntry string
	if id, found := params.Keys[requestIDKeyName].(string); found {
		reqIDEntry = fmt.Sprintf(`,"%s":%s`, requestIDKeyName, id)
	}

	return fmt.Sprintf(`{"ts":"%s","method":"%s","path":"%s","code":%d,"duration":%d%s%s}`,
		params.TimeStamp,
		params.Method,
		params.Path,
		params.StatusCode,
		params.Latency,
		reqIDEntry,
		errorMessageEntry) + "\n"
}
