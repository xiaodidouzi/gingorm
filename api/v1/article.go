package v1

import (
	"awesomeProject/global"
	"awesomeProject/models"
	"awesomeProject/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var cacheKey = "articles"

func CreateArticle(ctx *gin.Context) {
	var article models.Article
	if err := ctx.ShouldBindJSON(&article); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := global.DB.Create(&article).Error; err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if err := global.RedisDB.Del(ctx, cacheKey).Err(); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusOK, article)
}
func GetArticle(ctx *gin.Context) {
	cachedData, err := global.RedisDB.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		var articles []models.Article
		if err = global.DB.Find(&articles).Error; err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		articleJSON, err := json.Marshal(articles)
		if err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		if err := global.RedisDB.Set(ctx, cacheKey, articleJSON, 10*time.Minute).Err(); err != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, articles)
		return
	}
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	var articles []models.Article
	if err := json.Unmarshal([]byte(cachedData), &articles); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, articles)

}

func GetArticleByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var article models.Article
	if err := global.DB.Where("id=?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondError(ctx, http.StatusNotFound, "article not found")
		} else {
			utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, article)
}
