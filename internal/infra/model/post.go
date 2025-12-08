package model

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	// 1. 展开 gorm.Model，为了加 JSON 标签
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 2. 标题：核心字段，必须禁止空串和纯空格
	Title string `gorm:"not null;check:char_length(TRIM(title)) > 0" json:"title"`

	// 3. 别名 (Slug)：用于 URL (如 /post/my-first-post)
	// 必须唯一，且不能为空
	Slug string `gorm:"unique;not null;check:char_length(TRIM(slug)) > 0" json:"slug"`

	// 4. 内容：文章通常很长，显式指定 type:text 防止被截断
	// 内容允许为空（比如存草稿时可能只写了标题）
	Content string `gorm:"type:text" json:"content"`

	// 5. 封面图：存 URL，允许为空
	Cover string `json:"cover"`

	// 6. 状态：0=草稿, 1=发布, 2=归档
	// 给个默认值 0 (草稿)
	Status int `gorm:"default:0" json:"status"`

	// --- 关联关系 ---

	// 作者 (必填)
	AuthorID uint `gorm:"not null" json:"author_id"`
	
	Author User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`

	// 分类 (选填，使用指针 *uint 允许存 NULL)
	CategoryID *uint     `json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	// 标签 (多对多)
	Tags []Tag `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}