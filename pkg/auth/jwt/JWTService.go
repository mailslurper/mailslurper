package jwt

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	jwtrequest "github.com/dgrijalva/jwt-go/request"
	"github.com/mailslurper/mailslurper/pkg/mailslurper"
	"github.com/pkg/errors"
)

/*
JWTService provides methods for working with
JWTs in MailSlurper
*/
type JWTService struct {
	Config *mailslurper.Configuration
}

/*
CreateToken creates a new JWT token for use in
MailSlurper services
*/
func (s *JWTService) CreateToken(authSecret, user string) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    JWTIssuer,
		},
		User: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(authSecret))
}

/*
Parse extracts a token from an HTTP request
*/
func (s *JWTService) Parse(request *http.Request, authSecret string) (*jwt.Token, error) {
	var result *jwt.Token
	var err error

	if result, err = jwtrequest.ParseFromRequestWithClaims(request, jwtrequest.AuthorizationHeaderExtractor, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		var ok bool

		if _, ok = token.Method.(*jwt.SigningMethodHMAC); !ok {
			return result, ErrInvalidToken
		}

		return []byte(authSecret), nil
	}); err != nil {
		return result, errors.Wrapf(err, "Problem parsing JWT token")
	}

	if err = s.IsTokenValid(result); err != nil {
		return result, err
	}

	return result, nil
}

func (s *JWTService) IsTokenValid(token *jwt.Token) error {
	var claims *Claims
	var ok bool

	claims, ok = token.Claims.(*Claims)

	if !ok {
		return ErrTokenMissingClaims
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	if claims.Issuer != JWTIssuer {
		return ErrInvalidIssuer
	}

	if _, ok = s.Config.Credentials[claims.User]; !ok {
		return ErrInvalidUser
	}

	return nil
}

func (s *JWTService) GetUserFromToken(token *jwt.Token) string {
	var claims *Claims

	claims, _ = token.Claims.(*Claims)
	return claims.User
}
