// Package logging configures MylSlurper's structured logger.
package logging

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// New returns a logrus.Logger configured with the given level and format
// ("json" or "text", defaulting to text).
func New(level, format string) *logrus.Logger {
	log := logrus.New()

	if format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	parsed, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		parsed = logrus.InfoLevel
	}
	log.SetLevel(parsed)

	return log
}
