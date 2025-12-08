package model

import (
	"time"
	"gorm.io/gorm"
)

type Tag struct {
	// 1. 展开 gorm.Model
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 2. 标签名
	// 必须唯一，不能为空，不能全是空格
	Name string `gorm:"unique;not null;check:char_length(TRIM(name)) > 0" json:"name"`

	// 3. 别名
	// 用于 URL (/tag/golang)
	Slug string `gorm:"unique;not null;check:char_length(TRIM(slug)) > 0" json:"slug"`

	// 4. 反向关联文章 (多对多)
	//  重点：建议加 json:"-" 
	// 原因 A: 防止 JSON 序列化死循环 (Post -> Tags -> Posts -> Tags ...)
	// 原因 B: 获取标签列表时，通常不需要把每个标签下的几百篇文章都查出来，太慢
	Posts []Post `gorm:"many2many:post_tags;" json:"-"`
}