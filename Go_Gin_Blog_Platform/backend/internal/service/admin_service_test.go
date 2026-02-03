package service

import (
	"context"
	"errors"
	"testing"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/models"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/repository"
)

type fakeUserRepo struct {
	users         []models.User
	count         int64
	updateRoleErr error
}

func (f *fakeUserRepo) Create(context.Context, *models.User) error { return nil }
func (f *fakeUserRepo) GetByEmail(context.Context, string) (*models.User, error) {
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) GetByID(_ context.Context, id string) (*models.User, error) {
	for i := range f.users {
		if f.users[i].ID == id {
			copy := f.users[i]
			return &copy, nil
		}
	}
	return nil, repository.ErrNotFound
}
func (f *fakeUserRepo) List(_ context.Context, _ int, _ int) ([]models.User, error) {
	return f.users, nil
}
func (f *fakeUserRepo) Count(context.Context) (int64, error) { return f.count, nil }
func (f *fakeUserRepo) UpdateRole(_ context.Context, id string, role models.Role) error {
	if f.updateRoleErr != nil {
		return f.updateRoleErr
	}
	for i := range f.users {
		if f.users[i].ID == id {
			f.users[i].Role = role
			return nil
		}
	}
	return repository.ErrNotFound
}
func (f *fakeUserRepo) UpdatePasswordHash(context.Context, string, string) error { return nil }

func TestAdminServiceUpdateUserRoleValidation(t *testing.T) {
	svc := NewAdminService(&fakeUserRepo{})

	_, err := svc.UpdateUserRole(context.Background(), "u1", "invalid-role")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestAdminServiceUpdateUserRoleNotFound(t *testing.T) {
	svc := NewAdminService(&fakeUserRepo{updateRoleErr: repository.ErrNotFound})

	_, err := svc.UpdateUserRole(context.Background(), "missing", "reader")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestAdminServiceListUsersPagination(t *testing.T) {
	repo := &fakeUserRepo{
		users: []models.User{
			{ID: "u1", Email: "a@example.com", Role: models.RoleAuthor},
			{ID: "u2", Email: "b@example.com", Role: models.RoleReader},
		},
		count: 12,
	}
	svc := NewAdminService(repo)

	users, page, err := svc.ListUsers(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("expected list users success, got %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users in current fake response, got %d", len(users))
	}
	if page.Total != 12 || page.TotalPages != 2 {
		t.Fatalf("unexpected pagination metadata: %+v", page)
	}
}
