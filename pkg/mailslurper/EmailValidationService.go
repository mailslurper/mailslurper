// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "net/mail"

/*
EmailValidationService realizes the EmailValidationProvider
interface by offering functions for working with email validation
and manipulation.
*/
type EmailValidationService struct {
}

/*
GetEmailComponents returns an object with all the parts of a parsed email address
*/
func (service *EmailValidationService) GetEmailComponents(email string) (*mail.Address, error) {
	return mail.ParseAddress(email)
}

/*
IsValidEmail returns true if the provided email is valid and parses
*/
func (service *EmailValidationService) IsValidEmail(email string) bool {
	_, err := service.GetEmailComponents(email)
	return err == nil
}

/*
NewEmailValidationService creates a new object
*/
func NewEmailValidationService() *EmailValidationService {
	return &EmailValidationService{}
}
