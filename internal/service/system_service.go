package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SystemService struct {
	db         *gorm.DB
	systemRepo *repository.SystemRepository
	userSvc    core.UserService
}

func NewSystemService(db *gorm.DB, systemRepo *repository.SystemRepository, userSvc core.UserService) *SystemService {
	return &SystemService{db: db, systemRepo: systemRepo, userSvc: userSvc}
}

type SystemStatus struct {
	Installed bool
	SiteName  *string
}

func (s *SystemService) Status(ctx context.Context) (SystemStatus, error) {
	set, err := s.systemRepo.Get(ctx)
	if errors.Is(err, repository.ErrSystemSettingNotFound) {
		return SystemStatus{Installed: false, SiteName: nil}, nil
	}
	if err != nil {
		return SystemStatus{}, normalizeServiceErrorWithOpMsg("system.status", "load system status failed", err)
	}
	var siteName *string
	if strings.TrimSpace(set.SiteName) != "" {
		v := set.SiteName
		siteName = &v
	}
	return SystemStatus{Installed: set.Installed, SiteName: siteName}, nil
}

type SetupParams struct {
	SiteName      string
	AdminUsername string
	AdminEmail    string
	AdminPassword string
}

func (s *SystemService) SetupOnce(ctx context.Context, p SetupParams) error {
	p.SiteName = strings.TrimSpace(p.SiteName)
	p.AdminUsername = strings.TrimSpace(p.AdminUsername)
	p.AdminEmail = strings.TrimSpace(p.AdminEmail)

	if p.SiteName == "" || p.AdminUsername == "" || p.AdminEmail == "" || p.AdminPassword == "" {
		return core.ErrInvalidInput
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Ensure singleton row exists.
		if err := s.systemRepo.EnsureSingletonRow(ctx, tx); err != nil {
			return normalizeServiceErrorWithOpMsg("system.setup.ensure_singleton", "ensure singleton system row failed", err)
		}

		// Lock the row for update to serialize setup.
		set, err := s.systemRepo.LockForUpdate(ctx, tx)
		if err != nil {
			return normalizeServiceErrorWithOpMsg("system.setup.lock", "lock system setup row failed", err)
		}
		if set.Installed {
			return fmt.Errorf("%w: already installed", core.ErrConflict)
		}

		// Create the first super admin user.
		admin := entity.User{
			Username: p.AdminUsername,
			Email:    p.AdminEmail,
			Password: p.AdminPassword,
			Role:     "admin",
		}

		// Use the existing UserService so password hashing and duplicate handling stays consistent.
		// Note: UserService uses the repository which is bound to the root gorm.DB.
		// That means it won't automatically join this transaction.
		// To keep existing logic untouched AND keep setup atomic, we create the user via tx directly here.
		// We still reuse entity password hashing method.
		if err := admin.SetPassword(admin.Password); err != nil {
			return fmt.Errorf("%w: invalid admin password", core.ErrInvalidInput)
		}

		// Create user inside the same transaction.
		userRepo := repository.NewUserRepository(tx)
		if err := userRepo.Create(ctx, admin); err != nil {
			if errors.Is(err, core.ErrDuplicate) {
				return fmt.Errorf("%w: admin account already exists", core.ErrConflict)
			}
			return normalizeServiceErrorWithOpMsg("system.setup.create_admin", "create admin account during setup failed", err)
		}

		affected, err := s.systemRepo.MarkInstalled(ctx, tx, p.SiteName, time.Now())
		if err != nil {
			return normalizeServiceErrorWithOpMsg("system.setup.mark_installed", "mark system installed failed", err)
		}
		if affected == 0 {
			return fmt.Errorf("%w: already installed", core.ErrConflict)
		}

		return nil
	})
}
