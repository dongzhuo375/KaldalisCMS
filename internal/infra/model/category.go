package model

import (
	"time"
	"gorm.io/gorm"
)

type Category struct {
	// 1. 展开 gorm.Model，为了统一 JSON 风格 (id, created_at)
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 2. 分类名称
	// check: 禁止空名字
	Name string `gorm:"unique;not null;check:char_length(TRIM(name)) > 0" json:"name"`

	// 3. 别名/路由 (Slug)
	// 用于 URL (例如: /category/tech)
	// check: 禁止空 Slug
	Slug string `gorm:"unique;not null;check:char_length(TRIM(slug)) > 0" json:"slug"`
    
	// 4. (可选建议) 关联文章
	// 如果你需要查询 "这个分类下有哪些文章"，可以加这个字段
	// json:"-" 建议默认隐藏，防止查询分类列表时把所有文章都带出来，导致数据量爆炸
	Posts []Post `gorm:"foreignKey:CategoryID" json:"-"`
}