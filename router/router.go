package router

import (
	v1 "awesomeProject/api/v1"
	v2 "awesomeProject/api/v2"
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
	// v1 auth
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", v1.Login)
		auth.POST("/register", v1.Register)
	}
	// v1 article
	article := r.Group("api/v1/articles")
	{
		article.POST("", v1.CreateArticle)
		article.GET("/:id", v1.GetArticleByID)
		article.GET("", v1.GetArticle)
		like := article.Group("/:id/like")
		{
			like.POST("", v1.LikeArticle)
			like.GET("", v1.GetArticleLikes)
		}
	}
	//v2 article likes kafka
	v2Article := r.Group("api/v2/articles")
	v2Article.Use(middlewares.AuthMiddleWare())
	{
		v2Article.POST("/:id/like", v2.LikeArticle)
	}
	// v1 exchange rates
	exchange := r.Group("/api/v1/exchangeRates")
	{
		exchange.GET("", v1.GetExchange)
		// 需要鉴权
		exchangeAuth := exchange.Use(middlewares.AuthMiddleWare())
		{
			exchangeAuth.POST("", v1.CreateExchangeRate)
		}
	}
	return r
}
