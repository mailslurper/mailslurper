// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"strings"
)

/*
GetCommandValue splits an input by colon (:) and returns the right hand side.
If there isn't a split, or a missing colon, an InvalidCommandFormatError is
returned.
*/
func GetCommandValue(streamInput, command, delimiter string) (string, error) {
	split := strings.Split(streamInput, delimiter)

	if len(split) < 2 {
		return "", InvalidCommandFormat(command)
	}

	return strings.TrimSpace(strings.Join(split[1:], "")), nil
}
