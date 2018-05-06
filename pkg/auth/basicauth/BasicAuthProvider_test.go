package basicauth_test

import (
	"testing"

	"github.com/mailslurper/mailslurper/pkg/auth/auth"
	"github.com/mailslurper/mailslurper/pkg/auth/basicauth"
)

func TestLogin(t *testing.T) {
	userName := "adam"
	password := "password"
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
		Password:        password,
		UserName:        userName,
		PasswordService: mockPasswordService,
	}

	err = provider.Login()

	if err != nil {
		t.Errorf("Expected error to be nil, got %s instead", err.Error())
	}
}

func TestLoginInvalidUserName(t *testing.T) {
	userName := "bob"
	password := "password"
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
		Password:        password,
		UserName:        userName,
		PasswordService: mockPasswordService,
	}

	err = provider.Login()

	if err != auth.ErrInvalidUserName {
		t.Errorf("Expected error to be ErrInvalidUserName, got %s instead", err.Error())
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	userName := "adam"
	password := "password"
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
		Password:        password,
		UserName:        userName,
		PasswordService: mockPasswordService,
	}

	err = provider.Login()

	if err != auth.ErrInvalidPassword {
		t.Errorf("Expected error to be ErrInvalidPassword, got %s instead", err.Error())
	}
}
