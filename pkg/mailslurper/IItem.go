// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
IItem represents a single header entry. Headers describe emails, bodies,
and attachments. They are in the form of "Key: Value".
*/
type IItem interface {
	GetKey() string
	GetValues() []string
	ParseHeaderString(header string) error
}
