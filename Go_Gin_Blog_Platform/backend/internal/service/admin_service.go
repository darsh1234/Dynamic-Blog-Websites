package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/models"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/repository"
)

type AdminService struct {
	users repository.UserRepository
}

func NewAdminService(users repository.UserRepository) *AdminService {
	return &AdminService{users: users}
}

type UserSummary struct {
	ID    string      `json:"id"`
	Email string      `json:"email"`
	Role  models.Role `json:"role"`
}

func (s *AdminService) ListUsers(ctx context.Context, page, limit int) ([]UserSummary, Pagination, error) {
	page, limit = normalizePagination(page, limit)
	offset := (page - 1) * limit

	users, err := s.users.List(ctx, limit, offset)
	if err != nil {
		return nil, Pagination{}, fmt.Errorf("list users: %w", err)
	}

	total, err := s.users.Count(ctx)
	if err != nil {
		return nil, Pagination{}, fmt.Errorf("count users: %w", err)
	}

	summaries := make([]UserSummary, 0, len(users))
	for _, user := range users {
		summaries = append(summaries, UserSummary{ID: user.ID, Email: user.Email, Role: user.Role})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return summaries, Pagination{Page: page, Limit: limit, Total: total, TotalPages: totalPages}, nil
}

func (s *AdminService) UpdateUserRole(ctx context.Context, userID, role string) (UserSummary, error) {
	normalizedRole, err := normalizeRole(role)
	if err != nil {
		return UserSummary{}, err
	}

	if strings.TrimSpace(userID) == "" {
		return UserSummary{}, fmt.Errorf("user id is required: %w", ErrValidation)
	}

	if err := s.users.UpdateRole(ctx, userID, normalizedRole); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return UserSummary{}, ErrUserNotFound
		}
		return UserSummary{}, fmt.Errorf("update user role: %w", err)
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return UserSummary{}, ErrUserNotFound
		}
		return UserSummary{}, fmt.Errorf("load updated user: %w", err)
	}

	return UserSummary{ID: user.ID, Email: user.Email, Role: user.Role}, nil
}

func normalizeRole(role string) (models.Role, error) {
	value := strings.ToLower(strings.TrimSpace(role))
	parsed := models.Role(value)
	switch parsed {
	case models.RoleAdmin, models.RoleAuthor, models.RoleReader:
		return parsed, nil
	default:
		return "", fmt.Errorf("role must be admin, author, or reader: %w", ErrValidation)
	}
}
