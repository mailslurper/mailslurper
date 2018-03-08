// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/adampresley/webframework/sanitizer"
	"github.com/sirupsen/logrus"
)

/*
An SMTPWorker is responsible for executing, parsing, and processing a single
TCP connection's email.
*/
type SMTPWorker struct {
	Connection             net.Conn
	EmailValidationService EmailValidationProvider
	Mail                   *MailItem
	Reader                 *SMTPReader
	Receiver               chan *MailItem
	State                  SMTPWorkerState
	WorkerID               int
	Writer                 *SMTPWriter
	XSSService             sanitizer.IXSSServiceProvider

	connectionCloseChannel chan string
	ctx                    context.Context
	pool                   *ServerPool
	logger                 *logrus.Entry
}

/*
InitializeMailItem initializes the mail item structure that will eventually
be written to the receiver channel.
*/
func (smtpWorker *SMTPWorker) InitializeMailItem() {
	smtpWorker.Mail = &MailItem{}

	smtpWorker.Mail.ToAddresses = NewMailAddressCollection()
	smtpWorker.Mail.Attachments = make([]*Attachment, 0, 5)
	smtpWorker.Mail.Message = NewSMTPMessagePart(smtpWorker.logger)

	/*
	 * IDs are generated ahead of time because
	 * we do not know what order recievers
	 * get the mail item once it is retrieved from the TCP socket.
	 */
	id, _ := GenerateID()
	smtpWorker.Mail.ID = id
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
) *SMTPWorker {
	return &SMTPWorker{
		EmailValidationService: emailValidationService,
		WorkerID:               workerID,
		State:                  SMTP_WORKER_IDLE,
		XSSService:             xssService,

		pool:   pool,
		logger: logger,
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
	ctx context.Context,
	connectionCloseChannel chan string,
) {
	smtpWorker.State = SMTP_WORKER_WORKING

	smtpWorker.Connection = connection
	smtpWorker.Receiver = receiver

	smtpWorker.Reader = reader
	smtpWorker.Writer = writer

	smtpWorker.connectionCloseChannel = connectionCloseChannel
	smtpWorker.ctx = ctx
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
	var streamInput string
	var command SMTPCommand
	var err error

	smtpWorker.InitializeMailItem()
	smtpWorker.Writer.SayHello()

	/*
	 * Read from the connection until we receive a QUIT command
	 * or some critical error occurs and we force quit.
	 */
	startTime := time.Now()

	for smtpWorker.State != SMTP_WORKER_DONE && smtpWorker.State != SMTP_WORKER_ERROR {
		streamInput = smtpWorker.Reader.Read()

		if command, err = GetCommandFromString(streamInput); err != nil {
			smtpWorker.logger.Errorf("Problem finding command from input %s: %s", streamInput, err.Error())
			smtpWorker.State = SMTP_WORKER_ERROR
			continue
		}

		if command == NONE {
			smtpWorker.logger.Debugf("No command...")
			continue
		}

		smtpWorker.logger.Debugf("Command: %s", command.String())

		if command == QUIT {
			smtpWorker.State = SMTP_WORKER_DONE
			smtpWorker.logger.Infof("QUIT command received. Closing connection")
		} else {
			executor := smtpWorker.getExecutorFromCommand(command)
			streamInput = strings.TrimSpace(streamInput)

			if err = executor.Process(streamInput, smtpWorker.Mail); err != nil {
				smtpWorker.State = SMTP_WORKER_ERROR
				smtpWorker.logger.Errorf("Problem executing command %s (stream input == '%s'): %s", command.String(), streamInput, err.Error())
				continue
			}
		}

		if smtpWorker.TimeoutHasExpired(startTime) {
			smtpWorker.logger.Infof("Connection timeout. Terminating client connection")
			smtpWorker.State = SMTP_WORKER_ERROR
			continue
		}
	}

	smtpWorker.Writer.SayGoodbye()
	smtpWorker.Connection.Close()

	if smtpWorker.State != SMTP_WORKER_ERROR {
		smtpWorker.Receiver <- smtpWorker.Mail
	}

	smtpWorker.State = SMTP_WORKER_IDLE
	smtpWorker.rejoinWorkerQueue()
}

func (smtpWorker *SMTPWorker) getExecutorFromCommand(command SMTPCommand) ICommandExecutor {
	switch command {
	case MAIL:
		return NewMailCommandExecutor(smtpWorker.logger, smtpWorker.Reader, smtpWorker.Writer, smtpWorker.EmailValidationService, smtpWorker.XSSService)

	case RCPT:
		return NewRcptCommandExecutor(smtpWorker.logger, smtpWorker.Reader, smtpWorker.Writer, smtpWorker.EmailValidationService, smtpWorker.XSSService)

	case DATA:
		return NewDataCommandExecutor(smtpWorker.logger, smtpWorker.Reader, smtpWorker.Writer, smtpWorker.EmailValidationService, smtpWorker.XSSService)

	default:
		return NewHelloCommandExecutor(smtpWorker.logger, smtpWorker.Reader, smtpWorker.Writer)
	}
}

/*
TimeoutHasExpired determines if the time elapsed since a start time has exceeded
the command timeout.
*/
func (smtpWorker *SMTPWorker) TimeoutHasExpired(startTime time.Time) bool {
	return int(time.Since(startTime).Seconds()) > COMMAND_TIMEOUT_SECONDS
}
