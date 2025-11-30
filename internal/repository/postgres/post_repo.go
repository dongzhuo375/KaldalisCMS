package repository

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"KaldalisCMS/internal/model"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// --- Mapper Functions ---

// model转换成entity
func toEntity(m model.Post) entity.Post {
	// Convert nested models to entities
	var authorEntity entity.User
	if m.Author.ID != 0 {
		authorEntity = entity.User{ID: m.Author.ID, Name: m.Author.Username}
	}

	var categoryEntity entity.Category
	if m.Category.ID != 0 {
		categoryEntity = entity.Category{ID: m.Category.ID, Name: m.Category.Name}
	}

	var tagsEntity []entity.Tag
	for _, tagModel := range m.Tags {
		tagsEntity = append(tagsEntity, entity.Tag{ID: tagModel.ID, Name: tagModel.Name})
	}

	return entity.Post{
		ID:         int(m.ID),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		Title:      m.Title,
		Slug:       m.Slug,
		Content:    m.Content,
		Cover:      m.Cover,
		AuthorID:   m.AuthorID,
		Author:     authorEntity,
		CategoryID: m.CategoryID,
		Category:   categoryEntity,
		Tags:       tagsEntity,
		Status:     m.Status,
	}
}

// entity转换成model
func toModel(e entity.Post) model.Post {
	return model.Post{
		Model:      gorm.Model{ID: uint(e.ID), CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt},
		Title:      e.Title,
		Slug:       e.Slug,
		Content:    e.Content,
		Cover:      e.Cover,
		AuthorID:   e.AuthorID,
		CategoryID: e.CategoryID,
		Status:     e.Status,
	}
}

type PostRepository struct {
	db *gorm.DB
}

func (r *PostRepository) IsSlugExists(slug string) (bool, error) {
	//TODO implement me
	panic("implement me")
	//AMBER你好，我已经在service实现了slug的基本功能，但是在数据库上需要判断slug是否存在，所以给你定义了个新的接口，望你实现
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

//根据core/error.go,主要处理未找到数据和重复错误，未知错误返回InternalError

func (r *PostRepository) GetAll() ([]entity.Post, error) {
	var postModels []model.Post
	if err := r.db.Preload("Author").Preload("Category").Preload("Tags").Find(&postModels).Error; err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}
	var postEntities []entity.Post
	for _, pm := range postModels {
		postEntities = append(postEntities, toEntity(pm))
	}
	return postEntities, nil
}

func (r *PostRepository) GetByID(id int) (entity.Post, error) {
	var postModel model.Post
	if err := r.db.Preload("Author").Preload("Category").Preload("Tags").First(&postModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Post{}, core.ErrNotFound
		}
		return entity.Post{}, fmt.Errorf("repository.GetByID: %w", err)
	}
	return toEntity(postModel), nil
}

func (r *PostRepository) Create(post entity.Post) error {
	postModel := toModel(post)
	if err := r.db.Create(&postModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // 23505 is the SQLSTATE for unique_violation
				return core.ErrDuplicate
			}
		}
		return fmt.Errorf("repository.Create: %w", err)
	}

	return nil
}

func (r *PostRepository) Update(post entity.Post) error {
	postModel := toModel(post)

	if err := r.db.Save(&postModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // 23505 is the SQLSTATE for unique_violation
				return core.ErrDuplicate
			}
		}
	}
	return nil
}

func (r *PostRepository) Delete(id int) error {
	if err := r.db.Delete(&model.Post{}, id).Error; err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}

	return nil
}
