package router

import (
	"awesomeProject/controllers"
	"awesomeProject/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.LoggerMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://"
		//},
		MaxAge: 12 * time.Hour,
	}))

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}
	api := r.Group("/api")
	api.GET("/exchangeRates", controllers.GetExchange)
	article := api.Group("/article")
	{
		article.POST("", controllers.CreateArticle)
		article.GET("/:id", controllers.GetArticleByID)
		article.GET("", controllers.GetArticle)
		like := article.Group("/:id/like")
		{
			like.POST("", controllers.LikeArticle)
			like.GET("", controllers.GetArticleLikes)
		}
	}

	authRequest := api.Group("")
	authRequest.Use(middlewares.AuthMiddleWare())
	{
		authRequest.POST("/exchangeRates", controllers.CreateExchangeRate)
	}
	return r
}
