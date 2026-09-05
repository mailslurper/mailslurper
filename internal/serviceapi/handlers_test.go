package serviceapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"
	"github.com/p0vidl0/mylslurper/internal/storage"
	"github.com/sirupsen/logrus"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	store := storage.NewSQLiteStorage(filepath.Join(t.TempDir(), "test.db"))
	if err := store.Connect(context.Background()); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)
	return &Server{Store: store, Log: log}
}

func TestMailCountAndCollection(t *testing.T) {
	s := newTestServer(t)
	router := s.Router()

	sent := time.Date(2026, 3, 1, 12, 30, 0, 0, time.UTC)
	if _, err := s.Store.StoreMail(context.Background(), &mail.Item{
		ID: "m1", From: "a@b.com", To: []string{"user@example.com"},
		Subject: "Your log in code is 123456", DateSent: sent,
		TextBody: "Use code 123456", ContentType: "text/plain",
	}); err != nil {
		t.Fatalf("StoreMail: %v", err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/mailcount", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /mailcount = %d", rec.Code)
	}
	var countResp mailCountResponse
	if err := json.NewDecoder(rec.Body).Decode(&countResp); err != nil {
		t.Fatalf("decode mailcount: %v", err)
	}
	if countResp.MailCount != 1 {
		t.Fatalf("mailCount = %d", countResp.MailCount)
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/mail?to=user@example.com&pageNumber=1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /mail = %d: %s", rec.Code, rec.Body.String())
	}
	var listResp mailCollectionResponse
	if err := json.NewDecoder(rec.Body).Decode(&listResp); err != nil {
		t.Fatalf("decode mail list: %v", err)
	}
	if len(listResp.MailItems) != 1 {
		t.Fatalf("mailItems len = %d", len(listResp.MailItems))
	}
	item := listResp.MailItems[0]
	if item.Body != "Use code 123456" {
		t.Fatalf("body = %q", item.Body)
	}
	if item.DateSent != "2026-03-01 12:30:00" {
		t.Fatalf("dateSent = %q", item.DateSent)
	}
}

func TestGetMailByID(t *testing.T) {
	s := newTestServer(t)
	router := s.Router()
	ctx := context.Background()

	if _, err := s.Store.StoreMail(ctx, &mail.Item{
		ID: "m1", From: "a@b.com", To: []string{"user@example.com"},
		Subject: "Detail", DateSent: time.Now(), TextBody: "body", ContentType: "text/plain",
	}); err != nil {
		t.Fatalf("StoreMail: %v", err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/mail/m1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /mail/m1 = %d", rec.Code)
	}
	var item mailItemResponse
	if err := json.NewDecoder(rec.Body).Decode(&item); err != nil {
		t.Fatal(err)
	}
	if item.Subject != "Detail" || item.Body != "body" {
		t.Fatalf("unexpected item: %+v", item)
	}

	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/mail/missing", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestMailCollectionBadPage(t *testing.T) {
	s := newTestServer(t)
	router := s.Router()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/mail?pageNumber=abc", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestMailCollectionDateFilters(t *testing.T) {
	s := newTestServer(t)
	router := s.Router()
	ctx := context.Background()
	middle := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)

	if _, err := s.Store.StoreMail(ctx, &mail.Item{
		ID: "old", From: "a@b.com", Subject: "Old", DateSent: middle.Add(-48 * time.Hour), TextBody: "x",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Store.StoreMail(ctx, &mail.Item{
		ID: "new", From: "a@b.com", Subject: "New", DateSent: middle, TextBody: "y",
	}); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	url := "/mail?start=2026-04-01%2010:00:00&end=2026-04-01%2014:00:00"
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, url, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /mail = %d: %s", rec.Code, rec.Body.String())
	}
	var list mailCollectionResponse
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if len(list.MailItems) != 1 || list.MailItems[0].Subject != "New" {
		t.Fatalf("filtered list = %+v", list.MailItems)
	}
}

func TestMailCollectionPagination(t *testing.T) {
	s := newTestServer(t)
	router := s.Router()
	ctx := context.Background()

	for i := 0; i < 51; i++ {
		if _, err := s.Store.StoreMail(ctx, &mail.Item{
			ID:       fmt.Sprintf("m%d", i),
			From:     "a@b.com",
			Subject:  fmt.Sprintf("Subj %d", i),
			DateSent: time.Now().Add(time.Duration(i) * time.Second),
			TextBody: "x",
		}); err != nil {
			t.Fatal(err)
		}
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/mail?pageNumber=1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /mail = %d", rec.Code)
	}
	var list mailCollectionResponse
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if list.TotalRecords != 51 || list.TotalPages != 2 || len(list.MailItems) != 50 {
		t.Fatalf("pagination = records:%d pages:%d items:%d", list.TotalRecords, list.TotalPages, len(list.MailItems))
	}
}
