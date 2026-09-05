package api

import (
	"net/http"

	"github.com/p0vidl0/mylslurper/internal/version"
)

func (a *API) handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, VersionResponse{Version: version.Current})
}

func (a *API) handleSettings(w http.ResponseWriter, r *http.Request) {
	scheme := a.Config.AuthenticationScheme
	if scheme == "" {
		scheme = "none"
	}
	writeJSON(w, http.StatusOK, SettingsResponse{
		AuthenticationScheme: scheme,
		ServerVersion:        version.Current,
		PublicURL:            a.Config.PublicURL,
	})
}
