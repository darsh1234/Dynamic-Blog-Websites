package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminService interface {
	ListUsers(ctx context.Context, page, limit int) ([]service.UserSummary, service.Pagination, error)
	UpdateUserRole(ctx context.Context, userID, role string) (service.UserSummary, error)
}

type AdminHandler struct {
	adminService AdminService
}

func NewAdminHandler(adminService AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

type updateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	users, pagination, err := h.adminService.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		handleAdminError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users, "meta": pagination})
}

func (h *AdminHandler) UpdateRole(c *gin.Context) {
	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	user, err := h.adminService.UpdateUserRole(c.Request.Context(), c.Param("id"), req.Role)
	if err != nil {
		handleAdminError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func handleAdminError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		writeError(c, http.StatusBadRequest, "validation_error", "Request validation failed", gin.H{"reason": err.Error()})
	case errors.Is(err, service.ErrUserNotFound):
		writeError(c, http.StatusNotFound, "user_not_found", "User was not found", nil)
	default:
		writeError(c, http.StatusInternalServerError, "internal_error", "Unexpected server error", nil)
	}
}
