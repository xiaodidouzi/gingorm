package controllers

import (
	"awesomeProject/global"
	"awesomeProject/models"
	"awesomeProject/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate
	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	exchangeRate.Date = time.Now()
	if err := global.DB.Create(&exchangeRate).Error; err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, exchangeRate)
}
func GetExchange(ctx *gin.Context) {
	var exchangeRates []models.ExchangeRate
	if err := global.DB.Find(&exchangeRates).Error; err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, exchangeRates)
}
