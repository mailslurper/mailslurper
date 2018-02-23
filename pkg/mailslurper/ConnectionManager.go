package mailslurper

import "net"

/*
A ConnectionManager is responsible for maintaining, closing, and cleaning
client connections. For every connection there is a worker. After an idle
timeout period the manager will forceably close a client connection.
*/
type ConnectionManager struct {
}

/*
NewConnectionManager creates a new struct
*/
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{}
}

/*
CleanIdle cleans up client connections that have exceeded the
idle timeout period. It does this by forcably closing them
*/
func (m *ConnectionManager) CleanIdle() error {

}

/*
Close will close a client connection. The state of the worker
is only used for logging purposes
*/
func (m *ConnectionManager) Close(state SMTPWorkerState) error {

}

/*
New attempts to track a new client connection. The SMTPListener will
use this to track a client connection and its worker
*/
func (m *ConnectionManager) New(connection net.Conn, worker *SMTPWorker) error {

}
