// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package auth

/*
IAuthProvider describes a provider of authentication services, such
as Basic, LDAP, etc...
*/
type IAuthProvider interface {
	Login(credentials *AuthCredentials) error
}
