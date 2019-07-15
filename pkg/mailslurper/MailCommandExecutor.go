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
MailCommandExecutor process the MAIL FROM
*/
type MailCommandExecutor struct {
	emailValidationService EmailValidationProvider
	logger                 *logrus.Entry
	reader                 *SMTPReader
	writer                 *SMTPWriter
	xssService             sanitizer.IXSSServiceProvider
}

/*
NewMailCommandExecutor creates a new struct
*/
func NewMailCommandExecutor(
	logger *logrus.Entry,
	reader *SMTPReader,
	writer *SMTPWriter,
	emailValidationService EmailValidationProvider,
	xssService sanitizer.IXSSServiceProvider,
) *MailCommandExecutor {
	return &MailCommandExecutor{
		emailValidationService: emailValidationService,
		logger:                 logger,
		reader:                 reader,
		writer:                 writer,
		xssService:             xssService,
	}
}

/*
Process handles the MAIL FROM command. This command tells us who
the sender is
*/
func (e *MailCommandExecutor) Process(streamInput string, mailItem *MailItem) error {
	var err error
	var from string
	var fromComponents *mail.Address

	if err = IsValidCommand(streamInput, "MAIL FROM"); err != nil {
		return err
	}

	if from, err = GetCommandValue(streamInput, "MAIL FROM", ":"); err != nil {
		return err
	}

	// For all we know, <> is a valid email address (RFC 2821, Section 6.1 & 3.7; NULL return path)
	if from != "<>" {
		if fromComponents, err = e.emailValidationService.GetEmailComponents(from); err != nil {
			return InvalidEmail(from)
		}

		from = e.xssService.SanitizeString(fromComponents.Address)

		if !e.emailValidationService.IsValidEmail(from) {
			return InvalidEmail(from)
		}
	}

	mailItem.FromAddress = from
	e.writer.SendOkResponse()

	return nil
}
