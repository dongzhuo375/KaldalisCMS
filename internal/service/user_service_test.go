package service

import (
	"context"
	"errors"
	"testing"

	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
)

type fakeUserRepo struct {
	getAllFn        func(ctx context.Context) ([]entity.User, error)
	getByIDFn       func(ctx context.Context, id uint) (entity.User, error)
	getByUsernameFn func(ctx context.Context, username string) (entity.User, error)
	createFn        func(ctx context.Context, user entity.User) error
	updateFn        func(ctx context.Context, user entity.User) error
	deleteFn        func(ctx context.Context, id uint) error
}

func (f *fakeUserRepo) GetAll(ctx context.Context) ([]entity.User, error) { return f.getAllFn(ctx) }
func (f *fakeUserRepo) GetByID(ctx context.Context, id uint) (entity.User, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakeUserRepo) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	return f.getByUsernameFn(ctx, username)
}
func (f *fakeUserRepo) Create(ctx context.Context, user entity.User) error {
	return f.createFn(ctx, user)
}
func (f *fakeUserRepo) Update(ctx context.Context, user entity.User) error {
	return f.updateFn(ctx, user)
}
func (f *fakeUserRepo) Delete(ctx context.Context, id uint) error { return f.deleteFn(ctx, id) }

func TestUserService_CreateUser_HashesPassword(t *testing.T) {
	ctx := context.Background()
	var captured entity.User
	repo := &fakeUserRepo{createFn: func(ctx context.Context, user entity.User) error {
		captured = user
		return nil
	}}

	err := NewUserService(repo).CreateUser(ctx, entity.User{
		Username: "alice",
		Password: "s3cret!",
	})
	if err != nil {
		t.Fatal(err)
	}
	if captured.Password == "s3cret!" || captured.Password == "" {
		t.Fatalf("password was not hashed: %q", captured.Password)
	}
	if !captured.CheckPassword("s3cret!") {
		t.Fatal("hashed password does not verify")
	}
}

func TestUserService_CreateUser_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepo{createFn: func(ctx context.Context, user entity.User) error {
		return core.ErrDuplicate
	}}
	err := NewUserService(repo).CreateUser(ctx, entity.User{Username: "x", Password: "p"})
	if !errors.Is(err, core.ErrDuplicate) {
		t.Fatalf("want ErrDuplicate, got %v", err)
	}
}

func TestUserService_VerifyUser(t *testing.T) {
	ctx := context.Background()

	user := entity.User{ID: 1, Username: "alice"}
	if err := user.SetPassword("pw123"); err != nil {
		t.Fatal(err)
	}

	t.Run("user not found -> invalid credentials", func(t *testing.T) {
		repo := &fakeUserRepo{getByUsernameFn: func(ctx context.Context, username string) (entity.User, error) {
			return entity.User{}, core.ErrNotFound
		}}
		_, err := NewUserService(repo).VerifyUser(ctx, "alice", "pw123")
		if !errors.Is(err, core.ErrInvalidCredentials) {
			t.Fatalf("want ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := &fakeUserRepo{getByUsernameFn: func(ctx context.Context, username string) (entity.User, error) {
			return user, nil
		}}
		_, err := NewUserService(repo).VerifyUser(ctx, "alice", "wrong")
		if !errors.Is(err, core.ErrInvalidCredentials) {
			t.Fatalf("want ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeUserRepo{getByUsernameFn: func(ctx context.Context, username string) (entity.User, error) {
			return user, nil
		}}
		got, err := NewUserService(repo).VerifyUser(ctx, "alice", "pw123")
		if err != nil {
			t.Fatal(err)
		}
		if got.ID != 1 {
			t.Fatalf("unexpected user: %+v", got)
		}
	})

	t.Run("repo error normalized", func(t *testing.T) {
		repo := &fakeUserRepo{getByUsernameFn: func(ctx context.Context, username string) (entity.User, error) {
			return entity.User{}, errors.New("db down")
		}}
		_, err := NewUserService(repo).VerifyUser(ctx, "alice", "pw123")
		if !errors.Is(err, core.ErrInternalError) {
			t.Fatalf("want ErrInternalError, got %v", err)
		}
	})
}

func TestUserService_Login_DelegatesToVerify(t *testing.T) {
	ctx := context.Background()
	user := entity.User{ID: 42, Username: "bob"}
	_ = user.SetPassword("pw")
	repo := &fakeUserRepo{getByUsernameFn: func(ctx context.Context, username string) (entity.User, error) {
		return user, nil
	}}
	got, err := NewUserService(repo).Login(ctx, "bob", "pw")
	if err != nil || got.ID != 42 {
		t.Fatalf("unexpected: %+v %v", got, err)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepo{getByIDFn: func(ctx context.Context, id uint) (entity.User, error) {
		return entity.User{ID: id}, nil
	}}
	got, err := NewUserService(repo).GetUserByID(ctx, 7)
	if err != nil || got.ID != 7 {
		t.Fatalf("unexpected: %+v %v", got, err)
	}

	repoErr := &fakeUserRepo{getByIDFn: func(ctx context.Context, id uint) (entity.User, error) {
		return entity.User{}, core.ErrNotFound
	}}
	if _, err := NewUserService(repoErr).GetUserByID(ctx, 7); !errors.Is(err, core.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
