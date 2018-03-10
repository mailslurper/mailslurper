// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
An ISet is a set of header items. Most emails, bodies, and attachments have
more than one header to describe what the content is and how to handle it.
*/
type ISet interface {
	Get(headerName string) (IItem, error)
	ParseHeaderString(headers string) error
	ToMap() map[string][]string
	UnfoldHeaders(headers string) string
}
