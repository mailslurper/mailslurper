// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
MailSearch is a set of criteria used to filter a mail collection
*/
type MailSearch struct {
	Message string
	Start   string
	End     string
	From    string
	To      string

	OrderByField     string
	OrderByDirection string
}
