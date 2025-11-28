package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
	Slug string `gorm:"unique;not null"`
}
