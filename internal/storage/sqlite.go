package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/p0vidl0/mylslurper/internal/mail"

	_ "modernc.org/sqlite"
)

const timeFormat = time.RFC3339Nano

// SQLiteStorage is the sole Storage implementation, backed by a pure-Go
// SQLite driver so the binary never needs CGO.
type SQLiteStorage struct {
	path string
	db   *sql.DB
}

// NewSQLiteStorage returns a Storage backed by the SQLite file at path.
// The file is created if it doesn't exist and left untouched if it does.
func NewSQLiteStorage(path string) *SQLiteStorage {
	return &SQLiteStorage{path: path}
}

func (s *SQLiteStorage) Connect(ctx context.Context) error {
	db, err := sql.Open("sqlite", s.path)
	if err != nil {
		return fmt.Errorf("opening sqlite database %q: %w", s.path, err)
	}

	// A single connection avoids any ambiguity about which pooled connection
	// a PRAGMA applies to, which is plenty for this tool's write volume.
	db.SetMaxOpenConns(1)

	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		return fmt.Errorf("enabling foreign keys: %w", err)
	}
	if _, err := db.ExecContext(ctx, `PRAGMA journal_mode = WAL`); err != nil {
		return fmt.Errorf("enabling WAL mode: %w", err)
	}

	if err := ensureSchema(ctx, db); err != nil {
		db.Close()
		return err
	}

	s.db = db
	return nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

func (s *SQLiteStorage) StoreMail(ctx context.Context, item *mail.Item) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO mail_item (id, date_sent, from_address, to_addresses, subject, xmailer, content_type, boundary, text_body, html_body, raw_message)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		item.ID, item.DateSent.UTC().Format(timeFormat), item.From, strings.Join(item.To, "; "),
		item.Subject, item.XMailer, item.ContentType, item.Boundary, item.TextBody, item.HTMLBody, item.RawMessage,
	)
	if err != nil {
		return "", fmt.Errorf("inserting mail item: %w", err)
	}

	for _, att := range item.Attachments {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO attachment (id, mail_item_id, file_name, content_type, content)
			VALUES (?, ?, ?, ?, ?)
		`, att.ID, item.ID, att.FileName, att.ContentType, att.Content)
		if err != nil {
			return "", fmt.Errorf("inserting attachment: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return item.ID, nil
}

func (s *SQLiteStorage) GetMailByID(ctx context.Context, id string) (*mail.Item, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, date_sent, from_address, to_addresses, subject, xmailer, content_type, boundary, text_body, html_body
		FROM mail_item WHERE id = ?
	`, id)

	item, err := scanItem(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	attachments, err := s.attachmentsFor(ctx, []string{id})
	if err != nil {
		return nil, err
	}
	item.Attachments = attachments[id]

	return item, nil
}

func (s *SQLiteStorage) GetMailRawByID(ctx context.Context, id string) (string, error) {
	var raw string
	err := s.db.QueryRowContext(ctx, `SELECT raw_message FROM mail_item WHERE id = ?`, id).Scan(&raw)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return raw, err
}

