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

type fakeAuthService struct{}

func (f fakeAuthService) Register(_ context.Context, _ service.RegisterInput) (service.RegisteredUser, service.TokenPair, error) {
	return service.RegisteredUser{ID: "u1", Email: "test@example.com", Role: models.RoleAuthor}, service.TokenPair{AccessToken: "a", RefreshToken: "r", TokenType: "Bearer", AccessExpiresAt: time.Now(), RefreshExpiresAt: time.Now().Add(time.Hour)}, nil
}

func (f fakeAuthService) Login(_ context.Context, _ service.LoginInput) (service.RegisteredUser, service.TokenPair, error) {
	return service.RegisteredUser{}, service.TokenPair{}, service.ErrInvalidCredentials
}

func (f fakeAuthService) Refresh(_ context.Context, _ service.RefreshInput) (service.TokenPair, error) {
	return service.TokenPair{}, nil
}

func (f fakeAuthService) Logout(_ context.Context, _ service.LogoutInput) error {
	return nil
}

func (f fakeAuthService) RequestPasswordReset(_ context.Context, _ service.RequestResetInput) error {
	return nil
}

func (f fakeAuthService) ConfirmPasswordReset(_ context.Context, _ service.ConfirmResetInput) error {
	return nil
}

func TestAuthRegisterSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(fakeAuthService{})
	r.POST("/register", h.Register)

	body := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}
	if _, ok := payload["tokens"]; !ok {
		t.Fatalf("expected tokens in response")
	}
}

func TestAuthLoginInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(fakeAuthService{})
	r.POST("/login", h.Login)

	body := `{"email":"test@example.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}
