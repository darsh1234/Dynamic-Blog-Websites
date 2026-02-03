package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type fakeVerifier struct {
	claims *auth.AccessClaims
	err    error
}

func (f fakeVerifier) ParseAccessToken(_ string) (*auth.AccessClaims, error) {
	return f.claims, f.err
}

func TestAuthRequiredRejectsMissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", AuthRequired(fakeVerifier{}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestRequireRolesAcceptsAllowedRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	claims := &auth.AccessClaims{Role: "author"}
	r.GET("/protected", AuthRequired(fakeVerifier{claims: claims}), RequireRoles("author"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer test")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
