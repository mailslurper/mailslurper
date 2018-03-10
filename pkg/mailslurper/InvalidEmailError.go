// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "fmt"

/*
An InvalidEmailError is used to alert a client that an email address is invalid
*/
type InvalidEmailError struct {
	Email string
}

/*
InvalidEmail returns a new error object
*/
func InvalidEmail(email string) *InvalidEmailError {
	return &InvalidEmailError{
		Email: email,
	}
}

func (err *InvalidEmailError) Error() string {
	return fmt.Sprintf("The provided email address, '%s', is invalid", err.Email)
}
