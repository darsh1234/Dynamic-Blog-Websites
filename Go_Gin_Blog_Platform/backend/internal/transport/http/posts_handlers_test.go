package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/models"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type fakePostService struct{}

func (f fakePostService) Create(_ context.Context, input service.CreatePostInput) (service.PostItem, error) {
	return service.PostItem{ID: "p1", AuthorID: input.ActorID, Title: input.Title, Content: input.Content, Status: models.PostStatusPublished}, nil
}

func (f fakePostService) GetByID(_ context.Context, postID string) (service.PostItem, error) {
	return service.PostItem{ID: postID, AuthorID: "u1", Title: "Hello", Content: "World", Status: models.PostStatusPublished}, nil
}

func (f fakePostService) List(_ context.Context, _ service.ListPostsInput) ([]service.PostItem, service.Pagination, error) {
	return []service.PostItem{{ID: "p1", Title: "A", Content: "B", Status: models.PostStatusPublished}}, service.Pagination{Page: 1, Limit: 10, Total: 1, TotalPages: 1}, nil
}

func (f fakePostService) Update(_ context.Context, input service.UpdatePostInput) (service.PostItem, error) {
	return service.PostItem{ID: input.PostID, AuthorID: input.ActorID, Title: "Updated", Content: "Updated", Status: models.PostStatusPublished, UpdatedAt: time.Now()}, nil
}

func (f fakePostService) Delete(_ context.Context, _ service.DeletePostInput) error {
	return nil
}

func TestPostsListSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewPostHandler(fakePostService{})
	r.GET("/posts", h.List)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/posts?page=1&limit=10", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected json response: %v", err)
	}
	if _, ok := payload["meta"]; !ok {
		t.Fatalf("expected meta object in response")
	}
}

func TestPostsCreateUnauthorizedWhenContextMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewPostHandler(fakePostService{})
	r.POST("/posts", h.Create)

	body := `{"title":"A","content":"B"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}
