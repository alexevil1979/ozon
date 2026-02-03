package jwt

import (
	"testing"
	"time"
)

func TestJWT_IssueAndValidate(t *testing.T) {
	j := New("test_secret", "1h")
	userID := "user-123"
	token, err := j.Issue(userID)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if token == "" {
		t.Fatal("token empty")
	}
	sub, err := j.Validate(token)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if sub != userID {
		t.Errorf("subject = %q, want %q", sub, userID)
	}
}

func TestJWT_ValidateInvalid(t *testing.T) {
	j := New("test_secret", "1h")
	_, err := j.Validate("invalid.token.here")
	if err != ErrInvalidToken {
		t.Errorf("Validate invalid: got %v", err)
	}
}

func TestJWT_Expiry(t *testing.T) {
	j := New("test_secret", "1ms")
	token, err := j.Issue("user-1")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond)
	_, err = j.Validate(token)
	if err == nil {
		t.Error("expired token should be invalid")
	}
}
