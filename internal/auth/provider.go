package auth

import "fmt"

// ErrInvalidCredentials is returned by CheckCredentials when the username
// doesn't exist or the password doesn't match.
var ErrInvalidCredentials = fmt.Errorf("invalid username or password")

// CheckCredentials validates userName/password against the bcrypt hashes in
// credentials (config.json's "credentials" map: userName -> bcrypt hash).
func CheckCredentials(credentials map[string]string, userName, password string) error {
	hash, ok := credentials[userName]
	if !ok {
		return ErrInvalidCredentials
	}
	if !CheckPassword([]byte(hash), []byte(password)) {
		return ErrInvalidCredentials
	}
	return nil
}
