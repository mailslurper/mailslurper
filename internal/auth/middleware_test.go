package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserFromContext(t *testing.T) {
	ctx := ContextWithUser(t.Context(), "alex")
	user, ok := UserFromContext(ctx)
	if !ok || user != "alex" {
		t.Fatalf("UserFromContext = %q, ok = %v", user, ok)
	}

	user, ok = UserFromContext(t.Context())
	if ok || user != "" {
		t.Fatalf("expected empty context, got %q ok=%v", user, ok)
	}
}

func TestRequireAuthBearerValid(t *testing.T) {
	jwtSvc := NewJWTService("test-secret", time.Hour)
	sessionSvc := NewSessionService("test-secret", time.Hour)
	token, err := jwtSvc.Issue("alex")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	handler := RequireAuth(jwtSvc, sessionSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok || user != "alex" {
			t.Errorf("context user = %q, ok = %v", user, ok)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestRequireAuthBearerInvalid(t *testing.T) {
	jwtSvc := NewJWTService("test-secret", time.Hour)
	sessionSvc := NewSessionService("test-secret", time.Hour)

	handler := RequireAuth(jwtSvc, sessionSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["error"] != "unauthorized" {
		t.Fatalf("error = %q", body["error"])
	}
}

func TestRequireAuthBearerExpired(t *testing.T) {
	jwtSvc := NewJWTService("test-secret", -time.Hour)
	sessionSvc := NewSessionService("test-secret", time.Hour)
	token, err := jwtSvc.Issue("alex")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	handler := RequireAuth(jwtSvc, sessionSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestRequireAuthSessionCookie(t *testing.T) {
	jwtSvc := NewJWTService("test-secret", time.Hour)
	sessionSvc := NewSessionService("test-secret", time.Hour)
	cookieValue := sessionSvc.Create("alex")

	handler := RequireAuth(jwtSvc, sessionSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok || user != "alex" {
			t.Errorf("context user = %q, ok = %v", user, ok)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: CookieName, Value: cookieValue})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestRequireAuthMissing(t *testing.T) {
	jwtSvc := NewJWTService("test-secret", time.Hour)
	sessionSvc := NewSessionService("test-secret", time.Hour)

	handler := RequireAuth(jwtSvc, sessionSvc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", rec.Code)
	}
}
