package api

import (
	"bookstore-api/internal/observability"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Correlation-ID")
		if id == "" {
			id = generateID()
		}

		logger := slog.Default().With(slog.String("correlation_id", id))
		ctx := observability.WithLogger(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Correlation-ID", id)
		c.Next()
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		observability.FromContext(c.Request.Context()).Info("request started",
			slog.String("method", method),
			slog.String("path", path),
		)

		c.Next()

		observability.FromContext(c.Request.Context()).Info("request completed",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
		)
	}
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
