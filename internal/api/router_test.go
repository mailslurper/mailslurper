package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/p0vidl0/mylslurper/internal/auth"
	"github.com/p0vidl0/mylslurper/internal/config"
	"github.com/p0vidl0/mylslurper/internal/mail"
	"github.com/p0vidl0/mylslurper/internal/storage"
	"github.com/sirupsen/logrus"
)

func jsonBody(s string) io.Reader {
	return strings.NewReader(s)
}

func newTestAPI(t *testing.T, cfg *config.Config) *API {
	t.Helper()
	store := storage.NewSQLiteStorage(filepath.Join(t.TempDir(), "test.db"))
	if err := store.Connect(context.Background()); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)

	return &API{
		Store:    store,
		Config:   cfg,
		JWT:      auth.NewJWTService(cfg.AuthSecret, time.Hour),
		Sessions: auth.NewSessionService(cfg.AuthSecret, time.Hour),
		Log:      log,
	}
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(v); err != nil {
		t.Fatalf("decoding response body %q: %v", rec.Body.String(), err)
	}
}

func TestVersionAndSettingsArePublic(t *testing.T) {
	cfg := config.Default()
	a := newTestAPI(t, cfg)
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/version", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/version = %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/settings", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/settings = %d", rec.Code)
	}
	var settings SettingsResponse
	decodeJSON(t, rec, &settings)
	if settings.AuthenticationScheme != "none" {
		t.Errorf("AuthenticationScheme = %q", settings.AuthenticationScheme)
	}
}

func TestMailRoutesRequireAuthWhenEnabled(t *testing.T) {
	cfg := config.Default()
	cfg.AuthenticationScheme = config.AuthSchemeBasic
	cfg.AuthSecret = "test-secret"
	hash, _ := auth.HashPassword([]byte("s3cret"))
	cfg.Credentials = map[string]string{"alex": string(hash)}

	a := newTestAPI(t, cfg)
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/mail", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without credentials, got %d", rec.Code)
	}

	// Login, then use the returned bearer token.
	loginBody := `{"userName":"alex","password":"s3cret"}`
	loginReq := httptest.NewRequest(http.MethodPost, "/api/login", jsonBody(loginBody))
	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, loginReq)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login = %d: %s", loginRec.Code, loginRec.Body.String())
	}
	var loginResp LoginResponse
	decodeJSON(t, loginRec, &loginResp)
	if loginResp.Token == "" {
		t.Fatal("expected a token from login")
	}

	authedReq := httptest.NewRequest(http.MethodGet, "/api/mail", nil)
	authedReq.Header.Set("Authorization", "Bearer "+loginResp.Token)
	authedRec := httptest.NewRecorder()
	router.ServeHTTP(authedRec, authedReq)
	if authedRec.Code != http.StatusOK {
		t.Fatalf("expected 200 with bearer token, got %d: %s", authedRec.Code, authedRec.Body.String())
	}
}

func TestLoginRejectsBadCredentials(t *testing.T) {
	cfg := config.Default()
	cfg.AuthenticationScheme = config.AuthSchemeBasic
	cfg.AuthSecret = "test-secret"
	hash, _ := auth.HashPassword([]byte("s3cret"))
	cfg.Credentials = map[string]string{"alex": string(hash)}

	a := newTestAPI(t, cfg)
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/login", jsonBody(`{"userName":"alex","password":"wrong"}`)))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for bad credentials, got %d", rec.Code)
	}
}

func TestPruneOptionsAndPrune(t *testing.T) {
	cfg := config.Default()
	a := newTestAPI(t, cfg)
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/prune-options", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/prune-options = %d", rec.Code)
	}
	var options []PruneOptionDTO
	decodeJSON(t, rec, &options)
	if len(options) != 4 {
		t.Fatalf("expected 4 prune options, got %d", len(options))
	}

	if _, err := a.Store.StoreMail(context.Background(), &mail.Item{ID: "m1", From: "a@b.com", Subject: "x", DateSent: time.Now()}); err != nil {
		t.Fatalf("StoreMail: %v", err)
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/mail/prune", jsonBody(`{"code":"all"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("prune = %d: %s", rec.Code, rec.Body.String())
	}
	var pruneResp PruneResponse
	decodeJSON(t, rec, &pruneResp)
	if pruneResp.DeletedCount != 1 {
		t.Fatalf("DeletedCount = %d", pruneResp.DeletedCount)
	}
}

