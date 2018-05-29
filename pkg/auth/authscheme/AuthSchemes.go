// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package authscheme

const (
	NONE  string = ""
	BASIC string = "basic"
)

func IsValidAuthScheme(authScheme string) bool {
	if authScheme != BASIC {
		return false
	}

	return true
}
