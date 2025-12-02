package entity

import (
	"errors"
	"fmt"
	"time"
)

const (
	StatusDraft     = 0 // 草稿
	StatusPublished = 1 // 已发布
)

type Post struct {
	ID         int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string
	Slug       string
	Content    string
	Cover      string
	AuthorID   uint
	Author     User // 这里可能需要根据实际情况调整
	CategoryID *uint
	Category   Category // 这里可能需要根据实际情况调整
	Tags       []Tag    // 这里可能需要根据实际情况调整
	Status     int      // 文章状态
}



// Category 结构体，简化版
type Category struct {
	ID   uint
	Name string
}

// Tag 结构体，简化版
type Tag struct {
	ID   uint
	Name string
}

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

	// 发布前的业务规则校验
	if err := p.CheckValidity(); err != nil {
		return fmt.Errorf("文章发布失败，校验未通过: %w", err)
	}

	// 状态流转
	p.Status = StatusPublished
	p.UpdatedAt = time.Now()

	return nil
}

// Draft 设置文章状态为草稿
func (p *Post) Draft() error {
	if p.Status == StatusDraft {
		return errors.New("文章已是草稿状态")
	}
	p.Status = StatusDraft
	return nil
}
