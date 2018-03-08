package mailslurper

import (
	"net"
)

/*
IConnectionManager describes an iterface for managing TCP connections
*/
type IConnectionManager interface {
	CleanIdle() error
	Close(connection net.Conn, state SMTPWorkerState) error
	New(connection net.Conn) error
}
