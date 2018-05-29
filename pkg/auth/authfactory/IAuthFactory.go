// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package authfactory

import "github.com/mailslurper/mailslurper/pkg/auth/auth"

/*
IAuthFactory, when implmented, returns the correction authorization
provider based on a configuration setting
*/
type IAuthFactory interface {
	Get() auth.IAuthProvider
}
