// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"strings"
)

/*
An Item represents a single header entry. Headers describe emails, bodies,
and attachments. They are in the form of "Key: Value".
*/
type Item struct {
	Key    string
	Values []string
}

/*
GetKey returns the key for this header item
*/
func (headerItem *Item) GetKey() string {
	return headerItem.Key
}

/*
GetValues returns the value for this header item
*/
func (headerItem *Item) GetValues() []string {
	return headerItem.Values
}

/*
ParseHeaderString takes a header string and parses it into a key and value(s)
*/
func (headerItem *Item) ParseHeaderString(header string) error {
	headerSplit := strings.Split(header, ":")
	if len(headerSplit) < 2 {
		return InvalidHeader(header)
	}

	headerItem.Key = strings.TrimSpace(headerSplit[0])
	headerItem.Values = append(headerItem.Values, strings.TrimSpace(strings.Join(headerSplit[1:], "")))
	return nil
}
