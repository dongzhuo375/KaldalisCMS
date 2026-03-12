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

	// 6. 状态：0=草稿/下线, 1=已发布。
	// 公共读取与作者草稿查询都会基于该字段过滤，并参与 author+status 复合索引。
	//这个索引是为了支持 GetDraftsByAuthor 中使用的 “按作者获取草稿” 查询模式
	Status int `gorm:"not null;default:0;index:idx_posts_author_status,priority:2" json:"status"`

	// --- 关联关系 ---

	// 作者 (必填)
	// 作者后台会频繁按 (author_id, status=draft) 检索自己的草稿，因此加入复合索引。
	AuthorID uint `gorm:"not null;index:idx_posts_author_status,priority:1" json:"author_id"`

	Author User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`

	// 分类 (选填，使用指针 *uint 允许存 NULL)
	CategoryID *uint     `json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	// 标签 (多对多)
	Tags []Tag `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}
