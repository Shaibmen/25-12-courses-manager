package middleware

import (
	"log/slog"
	"online-courses/internal/server/http/handlers/handlers_utils"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		var err *gin.Error
		if len(c.Errors) > 0 {
			err = c.Errors[0]
			response := handlers_utils.DefineError(err)

			c.JSON(response.Status, response.Message)
		}

		duration := time.Since(start)

		if duration > 1*time.Second {
			logger.Warn("too long duration",
				"duration", duration,
			)
		}

		clientIP := c.ClientIP()
		urlPath := c.Request.URL.Path
		method := c.Request.Method
		requestProto := c.Request.Proto

		userAgent := c.Request.UserAgent()

		logger.Info("incoming request",
			"ip", clientIP,
			"url-path", urlPath,
			"method", method,
			"request-proto", requestProto,
			"duration", duration.Seconds(),
			"user-agent", userAgent,
			"error", err,
		)

	}
}
