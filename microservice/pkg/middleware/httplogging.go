package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
	"time"
)

// HttpLoggingMiddleware that logs incoming http requests and their latencies to the logger instance
func HttpLoggingMiddleware(logger *slog.Logger, app *fiber.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		if err != nil {
			c.Status(http.StatusInternalServerError)
		}
		logger.Info("http request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency", time.Since(start).String(),
		)
		return nil
	}
}
