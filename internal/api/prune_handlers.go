package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/p0vidl0/mylslurper/internal/events"
)

var pruneOptions = []PruneOptionDTO{
	{Code: "60plus", Description: "Older than 60 days"},
	{Code: "30plus", Description: "Older than 30 days"},
	{Code: "2wksplus", Description: "Older than 2 weeks"},
	{Code: "all", Description: "All emails"},
}

// pruneCutoff returns the cutoff time for a prune code, or all=true if the
// code means "delete everything". ok is false for an unrecognized code.
func pruneCutoff(code string) (cutoff time.Time, all bool, ok bool) {
	now := time.Now()
	switch code {
	case "60plus":
		return now.AddDate(0, 0, -60), false, true
	case "30plus":
		return now.AddDate(0, 0, -30), false, true
	case "2wksplus":
		return now.AddDate(0, 0, -14), false, true
	case "all":
		return time.Time{}, true, true
	default:
		return time.Time{}, false, false
	}
}

func (a *API) handlePruneOptions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, pruneOptions)
}

func (a *API) handlePruneMail(w http.ResponseWriter, r *http.Request) {
	var req PruneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cutoff, all, ok := pruneCutoff(req.Code)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid prune code")
		return
	}

	var deleted int64
	var err error
	if all {
		deleted, err = a.Store.DeleteAllMail(r.Context())
	} else {
		deleted, err = a.Store.DeleteMailsOlderThan(r.Context(), cutoff)
	}
	if err != nil {
		a.Log.WithError(err).Error("failed to prune mail")
		writeError(w, http.StatusInternalServerError, "could not prune mail")
		return
	}

	if a.Events != nil {
		a.Events.Publish(events.MailPruned(deleted))
	}

	writeJSON(w, http.StatusOK, PruneResponse{DeletedCount: deleted})
}
