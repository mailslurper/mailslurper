package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/p0vidl0/mylslurper/internal/storage"
)

const (
	defaultPageSize = 50
	maxPageSize     = 200
)

func (a *API) handleListMail(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	search := &storage.Search{
		Query:     q.Get("q"),
		From:      q.Get("from"),
		To:        q.Get("to"),
		SortField: q.Get("sort"),
		SortDir:   q.Get("dir"),
	}
	if start, ok := parseDateParam(q.Get("start")); ok {
		search.Start = start
	}
	if end, ok := parseDateParam(q.Get("end")); ok {
		search.End = end
	}

	page := parseIntParam(q.Get("page"), 1)
	if page < 1 {
		page = 1
	}
	pageSize := parseIntParam(q.Get("pageSize"), defaultPageSize)
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	ctx := r.Context()

	total, err := a.Store.GetMailCount(ctx, search)
	if err != nil {
		a.Log.WithError(err).Error("failed to count mail")
		writeError(w, http.StatusInternalServerError, "could not count mail")
		return
	}

	items, err := a.Store.GetMailCollection(ctx, (page-1)*pageSize, pageSize, search)
	if err != nil {
		a.Log.WithError(err).Error("failed to list mail")
		writeError(w, http.StatusInternalServerError, "could not list mail")
		return
	}

	dtoItems := make([]MailListItemDTO, 0, len(items))
	for _, item := range items {
		dtoItems = append(dtoItems, MailListItemDTO{
			ID:              item.ID,
			DateSent:        formatTime(item.DateSent),
			FromAddress:     item.From,
			ToAddresses:     item.To,
			Subject:         item.Subject,
			XMailer:         item.XMailer,
			ContentType:     item.ContentType,
			AttachmentCount: len(item.Attachments),
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	writeJSON(w, http.StatusOK, MailListResponse{
		MailItems:    dtoItems,
		Page:         page,
		PageSize:     pageSize,
		TotalRecords: total,
		TotalPages:   totalPages,
	})
}

func (a *API) handleMailCount(w http.ResponseWriter, r *http.Request) {
	count, err := a.Store.GetMailCount(r.Context(), nil)
	if err != nil {
		a.Log.WithError(err).Error("failed to count mail")
		writeError(w, http.StatusInternalServerError, "could not count mail")
		return
	}
	writeJSON(w, http.StatusOK, CountResponse{Count: count})
}

func (a *API) handleGetMail(w http.ResponseWriter, r *http.Request) {
	item, err := a.Store.GetMailByID(r.Context(), r.PathValue("id"))
	if err != nil {
		a.respondItemError(w, err)
		return
	}

	attachments := make([]AttachmentDTO, 0, len(item.Attachments))
	for _, att := range item.Attachments {
		attachments = append(attachments, AttachmentDTO{
			ID:          att.ID,
			FileName:    att.FileName,
			ContentType: att.ContentType,
			SizeBytes:   att.Size,
		})
	}

	writeJSON(w, http.StatusOK, MailDetailDTO{
		ID:          item.ID,
		DateSent:    formatTime(item.DateSent),
		FromAddress: item.From,
		ToAddresses: item.To,
		Subject:     item.Subject,
		XMailer:     item.XMailer,
		ContentType: item.ContentType,
		TextBody:    item.TextBody,
		HTMLBody:    sanitizeHTML(item.HTMLBody),
		Attachments: attachments,
	})
}

func (a *API) handleGetMailMessage(w http.ResponseWriter, r *http.Request) {
	item, err := a.Store.GetMailByID(r.Context(), r.PathValue("id"))
	if err != nil {
		a.respondItemError(w, err)
		return
	}

	if item.HTMLBody != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(sanitizeHTML(item.HTMLBody)))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(item.TextBody))
}

func (a *API) handleGetMailRaw(w http.ResponseWriter, r *http.Request) {
	raw, err := a.Store.GetMailRawByID(r.Context(), r.PathValue("id"))
	if err != nil {
		a.respondItemError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(raw))
}

func (a *API) respondItemError(w http.ResponseWriter, err error) {
	if err == storage.ErrNotFound {
		writeError(w, http.StatusNotFound, "mail item not found")
		return
	}
	a.Log.WithError(err).Error("mail lookup failed")
	writeError(w, http.StatusInternalServerError, "could not retrieve mail item")
}

func parseDateParam(raw string) (time.Time, bool) {
	if raw == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func parseIntParam(raw string, def int) int {
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return n
}
