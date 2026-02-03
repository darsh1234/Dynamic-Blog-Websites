package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	GetActiveByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeByHash(ctx context.Context, tokenHash string) error
}

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *models.PasswordResetToken) error
	GetActiveByHash(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error)
	MarkUsedByID(ctx context.Context, tokenID string) error
}

type GormRefreshTokenRepository struct {
	db *gorm.DB
}

type GormPasswordResetTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *GormRefreshTokenRepository {
	return &GormRefreshTokenRepository{db: db}
}

func NewPasswordResetTokenRepository(db *gorm.DB) *GormPasswordResetTokenRepository {
	return &GormPasswordResetTokenRepository{db: db}
}

func (r *GormRefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

func (r *GormRefreshTokenRepository) GetActiveByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Where("expires_at > ?", time.Now().UTC()).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get active refresh token by hash: %w", err)
	}
	return &token, nil
}

func (r *GormRefreshTokenRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Updates(map[string]any{"revoked_at": &now})

	if result.Error != nil {
		return fmt.Errorf("revoke refresh token by hash: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GormPasswordResetTokenRepository) Create(ctx context.Context, token *models.PasswordResetToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("create password reset token: %w", err)
	}
	return nil
}

func (r *GormPasswordResetTokenRepository) GetActiveByHash(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		Where("used_at IS NULL").
		Where("expires_at > ?", time.Now().UTC()).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get active password reset token by hash: %w", err)
	}
	return &token, nil
}

func (r *GormPasswordResetTokenRepository) MarkUsedByID(ctx context.Context, tokenID string) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&models.PasswordResetToken{}).
		Where("id = ?", tokenID).
		Where("used_at IS NULL").
		Updates(map[string]any{"used_at": &now})

	if result.Error != nil {
		return fmt.Errorf("mark password reset token used: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
