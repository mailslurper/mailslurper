package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CookieName is the name of the signed session cookie used by the browser
// shell.
const CookieName = "mylslurper_session"

var ErrInvalidSession = errors.New("invalid or expired session")

// SessionService issues and verifies small HMAC-signed cookie values
// ("userName|expiry" + signature) without depending on a third-party
// session library.
type SessionService struct {
	secret []byte
	ttl    time.Duration
}

// NewSessionService returns a SessionService signing cookies with secret and
// expiring them after ttl.
func NewSessionService(secret string, ttl time.Duration) *SessionService {
	return &SessionService{secret: []byte(secret), ttl: ttl}
}

// Create returns a signed cookie value authenticating userName.
func (s *SessionService) Create(userName string) string {
	payload := fmt.Sprintf("%s|%d", userName, time.Now().Add(s.ttl).Unix())
	sig := s.sign(payload)
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + base64.RawURLEncoding.EncodeToString(sig)
}

// Verify checks value's signature and expiry, returning the authenticated
// username on success.
func (s *SessionService) Verify(value string) (string, error) {
	parts := strings.SplitN(value, ".", 2)
	if len(parts) != 2 {
		return "", ErrInvalidSession
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", ErrInvalidSession
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ErrInvalidSession
	}

	if !hmac.Equal(sig, s.sign(string(payloadBytes))) {
		return "", ErrInvalidSession
	}

	fields := strings.SplitN(string(payloadBytes), "|", 2)
	if len(fields) != 2 {
		return "", ErrInvalidSession
	}
	expiry, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return "", ErrInvalidSession
	}
	if time.Now().Unix() > expiry {
		return "", ErrInvalidSession
	}

	return fields[0], nil
}

func (s *SessionService) sign(payload string) []byte {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(payload))
	return mac.Sum(nil)
}
