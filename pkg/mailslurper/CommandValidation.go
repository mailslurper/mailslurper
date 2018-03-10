// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"strings"
)

/*
IsValidCommand returns an error if the input stream does not contain the expected command.
The input and expected commands are lower cased, as we do not care about
case when comparing.
*/
func IsValidCommand(streamInput, expectedCommand string) error {
	check := strings.Index(strings.ToLower(streamInput), strings.ToLower(expectedCommand))

	if check < 0 {
		return InvalidCommand(expectedCommand)
	}

	return nil
}
