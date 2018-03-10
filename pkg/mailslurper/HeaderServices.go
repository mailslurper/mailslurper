// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"regexp"
)

/*
UnfoldHeaders follows the RFC-2822, which defines "folding" as the process of breaking up large
header lines into multiple lines. Long Subject lines or Content-Type
lines (with boundaries) sometimes do this. This function will "unfold"
them into a single line.
*/
func UnfoldHeaders(contents string) string {
	headerUnfolderRegex := regexp.MustCompile("(.*?)\r\n\\s{1}(.*?)\r\n")
	return headerUnfolderRegex.ReplaceAllString(contents, "$1 $2\r\n")
}
