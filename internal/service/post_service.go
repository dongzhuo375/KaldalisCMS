package service

import (
	"fmt"

	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity" // Service 只能使用 Entity
)

type PostService struct {
	repo core.PostRepository
}

func NewPostService(repo core.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) CreatePost(post entity.Post) (entity.Post, error) {
	// 进行业务逻辑验证
	if err := post.CheckValidity(); err != nil {
		return entity.Post{}, fmt.Errorf("文章数据校验失败: %w", err)
	}

	// 调用 Repository 执行持久化
	savedPost, err := s.repo.Create(post)
	if err != nil {
		// 封装错误
		return entity.Post{}, fmt.Errorf("保存文章失败: %w", err)
	}
	return savedPost, nil
}

func (s *PostService) UpdatePost(id int, updatedEntity entity.Post) (entity.Post, error) {
	// 获取现有 Entity (Repo 返回的就是 Entity)
	existingEntity, err := s.repo.GetByID(id)
	if err != nil {
		// 错误封装（假设已定义 core.ErrNotFound）
		return entity.Post{}, fmt.Errorf("更新失败，文章不存在: %w", err)
	}

	// 状态合并（直接更新 Entity 字段）
	existingEntity.Title = updatedEntity.Title
	existingEntity.Content = updatedEntity.Content
	existingEntity.ID = id // 确保 ID 被携带

	// Entity 检查自身完整性
	if err := existingEntity.CheckValidity(); err != nil {
		return entity.Post{}, fmt.Errorf("更新后的数据校验失败: %w", err)
	}

	// 调用 Repository 执行更新，传递 Entity
	savedPost, err := s.repo.Update(existingEntity)
	if err != nil {
		return entity.Post{}, fmt.Errorf("更新文章失败: %w", err)
	}

	return savedPost, nil
}

func (s *PostService) PublishPost(id int) error {
	// 1. 获取 Entity
	post, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("发布失败，文章不存在: %w", err)
	}

	// 2. 调用 Entity 的核心业务行为 (状态流转)
	if err := post.Publish(); err != nil {
		return fmt.Errorf("发布文章失败: %w", err)
	}

	// 3. Service 协调：将已修改的 Entity 传递给 Repo 持久化
	_, err = s.repo.Update(post)
	if err != nil {
		return fmt.Errorf("更新发布状态失败: %w", err)
	}

	return nil
}
