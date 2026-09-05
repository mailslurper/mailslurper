package mail

import (
	"strings"
	"testing"
)

func TestParsePlainText(t *testing.T) {
	raw := "From: sender@example.com\r\n" +
		"To: recipient@example.com\r\n" +
		"Subject: Hello\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Just a plain message.\r\n"

	item, err := Parse(raw, "sender@example.com", []string{"recipient@example.com"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if item.Subject != "Hello" {
		t.Errorf("Subject = %q", item.Subject)
	}
	if !strings.Contains(item.TextBody, "Just a plain message.") {
		t.Errorf("TextBody = %q", item.TextBody)
	}
	if item.HTMLBody != "" {
		t.Errorf("expected empty HTMLBody, got %q", item.HTMLBody)
	}
}

func TestParseMultipartAlternative(t *testing.T) {
	raw := "From: sender@example.com\r\n" +
		"To: recipient@example.com\r\n" +
		"Subject: Alt\r\n" +
		"Content-Type: multipart/alternative; boundary=\"B1\"\r\n" +
		"\r\n" +
		"--B1\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"plain version\r\n" +
		"--B1\r\n" +
		"Content-Type: text/html\r\n" +
		"\r\n" +
		"<p>html version</p>\r\n" +
		"--B1--\r\n"

	item, err := Parse(raw, "sender@example.com", []string{"recipient@example.com"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if !strings.Contains(item.TextBody, "plain version") {
		t.Errorf("TextBody = %q", item.TextBody)
	}
	if !strings.Contains(item.HTMLBody, "html version") {
		t.Errorf("HTMLBody = %q", item.HTMLBody)
	}
}

func TestParseMultipartMixedWithAttachment(t *testing.T) {
	raw := "From: sender@example.com\r\n" +
		"To: recipient@example.com\r\n" +
		"Subject: With attachment\r\n" +
		"Content-Type: multipart/mixed; boundary=\"B2\"\r\n" +
		"\r\n" +
		"--B2\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"body text\r\n" +
		"--B2\r\n" +
		"Content-Type: text/plain; name=\"note.txt\"\r\n" +
		"Content-Disposition: attachment; filename=\"note.txt\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" +
		"aGVsbG8gd29ybGQ=\r\n" +
		"--B2--\r\n"

	item, err := Parse(raw, "sender@example.com", []string{"recipient@example.com"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if !strings.Contains(item.TextBody, "body text") {
		t.Errorf("TextBody = %q", item.TextBody)
	}
	if len(item.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(item.Attachments))
	}
	att := item.Attachments[0]
	if att.FileName != "note.txt" {
		t.Errorf("FileName = %q", att.FileName)
	}
	if string(att.Content) != "hello world" {
		t.Errorf("Content = %q", string(att.Content))
	}
}

func TestParseMissingContentType(t *testing.T) {
	raw := "From: sender@example.com\r\n" +
		"To: recipient@example.com\r\n" +
		"Subject: No content type\r\n" +
		"\r\n" +
		"raw body\r\n"

	item, err := Parse(raw, "sender@example.com", []string{"recipient@example.com"})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if !strings.Contains(item.TextBody, "raw body") {
		t.Errorf("TextBody = %q", item.TextBody)
	}
}
