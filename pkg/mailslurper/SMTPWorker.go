// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/adampresley/webframework/sanitizer"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

/*
An SMTPWorker is responsible for executing, parsing, and processing a single
TCP connection's email.
*/
type SMTPWorker struct {
	Connection             net.Conn
	EmailValidationService EmailValidationProvider
	Error                  error
	Reader                 *SMTPReader
	Receiver               chan *MailItem
	State                  SMTPWorkerState
	WorkerID               int
	Writer                 *SMTPWriter
	XSSService             sanitizer.IXSSServiceProvider

	connectionCloseChannel chan net.Conn
	killServerContext      context.Context
	pool                   *ServerPool
	logger                 *logrus.Entry
	logLevel               string
	logFormat              string
}

type smtpCommand struct {
	Command     SMTPCommand
	StreamInput string
}

/*
NewSMTPWorker creates a new SMTP worker. An SMTP worker is
responsible for parsing and working with SMTP mail data.
*/
func NewSMTPWorker(
	workerID int,
	pool *ServerPool,
	emailValidationService EmailValidationProvider,
	xssService sanitizer.IXSSServiceProvider,
	logger *logrus.Entry,
	logLevel string,
	logFormat string,
) *SMTPWorker {
	return &SMTPWorker{
		EmailValidationService: emailValidationService,
		WorkerID:               workerID,
		State:                  SMTP_WORKER_IDLE,
		XSSService:             xssService,

		pool:      pool,
		logger:    logger,
		logLevel:  logLevel,
		logFormat: logFormat,
	}
}

/*
Prepare tells a worker about the TCP connection they will work
with, the IO handlers, and sets their state.
*/
func (smtpWorker *SMTPWorker) Prepare(
	connection net.Conn,
	receiver chan *MailItem,
	reader *SMTPReader,
	writer *SMTPWriter,
	killServerContext context.Context,
	connectionCloseChannel chan net.Conn,
) {
	smtpWorker.State = SMTP_WORKER_WORKING

	smtpWorker.Connection = connection
	smtpWorker.Receiver = receiver

	smtpWorker.Reader = reader
	smtpWorker.Writer = writer

	smtpWorker.connectionCloseChannel = connectionCloseChannel
	smtpWorker.killServerContext = killServerContext
}

func (smtpWorker *SMTPWorker) rejoinWorkerQueue() {
	smtpWorker.pool.JoinQueue(smtpWorker)
}

