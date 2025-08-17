package controllers

import (
	"awesomeProject/global"
	"awesomeProject/models"
	"awesomeProject/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
)

func LikeArticle(ctx *gin.Context) {
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + "likes"
	exists, _ := global.RedisDB.Exists(ctx, likeKey).Result()
	if exists == 0 {
		var article models.Article
		if err := global.DB.First(&article, articleID).Error; err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		global.RedisDB.Set(ctx, likeKey, article.Likes, 0)
	}
	if err := global.RedisDB.Incr(ctx, likeKey).Err(); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "successful"})
}
func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + "likes"
	likes, err := global.RedisDB.Get(ctx, likeKey).Result()
	if err == redis.Nil {
		var article models.Article
		if err := global.DB.First(&article, articleID).Error; err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		likes = strconv.FormatInt(article.Likes, 10)
		global.RedisDB.Set(ctx, likeKey, article.Likes, 0)
	} else if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"likes": likes})
}
