package repository

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"KaldalisCMS/internal/infra/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// --- Mapper Functions ---

// model转换成entity
func userToEntity(m model.User) entity.User {
	return entity.User{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		Role:      m.Role,
	}
}

// entity转换成model
func userToModel(e entity.User) model.User {
	return model.User{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		Username:  e.Username,
		Email:     e.Email,
		Password:  e.Password,
		Role:      e.Role,
	}
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	var userModels []model.User
	if err := r.db.WithContext(ctx).Find(&userModels).Error; err != nil {
		return nil, fmt.Errorf("user_repository.GetAll: %w", err)
	}
	var userEntities []entity.User
	for _, um := range userModels {
		userEntities = append(userEntities, userToEntity(um))
	}
	return userEntities, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (entity.User, error) {
	var userModel model.User
	if err := r.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, core.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("user_repository.GetByID: %w", err)
	}
	return userToEntity(userModel), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	var userModel model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, core.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("user_repository.GetByUsername: %w", err)
	}
	return userToEntity(userModel), nil
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) error {
	userModel := userToModel(user)
	if err := r.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return core.ErrDuplicate
			}
		}
		return fmt.Errorf("user_repository.CreateUser: %w", err)
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user entity.User) error {
	userModel := userToModel(user)
	if err := r.db.WithContext(ctx).Save(&userModel).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return core.ErrDuplicate
			}
		}
		return fmt.Errorf("user_repository.Update: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("repository.DeleteUser: %w", err)
	}
	return nil
}
