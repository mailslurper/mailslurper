// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "github.com/sirupsen/logrus"

/*
GetLogger returns a logger with who owns the logger
attached to it
*/
func GetLogger(logLevel, logFormat, who string) *logrus.Entry {
	l := logrus.New()

	switch logLevel {
	case "debug":
		l.SetLevel(logrus.DebugLevel)

	case "error":
		l.SetLevel(logrus.ErrorLevel)

	default:
		l.SetLevel(logrus.InfoLevel)
	}

	switch logFormat {
	case "json":
		l.Formatter = &logrus.JSONFormatter{}

	default:
		l.Formatter = &logrus.TextFormatter{}
	}

	e := l.WithField("who", who)
	return e
}
