package basicauth

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

/*
PasswordService provides methods for working with passwords
*/
type PasswordService struct{}

/*
HashPassword hashes a password
*/
func (s *PasswordService) HashPassword(password []byte) ([]byte, error) {
	var result []byte
	var err error

	if result, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost); err != nil {
		return result, errors.Wrapf(err, "Problem hashing password using BCRYPT")
	}

	return result, nil
}

/*
IsPasswordValid returns true if the provided password matches the
hashed password
*/
func (s *PasswordService) IsPasswordValid(password, hashedPassword []byte) bool {
	var err error

	if err = bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
		return false
	}

	return true

}
