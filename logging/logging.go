package logging

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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
