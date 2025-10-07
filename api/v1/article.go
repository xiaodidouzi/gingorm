package v1

import (
	"gingorm/service"
	"gingorm/service/dto"
	"gingorm/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ArticleAPI struct {
	articleService *service.ArticleService
}

func NewArticleAPI(articleService *service.ArticleService) *ArticleAPI {
	return &ArticleAPI{
		articleService: articleService}
}

// 创建文章
func (a *ArticleAPI) CreateArticle(ctx *gin.Context) {
	var req dto.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	//获取当前用户ID
	authorID := ctx.MustGet("userID").(uint)
	res, err := a.articleService.CreateArticle(req, authorID)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, res)
}

// 获取文章详情
func (a *ArticleAPI) GetArticle(ctx *gin.Context) {
	//获取文章ID
	idStr := ctx.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := a.articleService.GetArticle(uint(articleID))
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, res)
}

// 分页查询文章列表
func (a *ArticleAPI) ListArticles(ctx *gin.Context) {
	var req dto.ListArticlesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := a.articleService.ListArticles(req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, res)
}

// 更新文章
func (a *ArticleAPI) UpdateArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req dto.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userID := ctx.MustGet("userID").(uint)
	res, err := a.articleService.UpdateArticle(uint(articleID), userID, req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, res)
}

// 删除文章
func (a *ArticleAPI) DeleteArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userID := ctx.MustGet("userID").(uint)
	res, err := a.articleService.DeleteArticle(uint(articleID), userID)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, res)
}
