// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"fmt"
	"strings"
)

/*
An SMTPCommand represents a command issued over a TCP connection
*/
type SMTPCommand int

/*
Constants representing the commands that an SMTP client will
send during the course of communicating with our server.
*/
const (
	NONE SMTPCommand = iota
	RCPT SMTPCommand = iota
	MAIL SMTPCommand = iota
	HELO SMTPCommand = iota
	RSET SMTPCommand = iota
	DATA SMTPCommand = iota
	QUIT SMTPCommand = iota
)

/*
SMTPCommands is a map of SMTP command strings to their int
representation. This is primarily used because there can
be more than one command to do the same things. For example,
a client can send "helo" or "ehlo" to initiate the handshake.
*/
var SMTPCommands = map[string]SMTPCommand{
	"helo":      HELO,
	"ehlo":      HELO,
	"rcpt to":   RCPT,
	"mail from": MAIL,
	"send":      MAIL,
	"rset":      RSET,
	"quit":      QUIT,
	"data":      DATA,
}

/*
SMTPCommandsToStrings is a friendly string representations of commands. Useful in error
reporting.
*/
var SMTPCommandsToStrings = map[SMTPCommand]string{
	HELO: "HELO",
	RCPT: "RCPT TO",
	MAIL: "SEND",
	RSET: "RSET",
	QUIT: "QUIT",
	DATA: "DATA",
}

/*
GetCommandFromString takes a string and returns the integer command representation. For example
if the string contains "DATA" then the value 1 (the constant DATA) will be returned.
*/
func GetCommandFromString(input string) (SMTPCommand, error) {
	result := NONE
	input = strings.ToLower(input)

	if input == "" {
		return result, nil
	}

	for key, value := range SMTPCommands {
		if strings.Index(input, key) == 0 {
			result = value
			break
		}
	}

	if result == NONE {
		return result, fmt.Errorf("Command '%s' not found", input)
	}

	return result, nil
}

/*
Returns the string representation of a command.
*/
func (smtpCommand SMTPCommand) String() string {
	return SMTPCommandsToStrings[smtpCommand]
}
