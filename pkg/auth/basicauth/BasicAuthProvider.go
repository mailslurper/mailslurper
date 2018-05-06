package basicauth

import (
	"github.com/mailslurper/mailslurper/pkg/auth/auth"
)

/*
BasicAuthProvider offers in interface for authenticating
users with basic user name and password. These credentials
are stored in the config file. They are hashed for
security reasons.
*/
type BasicAuthProvider struct {
	CredentialMap   map[string]string
	PasswordService auth.IPasswordService
}

/*
Login returns an error if the credential provided are invalid
*/
func (p *BasicAuthProvider) Login(credentials *auth.AuthCredentials) error {
	var ok bool

	if _, ok = p.CredentialMap[credentials.UserName]; !ok {
		return auth.ErrInvalidUserName
	}

	if !p.PasswordService.IsPasswordValid([]byte(credentials.Password), []byte(p.CredentialMap[credentials.UserName])) {
		return auth.ErrInvalidPassword
	}

	return nil
}