func (s *SQLiteStorage) GetMailCollection(ctx context.Context, offset, limit int, search *Search) ([]*mail.Item, error) {
	where, args := buildWhere(search)
	orderBy := buildOrderBy(search)

	query := fmt.Sprintf(`
		SELECT id, date_sent, from_address, to_addresses, subject, xmailer, content_type
		FROM mail_item %s %s LIMIT ? OFFSET ?
	`, where, orderBy)
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying mail collection: %w", err)
	}
	defer rows.Close()

	var items []*mail.Item
	var ids []string
	for rows.Next() {
		var item mail.Item
		var dateSent, toAddresses string
		if err := rows.Scan(&item.ID, &dateSent, &item.From, &toAddresses, &item.Subject, &item.XMailer, &item.ContentType); err != nil {
			return nil, err
		}
		item.DateSent, _ = time.Parse(timeFormat, dateSent)
		item.To = splitAddresses(toAddresses)
		items = append(items, &item)
		ids = append(ids, item.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	attachments, err := s.attachmentsFor(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		item.Attachments = attachments[item.ID]
	}

	return items, nil
}

func (s *SQLiteStorage) GetMailCollectionWithBodies(ctx context.Context, offset, limit int, search *Search) ([]*mail.Item, error) {
	where, args := buildWhere(search)
	orderBy := buildOrderBy(search)

	query := fmt.Sprintf(`
		SELECT id, date_sent, from_address, to_addresses, subject, xmailer, content_type, text_body, html_body
		FROM mail_item %s %s LIMIT ? OFFSET ?
	`, where, orderBy)
	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying mail collection: %w", err)
	}
	defer rows.Close()

	var items []*mail.Item
	for rows.Next() {
		var item mail.Item
		var dateSent, toAddresses string
		if err := rows.Scan(&item.ID, &dateSent, &item.From, &toAddresses, &item.Subject, &item.XMailer, &item.ContentType, &item.TextBody, &item.HTMLBody); err != nil {
			return nil, err
		}
		item.DateSent, _ = time.Parse(timeFormat, dateSent)
		item.To = splitAddresses(toAddresses)
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (s *SQLiteStorage) GetMailCount(ctx context.Context, search *Search) (int, error) {
	where, args := buildWhere(search)
	query := fmt.Sprintf(`SELECT COUNT(*) FROM mail_item %s`, where)

	var count int
	if err := s.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("counting mail items: %w", err)
	}
	return count, nil
}

func (s *SQLiteStorage) GetAttachment(ctx context.Context, mailID, attachmentID string) (*mail.Attachment, error) {
	var att mail.Attachment
	att.MailItemID = mailID
	err := s.db.QueryRowContext(ctx, `
		SELECT id, file_name, content_type, content FROM attachment WHERE id = ? AND mail_item_id = ?
	`, attachmentID, mailID).Scan(&att.ID, &att.FileName, &att.ContentType, &att.Content)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	att.Size = len(att.Content)
	return &att, nil
}

func (s *SQLiteStorage) DeleteMailsOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	return s.deleteWhere(ctx, "date_sent < ?", cutoff.UTC().Format(timeFormat))
}

func (s *SQLiteStorage) DeleteAllMail(ctx context.Context) (int64, error) {
	return s.deleteWhere(ctx, "1 = 1")
}

func (s *SQLiteStorage) deleteWhere(ctx context.Context, condition string, args ...any) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, fmt.Sprintf(
		`DELETE FROM attachment WHERE mail_item_id IN (SELECT id FROM mail_item WHERE %s)`, condition,
	), args...); err != nil {
		return 0, fmt.Errorf("deleting attachments: %w", err)
	}

	res, err := tx.ExecContext(ctx, fmt.Sprintf(`DELETE FROM mail_item WHERE %s`, condition), args...)
	if err != nil {
		return 0, fmt.Errorf("deleting mail items: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// attachmentsFor batch-loads attachment metadata (no content bytes) for the
// given mail item IDs, grouped by mail item ID.
func (s *SQLiteStorage) attachmentsFor(ctx context.Context, ids []string) (map[string][]mail.Attachment, error) {
	result := map[string][]mail.Attachment{}
	if len(ids) == 0 {
		return result, nil
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(
		`SELECT id, mail_item_id, file_name, content_type, LENGTH(content) FROM attachment WHERE mail_item_id IN (%s)`, placeholders,
	), args...)
	if err != nil {
		return nil, fmt.Errorf("querying attachments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var att mail.Attachment
		if err := rows.Scan(&att.ID, &att.MailItemID, &att.FileName, &att.ContentType, &att.Size); err != nil {
			return nil, err
		}
		result[att.MailItemID] = append(result[att.MailItemID], att)
	}
	return result, rows.Err()
}

func scanItem(row *sql.Row) (*mail.Item, error) {
	var item mail.Item
	var dateSent, toAddresses string
	err := row.Scan(&item.ID, &dateSent, &item.From, &toAddresses, &item.Subject, &item.XMailer,
		&item.ContentType, &item.Boundary, &item.TextBody, &item.HTMLBody)
	if err != nil {
		return nil, err
	}
	item.DateSent, _ = time.Parse(timeFormat, dateSent)
	item.To = splitAddresses(toAddresses)
	return &item, nil
}

func splitAddresses(joined string) []string {
	if joined == "" {
		return nil
	}
	parts := strings.Split(joined, "; ")
	return parts
}
