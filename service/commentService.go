package service

import (
	"errors"
	"gingorm/global"
	"gingorm/models"
	"gingorm/service/dto"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CommentService struct {
	DB  *gorm.DB
	Rdb *redis.Client
}

func NewCommentService(db *gorm.DB, rdb *redis.Client) *CommentService {
	return &CommentService{
		DB:  db,
		Rdb: rdb,
	}
}

// 创建评论
func (s *CommentService) CreateComment(userID uint, req dto.CreateCommentRequest) (*dto.CreateCommentResponse, error) {
	// 检查文章是否存在
	var article models.Article
	if err := s.DB.First(&article, req.ArticleID).Error; err != nil {
		return nil, errors.New("文章不存在")
	}
	comment := models.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID,
		Content:   req.Content,
	}
	if err := global.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	res := &dto.CreateCommentResponse{
		ID:        comment.ID,
		ArticleID: comment.ArticleID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return res, nil
}

// 分页查询评论列表
func (s *CommentService) ListComments(req dto.ListCommentsRequest) (*dto.ListCommentsResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize

	query := global.DB.Model(&models.Comment{}).Where("article_id = ?", req.ArticleID)

	var total int64
	query.Count(&total)

	var comments []models.Comment
	if err := query.Order("created_at ASC").Offset(offset).Limit(req.PageSize).Find(&comments).Error; err != nil {
		return nil, err
	}

	list := make([]dto.CreateCommentResponse, 0, len(comments)) //
	for _, c := range comments {
		list = append(list, dto.CreateCommentResponse{
			ID:        c.ID,
			ArticleID: c.ArticleID,
			UserID:    c.UserID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	res := &dto.ListCommentsResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}

	return res, nil
}

// 删除评论
func (s *CommentService) DeleteComment(commentID, userID uint) (string, error) {
	var comment models.Comment
	if err := global.DB.First(&comment, commentID).Error; err != nil {
		return "", errors.New("评论不存在")
	}

	// 只能删除自己的评论
	if comment.UserID != userID {
		return "", errors.New("无权限删除")
	}
	if err := global.DB.Delete(&comment).Error; err != nil {
		return "", err
	}
	return "删除成功", nil
}
