package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PostService interface {
	Create(ctx context.Context, input service.CreatePostInput) (service.PostItem, error)
	GetByID(ctx context.Context, postID string) (service.PostItem, error)
	List(ctx context.Context, input service.ListPostsInput) ([]service.PostItem, service.Pagination, error)
	Update(ctx context.Context, input service.UpdatePostInput) (service.PostItem, error)
	Delete(ctx context.Context, input service.DeletePostInput) error
}

type PostHandler struct {
	postService PostService
}

func NewPostHandler(postService PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

type createPostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Status  string `json:"status"`
}

type updatePostRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  *string `json:"status"`
}

func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	posts, pagination, err := h.postService.List(c.Request.Context(), service.ListPostsInput{Page: page, Limit: limit})
	if err != nil {
		handlePostError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts, "meta": pagination})
}

func (h *PostHandler) GetByID(c *gin.Context) {
	post, err := h.postService.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handlePostError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": post})
}

func (h *PostHandler) Create(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	actorID, actorRole, ok := currentUserFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized", "Authentication context is missing", nil)
		return
	}

	post, err := h.postService.Create(c.Request.Context(), service.CreatePostInput{
		ActorID:   actorID,
		ActorRole: actorRole,
		Title:     req.Title,
		Content:   req.Content,
		Status:    req.Status,
	})
	if err != nil {
		handlePostError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": post})
}

func (h *PostHandler) Update(c *gin.Context) {
	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	actorID, actorRole, ok := currentUserFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized", "Authentication context is missing", nil)
		return
	}

	post, err := h.postService.Update(c.Request.Context(), service.UpdatePostInput{
		PostID:    c.Param("id"),
		ActorID:   actorID,
		ActorRole: actorRole,
		Title:     req.Title,
		Content:   req.Content,
		Status:    req.Status,
	})
	if err != nil {
		handlePostError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": post})
}

func (h *PostHandler) Delete(c *gin.Context) {
	actorID, actorRole, ok := currentUserFromContext(c)
	if !ok {
		writeError(c, http.StatusUnauthorized, "unauthorized", "Authentication context is missing", nil)
		return
	}

	if err := h.postService.Delete(c.Request.Context(), service.DeletePostInput{
		PostID:    c.Param("id"),
		ActorID:   actorID,
		ActorRole: actorRole,
	}); err != nil {
		handlePostError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handlePostError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		writeError(c, http.StatusBadRequest, "validation_error", "Request validation failed", gin.H{"reason": err.Error()})
	case errors.Is(err, service.ErrForbidden):
		writeError(c, http.StatusForbidden, "forbidden", "Insufficient permissions", nil)
	case errors.Is(err, service.ErrPostNotFound):
		writeError(c, http.StatusNotFound, "post_not_found", "Post was not found", nil)
	default:
		writeError(c, http.StatusInternalServerError, "internal_error", "Unexpected server error", nil)
	}
}
