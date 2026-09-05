package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"
)

func newTestStorage(t *testing.T) (*SQLiteStorage, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	s := NewSQLiteStorage(path)
	if err := s.Connect(context.Background()); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s, path
}

func sampleItem(subject, from string, sent time.Time) *mail.Item {
	return &mail.Item{
		ID:          from + "-" + subject,
		DateSent:    sent,
		From:        from,
		To:          []string{"recipient@example.com"},
		Subject:     subject,
		ContentType: "text/plain",
		TextBody:    "body of " + subject,
		RawMessage:  "raw " + subject,
	}
}

func TestConnectIsIdempotentAndPersists(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "persist.db")

	s1 := NewSQLiteStorage(path)
	if err := s1.Connect(ctx); err != nil {
		t.Fatalf("first Connect: %v", err)
	}
	if _, err := s1.StoreMail(ctx, sampleItem("Hi", "a@b.com", time.Now())); err != nil {
		t.Fatalf("StoreMail: %v", err)
	}
	if err := s1.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	s2 := NewSQLiteStorage(path)
	if err := s2.Connect(ctx); err != nil {
		t.Fatalf("second Connect: %v", err)
	}
	defer s2.Close()

	count, err := s2.GetMailCount(ctx, nil)
	if err != nil {
		t.Fatalf("GetMailCount: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected mail to survive reconnect, got count=%d", count)
	}
}

func TestStoreAndGetMailByID(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()

	item := sampleItem("Test", "sender@example.com", time.Now())
	item.Attachments = []mail.Attachment{{ID: "att1", FileName: "note.txt", ContentType: "text/plain", Content: []byte("hello")}}

	id, err := s.StoreMail(ctx, item)
	if err != nil {
		t.Fatalf("StoreMail: %v", err)
	}

	got, err := s.GetMailByID(ctx, id)
	if err != nil {
		t.Fatalf("GetMailByID: %v", err)
	}
	if got.Subject != "Test" || got.From != "sender@example.com" {
		t.Errorf("got %+v", got)
	}
	if len(got.Attachments) != 1 || got.Attachments[0].FileName != "note.txt" {
		t.Errorf("attachments = %+v", got.Attachments)
	}
	if got.Attachments[0].Size != len("hello") {
		t.Errorf("expected attachment metadata to carry Size without loading Content, got Size=%d", got.Attachments[0].Size)
	}

	att, err := s.GetAttachment(ctx, id, "att1")
	if err != nil {
		t.Fatalf("GetAttachment: %v", err)
	}
	if string(att.Content) != "hello" {
		t.Errorf("attachment content = %q", att.Content)
	}
}

