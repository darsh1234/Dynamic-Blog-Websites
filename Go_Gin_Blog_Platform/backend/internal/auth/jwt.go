package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

var ErrInvalidToken = errors.New("invalid token")

type AccessClaims struct {
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (m *TokenManager) GenerateAccessToken(userID, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(m.accessTTL)

	claims := AccessClaims{
		Role:      role,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        generateTokenID(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.accessSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}

	return token, expiresAt, nil
}

func (m *TokenManager) GenerateRefreshToken(userID string) (token string, tokenID string, expiresAt time.Time, err error) {
	now := time.Now().UTC()
	expiresAt = now.Add(m.refreshTTL)
	tokenID = generateTokenID()

	claims := RefreshClaims{
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        tokenID,
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.refreshSecret)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("sign refresh token: %w", err)
	}

	return token, tokenID, expiresAt, nil
}

func (m *TokenManager) ParseAccessToken(token string) (*AccessClaims, error) {
	parsedClaims := &AccessClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, parsedClaims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return m.accessSecret, nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	if parsedClaims.TokenType != TokenTypeAccess {
		return nil, ErrInvalidToken
	}

	if parsedClaims.Subject == "" || parsedClaims.Role == "" {
		return nil, ErrInvalidToken
	}

	return parsedClaims, nil
}

func (m *TokenManager) ParseRefreshToken(token string) (*RefreshClaims, error) {
	parsedClaims := &RefreshClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, parsedClaims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return m.refreshSecret, nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	if parsedClaims.TokenType != TokenTypeRefresh {
		return nil, ErrInvalidToken
	}

	if parsedClaims.Subject == "" || parsedClaims.ID == "" {
		return nil, ErrInvalidToken
	}

	return parsedClaims, nil
}

func generateTokenID() string {
	bytes := make([]byte, 18)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func GenerateRandomToken(length int) (string, error) {
	if length <= 0 {
		length = 32
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
