// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "fmt"

/*
An ConnectionNotExistsError is used to alert a client that the specified
connection is not in the ConnectionManager pool
*/
type ConnectionNotExistsError struct {
	Address string
}

/*
ConnectionNotExists returns a new error object
*/
func ConnectionNotExists(address string) *ConnectionNotExistsError {
	return &ConnectionNotExistsError{
		Address: address,
	}
}

func (err *ConnectionNotExistsError) Error() string {
	return fmt.Sprintf("Connection '%s' is not in the connection manager pool", err.Address)
}
