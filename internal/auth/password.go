// Package auth implements MylSlurper's authentication: bcrypt credential
// checks, JWT bearer tokens for the API, and signed session cookies for the
// browser shell.
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a bcrypt hash of password, suitable for storing in a
// config.json credentials map.
func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// CheckPassword reports whether password matches hash.
func CheckPassword(hash, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, password) == nil
}
