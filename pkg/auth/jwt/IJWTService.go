package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type IJWTService interface {
	CreateToken(authSecret string) *jwt.Token
}
