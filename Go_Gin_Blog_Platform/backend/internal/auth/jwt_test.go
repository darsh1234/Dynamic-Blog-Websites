package auth

import (
	"testing"
	"time"
)

func TestTokenManagerAccessAndRefreshLifecycle(t *testing.T) {
	m := NewTokenManager("access-secret", "refresh-secret", 15*time.Minute, 7*24*time.Hour)

	access, accessExp, err := m.GenerateAccessToken("user-1", "author")
	if err != nil {
		t.Fatalf("expected access token generation to succeed: %v", err)
	}
	if access == "" {
		t.Fatalf("expected non-empty access token")
	}
	if accessExp.Before(time.Now().UTC()) {
		t.Fatalf("expected access token expiry to be in the future")
	}

	accessClaims, err := m.ParseAccessToken(access)
	if err != nil {
		t.Fatalf("expected access token parsing to succeed: %v", err)
	}
	if accessClaims.Subject != "user-1" || accessClaims.Role != "author" {
		t.Fatalf("unexpected access token claims: %+v", accessClaims)
	}

	refresh, refreshID, _, err := m.GenerateRefreshToken("user-1")
	if err != nil {
		t.Fatalf("expected refresh token generation to succeed: %v", err)
	}
	if refresh == "" || refreshID == "" {
		t.Fatalf("expected non-empty refresh token and id")
	}

	refreshClaims, err := m.ParseRefreshToken(refresh)
	if err != nil {
		t.Fatalf("expected refresh token parsing to succeed: %v", err)
	}
	if refreshClaims.Subject != "user-1" || refreshClaims.ID == "" {
		t.Fatalf("unexpected refresh token claims: %+v", refreshClaims)
	}
}

func TestTokenManagerRejectsWrongTokenType(t *testing.T) {
	m := NewTokenManager("access-secret", "refresh-secret", 15*time.Minute, 7*24*time.Hour)

	refresh, _, _, err := m.GenerateRefreshToken("user-1")
	if err != nil {
		t.Fatalf("expected refresh token generation to succeed: %v", err)
	}

	if _, err := m.ParseAccessToken(refresh); err == nil {
		t.Fatalf("expected refresh token to fail access parsing")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	token, err := GenerateRandomToken(24)
	if err != nil {
		t.Fatalf("expected random token generation to succeed: %v", err)
	}
	if len(token) < 24 {
		t.Fatalf("expected generated token to have a reasonable length")
	}
}
