package auth

import "testing"

func TestHashPasswordAndVerify(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("expected hash generation to succeed: %v", err)
	}

	if hash == "password123" {
		t.Fatalf("expected hashed password to differ from plain text")
	}

	if !VerifyPassword(hash, "password123") {
		t.Fatalf("expected password verification to succeed")
	}

	if VerifyPassword(hash, "wrong-password") {
		t.Fatalf("expected wrong password verification to fail")
	}
}

func TestHashTokenDeterministic(t *testing.T) {
	left := HashToken("sample-token")
	right := HashToken("sample-token")
	other := HashToken("different-token")

	if left != right {
		t.Fatalf("expected same input tokens to hash identically")
	}
	if left == other {
		t.Fatalf("expected different input tokens to hash differently")
	}
}
