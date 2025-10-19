package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ArticleID  uint   `gorm:"index" json:"article_id"`
	UserID     uint   `gorm:"index" json:"user_id"`
	Content    string `gorm:"type:text;not null" json:"content"`
	AuditState string `gorm:"type:varchar(20);default:'pending'" json:"audit_state"`
	Reason     string `gorm:"type:text" json:"reason"`
	LikeCount  int64  `gorm:"default:0" json:"like_count"`
}
