// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "fmt"

/*
An InvalidCommandFormatError is used to alert a client that the command passed in
has an invalid format
*/
type InvalidCommandFormatError struct {
	InvalidCommand string
}

/*
InvalidCommandFormat returns a new error object
*/
func InvalidCommandFormat(command string) *InvalidCommandFormatError {
	return &InvalidCommandFormatError{
		InvalidCommand: command,
	}
}

func (err *InvalidCommandFormatError) Error() string {
	return fmt.Sprintf("%s command format is invalid", err.InvalidCommand)
}
