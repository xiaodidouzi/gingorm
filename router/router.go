package router

import (
	v1 "gingorm/api/v1"
	"gingorm/kafka"
	"gingorm/middlewares"
	"gingorm/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

import (
	"github.com/go-redis/redis/v8"
)

func SetupRouter(db *gorm.DB, rdb *redis.Client, producer *kafka.LikeProducer) *gin.Engine {
	r := gin.Default()
	// 注册 service
	userService := service.NewUserService(db, rdb)
	articleService := service.NewArticleService(db, rdb)
	commentService := service.NewCommentService(db, rdb)
	likeService := service.NewLikeService(db, rdb, producer)
	// 注册 API
	userAPI := v1.NewUserAPI(userService)
	articleAPI := v1.NewArticleAPI(articleService)
	commentAPI := v1.NewCommentAPI(commentService)
	likeAPI := v1.NewLikeAPI(likeService)

	api := r.Group("/api/v1")

	userGroup := api.Group("/user")
	{
		userGroup.POST("/register", userAPI.Register)
		userGroup.POST("/login", userAPI.Login)
	}
	auth := middlewares.AuthMiddleWare()
	articleGroup := api.Group("/article")
	{
		articleGroup.POST("", auth, articleAPI.CreateArticle)
		articleGroup.GET("/:id", articleAPI.GetArticle)
		articleGroup.GET("", articleAPI.ListArticles)
		articleGroup.PUT("/:id", auth, articleAPI.UpdateArticle)
		articleGroup.DELETE("/:id", auth, articleAPI.DeleteArticle)

		articleGroup.POST("/:id/comments", auth, commentAPI.CreateComment)
		articleGroup.GET("/:id/comments", commentAPI.ListComments)
		articleGroup.DELETE("/:id/comments/:comment_id", auth, commentAPI.DeleteComment)
	}

	likeGroup := api.Group("/like")
	{
		likeGroup.POST("", auth, likeAPI.Like)
		//likeGroup.GET("/count", likeAPI.GetLikeCount)
	}
	return r
}
