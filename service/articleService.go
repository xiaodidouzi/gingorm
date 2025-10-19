package service

import (
	"errors"
	"gingorm/global"
	"gingorm/models"
	"gingorm/service/dto"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ArticleService struct {
	DB  *gorm.DB
	Rdb *redis.Client
}

func NewArticleService(db *gorm.DB, rdb *redis.Client) *ArticleService {
	return &ArticleService{
		DB:  db,
		Rdb: rdb,
	}
}

// 创建文章
func (s *ArticleService) CreateArticle(req dto.CreateArticleRequest, authorID uint) (*dto.CreateArticleResponse, error) {

	article := models.Article{
		Title:      req.Title,
		Content:    req.Content,
		AuthorID:   authorID,
		Category:   req.Category,
		AuditState: "pending",
	}

	if err := s.DB.Create(&article).Error; err != nil {
		return nil, err
	}

	res := &dto.CreateArticleResponse{
		ID:         article.ID,
		Title:      article.Title,
		AuthorID:   article.AuthorID,
		CreatedAt:  article.CreatedAt.Format("2006-01-02 15:04:05"),
		AuditState: article.AuditState,
		Reason:     article.Reason,
	}

	return res, nil
}

// 获取文章详情
func (s *ArticleService) GetArticle(articleID uint) (*dto.GetArticleResponse, error) {
	//查询数据库
	var article models.Article
	if err := global.DB.Where("id=?", articleID).First(&article).Error; err != nil {
		return nil, errors.New("文章不存在")
	}
	//返回结果
	content := article.Content
	if len(content) > 20 {
		content = content[:20]
	}

	res := &dto.GetArticleResponse{
		ID:           article.ID,
		Title:        article.Title,
		Content:      content,
		AuthorID:     article.AuthorID,
		LikeCount:    article.LikeCount,
		CommentCount: article.CommentCount,
		CreatedAt:    article.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    article.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return res, nil
}

// 分页查询文章列表
func (s *ArticleService) ListArticles(req dto.ListArticlesRequest) (*dto.ListArticlesResponse, error) {
	//默认页数条数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	//偏移量
	offset := (req.Page - 1) * req.PageSize
	//查询对象query
	query := global.DB.Model(&models.Article{})
	//查询条件 只查一种条件
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.AuthorID != 0 {
		query = query.Where("author_id = ?", req.AuthorID)
	}
	//计数
	var total int64
	query.Count(&total)

	var articles []models.Article
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&articles).Error; err != nil {
		return nil, err
	}
	//查出的内容转换为DTO列表
	var list []dto.GetArticleResponse
	for _, a := range articles {
		list = append(list, dto.GetArticleResponse{
			ID:           a.ID,
			Title:        a.Title,
			Content:      a.Content,
			AuthorID:     a.AuthorID,
			LikeCount:    a.LikeCount,
			CommentCount: a.CommentCount,
			CreatedAt:    a.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    a.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	res := &dto.ListArticlesResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}
	return res, nil
}

// 更新文章
func (s *ArticleService) UpdateArticle(articleID, userID uint, req dto.UpdateArticleRequest) (*dto.UpdateArticleResponse, error) {
	//检查文章是否存在
	var article models.Article
	if err := global.DB.First(&article, articleID).Error; err != nil {
		return nil, errors.New("文章不存在")
	}
	// 校验是否作者
	if article.AuthorID != userID {
		return nil, errors.New("无权限修改")
	}
	//构建更新字段
	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if len(updates) > 0 {
		if err := global.DB.Model(&article).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	res := &dto.UpdateArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		AuthorID:  article.AuthorID,
		CreatedAt: article.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: article.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return res, nil
}

// 删除文章
func (s *ArticleService) DeleteArticle(articleID, userID uint) (string, error) {
	//检查文章是否存在
	var article models.Article
	if err := global.DB.First(&article, articleID).Error; err != nil {
		return "", errors.New("文章不存在")
	}
	// 校验是否作者
	if article.AuthorID != userID {
		return "", errors.New("无权限删除")
	}

	if err := global.DB.Delete(&article).Error; err != nil {
		return "", err
	}
	return "删除成功", nil
}
