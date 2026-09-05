// Package web embeds MylSlurper's static frontend (plain HTML/CSS/JS, no
// build step) into the compiled binary.
package web

import (
	"embed"
	"io/fs"
)

//go:embed static
var embedded embed.FS

// Assets returns the embedded frontend as an fs.FS rooted at its content
// (i.e. index.html is at the root, not under "static/").
func Assets() (fs.FS, error) {
	return fs.Sub(embedded, "static")
}
