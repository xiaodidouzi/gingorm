package dto

// 创建评论请求 DTO
type CreateCommentRequest struct {
	ArticleID uint   `json:"article_id"`
	UserID    uint   `json:"user_id" `
	Content   string `json:"content" binding:"required,min=1,max=500"`
}

// 创建评论响应 DTO
type CreateCommentResponse struct {
	ID        uint   `json:"id"`
	ArticleID uint   `json:"article_id"`
	UserID    uint   `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// 分页查询评论列表请求 DTO
type ListCommentsRequest struct {
	ArticleID uint `form:"article_id" `
	UserID    uint `form:"user_id"`
	Page      int  `form:"page" binding:"min=1"`
	PageSize  int  `form:"page_size" binding:"min=1,max=100"`
}

// 分页查询评论列表响应 DTO
type ListCommentsResponse struct {
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	List     []CreateCommentResponse `json:"list"`
}
