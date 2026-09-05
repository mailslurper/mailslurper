// Package idgen generates random identifiers used as primary keys.
package idgen

import (
	"crypto/rand"
	"fmt"
)

// New returns a random UUIDv4-formatted string.
func New() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err) // crypto/rand.Read failing means the OS's CSPRNG is broken
	}

	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