func TestGetMailByIDNotFound(t *testing.T) {
	s, _ := newTestStorage(t)
	if _, err := s.GetMailByID(context.Background(), "missing"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSearchSortAndPaging(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()
	base := time.Now().Add(-time.Hour)

	for i, subj := range []string{"Alpha", "Beta", "Gamma"} {
		item := sampleItem(subj, "sender@example.com", base.Add(time.Duration(i)*time.Minute))
		if _, err := s.StoreMail(ctx, item); err != nil {
			t.Fatalf("StoreMail: %v", err)
		}
	}

	items, err := s.GetMailCollection(ctx, 0, 10, &Search{SortField: "subject", SortDir: "asc"})
	if err != nil {
		t.Fatalf("GetMailCollection: %v", err)
	}
	if len(items) != 3 || items[0].Subject != "Alpha" || items[2].Subject != "Gamma" {
		t.Fatalf("unexpected order: %v", subjectsOf(items))
	}

	page, err := s.GetMailCollection(ctx, 1, 1, &Search{SortField: "subject", SortDir: "asc"})
	if err != nil {
		t.Fatalf("GetMailCollection page: %v", err)
	}
	if len(page) != 1 || page[0].Subject != "Beta" {
		t.Fatalf("unexpected page: %v", subjectsOf(page))
	}

	filtered, err := s.GetMailCollection(ctx, 0, 10, &Search{Query: "Gamma"})
	if err != nil {
		t.Fatalf("GetMailCollection filtered: %v", err)
	}
	if len(filtered) != 1 || filtered[0].Subject != "Gamma" {
		t.Fatalf("unexpected filtered result: %v", subjectsOf(filtered))
	}
}

func subjectsOf(items []*mail.Item) []string {
	var out []string
	for _, i := range items {
		out = append(out, i.Subject)
	}
	return out
}

func TestDeleteMailsOlderThan(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()
	now := time.Now()

	old := sampleItem("Old", "sender@example.com", now.Add(-48*time.Hour))
	old.Attachments = []mail.Attachment{{ID: "a1", FileName: "f.txt", ContentType: "text/plain", Content: []byte("x")}}
	recent := sampleItem("Recent", "sender@example.com", now)

	if _, err := s.StoreMail(ctx, old); err != nil {
		t.Fatal(err)
	}
	if _, err := s.StoreMail(ctx, recent); err != nil {
		t.Fatal(err)
	}

	deleted, err := s.DeleteMailsOlderThan(ctx, now.Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("DeleteMailsOlderThan: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected 1 deleted, got %d", deleted)
	}

	count, err := s.GetMailCount(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 remaining, got %d", count)
	}

	if _, err := s.GetAttachment(ctx, old.ID, "a1"); err != ErrNotFound {
		t.Fatalf("expected orphaned attachment to be deleted, got err=%v", err)
	}
}

func TestDeleteAllMail(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		item := sampleItem("Subj", "sender@example.com", time.Now())
		item.ID += string(rune('a' + i))
		if _, err := s.StoreMail(ctx, item); err != nil {
			t.Fatal(err)
		}
	}

	deleted, err := s.DeleteAllMail(ctx)
	if err != nil {
		t.Fatalf("DeleteAllMail: %v", err)
	}
	if deleted != 3 {
		t.Fatalf("expected 3 deleted, got %d", deleted)
	}
}

func TestGetMailRawByID(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()

	item := sampleItem("Raw", "sender@example.com", time.Now())
	if _, err := s.StoreMail(ctx, item); err != nil {
		t.Fatalf("StoreMail: %v", err)
	}

	raw, err := s.GetMailRawByID(ctx, item.ID)
	if err != nil {
		t.Fatalf("GetMailRawByID: %v", err)
	}
	if raw != item.RawMessage {
		t.Fatalf("raw = %q, want %q", raw, item.RawMessage)
	}
}

func TestGetMailCollectionWithBodies(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()

	item := sampleItem("Bodies", "sender@example.com", time.Now())
	item.HTMLBody = "<p>html</p>"
	if _, err := s.StoreMail(ctx, item); err != nil {
		t.Fatalf("StoreMail: %v", err)
	}

	items, err := s.GetMailCollectionWithBodies(ctx, 0, 10, nil)
	if err != nil {
		t.Fatalf("GetMailCollectionWithBodies: %v", err)
	}
	if len(items) != 1 || items[0].TextBody != item.TextBody {
		t.Fatalf("unexpected items: %+v", items)
	}
}

func TestSearchByFromAndTo(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()
	sent := time.Now()

	a := sampleItem("A", "alice@example.com", sent)
	a.To = []string{"bob@example.com"}
	b := sampleItem("B", "carol@example.com", sent)
	b.To = []string{"dave@example.com"}
	if _, err := s.StoreMail(ctx, a); err != nil {
		t.Fatal(err)
	}
	if _, err := s.StoreMail(ctx, b); err != nil {
		t.Fatal(err)
	}

	fromFiltered, err := s.GetMailCollection(ctx, 0, 10, &Search{From: "alice"})
	if err != nil {
		t.Fatal(err)
	}
	if len(fromFiltered) != 1 || fromFiltered[0].Subject != "A" {
		t.Fatalf("from filter = %+v", subjectsOf(fromFiltered))
	}

	toFiltered, err := s.GetMailCollection(ctx, 0, 10, &Search{To: "dave"})
	if err != nil {
		t.Fatal(err)
	}
	if len(toFiltered) != 1 || toFiltered[0].Subject != "B" {
		t.Fatalf("to filter = %+v", subjectsOf(toFiltered))
	}
}

func TestSearchByDateRange(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()
	middle := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)

	old := sampleItem("Old", "sender@example.com", middle.Add(-48*time.Hour))
	recent := sampleItem("Recent", "sender@example.com", middle)
	if _, err := s.StoreMail(ctx, old); err != nil {
		t.Fatal(err)
	}
	if _, err := s.StoreMail(ctx, recent); err != nil {
		t.Fatal(err)
	}

	items, err := s.GetMailCollection(ctx, 0, 10, &Search{
		Start: middle.Add(-time.Hour),
		End:   middle.Add(time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Subject != "Recent" {
		t.Fatalf("date filter = %+v", subjectsOf(items))
	}
}

func TestSearchInvalidSortField(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()
	base := time.Now()

	for i, subj := range []string{"First", "Second"} {
		item := sampleItem(subj, "sender@example.com", base.Add(time.Duration(i)*time.Minute))
		if _, err := s.StoreMail(ctx, item); err != nil {
			t.Fatal(err)
		}
	}

	items, err := s.GetMailCollection(ctx, 0, 10, &Search{SortField: "bogus", SortDir: "asc"})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].Subject != "First" {
		t.Fatalf("expected date_sent asc fallback order, got %+v", subjectsOf(items))
	}
}

func TestGetAttachmentNotFound(t *testing.T) {
	s, _ := newTestStorage(t)
	ctx := context.Background()

	item := sampleItem("NoAtt", "sender@example.com", time.Now())
	if _, err := s.StoreMail(ctx, item); err != nil {
		t.Fatal(err)
	}

	if _, err := s.GetAttachment(ctx, item.ID, "missing"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
