package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/models"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/repository"
)

type fakePostRepo struct {
	post       models.Post
	listTotal  int64
	listPosts  []models.Post
	lastLimit  int
	lastOffset int
}

func (f *fakePostRepo) Create(_ context.Context, post *models.Post) error {
	f.post = *post
	if f.post.ID == "" {
		f.post.ID = "post-1"
	}
	f.post.CreatedAt = time.Now().UTC()
	f.post.UpdatedAt = f.post.CreatedAt
	*post = f.post
	return nil
}

func (f *fakePostRepo) GetByID(_ context.Context, id string) (*models.Post, error) {
	if f.post.ID == "" || id != f.post.ID {
		return nil, repository.ErrNotFound
	}
	copy := f.post
	return &copy, nil
}

func (f *fakePostRepo) List(_ context.Context, limit, offset int) ([]models.Post, int64, error) {
	f.lastLimit = limit
	f.lastOffset = offset
	if f.listPosts != nil {
		return f.listPosts, f.listTotal, nil
	}
	return []models.Post{f.post}, 1, nil
}

func (f *fakePostRepo) Update(_ context.Context, id string, updates map[string]any) error {
	if f.post.ID == "" || id != f.post.ID {
		return repository.ErrNotFound
	}
	if v, ok := updates["title"].(string); ok {
		f.post.Title = v
	}
	if v, ok := updates["content"].(string); ok {
		f.post.Content = v
	}
	if v, ok := updates["status"].(models.PostStatus); ok {
		f.post.Status = v
	}
	f.post.UpdatedAt = time.Now().UTC()
	return nil
}

func (f *fakePostRepo) Delete(_ context.Context, id string) error {
	if f.post.ID == "" || id != f.post.ID {
		return repository.ErrNotFound
	}
	f.post = models.Post{}
	return nil
}

func TestPostServiceCreateRejectsReader(t *testing.T) {
	svc := NewPostService(&fakePostRepo{})

	_, err := svc.Create(context.Background(), CreatePostInput{
		ActorID:   "u1",
		ActorRole: "reader",
		Title:     "Hello",
		Content:   "World",
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestPostServiceUpdateEnforcesOwnership(t *testing.T) {
	repo := &fakePostRepo{post: models.Post{ID: "p1", AuthorID: "owner", Title: "Old", Content: "Text", Status: models.PostStatusPublished}}
	svc := NewPostService(repo)

	title := "New"
	_, err := svc.Update(context.Background(), UpdatePostInput{
		PostID:    "p1",
		ActorID:   "different-user",
		ActorRole: "author",
		Title:     &title,
	})

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden for non-owner author update, got %v", err)
	}
}

func TestPostServiceListAppliesPaginationDefaults(t *testing.T) {
	repo := &fakePostRepo{listPosts: []models.Post{}, listTotal: 120}
	svc := NewPostService(repo)

	_, meta, err := svc.List(context.Background(), ListPostsInput{Page: 0, Limit: 500})
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}

	if repo.lastLimit != 50 {
		t.Fatalf("expected capped limit 50, got %d", repo.lastLimit)
	}
	if repo.lastOffset != 0 {
		t.Fatalf("expected first-page offset 0, got %d", repo.lastOffset)
	}
	if meta.Page != 1 || meta.Limit != 50 || meta.TotalPages != 3 {
		t.Fatalf("unexpected pagination metadata: %+v", meta)
	}
}
