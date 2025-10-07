package v1

import (
	"gingorm/service"
	"gingorm/service/dto"
	"gingorm/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserAPI struct {
	userService *service.UserService
}

// NewUserAPI 构造函数
func NewUserAPI(userService *service.UserService) *UserAPI {
	return &UserAPI{userService: userService}
}

// Register 用户注册
func (u *UserAPI) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := u.userService.Register(req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, res)
}

// Login 用户登录
func (u *UserAPI) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := u.userService.Login(req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, res)
}
