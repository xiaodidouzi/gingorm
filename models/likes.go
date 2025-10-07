package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID     uint   `gorm:"not null;index:idx_user_target,unique" json:"user_id"`
	TargetID   uint   `gorm:"not null;index:idx_user_target,unique" json:"target_id"`                    // 目标ID（文章ID或评论ID）
	TargetType string `gorm:"type:varchar(20);not null;index:idx_user_target,unique" json:"target_type"` // like类型：article/comment
}