func TestPruneRejectsInvalidCode(t *testing.T) {
	cfg := config.Default()
	a := newTestAPI(t, cfg)
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/mail/prune", jsonBody(`{"code":"bogus"}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid prune code, got %d", rec.Code)
	}
}

func authEnabledConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.Default()
	cfg.AuthenticationScheme = config.AuthSchemeBasic
	cfg.AuthSecret = "test-secret"
	hash, err := auth.HashPassword([]byte("s3cret"))
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	cfg.Credentials = map[string]string{"alex": string(hash)}
	return cfg
}

func TestLoginSuccessSetsCookieAndToken(t *testing.T) {
	a := newTestAPI(t, authEnabledConfig(t))
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/login", jsonBody(`{"userName":"alex","password":"s3cret"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("login = %d: %s", rec.Code, rec.Body.String())
	}

	var resp LoginResponse
	decodeJSON(t, rec, &resp)
	if resp.Token == "" {
		t.Fatal("expected token in login response")
	}

	cookies := rec.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == auth.CookieName {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil || sessionCookie.Value == "" {
		t.Fatalf("expected session cookie, got cookies: %+v", cookies)
	}
}

func TestLogoutInvalidatesToken(t *testing.T) {
	a := newTestAPI(t, authEnabledConfig(t))
	router := a.Router()

	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, httptest.NewRequest(http.MethodPost, "/api/login", jsonBody(`{"userName":"alex","password":"s3cret"}`)))
	var loginResp LoginResponse
	decodeJSON(t, loginRec, &loginResp)

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	logoutReq.Header.Set("Authorization", "Bearer "+loginResp.Token)
	logoutRec := httptest.NewRecorder()
	router.ServeHTTP(logoutRec, logoutReq)
	if logoutRec.Code != http.StatusOK {
		t.Fatalf("logout = %d: %s", logoutRec.Code, logoutRec.Body.String())
	}

	authedReq := httptest.NewRequest(http.MethodGet, "/api/mail", nil)
	authedReq.Header.Set("Authorization", "Bearer "+loginResp.Token)
	authedRec := httptest.NewRecorder()
	router.ServeHTTP(authedRec, authedReq)
	if authedRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 after logout, got %d", authedRec.Code)
	}
}

func TestCookieAuthWorksForMail(t *testing.T) {
	a := newTestAPI(t, authEnabledConfig(t))
	router := a.Router()

	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, httptest.NewRequest(http.MethodPost, "/api/login", jsonBody(`{"userName":"alex","password":"s3cret"}`)))
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login = %d", loginRec.Code)
	}

	mailReq := httptest.NewRequest(http.MethodGet, "/api/mail", nil)
	for _, c := range loginRec.Result().Cookies() {
		mailReq.AddCookie(c)
	}
	mailRec := httptest.NewRecorder()
	router.ServeHTTP(mailRec, mailReq)
	if mailRec.Code != http.StatusOK {
		t.Fatalf("GET /api/mail with session cookie = %d: %s", mailRec.Code, mailRec.Body.String())
	}
}

func TestLoginWhenAuthDisabled(t *testing.T) {
	cfg := config.Default()
	cfg.AuthenticationScheme = config.AuthSchemeNone
	a := newTestAPI(t, cfg)
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/login", jsonBody(`{"userName":"alex","password":"s3cret"}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when auth disabled, got %d", rec.Code)
	}
}

func TestStaticHandlerServesIndexAndAssets(t *testing.T) {
	cfg := config.Default()
	a := newTestAPI(t, cfg)
	a.Assets = fstest.MapFS{
		"index.html":       {Data: []byte("<html>index</html>")},
		"css/app.css":      {Data: []byte("body{}")},
	}
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "index") {
		t.Fatalf("/ = %d %q", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/inbox/deep/link", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "index") {
		t.Fatalf("SPA fallback = %d %q", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/css/app.css", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "body{}" {
		t.Fatalf("static asset = %d %q", rec.Code, rec.Body.String())
	}

	a.Assets = fstest.MapFS{}
	router = a.Router()
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/missing-only", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing index.html = %d", rec.Code)
	}
}

func TestCORSPreflight(t *testing.T) {
	cfg := config.Default()
	cfg.DevCORSOrigin = "http://localhost:5173"
	a := newTestAPI(t, cfg)
	router := a.Router()

	req := httptest.NewRequest(http.MethodOptions, "/api/mail", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("OPTIONS = %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != cfg.DevCORSOrigin {
		t.Fatalf("Allow-Origin = %q", got)
	}
}
