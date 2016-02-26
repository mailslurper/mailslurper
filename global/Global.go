package global

import "github.com/mailslurper/libmailslurper/storage"

const (
	// Version of the MailSlurper Server application
	SERVER_VERSION string = "1.8"
)

var Database storage.IStorage
