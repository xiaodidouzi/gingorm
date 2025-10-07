package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username" binding:"required,min=3,max=20"`
	Password string `gorm:"type:varchar(255);not null" json:"-" binding:"required,min=6,max=20"`
}
