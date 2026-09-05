package mail

import "testing"

func TestItemBodyPrefersHTML(t *testing.T) {
	item := &Item{TextBody: "plain", HTMLBody: "<p>html</p>"}
	if got := item.Body(); got != "<p>html</p>" {
		t.Fatalf("Body() = %q", got)
	}
}

func TestItemBodyFallsBackToText(t *testing.T) {
	item := &Item{TextBody: "plain"}
	if got := item.Body(); got != "plain" {
		t.Fatalf("Body() = %q", got)
	}
}
