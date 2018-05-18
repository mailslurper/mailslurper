package jwt

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var JWTIssuer string = "mailslurper"
var ErrInvalidToken error = fmt.Errorf("Invalid token")
var ErrTokenMissingClaims error = fmt.Errorf("Token is missing claims")
var ErrInvalidUser error = fmt.Errorf("Invalid user")
var ErrInvalidIssuer error = fmt.Errorf("Invalid issuer")

type Claims struct {
	jwt.StandardClaims
	User string `json:"user"`
}
