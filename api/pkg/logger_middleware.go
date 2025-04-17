package pkg

import (
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

func GinLogger(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		// Construir campos de log
		fields := map[string]any{
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"status":     c.Writer.Status(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"latency_ms": time.Since(start).Milliseconds(),
		}

		logger.Info("Request handled", fields)
	}
}

func GinRecovery(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())

				fields := map[string]any{
					"path":  c.Request.URL.Path,
					"error": err,
					"stack": stack,
				}

				logger.Error("Recovered from panic", fields)

				c.AbortWithStatusJSON(500, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
