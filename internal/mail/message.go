// Package mail holds MylSlurper's mail domain model and MIME parser.
package mail

import "time"

// Item is a single received email, fully parsed and ready to persist.
type Item struct {
	ID          string
	DateSent    time.Time
	From        string
	To          []string
	Subject     string
	XMailer     string
	ContentType string
	Boundary    string
	TextBody    string
	HTMLBody    string
	RawMessage  string
	Attachments []Attachment
}

// Attachment is a single file attached to a mail Item.
type Attachment struct {
	ID          string
	MailItemID  string
	FileName    string
	ContentType string
	Content     []byte
	// Size is the attachment's byte size. It's always populated, even when
	// Content itself hasn't been loaded (e.g. for list/detail views that
	// only need metadata, not the attachment's bytes).
	Size int
}

// Body returns the best available body for display: HTML if present,
// otherwise plain text.
func (i *Item) Body() string {
	if i.HTMLBody != "" {
		return i.HTMLBody
	}
	return i.TextBody
}
