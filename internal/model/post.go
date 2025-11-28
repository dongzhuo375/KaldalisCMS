package model

import "gorm.io/gorm"

// Post represents the data structure for a blog post.
type Post struct {
	gorm.Model
	Title      string `gorm:"not null"`
	Slug       string `gorm:"unique;not null"`
	Content    string
	Cover      string
	AuthorID   uint
	Author     User `gorm:"foreignKey:AuthorID"`
	CategoryID uint
	Category   Category `gorm:"foreignKey:CategoryID"`
	Tags       []Tag    `gorm:"many2many:post_tags;"`
}
