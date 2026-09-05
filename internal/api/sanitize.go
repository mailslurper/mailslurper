package api

import "github.com/microcosm-cc/bluemonday"

// htmlSanitizer strips scripts, event handlers, and other unsafe markup
// from mail HTML bodies before they're ever sent to the browser. The
// frontend additionally renders HTML bodies inside a script-disabled
// sandboxed iframe — this sanitizer and that sandbox are complementary
// defenses, not substitutes for each other.
var htmlSanitizer = bluemonday.UGCPolicy()

func sanitizeHTML(html string) string {
	if html == "" {
		return html
	}
	return htmlSanitizer.Sanitize(html)
}
