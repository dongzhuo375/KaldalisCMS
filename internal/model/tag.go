package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Slug  string `gorm:"unique;not null"`
	Posts []Post `gorm:"many2many:post_tags;"`
}
