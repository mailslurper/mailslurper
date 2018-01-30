// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "github.com/sirupsen/logrus"

/*
MailItem is a struct describing a parsed mail item. This is
populated after an incoming client connection has finished
sending mail data to this server.
*/
type MailItem struct {
	ID          string                `json:"id"`
	DateSent    string                `json:"dateSent"`
	FromAddress string                `json:"fromAddress"`
	ToAddresses MailAddressCollection `json:"toAddresses"`
	Subject     string                `json:"subject"`
	XMailer     string                `json:"xmailer"`
	MIMEVersion string                `json:"mimeVersion"`
	Body        string                `json:"body"`
	ContentType string                `json:"contentType"`
	Boundary    string                `json:"boundary"`
	Attachments []*Attachment         `json:"attachments"`

	Message           *SMTPMessagePart
	InlineAttachments []*Attachment
	TextBody          string
	HTMLBody          string
}

/*
NewMailItem creates a new MailItem object
*/
func NewMailItem(id, dateSent string, fromAddress string, toAddresses MailAddressCollection, subject, xMailer, body, contentType, boundary string, attachments []*Attachment, logger *logrus.Entry) *MailItem {
	return &MailItem{
		ID:          id,
		DateSent:    dateSent,
		FromAddress: fromAddress,
		ToAddresses: toAddresses,
		Subject:     subject,
		XMailer:     xMailer,
		Body:        body,
		ContentType: contentType,
		Boundary:    boundary,
		Attachments: attachments,

		Message: NewSMTPMessagePart(logger),
	}
}
