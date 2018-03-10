// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "net/mail"

/*
EmailValidationProvider is an interface for describing an email
validation service.
*/
type EmailValidationProvider interface {
	GetEmailComponents(email string) (*mail.Address, error)
	IsValidEmail(email string) bool
}
