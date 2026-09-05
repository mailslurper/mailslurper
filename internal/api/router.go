// Package api implements MylSlurper's JSON HTTP API and serves the
// embedded single-page frontend from the same server/origin.
package api

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/p0vidl0/mylslurper/internal/auth"
	"github.com/p0vidl0/mylslurper/internal/config"
	"github.com/p0vidl0/mylslurper/internal/events"
	"github.com/p0vidl0/mylslurper/internal/storage"
	"github.com/sirupsen/logrus"
)

// API wires together storage, auth, and configuration to build the HTTP
// route table.
type API struct {
	Store    storage.Storage
	Config   *config.Config
	JWT      *auth.JWTService
	Sessions *auth.SessionService
	Log      *logrus.Logger
	Assets   fs.FS
	Events   *events.Hub
}

// Router builds the complete http.Handler: the JSON API under /api/, and
// the embedded SPA everywhere else.
func (a *API) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/version", a.handleVersion)
	mux.HandleFunc("GET /api/settings", a.handleSettings)
	mux.HandleFunc("POST /api/login", a.handleLogin)

	mux.Handle("GET /api/mail", a.protect(a.handleListMail))
	mux.Handle("GET /api/mail/count", a.protect(a.handleMailCount))
	mux.Handle("GET /api/mail/{id}", a.protect(a.handleGetMail))
	mux.Handle("GET /api/mail/{id}/message", a.protect(a.handleGetMailMessage))
	mux.Handle("GET /api/mail/{id}/messageraw", a.protect(a.handleGetMailRaw))
	mux.Handle("GET /api/mail/{id}/attachments/{attachmentId}", a.protect(a.handleGetAttachment))
	mux.Handle("POST /api/mail/prune", a.protect(a.handlePruneMail))
	mux.Handle("GET /api/prune-options", a.protect(a.handlePruneOptions))
	mux.Handle("POST /api/logout", a.protect(a.handleLogout))
	mux.Handle("GET /api/events", a.protect(a.handleEvents))

	if a.Assets != nil {
		mux.Handle("/", a.staticHandler())
	}

	return a.withLogging(a.withCORS(mux))
}

// protect wraps h with auth middleware unless authentication is disabled.
func (a *API) protect(h http.HandlerFunc) http.Handler {
	if !a.Config.IsAuthEnabled() {
		return h
	}
	return auth.RequireAuth(a.JWT, a.Sessions)(h)
}

// staticHandler serves the embedded frontend, falling back to index.html
// for any path that isn't a real file, so the SPA's client-side router
// receives every deep link.
func (a *API) staticHandler() http.Handler {
	fileServer := http.FileServerFS(a.Assets)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if clean == "" || clean == "." {
			clean = "index.html"
		}

		if _, err := fs.Stat(a.Assets, clean); err != nil {
			data, err := fs.ReadFile(a.Assets, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(data)
			return
		}

		fileServer.ServeHTTP(w, r)
	})
}
