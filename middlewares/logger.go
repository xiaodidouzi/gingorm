package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		status := ctx.Writer.Status()
		log.Printf("[INFO] %s | %d | %v | %s %s",
			time.Now().Format("2006-01-02 15:04:05"),
			status,
			duration,
			ctx.Request.Method,
			ctx.Request.URL.Path,
		)
	}
}
