package mailslurper

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/*
A ConnectionManager is responsible for maintaining, closing, and cleaning
client connections. For every connection there is a worker. After an idle
timeout period the manager will forceably close a client connection.
*/
type ConnectionManager struct {
	closeChannel    chan string
	config          *Configuration
	connectionPool  ConnectionPool
	ctx             context.Context
	logger          *logrus.Entry
	mailItemChannel chan *MailItem
	serverPool      *ServerPool
}

/*
NewConnectionManager creates a new struct
*/
func NewConnectionManager(logger *logrus.Entry, config *Configuration, ctx context.Context, mailItemChannel chan *MailItem, serverPool *ServerPool) *ConnectionManager {
	return &ConnectionManager{
		closeChannel:    make(chan string, 5),
		config:          config,
		connectionPool:  NewConnectionPool(),
		ctx:             ctx,
		logger:          logger,
		mailItemChannel: mailItemChannel,
		serverPool:      serverPool,
	}
}

/*
CleanIdle cleans up client connections that have exceeded the
idle timeout period. It does this by forcably closing them
*/
func (m *ConnectionManager) CleanIdle() error {
	return nil
}

/*
Close will close a client connection. The state of the worker
is only used for logging purposes
*/
func (m *ConnectionManager) Close(connection net.Conn, state SMTPWorkerState) error {
	if m.connectionExistsInPool(connection) {
		return m.connectionPool[connection.RemoteAddr().String()].Connection.Close()
	}

	return ConnectionNotExists(connection.RemoteAddr().String())
}

func (m *ConnectionManager) connectionExistsInPool(connection net.Conn) bool {
	if _, ok := m.connectionPool[connection.RemoteAddr().String()]; ok {
		return true
	}

	return false
}

/*
New attempts to track a new client connection. The SMTPListener will
use this to track a client connection and its worker
*/
func (m *ConnectionManager) New(connection net.Conn) error {
	var err error
	var worker *SMTPWorker

	if m.connectionExistsInPool(connection) {
		return ConnectionExists(connection.RemoteAddr().String())
	}

	if worker, err = m.serverPool.NextWorker(connection, m.mailItemChannel, m.ctx, m.closeChannel); err != nil {
		connection.Close()
		m.logger.WithError(err).Errorf("Error getting next SMTP worker")
		return errors.Wrapf(err, "Error getting work in ConnectionManager")
	}

	m.connectionPool[connection.RemoteAddr().String()] = NewConnectionPoolItem(connection, worker)
	go m.connectionPool[connection.RemoteAddr().String()].Worker.Work()

	return nil
}
