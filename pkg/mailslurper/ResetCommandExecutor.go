// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

/*
ResetCommandExecutor process the command RSET
*/
type ResetCommandExecutor struct {
	logger *logrus.Entry
	writer *SMTPWriter
}

/*
NewResetCommandExecutor creates a new struct
*/
func NewResetCommandExecutor(logger *logrus.Entry, writer *SMTPWriter) *ResetCommandExecutor {
	return &ResetCommandExecutor{
		logger: logger,
		writer: writer,
	}
}

/*
Process handles the RSET command
*/
func (e *ResetCommandExecutor) Process(streamInput string, mailItem *MailItem) error {
	if strings.ToLower(streamInput) != "rset" {
		return fmt.Errorf("Invalid RSET command")
	}

	// Overwrite current mail object with an empty one
	*mailItem = *NewEmptyMailItem(e.logger)

	return e.writer.SendOkResponse()
}
