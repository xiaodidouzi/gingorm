package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title        string `gorm:"type:varchar(50);not null" json:"title" binding:"required,min=3,max=50"`
	Content      string `gorm:"type:text;not null" json:"content" binding:"required"`
	AuthorID     uint   `gorm:"not null;index" json:"author_id"` //外键
	Category     string `gorm:"type:varchar(50)" json:"category"`
	Summary      string `gorm:"type:varchar(500)" json:"summary"`
	LikeCount    int64  `gorm:"default:0" json:"like_count"`
	CommentCount int64  `gorm:"default:0" json:"comment_count"`
}
