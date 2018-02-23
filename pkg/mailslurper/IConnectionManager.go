package mailslurper

import (
	"net"
)

/*
IConnectionManager describes an iterface for managing TCP connections
*/
type IConnectionManager interface {
	CleanIdle() error
	Close(state SMTPWorkerState) error
	New(connection net.Conn, worker *SMTPWorker) error
}
