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
