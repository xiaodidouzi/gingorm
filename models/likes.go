package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	ArticleID int `json:"article_id" gorm:"index:idx_article_user,unique"`
	UserID    int `json:"user_id" gorm:"index:idx_article_user,unique"`
}
