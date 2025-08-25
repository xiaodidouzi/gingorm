package v2

import (
	"awesomeProject/global"
	"awesomeProject/kafka"
	"awesomeProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LikeArticle(ctx *gin.Context) {
	articleIDStr := ctx.Param("id")
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, "未登录")
		return
	}
	userID := userIDVal.(int)
	likeKey := fmt.Sprintf("like:%d:%d", articleID, userID)
	rexists, err := global.RedisDB.Exists(ctx, likeKey).Result()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Redis错误")
		return
	}
	if rexists == 1 {
		ctx.JSON(http.StatusOK, gin.H{"message": "已点赞"})
		return
	}
	if err := kafka.SendLikeMessage(articleID, userID); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "发送消息失败")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}
