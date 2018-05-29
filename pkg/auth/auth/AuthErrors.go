// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package auth

import (
	"fmt"
)

var ErrInvalidUserName = fmt.Errorf("Invalid user name")
var ErrInvalidPassword = fmt.Errorf("Invalid password")
