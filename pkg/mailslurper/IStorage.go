// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
IStorage defines an interface for structures that need to connect to
storage engines. They store and retrieve data for MailSlurper
*/
type IStorage interface {
	Connect() error
	Disconnect()
	Create() error

	GetAttachment(mailID, attachmentID string) (*Attachment, error)
	GetMailByID(id string) (*MailItem, error)
	GetMailCollection(offset, length int, mailSearch *MailSearch) ([]*MailItem, error)
	GetMailCount(mailSearch *MailSearch) (int, error)

	DeleteMailsAfterDate(startDate string) (int64, error)
	StoreMail(mailItem *MailItem) (string, error)
}
