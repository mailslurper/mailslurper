// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"github.com/sirupsen/logrus"
)

/*
NoopCommandExecutor process the command NOOP
*/
type NoopCommandExecutor struct {
	logger *logrus.Entry
	writer *SMTPWriter
}

/*
NewNoopCommandExecutor creates a new struct
*/
func NewNoopCommandExecutor(logger *logrus.Entry, writer *SMTPWriter) *NoopCommandExecutor {
	return &NoopCommandExecutor{
		logger: logger,
		writer: writer,
	}
}

/*
Process handles the NOOP command
*/
func (e *NoopCommandExecutor) Process(streamInput string, mailItem *MailItem) error {
	var err error

	//validate
	if err = IsValidCommand(streamInput, "NOOP"); err != nil {
		return err
	}

	//log the command, and do nothing
	e.logger.Debugf("NOOP command received")

	//Response: 250 Ok
	return e.writer.SendOkResponse()
}
