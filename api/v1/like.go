package v1

import (
	"gingorm/service"
	"gingorm/service/dto"
	"gingorm/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LikeAPI struct {
	likeService *service.LikeService
}

func NewLikeAPI(likeService *service.LikeService) *LikeAPI {
	return &LikeAPI{likeService: likeService}
}

func (a *LikeAPI) Like(ctx *gin.Context) {
	var req dto.LikeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID := ctx.MustGet("userID").(uint)
	res, err := a.likeService.Like(ctx, userID, req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondSuccess(ctx, res)
}
