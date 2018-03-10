// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "fmt"

/*
An InvalidHeaderError is used to alert a client that a header is malformed.
*/
type InvalidHeaderError struct {
	InvalidHeader string
}

/*
InvalidHeader returns a new error object
*/
func InvalidHeader(header string) *InvalidHeaderError {
	return &InvalidHeaderError{
		InvalidHeader: header,
	}
}

func (err *InvalidHeaderError) Error() string {
	return fmt.Sprintf("Invalid header '%s'", err.InvalidHeader)
}
