package basicauth

import "github.com/mailslurper/mailslurper/pkg/auth/auth"

/*
BasicAuthProvider offers in interface for authenticating
users with basic user name and password. These credentials
are stored in the config file. They are hashed for
security reasons.
*/
type BasicAuthProvider struct {
	CredentialMap   map[string]string     `json:"-"`
	Password        string                `json:"-"`
	PasswordService auth.IPasswordService `json:"-"`
	UserName        string                `json:"userName"`
}

/*
Login returns an error if the credential provided are invalid
*/
func (p *BasicAuthProvider) Login() error {
	var ok bool

	if _, ok = p.CredentialMap[p.UserName]; !ok {
		return auth.ErrInvalidUserName
	}

	if !p.PasswordService.IsPasswordValid([]byte(p.Password), []byte(p.CredentialMap[p.UserName])) {
		return auth.ErrInvalidPassword
	}

	return nil
}
