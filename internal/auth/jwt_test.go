package auth

import (
	"testing"
	"time"
)

func TestJWTIssueAndParse(t *testing.T) {
	svc := NewJWTService("test-secret", time.Hour)

	token, err := svc.Issue("alex")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	userName, err := svc.Parse(token)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if userName != "alex" {
		t.Errorf("userName = %q", userName)
	}
}

func TestJWTRejectsExpired(t *testing.T) {
	svc := NewJWTService("test-secret", -time.Hour)
	token, err := svc.Issue("alex")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if _, err := svc.Parse(token); err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}

func TestJWTRejectsWrongSecret(t *testing.T) {
	svc1 := NewJWTService("secret-one", time.Hour)
	svc2 := NewJWTService("secret-two", time.Hour)

	token, err := svc1.Issue("alex")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if _, err := svc2.Parse(token); err == nil {
		t.Fatal("expected token signed with a different secret to be rejected")
	}
}

func TestJWTLogoutInvalidatesToken(t *testing.T) {
	svc := NewJWTService("test-secret", time.Hour)
	token, err := svc.Issue("alex")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	svc.Logout(token)

	if _, err := svc.Parse(token); err != ErrTokenInvalidated {
		t.Fatalf("expected ErrTokenInvalidated, got %v", err)
	}
}
