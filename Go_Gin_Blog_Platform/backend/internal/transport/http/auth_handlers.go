package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(ctx context.Context, input service.RegisterInput) (service.RegisteredUser, service.TokenPair, error)
	Login(ctx context.Context, input service.LoginInput) (service.RegisteredUser, service.TokenPair, error)
	Refresh(ctx context.Context, input service.RefreshInput) (service.TokenPair, error)
	Logout(ctx context.Context, input service.LogoutInput) error
	RequestPasswordReset(ctx context.Context, input service.RequestResetInput) error
	ConfirmPasswordReset(ctx context.Context, input service.ConfirmResetInput) error
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type requestResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type confirmResetRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	user, tokenPair, err := h.authService.Register(c.Request.Context(), service.RegisterInput{Email: req.Email, Password: req.Password})
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user, "tokens": tokenPair})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	user, tokenPair, err := h.authService.Login(c.Request.Context(), service.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "tokens": tokenPair})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	tokenPair, err := h.authService.Refresh(c.Request.Context(), service.RefreshInput{RefreshToken: req.RefreshToken})
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokenPair})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req logoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), service.LogoutInput{RefreshToken: req.RefreshToken}); err != nil {
		handleAuthError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req requestResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	if err := h.authService.RequestPasswordReset(c.Request.Context(), service.RequestResetInput{Email: req.Email}); err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link will be sent"})
}

func (h *AuthHandler) ConfirmPasswordReset(c *gin.Context) {
	var req confirmResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeValidationError(c, err)
		return
	}

	if err := h.authService.ConfirmPasswordReset(c.Request.Context(), service.ConfirmResetInput{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}); err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		writeError(c, http.StatusBadRequest, "validation_error", "Request validation failed", gin.H{"reason": err.Error()})
	case errors.Is(err, service.ErrEmailAlreadyUsed):
		writeError(c, http.StatusConflict, "email_already_exists", "Email is already registered", nil)
	case errors.Is(err, service.ErrInvalidCredentials):
		writeError(c, http.StatusUnauthorized, "invalid_credentials", "Email or password is incorrect", nil)
	case errors.Is(err, service.ErrInvalidToken):
		writeError(c, http.StatusUnauthorized, "invalid_token", "Token is invalid or expired", nil)
	default:
		writeError(c, http.StatusInternalServerError, "internal_error", "Unexpected server error", nil)
	}
}
