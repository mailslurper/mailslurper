// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package auth

/*
IPasswordService is an interface for validating passwords
*/
type IPasswordService interface {
	IsPasswordValid(password, storedPassword []byte) bool
}
