package controllers

import (
	"awesomeProject/global"
	"awesomeProject/models"
	"awesomeProject/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	user.Password = hashedPwd
	if err := global.DB.Create(&user).Error; err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "注册成功", "username": user.Username})
}
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user, err := getUserByUsername(input.Username)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	if !utils.CheckPassword(input.Password, user.Password) {
		utils.RespondError(ctx, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
