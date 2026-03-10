package entity

import (
	"errors"
	"fmt"
	"time"
)

const (
	// StatusDraft represents both a newly created draft and a post that has been taken offline.
	// The minimal workflow does not distinguish between "never published" and "unpublished" yet.
	StatusDraft = 0
	// StatusPublished marks content that is safe to expose on public read endpoints.
	StatusPublished = 1
)

// Post is the core publishing aggregate shared across service, repository, and API layers.
// Visibility is derived from Status rather than from routing or caller identity.
type Post struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string
	Slug       string
	Content    string
	Cover      string
	AuthorID   uint
	Author     User
	CategoryID *uint
	Category   Category
	Tags       []Tag
	Status     int // Draft or Published
}

// PostPatch models the editable subset of a post for management updates.
// Nil pointer fields mean "leave the current value unchanged".
// For Tags, nil means "do not touch tags", while an empty slice means "replace with empty".
type PostPatch struct {
	Title      *string
	Content    *string
	Cover      *string
	CategoryID *uint
	Tags       []Tag
}

// Category 结构体，简化版
type Category struct {
	ID   uint
	Name string
}

// Tag 结构体，简化版
//
// NOTE: Tag 已迁移到 entity/tag.go（避免在同一 package 内重复声明）。

// 检查实体自身是否满足发布或更新的基本业务要求
func (p *Post) CheckValidity() error {
	if p.Title == "" {
		return errors.New("标题不能为空")
	}

	// 内容最短设置
	//if len(p.Content) < 100 {
	//	return fmt.Errorf("内容长度必须大于 100 个字符 (当前 %d)", len(p.Content))
	//}

	//if p.Slug == "" {
	//	return errors.New("URL 标识符 (Slug) 不能为空")
	//}

	// 其余检查
	return nil
}

// Publish 封装了文章从草稿到发布的业务行为和状态流转
func (p *Post) Publish() error {
	if p.Status == StatusPublished {
		return errors.New("文章已是发布状态，无法重复发布")
	}

	if err := p.CheckValidity(); err != nil {
		return fmt.Errorf("文章发布失败，校验未通过: %w", err)
	}

	p.Status = StatusPublished
	p.UpdatedAt = time.Now()

	return nil
}

// Draft 将文章切回草稿状态，用于“下线”已发布内容。
func (p *Post) Draft() error {
	if p.Status == StatusDraft {
		return errors.New("文章已是草稿状态")
	}
	p.Status = StatusDraft
	p.UpdatedAt = time.Now()
	return nil
}
