package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/p0vidl0/mylslurper/internal/config"
	"github.com/p0vidl0/mylslurper/internal/mail"
)

func seedMail(t *testing.T, a *API, item *mail.Item) {
	t.Helper()
	if _, err := a.Store.StoreMail(context.Background(), item); err != nil {
		t.Fatalf("StoreMail: %v", err)
	}
}

func TestListGetAndCountMail(t *testing.T) {
	a := newTestAPI(t, config.Default())
	router := a.Router()

	seedMail(t, a, &mail.Item{
		ID: "m1", From: "sender@example.com", To: []string{"recipient@example.com"},
		Subject: "Hello", DateSent: time.Now(), TextBody: "hi", ContentType: "text/plain",
		Attachments: []mail.Attachment{{ID: "a1", FileName: "note.txt", ContentType: "text/plain", Content: []byte("hey")}},
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/mail", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/mail = %d: %s", rec.Code, rec.Body.String())
	}
	var list MailListResponse
	decodeJSON(t, rec, &list)
	if list.TotalRecords != 1 || len(list.MailItems) != 1 {
		t.Fatalf("unexpected list response: %+v", list)
	}
	if list.MailItems[0].AttachmentCount != 1 {
		t.Fatalf("AttachmentCount = %d", list.MailItems[0].AttachmentCount)
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/mail/count", nil))
	var count CountResponse
	decodeJSON(t, rec, &count)
	if count.Count != 1 {
		t.Fatalf("Count = %d", count.Count)
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/mail/m1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/mail/m1 = %d: %s", rec.Code, rec.Body.String())
	}
	var detail MailDetailDTO
	decodeJSON(t, rec, &detail)
	if detail.Subject != "Hello" || len(detail.Attachments) != 1 {
		t.Fatalf("unexpected detail: %+v", detail)
	}
	if detail.Attachments[0].SizeBytes != len("hey") {
		t.Fatalf("expected attachment SizeBytes to reflect content length, got %d", detail.Attachments[0].SizeBytes)
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/mail/m1/attachments/a1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET attachment = %d", rec.Code)
	}
	if rec.Body.String() != "hey" {
		t.Errorf("attachment body = %q", rec.Body.String())
	}
}

func TestGetMissingMailReturns404(t *testing.T) {
	a := newTestAPI(t, config.Default())
	router := a.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/mail/missing", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetMailMessageAndRaw(t *testing.T) {
	a := newTestAPI(t, config.Default())
	router := a.Router()

	seedMail(t, a, &mail.Item{
		ID: "m2", From: "sender@example.com", To: []string{"recipient@example.com"},
		Subject: "Raw test", DateSent: time.Now(), HTMLBody: "<p>hi</p>", RawMessage: "Subject: Raw test\r\n\r\nraw",
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/mail/m2/message", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "<p>hi</p>" {
		t.Fatalf("message = %d %q", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/mail/m2/messageraw", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "Subject: Raw test\r\n\r\nraw" {
		t.Fatalf("messageraw = %d %q", rec.Code, rec.Body.String())
	}
}

func TestListMailFiltersAndPagination(t *testing.T) {
	a := newTestAPI(t, config.Default())
	router := a.Router()
	sent := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)

	seedMail(t, a, &mail.Item{
		ID: "m-a", From: "alice@example.com", To: []string{"bob@example.com"},
		Subject: "Invoice Alpha", DateSent: sent, TextBody: "one",
	})
	seedMail(t, a, &mail.Item{
		ID: "m-b", From: "carol@example.com", To: []string{"dave@example.com"},
		Subject: "Beta note", DateSent: sent.Add(time.Hour), TextBody: "two",
	})

	rec := httptest.NewRecorder()
	url := "/api/mail?from=alice&to=bob&start=2026-05-10&end=2026-05-11&page=1&pageSize=200"
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, url, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET filtered list = %d: %s", rec.Code, rec.Body.String())
	}
	var list MailListResponse
	decodeJSON(t, rec, &list)
	if list.TotalRecords != 1 || len(list.MailItems) != 1 {
		t.Fatalf("unexpected filtered list: %+v", list)
	}
	if list.PageSize != 200 {
		t.Fatalf("PageSize = %d", list.PageSize)
	}
}

func TestGetMissingAttachmentReturns404(t *testing.T) {
	a := newTestAPI(t, config.Default())
	router := a.Router()

	seedMail(t, a, &mail.Item{
		ID: "m-att", From: "sender@example.com", To: []string{"recipient@example.com"},
		Subject: "No attachment", DateSent: time.Now(),
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/mail/m-att/attachments/missing", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
