package mailslurper

import (
	"strings"
)

/*
MailAddressCollection is a set of email address
*/
type MailAddressCollection []string

/*
NewMailAddressCollection returns a new MailAddressCollection
*/
func NewMailAddressCollection() MailAddressCollection {
	return make(MailAddressCollection, 0, 5)
}

/*
NewMailAddressCollectionFromStringList takes a list of delimited email address and
breaks it into a collection of mail addresses
*/
func NewMailAddressCollectionFromStringList(addresses string) MailAddressCollection {
	split := strings.Split(addresses, "; ")
	result := NewMailAddressCollection()

	for _, s := range split {
		result = append(result, s)
	}

	return result
}
