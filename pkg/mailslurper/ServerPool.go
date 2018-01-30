// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/adampresley/webframework/sanitizer"
)

/*
ServerPool represents a pool of SMTP workers. This will
manage how many workers may respond to SMTP client requests
and allocation of those workers.
*/
type ServerPool struct {
	logger *logrus.Entry
	pool   chan *SMTPWorker
}

/*
JoinQueue adds a worker to the queue.
*/
func (pool *ServerPool) JoinQueue(worker *SMTPWorker) {
	pool.pool <- worker
}

/*
NewServerPool creates a new server pool with a maximum number of SMTP
workers. An array of workers is initialized with an ID
and an initial state of SMTP_WORKER_IDLE.
*/
func NewServerPool(logger *logrus.Entry, maxWorkers int) *ServerPool {
	xssService := sanitizer.NewXSSService()
	emailValidationService := NewEmailValidationService()

	pool := &ServerPool{
		pool:   make(chan *SMTPWorker, maxWorkers),
		logger: logger,
	}

	for index := 0; index < maxWorkers; index++ {
		pool.JoinQueue(NewSMTPWorker(
			index+1,
			pool,
			emailValidationService,
			xssService,
			logger,
		))
	}

	logger.Infof("Worker pool configured for %d workers", maxWorkers)
	return pool
}

/*
NextWorker retrieves the next available worker from the queue.
*/
func (pool *ServerPool) NextWorker(connection net.Conn, receiver chan *MailItem) (*SMTPWorker, error) {
	select {
	case worker := <-pool.pool:
		worker.Prepare(
			connection,
			receiver,
			&SMTPReader{Connection: connection, logger: pool.logger},
			&SMTPWriter{Connection: connection, logger: pool.logger},
		)

		pool.logger.Infof("Worker %d queued to handle connections from %s", worker.WorkerID, connection.RemoteAddr().String())
		return worker, nil

	case <-time.After(time.Second * 2):
		return &SMTPWorker{}, NoWorkerAvailable()
	}
}
