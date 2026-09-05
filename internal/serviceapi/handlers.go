// Package serviceapi implements the classic MailSlurper service REST API
// consumed by integration tests and tooling (GET /mailcount, GET /mail, …).
package serviceapi

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"
	"github.com/p0vidl0/mylslurper/internal/storage"
	"github.com/sirupsen/logrus"
)

const legacyPageSize = 50

const legacyTimeFormat = "2006-01-02 15:04:05"

// Server serves the unauthenticated MailSlurper-compatible service API.
type Server struct {
	Store storage.Storage
	Log   *logrus.Logger
}

type mailCountResponse struct {
	MailCount int `json:"mailCount"`
}

type mailCollectionResponse struct {
	MailItems    []mailItemResponse `json:"mailItems"`
	TotalPages   int                `json:"totalPages"`
	TotalRecords int                `json:"totalRecords"`
}

type mailItemResponse struct {
	ID          string   `json:"id"`
	DateSent    string   `json:"dateSent"`
	FromAddress string   `json:"fromAddress"`
	ToAddresses []string `json:"toAddresses"`
	Subject     string   `json:"subject"`
	XMailer     string   `json:"xmailer,omitempty"`
	Body        string   `json:"body"`
	ContentType string   `json:"contentType,omitempty"`
}

// Router returns the service API handler table.
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /mailcount", s.handleMailCount)
	mux.HandleFunc("GET /mail", s.handleMailCollection)
	mux.HandleFunc("GET /mail/{id}", s.handleGetMail)
	return mux
}

func (s *Server) handleMailCount(w http.ResponseWriter, r *http.Request) {
	count, err := s.Store.GetMailCount(r.Context(), nil)
	if err != nil {
		s.Log.WithError(err).Error("service api: mail count failed")
		http.Error(w, "Problem getting mail count", http.StatusInternalServerError)
		return
	}
	writeJSON(w, mailCountResponse{MailCount: count})
}

func (s *Server) handleMailCollection(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page := 1
	if raw := q.Get("pageNumber"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			http.Error(w, "A valid page number is required", http.StatusBadRequest)
			return
		}
		page = n
	}

	search := &storage.Search{
		Query:     q.Get("message"),
		From:      q.Get("from"),
		To:        q.Get("to"),
		SortField: q.Get("orderby"),
		SortDir:   q.Get("dir"),
	}
	if start, ok := parseLegacyDateParam(q.Get("start")); ok {
		search.Start = start
	}
	if end, ok := parseLegacyDateParam(q.Get("end")); ok {
		search.End = end
	}
	if search.SortField == "" {
		search.SortField = "date"
	}
	if search.SortDir == "" {
		search.SortDir = "desc"
	}

	ctx := r.Context()
	offset := (page - 1) * legacyPageSize

	items, err := s.Store.GetMailCollectionWithBodies(ctx, offset, legacyPageSize, search)
	if err != nil {
		s.Log.WithError(err).Error("service api: list mail failed")
		http.Error(w, "Problem getting mail collection", http.StatusInternalServerError)
		return
	}

	total, err := s.Store.GetMailCount(ctx, search)
	if err != nil {
		s.Log.WithError(err).Error("service api: count mail failed")
		http.Error(w, "Error getting record count", http.StatusInternalServerError)
		return
	}

	respItems := make([]mailItemResponse, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, toMailItemResponse(item))
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / legacyPageSize))
	}

	writeJSON(w, mailCollectionResponse{
		MailItems:    respItems,
		TotalPages:   totalPages,
		TotalRecords: total,
	})
}

func (s *Server) handleGetMail(w http.ResponseWriter, r *http.Request) {
	item, err := s.Store.GetMailByID(r.Context(), r.PathValue("id"))
	if err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, "Mail item not found", http.StatusNotFound)
			return
		}
		s.Log.WithError(err).Error("service api: get mail failed")
		http.Error(w, "Problem getting mail item", http.StatusInternalServerError)
		return
	}
	writeJSON(w, toMailItemResponse(item))
}

func toMailItemResponse(item *mail.Item) mailItemResponse {
	return mailItemResponse{
		ID:          item.ID,
		DateSent:    formatLegacyTime(item.DateSent),
		FromAddress: item.From,
		ToAddresses: item.To,
		Subject:     item.Subject,
		XMailer:     item.XMailer,
		Body:        item.Body(),
		ContentType: item.ContentType,
	}
}

func formatLegacyTime(t time.Time) string {
	return t.UTC().Format(legacyTimeFormat)
}

func parseLegacyDateParam(raw string) (time.Time, bool) {
	if raw == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{legacyTimeFormat, time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}
