// Package storage persists parsed mail items and their attachments.
package storage

import (
	"context"
	"errors"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"
)

// ErrNotFound is returned when a mail item or attachment doesn't exist.
var ErrNotFound = errors.New("not found")

// Search describes filter and sort criteria for listing mail.
type Search struct {
	Query string // matches subject or body
	From  string
	To    string
	Start time.Time
	End   time.Time

	SortField string // "date" | "subject" | "from"
	SortDir   string // "asc" | "desc"
}

// Storage is the persistence contract for mail items and attachments.
type Storage interface {
	Connect(ctx context.Context) error
	Close() error

	StoreMail(ctx context.Context, item *mail.Item) (id string, err error)
	GetMailByID(ctx context.Context, id string) (*mail.Item, error)
	GetMailRawByID(ctx context.Context, id string) (string, error)
	GetMailCollection(ctx context.Context, offset, limit int, search *Search) ([]*mail.Item, error)
	GetMailCollectionWithBodies(ctx context.Context, offset, limit int, search *Search) ([]*mail.Item, error)
	GetMailCount(ctx context.Context, search *Search) (int, error)
	GetAttachment(ctx context.Context, mailID, attachmentID string) (*mail.Attachment, error)

	DeleteMailsOlderThan(ctx context.Context, cutoff time.Time) (deleted int64, err error)
	DeleteAllMail(ctx context.Context) (deleted int64, err error)
}
