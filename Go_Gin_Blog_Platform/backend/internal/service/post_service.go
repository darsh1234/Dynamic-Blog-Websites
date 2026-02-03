package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/models"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/repository"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("forbidden")
)

type CreatePostInput struct {
	ActorID   string
	ActorRole string
	Title     string
	Content   string
	Status    string
}

type UpdatePostInput struct {
	PostID    string
	ActorID   string
	ActorRole string
	Title     *string
	Content   *string
	Status    *string
}

type DeletePostInput struct {
	PostID    string
	ActorID   string
	ActorRole string
}

type ListPostsInput struct {
	Page  int
	Limit int
}

type PostItem struct {
	ID        string            `json:"id"`
	AuthorID  string            `json:"author_id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Status    models.PostStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PostService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(ctx context.Context, input CreatePostInput) (PostItem, error) {
	if !canWritePosts(input.ActorRole) {
		return PostItem{}, ErrForbidden
	}

	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)
	if title == "" || content == "" {
		return PostItem{}, fmt.Errorf("title and content are required: %w", ErrValidation)
	}

	status, err := normalizeStatus(input.Status)
	if err != nil {
		return PostItem{}, err
	}

	post := &models.Post{
		AuthorID: input.ActorID,
		Title:    title,
		Content:  content,
		Status:   status,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return PostItem{}, fmt.Errorf("create post: %w", err)
	}

	return toPostItem(*post), nil
}

func (s *PostService) GetByID(ctx context.Context, postID string) (PostItem, error) {
	post, err := s.repo.GetByID(ctx, strings.TrimSpace(postID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return PostItem{}, ErrPostNotFound
		}
		return PostItem{}, fmt.Errorf("get post by id: %w", err)
	}
	return toPostItem(*post), nil
}

func (s *PostService) List(ctx context.Context, input ListPostsInput) ([]PostItem, Pagination, error) {
	page, limit := normalizePagination(input.Page, input.Limit)
	offset := (page - 1) * limit

	posts, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, Pagination{}, fmt.Errorf("list posts: %w", err)
	}

	items := make([]PostItem, 0, len(posts))
	for _, post := range posts {
		items = append(items, toPostItem(post))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return items, Pagination{Page: page, Limit: limit, Total: total, TotalPages: totalPages}, nil
}

func (s *PostService) Update(ctx context.Context, input UpdatePostInput) (PostItem, error) {
	if !canWritePosts(input.ActorRole) {
		return PostItem{}, ErrForbidden
	}

	post, err := s.repo.GetByID(ctx, strings.TrimSpace(input.PostID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return PostItem{}, ErrPostNotFound
		}
		return PostItem{}, fmt.Errorf("get existing post: %w", err)
	}

	if !canModifyPost(input.ActorRole, input.ActorID, post.AuthorID) {
		return PostItem{}, ErrForbidden
	}

	updates := map[string]any{}
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return PostItem{}, fmt.Errorf("title cannot be empty: %w", ErrValidation)
		}
		updates["title"] = title
	}
	if input.Content != nil {
		content := strings.TrimSpace(*input.Content)
		if content == "" {
			return PostItem{}, fmt.Errorf("content cannot be empty: %w", ErrValidation)
		}
		updates["content"] = content
	}
	if input.Status != nil {
		status, err := normalizeStatus(*input.Status)
		if err != nil {
			return PostItem{}, err
		}
		updates["status"] = status
	}

	if len(updates) == 0 {
		return PostItem{}, fmt.Errorf("no update fields provided: %w", ErrValidation)
	}

	if err := s.repo.Update(ctx, post.ID, updates); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return PostItem{}, ErrPostNotFound
		}
		return PostItem{}, fmt.Errorf("update post: %w", err)
	}

	updated, err := s.repo.GetByID(ctx, post.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return PostItem{}, ErrPostNotFound
		}
		return PostItem{}, fmt.Errorf("reload updated post: %w", err)
	}
	return toPostItem(*updated), nil
}

func (s *PostService) Delete(ctx context.Context, input DeletePostInput) error {
	if !canWritePosts(input.ActorRole) {
		return ErrForbidden
	}

	post, err := s.repo.GetByID(ctx, strings.TrimSpace(input.PostID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrPostNotFound
		}
		return fmt.Errorf("get existing post: %w", err)
	}

	if !canModifyPost(input.ActorRole, input.ActorID, post.AuthorID) {
		return ErrForbidden
	}

	if err := s.repo.Delete(ctx, post.ID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrPostNotFound
		}
		return fmt.Errorf("delete post: %w", err)
	}

	return nil
}

func toPostItem(post models.Post) PostItem {
	return PostItem{
		ID:        post.ID,
		AuthorID:  post.AuthorID,
		Title:     post.Title,
		Content:   post.Content,
		Status:    post.Status,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func normalizePagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return page, limit
}

func normalizeStatus(status string) (models.PostStatus, error) {
	value := strings.ToLower(strings.TrimSpace(status))
	if value == "" {
		return models.PostStatusPublished, nil
	}

	s := models.PostStatus(value)
	if s != models.PostStatusDraft && s != models.PostStatusPublished {
		return "", fmt.Errorf("status must be draft or published: %w", ErrValidation)
	}
	return s, nil
}

func canWritePosts(role string) bool {
	role = strings.ToLower(strings.TrimSpace(role))
	return role == string(models.RoleAdmin) || role == string(models.RoleAuthor)
}

func canModifyPost(actorRole, actorID, authorID string) bool {
	if strings.EqualFold(strings.TrimSpace(actorRole), string(models.RoleAdmin)) {
		return true
	}
	return strings.TrimSpace(actorID) != "" && strings.TrimSpace(actorID) == strings.TrimSpace(authorID)
}