/*
Work is the function called by the SmtpListener when a client request
is received. This will start the process by responding to the client,
start processing commands, and finally close the connection.
*/
func (smtpWorker *SMTPWorker) Work() {
	var err error

	smtpWorker.Writer.SayHello()
	mailItem := NewEmptyMailItem(smtpWorker.logger)

	quitChannel := make(chan bool, 2)
	quitCommandChannel := make(chan bool, 2)
	workerErrorChannel := make(chan error, 2)
	commandChannel := make(chan smtpCommand)
	commandDoneChannel := make(chan error)

	/*
	 * This goroutine is the command processor
	 */
	go func() {
		var streamInput string
		var command SMTPCommand
		var err error
		var networkError net.Error
		var ok bool

		for {
			select {
			case <-quitCommandChannel:
				return

			default:
				if streamInput, err = smtpWorker.Reader.Read(); err != nil {
					if networkError, ok = err.(net.Error); ok {
						if networkError.Timeout() {
							smtpWorker.logger.WithField("connection", smtpWorker.Connection.RemoteAddr().String()).Infof("Connection inactivity timeout")

							quitCommandChannel <- true
							quitChannel <- true
							break
						}
					}

					workerErrorChannel <- err
					break
				}

				if command, err = GetCommandFromString(streamInput); err != nil {
					smtpWorker.logger.WithError(err).WithField("input", streamInput).Errorf("Problem finding command from input")
					workerErrorChannel <- errors.Wrapf(err, "Problem finding command from input %s", streamInput)
					break
				}

				if command == QUIT {
					quitCommandChannel <- true
					quitChannel <- true
					break
				}

				commandChannel <- smtpCommand{Command: command, StreamInput: streamInput}
				err = <-commandDoneChannel

				if err != nil {
					smtpWorker.logger.WithError(err).Errorf("Error executing command")
					quitCommandChannel <- true
				}
			}
		}
	}()

	for {
		select {
		case <-smtpWorker.killServerContext.Done():
			smtpWorker.State = SMTP_WORKER_DONE
			smtpWorker.Writer.SayGoodbye()
			smtpWorker.connectionCloseChannel <- smtpWorker.Connection
			break

		case <-quitChannel:
			smtpWorker.logger.WithField("connection", smtpWorker.Connection.RemoteAddr().String()).Infof("QUIT command received")
			smtpWorker.Writer.SayGoodbye()

			smtpWorker.State = SMTP_WORKER_DONE
			smtpWorker.connectionCloseChannel <- smtpWorker.Connection
			smtpWorker.rejoinWorkerQueue()

			break

		case workerError := <-workerErrorChannel:
			smtpWorker.State = SMTP_WORKER_ERROR
			smtpWorker.Error = workerError
			smtpWorker.Writer.SayGoodbye()

			smtpWorker.connectionCloseChannel <- smtpWorker.Connection
			smtpWorker.rejoinWorkerQueue()
			break

		case command := <-commandChannel:
			if command.Command == QUIT {
				quitChannel <- true
				continue
			}

			executor := smtpWorker.getExecutorFromCommand(command.Command)
			command.StreamInput = strings.TrimSpace(command.StreamInput)

			if err = executor.Process(command.StreamInput, mailItem); err != nil {
				smtpWorker.logger.WithError(err).WithFields(logrus.Fields{"command": command.Command.String(), "input": command.StreamInput}).Errorf("Problem executing command")
				workerErrorChannel <- errors.Wrapf(err, "Problem executing command %s (stream input == '%s')", command.Command.String(), command.StreamInput)

				commandDoneChannel <- err
				continue
			}

			if command.Command == DATA {
				copy := NewEmptyMailItem(smtpWorker.logger)
				copier.Copy(copy, mailItem)
				smtpWorker.Receiver <- copy

				mailItem = NewEmptyMailItem(smtpWorker.logger)
			}

			commandDoneChannel <- nil
		}
	}
}

func (smtpWorker *SMTPWorker) getExecutorFromCommand(command SMTPCommand) ICommandExecutor {
	switch command {
	case MAIL:
		return NewMailCommandExecutor(
			GetLogger(smtpWorker.logLevel, smtpWorker.logFormat, "MAIL Command Executor"),
			smtpWorker.Reader,
			smtpWorker.Writer,
			smtpWorker.EmailValidationService,
			smtpWorker.XSSService,
		)

	case RCPT:
		return NewRcptCommandExecutor(
			GetLogger(smtpWorker.logLevel, smtpWorker.logFormat, "RCPT TO Command Executor"),
			smtpWorker.Reader,
			smtpWorker.Writer,
			smtpWorker.EmailValidationService,
			smtpWorker.XSSService,
		)

	case DATA:
		return NewDataCommandExecutor(
			GetLogger(smtpWorker.logLevel, smtpWorker.logFormat, "DATA Command Executor"),
			smtpWorker.Reader,
			smtpWorker.Writer,
			smtpWorker.EmailValidationService,
			smtpWorker.XSSService,
		)

	default:
		return NewHelloCommandExecutor(
			GetLogger(smtpWorker.logLevel, smtpWorker.logFormat, "HELO Command Executor"),
			smtpWorker.Reader,
			smtpWorker.Writer,
		)
	}
}

/*
TimeoutHasExpired determines if the time elapsed since a start time has exceeded
the command timeout.
*/
func (smtpWorker *SMTPWorker) TimeoutHasExpired(startTime time.Time) bool {
	return int(time.Since(startTime).Seconds()) > COMMAND_TIMEOUT_SECONDS
}
