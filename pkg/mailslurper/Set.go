// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

/*
A Set is a set of header items. Most emails, bodies, and attachments have
more than one header to describe what the content is and how to handle it.
*/
type Set struct {
	HeaderItems []IItem
}

/*
NewHeaderSet takes a set of headers from a raw string, creates
a new Set, and returns it.
*/
func NewHeaderSet(headers string) (*Set, error) {
	result := &Set{}
	err := result.ParseHeaderString(headers)

	return result, err
}

/*
Get returns a header item by it's key name. If the key does not
exist in this set a MissingHeaderError is returned.
*/
func (set *Set) Get(headerName string) (IItem, error) {
	lowerCaseHeaderName := strings.ToLower(headerName)

	for _, header := range set.HeaderItems {
		if strings.ToLower(header.GetKey()) == lowerCaseHeaderName {
			return header, nil
		}
	}

	return nil, MissingHeader(headerName)
}

/*
ParseHeaderString will take a set of headers from a raw string
and parse them into a set of header items.
*/
func (set *Set) ParseHeaderString(headers string) error {
	var err error

	headerLines := strings.Split(strings.TrimSpace(set.UnfoldHeaders(headers)), "\r\n")
	set.HeaderItems = make([]IItem, len(headerLines))

	for index, item := range headerLines {
		set.HeaderItems[index] = &Item{}

		if err = set.HeaderItems[index].ParseHeaderString(item); err != nil {
			return errors.Wrapf(err, "Error parsing header #%d", index+1)
		}
	}

	return nil
}

/*
ToMap converts this structure to a map suiteable for use in a mail.Message
structure
*/
func (set *Set) ToMap() map[string][]string {
	result := make(map[string][]string)

	for _, header := range set.HeaderItems {
		result[header.GetKey()] = header.GetValues()
	}

	return result
}

/*
UnfoldHeaders "unfolds" headers that broken up into a single line.
The RFC-2822 defines "folding" as the process of breaking up large
header lines into multiple lines. Long Subject lines or Content-Type
lines (with boundaries) sometimes do this.
*/
func (set *Set) UnfoldHeaders(headers string) string {
	headerUnfolderRegex := regexp.MustCompile("(.*?)\r\n\\s{1}(.*?)\r\n")
	return headerUnfolderRegex.ReplaceAllString(headers, "$1 $2\r\n")
}
