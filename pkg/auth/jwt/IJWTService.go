package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type IJWTService interface {
	CreateToken(authSecret, user string) (string, error)
	DecryptToken(token string) (string, error)
	EncryptToken(token string) (string, error)
	GetUserFromToken(token *jwt.Token) string
	Parse(tokenFromHeader, authSecret string) (*jwt.Token, error)
	IsTokenValid(token *jwt.Token) error
}
