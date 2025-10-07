package dto

// 创建文章请求 DTO
type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required,min=3,max=50"`
	Content  string `json:"content" binding:"required"`
	Category string `json:"category"`
}

// 创建文章响应 DTO
type CreateArticleResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	AuthorID  uint   `json:"author_id"`
	CreatedAt string `json:"created_at"`
}

// 获取文章详情响应 DTO
type GetArticleResponse struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	AuthorID     uint   `json:"author_id"`
	Summary      string `json:"summary"`
	LikeCount    int64  `json:"like_count"`
	CommentCount int64  `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// 分页查询列表请求 DTO
type ListArticlesRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Category string `form:"category"`
	AuthorID uint   `form:"author_id"`
}

// 分页查询列表响应 DTO
type ListArticlesResponse struct {
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
	List     []GetArticleResponse `json:"list"`
}

// 更新文章请求 DTO
type UpdateArticleRequest struct {
	//指针确保空字段不更新
	Title    *string `json:"title,omitempty" `
	Content  *string `json:"content,omitempty"`
	Category *string `json:"category,omitempty"`
}

// 更新文章响应 DTO
type UpdateArticleResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	AuthorID  uint   `json:"author_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
