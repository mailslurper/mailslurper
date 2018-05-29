// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package jwt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"time"

	"golang.org/x/crypto/pbkdf2"

	"github.com/dgrijalva/jwt-go"
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
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.Config.AuthTimeoutInMinutes)).Unix(),
			Issuer:    JWTIssuer,
		},
		User: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(authSecret))
}

/*
DecryptToken takes a Base64 encoded token which has been encrypted
using AES-256 encryption. This returns the unencoded, unencrypted
token
*/
func (s *JWTService) DecryptToken(token string) (string, error) {
	var err error
	var aesBlock cipher.Block
	var unencodedToken []byte
	var gcm cipher.AEAD
	var nonce []byte
	var resultBytes []byte

	key := s.generateAESKey()

	if unencodedToken, err = base64.StdEncoding.DecodeString(token); err != nil {
		return "", errors.Wrapf(err, "Unable to base64 decode JWT token")
	}

	if aesBlock, err = aes.NewCipher(key); err != nil {
		return "", errors.Wrapf(err, "Unable to create AES cipher block")
	}

	if gcm, err = cipher.NewGCM(aesBlock); err != nil {
		return "", errors.Wrapf(err, "Problem creating GCM")
	}

	nonceSize := gcm.NonceSize()
	if len(unencodedToken) < nonceSize {
		return "", errors.Wrapf(err, "Ciphertext too short")
	}

	nonce, cipherText := unencodedToken[:nonceSize], unencodedToken[nonceSize:]

	if resultBytes, err = gcm.Open(nil, nonce, cipherText, nil); err != nil {
		return "", errors.Wrapf(err, "Problem decrypting token")
	}

	return string(resultBytes), nil
}

/*
EncryptToken takes a token string, encrypts it using AES-256,
then encodes it in Base64.
*/
func (s *JWTService) EncryptToken(token string) (string, error) {
	var err error
	var aesBlock cipher.Block
	var gcm cipher.AEAD
	var nonce []byte
	var encryptedResult []byte

	key := s.generateAESKey()

	if aesBlock, err = aes.NewCipher(key); err != nil {
		return "", errors.Wrapf(err, "Unable to create AES cipher block")
	}

	if gcm, err = cipher.NewGCM(aesBlock); err != nil {
		return "", errors.Wrapf(err, "Problem creating GCM")
	}

	nonce = make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	encryptedResult = gcm.Seal(nonce, nonce, []byte(token), nil)
	encodedResult := base64.StdEncoding.EncodeToString(encryptedResult)

	return encodedResult, nil
}

/*
GetUserFromToken retrieves the user from the claims in
a JWT token
*/
func (s *JWTService) GetUserFromToken(token *jwt.Token) string {
	var claims *Claims

	claims, _ = token.Claims.(*Claims)
	return claims.User
}

/*
Parse decrypts the provided token and returns a JWT token object
*/
func (s *JWTService) Parse(tokenFromHeader, authSecret string) (*jwt.Token, error) {
	var result *jwt.Token
	var decryptedToken string
	var err error

	/*
	 * Decrypt token first
	 */
	if decryptedToken, err = s.DecryptToken(tokenFromHeader); err != nil {
		return result, errors.Wrapf(err, "Problem decrypting JWT token in Parse")
	}

	if result, err = jwt.ParseWithClaims(decryptedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

/*
IsTokenValid returns an error if there are any issues with the
provided JWT token. Possible issues include:
	* Missing claims
	* Invalid token format
	* Invalid issuer
	* User doesn't have a corresponding entry in the credentials table
*/
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

func (s *JWTService) generateAESKey() []byte {
	return pbkdf2.Key([]byte(s.Config.AuthSecret), []byte(s.Config.AuthSalt), 4096, 32, sha1.New)
}

func (s *JWTService) pkcs5Padding(content []byte) []byte {
	padding := aes.BlockSize - len(content)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(content, padtext...)
}
