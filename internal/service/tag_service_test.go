package service

import (
	"context"
	"errors"
	"testing"

	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
)

type fakeTagRepo struct {
	createFn    func(ctx context.Context, tag entity.Tag) (entity.Tag, error)
	getAllFn    func(ctx context.Context) ([]entity.Tag, error)
	getByIDFn   func(ctx context.Context, id uint) (entity.Tag, error)
	getByNameFn func(ctx context.Context, name string) (entity.Tag, error)
	updateFn    func(ctx context.Context, tag entity.Tag) (entity.Tag, error)
	deleteFn    func(ctx context.Context, id uint) error
}

func (f *fakeTagRepo) Create(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	return f.createFn(ctx, tag)
}
func (f *fakeTagRepo) GetAll(ctx context.Context) ([]entity.Tag, error) { return f.getAllFn(ctx) }
func (f *fakeTagRepo) GetByID(ctx context.Context, id uint) (entity.Tag, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakeTagRepo) GetByName(ctx context.Context, name string) (entity.Tag, error) {
	return f.getByNameFn(ctx, name)
}
func (f *fakeTagRepo) Update(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	return f.updateFn(ctx, tag)
}
func (f *fakeTagRepo) Delete(ctx context.Context, id uint) error { return f.deleteFn(ctx, id) }

func TestTagService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("empty name rejected", func(t *testing.T) {
		svc := NewTagService(&fakeTagRepo{})
		_, err := svc.Create(ctx, entity.Tag{Name: "   "})
		if !errors.Is(err, core.ErrInvalidInput) {
			t.Fatalf("want ErrInvalidInput, got %v", err)
		}
	})

	t.Run("duplicate rejected", func(t *testing.T) {
		repo := &fakeTagRepo{
			getByNameFn: func(ctx context.Context, name string) (entity.Tag, error) {
				return entity.Tag{ID: 7, Name: name}, nil
			},
		}
		svc := NewTagService(repo)
		_, err := svc.Create(ctx, entity.Tag{Name: "go"})
		if !errors.Is(err, core.ErrDuplicate) {
			t.Fatalf("want ErrDuplicate, got %v", err)
		}
	})

	t.Run("creates when not exists", func(t *testing.T) {
		called := false
		repo := &fakeTagRepo{
			getByNameFn: func(ctx context.Context, name string) (entity.Tag, error) {
				return entity.Tag{}, core.ErrNotFound
			},
			createFn: func(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
				called = true
				tag.ID = 1
				return tag, nil
			},
		}
		svc := NewTagService(repo)
		got, err := svc.Create(ctx, entity.Tag{Name: "  go "})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if !called {
			t.Fatal("repo.Create was not called")
		}
		if got.Name != "go" {
			t.Fatalf("name not trimmed: %q", got.Name)
		}
		if got.ID != 1 {
			t.Fatalf("want ID 1, got %d", got.ID)
		}
	})

	t.Run("lookup repo error normalized", func(t *testing.T) {
		repo := &fakeTagRepo{
			getByNameFn: func(ctx context.Context, name string) (entity.Tag, error) {
				return entity.Tag{}, errors.New("db boom")
			},
		}
		svc := NewTagService(repo)
		_, err := svc.Create(ctx, entity.Tag{Name: "go"})
		if !errors.Is(err, core.ErrInternalError) {
			t.Fatalf("want ErrInternalError, got %v", err)
		}
	})
}

func TestTagService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("zero id invalid", func(t *testing.T) {
		svc := NewTagService(&fakeTagRepo{})
		_, err := svc.GetByID(ctx, 0)
		if !errors.Is(err, core.ErrInvalidInput) {
			t.Fatalf("want ErrInvalidInput, got %v", err)
		}
	})

	t.Run("not found propagated", func(t *testing.T) {
		repo := &fakeTagRepo{
			getByIDFn: func(ctx context.Context, id uint) (entity.Tag, error) {
				return entity.Tag{}, core.ErrNotFound
			},
		}
		svc := NewTagService(repo)
		_, err := svc.GetByID(ctx, 9)
		if !errors.Is(err, core.ErrNotFound) {
			t.Fatalf("want ErrNotFound, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeTagRepo{
			getByIDFn: func(ctx context.Context, id uint) (entity.Tag, error) {
				return entity.Tag{ID: id, Name: "go"}, nil
			},
		}
		got, err := NewTagService(repo).GetByID(ctx, 3)
		if err != nil || got.ID != 3 {
			t.Fatalf("unexpected: %+v %v", got, err)
		}
	})
}

func TestTagService_GetByName(t *testing.T) {
	ctx := context.Background()
	svc := NewTagService(&fakeTagRepo{})
	if _, err := svc.GetByName(ctx, "  "); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}

	repo := &fakeTagRepo{getByNameFn: func(ctx context.Context, name string) (entity.Tag, error) {
		return entity.Tag{Name: name}, nil
	}}
	got, err := NewTagService(repo).GetByName(ctx, " go ")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "go" {
		t.Fatalf("name not trimmed before lookup: %q", got.Name)
	}
}

func TestTagService_Update(t *testing.T) {
	ctx := context.Background()

	if _, err := NewTagService(&fakeTagRepo{}).Update(ctx, entity.Tag{}); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput for zero id, got %v", err)
	}

	repo := &fakeTagRepo{updateFn: func(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
		return tag, nil
	}}
	got, err := NewTagService(repo).Update(ctx, entity.Tag{ID: 1, Name: "  golang "})
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "golang" {
		t.Fatalf("name not trimmed: %q", got.Name)
	}
}

func TestTagService_Delete(t *testing.T) {
	ctx := context.Background()
	if err := NewTagService(&fakeTagRepo{}).Delete(ctx, 0); !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}

	var deleted uint
	repo := &fakeTagRepo{deleteFn: func(ctx context.Context, id uint) error {
		deleted = id
		return nil
	}}
	if err := NewTagService(repo).Delete(ctx, 5); err != nil {
		t.Fatal(err)
	}
	if deleted != 5 {
		t.Fatalf("want id 5, got %d", deleted)
	}
}

func TestTagService_GetAll(t *testing.T) {
	ctx := context.Background()
	repo := &fakeTagRepo{getAllFn: func(ctx context.Context) ([]entity.Tag, error) {
		return []entity.Tag{{ID: 1}, {ID: 2}}, nil
	}}
	got, err := NewTagService(repo).GetAll(ctx)
	if err != nil || len(got) != 2 {
		t.Fatalf("unexpected: %+v %v", got, err)
	}
}
