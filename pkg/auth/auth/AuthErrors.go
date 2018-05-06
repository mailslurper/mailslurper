package auth

import (
	"fmt"
)

var ErrInvalidUserName = fmt.Errorf("Invalid user name")
var ErrInvalidPassword = fmt.Errorf("Invalid password")
