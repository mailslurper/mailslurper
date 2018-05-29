// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"net"
)

/*
IConnectionManager describes an iterface for managing TCP connections
*/
type IConnectionManager interface {
	Close(connection net.Conn) error
	New(connection net.Conn) error
}
