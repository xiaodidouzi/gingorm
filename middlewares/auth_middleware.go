package middlewares

import (
	"gingorm/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			utils.RespondError(ctx, http.StatusUnauthorized, "Missing Authorization")
			ctx.Abort()
			return
		}
		userID, username, err := utils.ParseJWT(token)
		if err != nil {
			utils.RespondError(ctx, http.StatusUnauthorized, "Invalid token")
			ctx.Abort()
			return
		}
		ctx.Set("username", username)
		ctx.Set("userID", uint(userID))
		ctx.Next()
	}
}
