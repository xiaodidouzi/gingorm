package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ArticleID uint   `gorm:"index" json:"article_id"`
	UserID    uint   `gorm:"index" json:"user_id"`
	Content   string `gorm:"type:text;not null" json:"content"`
	LikeCount int64  `gorm:"default:0" json:"like_count"`
}
