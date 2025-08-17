package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func RespondError(ctx *gin.Context, code int, msg string) {
	log.Printf("[ERROR] %s | %d | %s | Path: %s",
		time.Now().Format("2006-01-02 15:04:05"),
		code,
		msg,
		ctx.Request.URL.Path,
	)
	ctx.JSON(code, gin.H{
		"code":    code,
		"message": msg,
		"data":    nil,
	})
}
