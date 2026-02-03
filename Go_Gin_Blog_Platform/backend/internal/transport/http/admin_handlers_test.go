package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/models"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type fakeAdminService struct{}

func (f fakeAdminService) ListUsers(_ context.Context, _, _ int) ([]service.UserSummary, service.Pagination, error) {
	return []service.UserSummary{{ID: "u1", Email: "a@example.com", Role: models.RoleAuthor}}, service.Pagination{Page: 1, Limit: 10, Total: 1, TotalPages: 1}, nil
}

func (f fakeAdminService) UpdateUserRole(_ context.Context, userID, role string) (service.UserSummary, error) {
	return service.UserSummary{ID: userID, Email: "a@example.com", Role: models.Role(role)}, nil
}

func TestAdminListUsersSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAdminHandler(fakeAdminService{})
	r.GET("/admin/users", h.ListUsers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected json response: %v", err)
	}
	if _, ok := payload["data"]; !ok {
		t.Fatalf("expected data in response")
	}
}

func TestAdminUpdateRoleValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAdminHandler(fakeAdminService{})
	r.PATCH("/admin/users/:id/role", h.UpdateRole)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/admin/users/u1/role", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
