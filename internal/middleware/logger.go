package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/welcomemonth/dancer-elite/internal/pkg/logger"
)

func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		ctx.Next()

		cost := time.Since(start)
		requestID, _ := ctx.Get("request_id")

		logger.Infow("request",
			"request_id", requestID,
			"status", ctx.Writer.Status(),
			"method", ctx.Request.Method,
			"path", path,
			"query", query,
			"ip", ctx.ClientIP(),
			"user-agent", ctx.Request.UserAgent(),
			"cost", cost.String(), // 或直接传 cost
			"body_size", ctx.Writer.Size(),
		)
	}
}
