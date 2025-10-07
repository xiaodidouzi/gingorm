package utils

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// 统一错误响应
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

// 统一成功响应
func RespondSuccess(ctx *gin.Context, data interface{}) {
	log.Printf("[INFO] %s | 200 | Path: %s",
		time.Now().Format("2006-01-02 15:04:05"),
		ctx.Request.URL.Path,
	)
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}
