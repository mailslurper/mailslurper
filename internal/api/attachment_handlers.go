package api

import (
	"fmt"
	"net/http"

	"github.com/p0vidl0/mylslurper/internal/storage"
)

func (a *API) handleGetAttachment(w http.ResponseWriter, r *http.Request) {
	att, err := a.Store.GetAttachment(r.Context(), r.PathValue("id"), r.PathValue("attachmentId"))
	if err != nil {
		if err == storage.ErrNotFound {
			writeError(w, http.StatusNotFound, "attachment not found")
			return
		}
		a.Log.WithError(err).Error("attachment lookup failed")
		writeError(w, http.StatusInternalServerError, "could not retrieve attachment")
		return
	}

	contentType := att.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", att.FileName))
	_, _ = w.Write(att.Content)
}
