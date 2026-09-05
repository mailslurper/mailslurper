package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey int

const userContextKey contextKey = iota

// ContextWithUser returns a context carrying the authenticated username.
func ContextWithUser(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, userContextKey, userName)
}

// UserFromContext returns the authenticated username stored by RequireAuth,
// if any.
func UserFromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(userContextKey).(string)
	return u, ok
}

// RequireAuth returns middleware that accepts either a JWT bearer token
// (for fetch calls) or the signed session cookie (for plain-link
// navigations, e.g. attachment downloads), rejecting the request with 401
// if neither is present and valid.
func RequireAuth(jwtSvc *JWTService, sessionSvc *SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userName, ok := authenticate(r, jwtSvc, sessionSvc)
			if !ok {
				writeUnauthorized(w)
				return
			}
			next.ServeHTTP(w, r.WithContext(ContextWithUser(r.Context(), userName)))
		})
	}
}

func authenticate(r *http.Request, jwtSvc *JWTService, sessionSvc *SessionService) (string, bool) {
	if header := r.Header.Get("Authorization"); header != "" {
		if token, ok := strings.CutPrefix(header, "Bearer "); ok {
			if userName, err := jwtSvc.Parse(token); err == nil {
				return userName, true
			}
			return "", false
		}
	}

	if cookie, err := r.Cookie(CookieName); err == nil {
		if userName, err := sessionSvc.Verify(cookie.Value); err == nil {
			return userName, true
		}
	}

	return "", false
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
}
