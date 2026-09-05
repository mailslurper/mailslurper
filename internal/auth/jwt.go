package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrTokenInvalidated is returned by JWTService.Parse for a token that was
// explicitly logged out before it expired.
var ErrTokenInvalidated = errors.New("token has been invalidated")

// Claims is the JWT payload MylSlurper issues: just the authenticated
// username plus the standard expiry/issued-at claims.
type Claims struct {
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

// JWTService issues and validates bearer tokens for the JSON API.
type JWTService struct {
	secret []byte
	ttl    time.Duration
	cache  *invalidationCache
}

// NewJWTService returns a JWTService signing tokens with secret and expiring
// them after ttl.
func NewJWTService(secret string, ttl time.Duration) *JWTService {
	return &JWTService{secret: []byte(secret), ttl: ttl, cache: newInvalidationCache()}
}

// Issue returns a signed bearer token for userName.
func (j *JWTService) Issue(userName string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// Parse validates tokenString's signature and expiry, and rejects it if it
// was explicitly invalidated by Logout, returning the authenticated
// username on success.
func (j *JWTService) Parse(tokenString string) (string, error) {
	if j.cache.isInvalidated(tokenString) {
		return "", ErrTokenInvalidated
	}

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return j.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return "", fmt.Errorf("parsing token: %w", err)
	}

	return claims.UserName, nil
}

// Logout invalidates tokenString for the remainder of its natural lifetime.
func (j *JWTService) Logout(tokenString string) {
	j.cache.invalidate(tokenString, j.ttl)
}
