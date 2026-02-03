package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID = "auth_user_id"
	ContextKeyRole   = "auth_role"
)

type AccessTokenVerifier interface {
	ParseAccessToken(token string) (*auth.AccessClaims, error)
}

func AuthRequired(verifier AccessTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		if verifier == nil {
			writeError(c, http.StatusServiceUnavailable, "auth_unavailable", "Authentication is not configured", nil)
			c.Abort()
			return
		}

		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
			writeError(c, http.StatusUnauthorized, "missing_token", "Authorization token is required", nil)
			c.Abort()
			return
		}

		rawToken := strings.TrimSpace(header[len("Bearer "):])
		if rawToken == "" {
			writeError(c, http.StatusUnauthorized, "missing_token", "Authorization token is required", nil)
			c.Abort()
			return
		}

		claims, err := verifier.ParseAccessToken(rawToken)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) {
				writeError(c, http.StatusUnauthorized, "invalid_token", "Access token is invalid or expired", nil)
			} else {
				writeError(c, http.StatusUnauthorized, "invalid_token", "Access token is invalid", nil)
			}
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.Subject)
		c.Set(ContextKeyRole, claims.Role)
		c.Next()
	}
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[strings.ToLower(strings.TrimSpace(role))] = struct{}{}
	}

	return func(c *gin.Context) {
		value, exists := c.Get(ContextKeyRole)
		if !exists {
			writeError(c, http.StatusForbidden, "forbidden", "Insufficient permissions", nil)
			c.Abort()
			return
		}

		role, ok := value.(string)
		if !ok {
			writeError(c, http.StatusForbidden, "forbidden", "Insufficient permissions", nil)
			c.Abort()
			return
		}

		if _, ok := allowed[strings.ToLower(role)]; !ok {
			writeError(c, http.StatusForbidden, "forbidden", "Insufficient permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
