package dto

// 点赞请求
type LikeRequest struct {
	TargetID   uint   `json:"target_id" binding:"required"`
	TargetType string `json:"target_type" binding:"required,oneof=article comment"` // 只允许两种类型
}
type LikeMessage struct {
	UserID     uint   `json:"user_id"`
	TargetID   uint   `json:"target_id"`
	TargetType string `json:"target_type"`
	Action     string `json:"action"` // "like" or "unlike"
}
