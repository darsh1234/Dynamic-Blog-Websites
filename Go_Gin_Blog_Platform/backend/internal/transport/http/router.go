package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type RouterDependencies struct {
	HealthChecker       HealthChecker
	HealthCheckTimeout  time.Duration
	AuthHandler         *AuthHandler
	PostHandler         *PostHandler
	AdminHandler        *AdminHandler
	AccessTokenVerifier AccessTokenVerifier
}

func NewRouter(logger *slog.Logger, deps RouterDependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	api := router.Group("/api/v1")
	{
		api.GET("/healthz", func(c *gin.Context) {
			if deps.HealthChecker != nil {
				ctx, cancel := context.WithTimeout(c.Request.Context(), deps.HealthCheckTimeout)
				defer cancel()

				if err := deps.HealthChecker.Ping(ctx); err != nil {
					c.JSON(http.StatusServiceUnavailable, gin.H{
						"error": gin.H{
							"code":    "database_unavailable",
							"message": "Database health check failed",
							"details": gin.H{"reason": err.Error()},
						},
					})
					return
				}
			}

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		auth := api.Group("/auth")
		{
			if deps.AuthHandler != nil {
				auth.POST("/register", deps.AuthHandler.Register)
				auth.POST("/login", deps.AuthHandler.Login)
				auth.POST("/refresh", deps.AuthHandler.Refresh)
				auth.POST("/logout", deps.AuthHandler.Logout)
				auth.POST("/password-reset/request", deps.AuthHandler.RequestPasswordReset)
				auth.POST("/password-reset/confirm", deps.AuthHandler.ConfirmPasswordReset)
			} else {
				auth.POST("/register", notImplemented(canonicalRoute("POST /auth/register")))
				auth.POST("/login", notImplemented(canonicalRoute("POST /auth/login")))
				auth.POST("/refresh", notImplemented(canonicalRoute("POST /auth/refresh")))
				auth.POST("/logout", notImplemented(canonicalRoute("POST /auth/logout")))
				auth.POST("/password-reset/request", notImplemented(canonicalRoute("POST /auth/password-reset/request")))
				auth.POST("/password-reset/confirm", notImplemented(canonicalRoute("POST /auth/password-reset/confirm")))
			}
		}

		if deps.PostHandler != nil {
			api.GET("/posts", deps.PostHandler.List)
			api.GET("/posts/:id", deps.PostHandler.GetByID)
		} else {
			api.GET("/posts", notImplemented(canonicalRoute("GET /posts")))
			api.GET("/posts/:id", notImplemented(canonicalRoute("GET /posts/:id")))
		}

		postsWrite := api.Group("/posts")
		postsWrite.Use(AuthRequired(deps.AccessTokenVerifier), RequireRoles("author", "admin"))
		{
			if deps.PostHandler != nil {
				postsWrite.POST("", deps.PostHandler.Create)
				postsWrite.PATCH("/:id", deps.PostHandler.Update)
				postsWrite.DELETE("/:id", deps.PostHandler.Delete)
			} else {
				postsWrite.POST("", notImplemented(canonicalRoute("POST /posts")))
				postsWrite.PATCH("/:id", notImplemented(canonicalRoute("PATCH /posts/:id")))
				postsWrite.DELETE("/:id", notImplemented(canonicalRoute("DELETE /posts/:id")))
			}
		}

		admin := api.Group("/admin")
		admin.Use(AuthRequired(deps.AccessTokenVerifier), RequireRoles("admin"))
		{
			if deps.AdminHandler != nil {
				admin.GET("/users", deps.AdminHandler.ListUsers)
				admin.PATCH("/users/:id/role", deps.AdminHandler.UpdateRole)
			} else {
				admin.GET("/users", notImplemented(canonicalRoute("GET /admin/users")))
				admin.PATCH("/users/:id/role", notImplemented(canonicalRoute("PATCH /admin/users/:id/role")))
			}
		}
	}

	logger.Info("router initialized", "base_path", "/api/v1")
	return router
}

func canonicalRoute(route string) string {
	return route
}

func notImplemented(route string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": gin.H{
				"code":    "not_implemented",
				"message": "Endpoint is not available in this build",
				"details": gin.H{"route": route},
			},
		})
	}
}
