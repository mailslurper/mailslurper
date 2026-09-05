package auth

import (
	"testing"
	"time"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword([]byte("s3cret"))
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !CheckPassword(hash, []byte("s3cret")) {
		t.Error("expected correct password to check out")
	}
	if CheckPassword(hash, []byte("wrong")) {
		t.Error("expected wrong password to fail")
	}
}

func TestCheckCredentials(t *testing.T) {
	hash, _ := HashPassword([]byte("s3cret"))
	creds := map[string]string{"alex": string(hash)}

	if err := CheckCredentials(creds, "alex", "s3cret"); err != nil {
		t.Errorf("expected valid credentials to pass, got %v", err)
	}
	if err := CheckCredentials(creds, "alex", "wrong"); err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
	if err := CheckCredentials(creds, "nobody", "s3cret"); err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials for unknown user, got %v", err)
	}
}

func TestSessionServiceRoundTrip(t *testing.T) {
	svc := NewSessionService("test-secret", time.Hour)
	value := svc.Create("alex")

	userName, err := svc.Verify(value)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if userName != "alex" {
		t.Errorf("userName = %q", userName)
	}
}

func TestSessionServiceRejectsTampering(t *testing.T) {
	svc := NewSessionService("test-secret", time.Hour)
	value := svc.Create("alex")

	tampered := value + "x"
	if _, err := svc.Verify(tampered); err == nil {
		t.Fatal("expected tampered cookie to be rejected")
	}
}

func TestSessionServiceRejectsExpired(t *testing.T) {
	svc := NewSessionService("test-secret", -time.Hour)
	value := svc.Create("alex")

	if _, err := svc.Verify(value); err != ErrInvalidSession {
		t.Fatalf("expected ErrInvalidSession, got %v", err)
	}
}
