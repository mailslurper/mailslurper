// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "fmt"

/*
An ConnectionExistsError is used to alert a client that there is already a
connection by this address cached
*/
type ConnectionExistsError struct {
	Address string
}

/*
ConnectionExists returns a new error object
*/
func ConnectionExists(address string) *ConnectionExistsError {
	return &ConnectionExistsError{
		Address: address,
	}
}

func (err *ConnectionExistsError) Error() string {
	return fmt.Sprintf("Connection on '%s' already exists", err.Address)
}
