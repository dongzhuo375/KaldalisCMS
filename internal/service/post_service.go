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

// --- Write Operations ---

func (s *PostService) CreatePost(post entity.Post) ( error) {
	// 进行业务逻辑验证 (Entity 自身校验)
	if err := post.CheckValidity(); err != nil {
		return  fmt.Errorf("文章数据校验失败: %w", err)
	}

	// 调用 Repository 执行持久化
	 err := s.repo.Create(post)
	if err != nil {
		// 封装错误
		return  fmt.Errorf("保存文章失败: %w", err)
	}
	return  nil
}

func (s *PostService) UpdatePost(id int, updatedEntity entity.Post) ( error) {
	// 获取现有 Entity
	existingEntity, err := s.repo.GetByID(id)
	if err != nil {
		// 错误检查 (假设 core.ErrNotFound 已定义)
		return fmt.Errorf("更新失败，文章不存在: %w", err)
	}

	// 状态合并
	existingEntity.Title = updatedEntity.Title
	existingEntity.Content = updatedEntity.Content
	existingEntity.ID = id

	// Entity 检查自身完整性
	if err := existingEntity.CheckValidity(); err != nil {
		return fmt.Errorf("更新后的数据校验失败: %w", err)
	}

	// 调用 Repository 执行更新
	 	err = s.repo.Update(existingEntity)
	if err != nil {
		return fmt.Errorf("更新文章失败: %w", err)
	}

	return  nil
}

// --- Read Operations ---

// 补充：根据 ID 获取文章
func (s *PostService) GetPostByID(id int) (entity.Post, error) {
	post, err := s.repo.GetByID(id)

	// 检查核心层抛出的契约错误
	//if errors.Is(err, core.ErrNotFound) {
	//	// 转换为 Service 层的语义错误或直接返回封装错误
	//	return entity.Post{}, fmt.Errorf("文章查找失败: %w", err)
	//}
	if err != nil {
		return entity.Post{}, fmt.Errorf("获取文章失败: %w", err)
	}
	return post, nil
}

// 补充：获取所有文章列表
func (s *PostService) GetAllPosts() ([]entity.Post, error) {
	// 业务逻辑 (例如：分页参数处理、权限筛选等) 可以在这里添加

	posts, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("获取所有文章列表失败: %w", err)
	}
	return posts, nil
}

//待repo补充错误处理

// --- Status Operations ---

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
	  err = s.repo.Update(post)
	if err != nil {
		return fmt.Errorf("更新发布状态失败: %w", err)
	}

	return nil
}

// --- Delete Operations ---

// 补充：删除文章
func (s *PostService) DeletePost(id int) error {
	// 可以在这里添加业务逻辑 (例如：权限检查、存档/软删除逻辑)

	err := s.repo.Delete(id)

	// 检查核心层抛出的契约错误
	//if errors.Is(err, core.ErrNotFound) {
	//	return fmt.Errorf("删除失败，文章不存在: %w", err)
	//}
	if err != nil {
		return fmt.Errorf("删除文章失败: %w", err)
	}
	return nil
}
