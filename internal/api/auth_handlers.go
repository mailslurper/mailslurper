package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/p0vidl0/mylslurper/internal/auth"
)

func (a *API) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !a.Config.IsAuthEnabled() {
		writeError(w, http.StatusBadRequest, "authentication is disabled")
		return
	}

	if err := auth.CheckCredentials(a.Config.Credentials, req.UserName, req.Password); err != nil {
		writeError(w, http.StatusForbidden, "invalid credentials")
		return
	}

	token, err := a.JWT.Issue(req.UserName)
	if err != nil {
		a.Log.WithError(err).Error("failed to issue token")
		writeError(w, http.StatusInternalServerError, "could not issue token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieName,
		Value:    a.Sessions.Create(req.UserName),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   int(time.Duration(a.Config.AuthTimeoutInMinutes) * time.Minute / time.Second),
	})

	writeJSON(w, http.StatusOK, LoginResponse{Token: token})
}

func (a *API) handleLogout(w http.ResponseWriter, r *http.Request) {
	if header := r.Header.Get("Authorization"); header != "" {
		if token, ok := cutBearer(header); ok {
			a.JWT.Logout(token)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func cutBearer(header string) (string, bool) {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
		return "", false
	}
	return header[len(prefix):], true
}
