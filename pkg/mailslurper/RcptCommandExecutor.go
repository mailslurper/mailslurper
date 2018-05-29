// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import (
	"net/mail"

	"github.com/adampresley/webframework/sanitizer"
	"github.com/sirupsen/logrus"
)

/*
RcptCommandExecutor process the RCPT TO command
*/
type RcptCommandExecutor struct {
	emailValidationService EmailValidationProvider
	logger                 *logrus.Entry
	reader                 *SMTPReader
	writer                 *SMTPWriter
	xssService             sanitizer.IXSSServiceProvider
}

/*
NewRcptCommandExecutor creates a new struct
*/
func NewRcptCommandExecutor(
	logger *logrus.Entry,
	reader *SMTPReader,
	writer *SMTPWriter,
	emailValidationService EmailValidationProvider,
	xssService sanitizer.IXSSServiceProvider,
) *RcptCommandExecutor {
	return &RcptCommandExecutor{
		emailValidationService: emailValidationService,
		logger:                 logger,
		reader:                 reader,
		writer:                 writer,
		xssService:             xssService,
	}
}

/*
Process handles the RCPT TO command. This command tells us who
the recipient is
*/
func (e *RcptCommandExecutor) Process(streamInput string, mailItem *MailItem) error {
	var err error
	var to string
	var toComponents *mail.Address

	if err = IsValidCommand(streamInput, "RCPT TO"); err != nil {
		return err
	}

	if to, err = GetCommandValue(streamInput, "RCPT TO", ":"); err != nil {
		return err
	}

	if toComponents, err = e.emailValidationService.GetEmailComponents(to); err != nil {
		return InvalidEmail(to)
	}

	to = e.xssService.SanitizeString(toComponents.Address)

	if !e.emailValidationService.IsValidEmail(to) {
		return InvalidEmail(to)
	}

	mailItem.ToAddresses = append(mailItem.ToAddresses, to)
	e.writer.SendOkResponse()
	return nil
}
