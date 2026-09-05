package storage

import (
	"context"
	"database/sql"
	"fmt"
)

const currentSchemaVersion = 1

// ensureSchema creates the database schema if it doesn't already exist and
// runs any forward migrations needed to reach currentSchemaVersion. It never
// drops or recreates existing tables, so data survives restarts.
func ensureSchema(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL)
	`); err != nil {
		return fmt.Errorf("creating schema_version table: %w", err)
	}

	version, err := readSchemaVersion(ctx, db)
	if err != nil {
		return err
	}

	if version < currentSchemaVersion {
		if err := migrateToV1(ctx, db); err != nil {
			return fmt.Errorf("migrating to schema v1: %w", err)
		}
		if err := setSchemaVersion(ctx, db, currentSchemaVersion); err != nil {
			return err
		}
	}

	return nil
}

func migrateToV1(ctx context.Context, db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS mail_item (
			id           TEXT PRIMARY KEY,
			date_sent    TEXT NOT NULL,
			from_address TEXT NOT NULL,
			to_addresses TEXT NOT NULL,
			subject      TEXT NOT NULL DEFAULT '',
			xmailer      TEXT NOT NULL DEFAULT '',
			content_type TEXT NOT NULL DEFAULT '',
			boundary     TEXT NOT NULL DEFAULT '',
			text_body    TEXT NOT NULL DEFAULT '',
			html_body    TEXT NOT NULL DEFAULT '',
			raw_message  TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE INDEX IF NOT EXISTS idx_mail_item_date_sent ON mail_item(date_sent)`,
		`CREATE INDEX IF NOT EXISTS idx_mail_item_subject   ON mail_item(subject)`,
		`CREATE INDEX IF NOT EXISTS idx_mail_item_from      ON mail_item(from_address)`,
		`CREATE TABLE IF NOT EXISTS attachment (
			id            TEXT PRIMARY KEY,
			mail_item_id  TEXT NOT NULL REFERENCES mail_item(id) ON DELETE CASCADE,
			file_name     TEXT NOT NULL DEFAULT '',
			content_type  TEXT NOT NULL DEFAULT '',
			content       BLOB NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_attachment_mail_item_id ON attachment(mail_item_id)`,
	}

	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func readSchemaVersion(ctx context.Context, db *sql.DB) (int, error) {
	var version int
	err := db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("reading schema version: %w", err)
	}
	return version, nil
}

func setSchemaVersion(ctx context.Context, db *sql.DB, version int) error {
	if _, err := db.ExecContext(ctx, `DELETE FROM schema_version`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (?)`, version); err != nil {
		return fmt.Errorf("writing schema version: %w", err)
	}
	return nil
}
