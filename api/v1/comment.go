package v1

import (
	"gingorm/service"
	"gingorm/service/dto"
	"gingorm/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentAPI struct {
	commentService *service.CommentService
}

func NewCommentAPI(commentService *service.CommentService) *CommentAPI {
	return &CommentAPI{
		commentService: commentService}
}

// 创建评论
func (c *CommentAPI) CreateComment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req dto.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	req.ArticleID = uint(articleID)
	userID := ctx.MustGet("userID").(uint)
	res, err := c.commentService.CreateComment(userID, req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, res)
}

// 获取文章评论列表
func (c *CommentAPI) ListComments(ctx *gin.Context) {
	idStr := ctx.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req dto.ListCommentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	req.ArticleID = uint(articleID)
	res, err := c.commentService.ListComments(req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, res)
}

// 删除评论
func (c *CommentAPI) DeleteComment(ctx *gin.Context) {
	commentIDStr := ctx.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID := ctx.MustGet("userID").(uint)
	msg, err := c.commentService.DeleteComment(uint(commentID), userID)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, msg)
}
