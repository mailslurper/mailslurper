// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
A Page is the basis for all displayed HTML pages
*/
type Page struct {
	PublicWWWURL string
	Error        bool
	Message      string
	Theme        string
	Title        string
	User         string
}
