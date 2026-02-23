package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gosimple/slug"

	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity" // Service 只能使用 Entity
)

type PostService struct {
	repo core.PostRepository
	// media is optional; when nil, reference sync is skipped.
	media *MediaService
}

func NewPostService(repo core.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

// NewPostServiceWithMedia wires an optional MediaService to keep post_assets in sync.
func NewPostServiceWithMedia(repo core.PostRepository, media *MediaService) *PostService {
	return &PostService{repo: repo, media: media}
}

func (s *PostService) DraftPost(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}

func (s *PostService) CreatePost(ctx context.Context, post entity.Post) error {
	// 进行业务逻辑验证 (Entity 自身校验)
	if err := post.CheckValidity(); err != nil {
		return fmt.Errorf("文章数据校验失败: %w", err)
	}

	generatedSlug := slug.Make(post.Title)

	if generatedSlug == "" {
		return fmt.Errorf("标题无法生成有效的URL标识符")
	}

	finalSlug, err := s.generateUniqueSlug(ctx, generatedSlug)
	if err != nil {
		return err // 无法生成唯一 Slug
	}

	post.Slug = finalSlug

	created, err := s.repo.Create(ctx, post)
	if err != nil {
		// 封装错误
		return fmt.Errorf("保存文章失败: %w", err)
	}

	// Sync media references.
	if s.media != nil {
		// If sync fails, we must rollback the post creation to maintain consistency
		// because we don't have a distributed transaction manager here.
		if err := s.media.SyncPostReferences(ctx, created.ID, created.Content, created.Cover); err != nil {
			// Rollback: delete the created post
			deleteErr := s.repo.Delete(ctx, created.ID)
			if deleteErr != nil {
				// Both sync and rollback failed. This is a critical error.
				// In a real system, we might log this to a special audit log or alert system.
				return fmt.Errorf("media sync failed: %v; ROLLBACK FAILED: %v", err, deleteErr)
			}
			return fmt.Errorf("media sync failed (post rolled back): %w", err)
		}
	}
	return nil
}

func (s *PostService) generateUniqueSlug(ctx context.Context, initialSlug string) (string, error) {
	currentSlug := initialSlug
	counter := 0
	maxAttempts := 100 // 最大尝试次数

	for {
		exists, err := s.repo.IsSlugExists(ctx, currentSlug)
		if err != nil {
			return "", fmt.Errorf("检查Slug唯一性失败: %w", err)
		}

		if !exists {
			return currentSlug, nil
		}

		counter++
		if counter >= maxAttempts {
			return "", errors.New("无法在合理尝试次数内生成唯一的URL标识符")
		}

		currentSlug = fmt.Sprintf("%s-%d", initialSlug, counter)
	}
}

func (s *PostService) UpdatePost(ctx context.Context, id uint, updatedEntity entity.Post) error {
	// 获取现有 Entity
	existingEntity, err := s.repo.GetByID(ctx, id)
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
	err = s.repo.Update(ctx, existingEntity)
	if err != nil {
		return fmt.Errorf("更新文章失败: %w", err)
	}

	if s.media != nil {
		// For update, we probably don't want to rollback the whole update if sync fails,
		// but it's arguable. The user specific concerns were about Create.
		// For now, let's just surface the error significantly.
		if err := s.media.SyncPostReferences(ctx, id, existingEntity.Content, existingEntity.Cover); err != nil {
			return fmt.Errorf("文章更新成功 but 媒体引用同步失败: %w", err)
		}
	}

	return nil
}

// --- Read Operations ---

// 补充：根据 ID 获取文章
func (s *PostService) GetPostByID(ctx context.Context, id uint) (entity.Post, error) {
	post, err := s.repo.GetByID(ctx, id)

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
func (s *PostService) GetAllPosts(ctx context.Context) ([]entity.Post, error) {
	// 业务逻辑 (例如：分页参数处理、权限筛选等) 可以在这里添加

	posts, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取所有文章列表失败: %w", err)
	}
	return posts, nil
}

// --- Status Operations ---

func (s *PostService) PublishPost(ctx context.Context, id uint) error {
	// 1. 获取 Entity
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("发布失败，文章不存在: %w", err)
	}

	// 2. 调用 Entity 的核心业务行为 (状态流转)
	if err := post.Publish(); err != nil {
		return fmt.Errorf("发布文章失败: %w", err)
	}

	// 3. Service 协调：将已修改的 Entity 传递给 Repo 持久化
	err = s.repo.Update(ctx, post)
	if err != nil {
		return fmt.Errorf("更新发布状态失败: %w", err)
	}

	return nil
}

// --- Delete Operations ---

// 补充：删除文章
func (s *PostService) DeletePost(ctx context.Context, id uint) error {
	// 可以在这里添加业务逻辑 (例如：权限检查、存档/软删除逻辑)

	err := s.repo.Delete(ctx, id)

	// 检查核心层抛出的契约错误
	//if errors.Is(err, core.ErrNotFound) {
	//	return fmt.Errorf("删除失败，文章不存在: %w", err)
	//}
	if err != nil {
		return fmt.Errorf("删除文章失败: %w", err)
	}
	return nil
}
