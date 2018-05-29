// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package basicauth_test

import (
	"testing"

	"github.com/mailslurper/mailslurper/pkg/auth/auth"
	"github.com/mailslurper/mailslurper/pkg/auth/basicauth"
)

func TestLogin(t *testing.T) {
	credentials := &auth.AuthCredentials{
		UserName: "adam",
		Password: "password",
	}
	hashedPassword := "abcdefg"
	var err error

	mockPasswordService := &basicauth.MockPasswordService{
		FnIsPasswordValid: func(password, storedPassword []byte) bool {
			return true
		},
	}

	provider := basicauth.BasicAuthProvider{
		CredentialMap: map[string]string{
			"adam": hashedPassword,
		},
		PasswordService: mockPasswordService,
	}

	err = provider.Login(credentials)

	if err != nil {
		t.Errorf("Expected error to be nil, got %s instead", err.Error())
	}
}

func TestLoginInvalidUserName(t *testing.T) {
	credentials := &auth.AuthCredentials{
		UserName: "bob",
		Password: "password",
	}
	hashedPassword := "abcdefg"
	var err error

	mockPasswordService := &basicauth.MockPasswordService{
		FnIsPasswordValid: func(password, storedPassword []byte) bool {
			return true
		},
	}

	provider := basicauth.BasicAuthProvider{
		CredentialMap: map[string]string{
			"adam": hashedPassword,
		},
		PasswordService: mockPasswordService,
	}

	err = provider.Login(credentials)

	if err != auth.ErrInvalidUserName {
		t.Errorf("Expected error to be ErrInvalidUserName, got %s instead", err.Error())
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	credentials := &auth.AuthCredentials{
		UserName: "adam",
		Password: "password",
	}
	hashedPassword := "abcdefg"
	var err error

	mockPasswordService := &basicauth.MockPasswordService{
		FnIsPasswordValid: func(password, storedPassword []byte) bool {
			return false
		},
	}

	provider := basicauth.BasicAuthProvider{
		CredentialMap: map[string]string{
			"adam": hashedPassword,
		},
		PasswordService: mockPasswordService,
	}

	err = provider.Login(credentials)

	if err != auth.ErrInvalidPassword {
		t.Errorf("Expected error to be ErrInvalidPassword, got %s instead", err.Error())
	}
}
