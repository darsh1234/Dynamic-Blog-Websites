package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/auth"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/email"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/models"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/repository"
)

var (
	ErrValidation         = errors.New("validation error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyUsed   = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
)

type RegisterInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type RefreshInput struct {
	RefreshToken string
}

type LogoutInput struct {
	RefreshToken string
}

type RequestResetInput struct {
	Email string
}

type ConfirmResetInput struct {
	Token       string
	NewPassword string
}

type RegisteredUser struct {
	ID    string      `json:"id"`
	Email string      `json:"email"`
	Role  models.Role `json:"role"`
}

type TokenPair struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenType        string    `json:"token_type"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

type AuthService struct {
	logger           *slog.Logger
	users            repository.UserRepository
	refreshTokens    repository.RefreshTokenRepository
	resetTokens      repository.PasswordResetTokenRepository
	tokenManager     *auth.TokenManager
	emailSender      email.Sender
	defaultRole      models.Role
	passwordResetTTL time.Duration
	frontendBaseURL  string
}

func NewAuthService(
	logger *slog.Logger,
	users repository.UserRepository,
	refreshTokens repository.RefreshTokenRepository,
	resetTokens repository.PasswordResetTokenRepository,
	tokenManager *auth.TokenManager,
	emailSender email.Sender,
	passwordResetTTL time.Duration,
	frontendBaseURL string,
) *AuthService {
	return &AuthService{
		logger:           logger,
		users:            users,
		refreshTokens:    refreshTokens,
		resetTokens:      resetTokens,
		tokenManager:     tokenManager,
		emailSender:      emailSender,
		defaultRole:      models.RoleAuthor,
		passwordResetTTL: passwordResetTTL,
		frontendBaseURL:  strings.TrimRight(frontendBaseURL, "/"),
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (RegisteredUser, TokenPair, error) {
	email := normalizeEmail(input.Email)
	if !isValidEmail(email) || len(input.Password) < 8 {
		return RegisteredUser{}, TokenPair{}, fmt.Errorf("register input invalid: %w", ErrValidation)
	}

	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return RegisteredUser{}, TokenPair{}, ErrEmailAlreadyUsed
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return RegisteredUser{}, TokenPair{}, fmt.Errorf("check existing user: %w", err)
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return RegisteredUser{}, TokenPair{}, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         s.defaultRole,
	}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return RegisteredUser{}, TokenPair{}, ErrEmailAlreadyUsed
		}
		return RegisteredUser{}, TokenPair{}, fmt.Errorf("create user: %w", err)
	}

	pair, err := s.issueTokenPair(ctx, user.ID, string(user.Role))
	if err != nil {
		return RegisteredUser{}, TokenPair{}, err
	}

	return RegisteredUser{ID: user.ID, Email: user.Email, Role: user.Role}, pair, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (RegisteredUser, TokenPair, error) {
	email := normalizeEmail(input.Email)
	if !isValidEmail(email) || strings.TrimSpace(input.Password) == "" {
		return RegisteredUser{}, TokenPair{}, fmt.Errorf("login input invalid: %w", ErrValidation)
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return RegisteredUser{}, TokenPair{}, ErrInvalidCredentials
		}
		return RegisteredUser{}, TokenPair{}, fmt.Errorf("fetch user by email: %w", err)
	}

	if !auth.VerifyPassword(user.PasswordHash, input.Password) {
		return RegisteredUser{}, TokenPair{}, ErrInvalidCredentials
	}

	pair, err := s.issueTokenPair(ctx, user.ID, string(user.Role))
	if err != nil {
		return RegisteredUser{}, TokenPair{}, err
	}

	return RegisteredUser{ID: user.ID, Email: user.Email, Role: user.Role}, pair, nil
}

func (s *AuthService) Refresh(ctx context.Context, input RefreshInput) (TokenPair, error) {
	rawToken := strings.TrimSpace(input.RefreshToken)
	if rawToken == "" {
		return TokenPair{}, fmt.Errorf("refresh input invalid: %w", ErrValidation)
	}

	claims, err := s.tokenManager.ParseRefreshToken(rawToken)
	if err != nil {
		return TokenPair{}, ErrInvalidToken
	}

	hashed := auth.HashToken(rawToken)
	storedToken, err := s.refreshTokens.GetActiveByHash(ctx, hashed)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return TokenPair{}, ErrInvalidToken
		}
		return TokenPair{}, fmt.Errorf("lookup refresh token: %w", err)
	}

	if storedToken.UserID != claims.Subject {
		return TokenPair{}, ErrInvalidToken
	}

	user, err := s.users.GetByID(ctx, claims.Subject)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return TokenPair{}, ErrInvalidToken
		}
		return TokenPair{}, fmt.Errorf("get user for refresh: %w", err)
	}

	if err := s.refreshTokens.RevokeByHash(ctx, hashed); err != nil && !errors.Is(err, repository.ErrNotFound) {
		return TokenPair{}, fmt.Errorf("revoke old refresh token: %w", err)
	}

	return s.issueTokenPair(ctx, user.ID, string(user.Role))
}

func (s *AuthService) Logout(ctx context.Context, input LogoutInput) error {
	rawToken := strings.TrimSpace(input.RefreshToken)
	if rawToken == "" {
		return fmt.Errorf("logout input invalid: %w", ErrValidation)
	}

	if _, err := s.tokenManager.ParseRefreshToken(rawToken); err != nil {
		return ErrInvalidToken
	}

	hashed := auth.HashToken(rawToken)
	err := s.refreshTokens.RevokeByHash(ctx, hashed)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, input RequestResetInput) error {
	emailAddress := normalizeEmail(input.Email)
	if !isValidEmail(emailAddress) {
		return fmt.Errorf("password reset request invalid: %w", ErrValidation)
	}

	user, err := s.users.GetByEmail(ctx, emailAddress)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("lookup user for reset: %w", err)
	}

	rawToken, err := auth.GenerateRandomToken(32)
	if err != nil {
		return fmt.Errorf("generate password reset token: %w", err)
	}

	hashedToken := auth.HashToken(rawToken)
	expiresAt := time.Now().UTC().Add(s.passwordResetTTL)
	if err := s.resetTokens.Create(ctx, &models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: expiresAt,
	}); err != nil {
		return fmt.Errorf("store password reset token: %w", err)
	}

	resetURL := s.frontendBaseURL + "/reset-password?token=" + url.QueryEscape(rawToken)
	message := email.Message{
		To:      user.Email,
		Subject: "Password Reset Request",
		Body:    "Use this link to reset your password: " + resetURL,
	}
	if err := s.emailSender.Send(ctx, message); err != nil {
		s.logger.Error("password reset email send failed", "error", err, "email", user.Email)
	}

	return nil
}

func (s *AuthService) ConfirmPasswordReset(ctx context.Context, input ConfirmResetInput) error {
	rawToken := strings.TrimSpace(input.Token)
	if rawToken == "" || len(input.NewPassword) < 8 {
		return fmt.Errorf("password reset confirm invalid: %w", ErrValidation)
	}

	hashedToken := auth.HashToken(rawToken)
	storedToken, err := s.resetTokens.GetActiveByHash(ctx, hashedToken)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrInvalidToken
		}
		return fmt.Errorf("load reset token: %w", err)
	}

	hashedPassword, err := auth.HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	if err := s.users.UpdatePasswordHash(ctx, storedToken.UserID, hashedPassword); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrInvalidToken
		}
		return fmt.Errorf("update password hash: %w", err)
	}

	if err := s.resetTokens.MarkUsedByID(ctx, storedToken.ID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrInvalidToken
		}
		return fmt.Errorf("mark reset token used: %w", err)
	}

	return nil
}

func (s *AuthService) issueTokenPair(ctx context.Context, userID, role string) (TokenPair, error) {
	accessToken, accessExpiresAt, err := s.tokenManager.GenerateAccessToken(userID, role)
	if err != nil {
		return TokenPair{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, _, refreshExpiresAt, err := s.tokenManager.GenerateRefreshToken(userID)
	if err != nil {
		return TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := s.refreshTokens.Create(ctx, &models.RefreshToken{
		UserID:    userID,
		TokenHash: auth.HashToken(refreshToken),
		ExpiresAt: refreshExpiresAt,
	}); err != nil {
		return TokenPair{}, fmt.Errorf("store refresh token: %w", err)
	}

	return TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && len(email) >= 5
}
